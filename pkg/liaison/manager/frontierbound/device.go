package frontierbound

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
	"github.com/singchia/liaison/pkg/proto"
	"gorm.io/gorm"
)

// 上报设备使用情况
func (fb *frontierBound) reportDeviceUsage(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	var usage proto.DeviceUsage
	if err := json.Unmarshal(req.Data(), &usage); err != nil {
		rsp.SetError(err)
		return
	}

	// 客户端ID就是EdgeID
	err := fb.repo.UpdateDeviceUsage(uint(req.ClientID()), usage.CPUUsage, usage.MemoryUsage, usage.DiskUsage)
	if err != nil {
		rsp.SetError(err)
		return
	}
}

// 上报设备
func (fb *frontierBound) reportDevice(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	var device proto.Device
	if err := json.Unmarshal(req.Data(), &device); err != nil {
		rsp.SetError(err)
		return
	}

	tx := fb.repo.Begin()
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 根据设备指纹查找设备
	deviceModel, err := tx.GetDeviceByFingerprint(device.Fingerprint)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("get device by fingerprint error: %s", err)
		rsp.SetError(err)
		return
	}

	if deviceModel == nil {
		deviceModel = &model.Device{
			Name:        device.HostName,
			Fingerprint: device.Fingerprint,
			HostName:    device.HostName,
			CPU:         device.CPU,
			Memory:      device.Memory,
			Disk:        device.Disk,
			OS:          device.OS,
			OSVersion:   device.OSVersion,
			CPUUsage:    device.DeviceUsage.CPUUsage,
			MemoryUsage: device.DeviceUsage.MemoryUsage,
			DiskUsage:   device.DeviceUsage.DiskUsage,
		}
		if err := tx.CreateDevice(deviceModel); err != nil {
			log.Errorf("create device error: %s", err)
			rsp.SetError(err)
			return
		}
	} else {
		updates := &model.Device{
			Fingerprint: device.Fingerprint,
			HostName:    device.HostName,
			CPU:         device.CPU,
			Memory:      device.Memory,
			Disk:        device.Disk,
			OS:          device.OS,
			OSVersion:   device.OSVersion,
			CPUUsage:    device.DeviceUsage.CPUUsage,
			MemoryUsage: device.DeviceUsage.MemoryUsage,
			DiskUsage:   device.DeviceUsage.DiskUsage,
		}
		updates.ID = deviceModel.ID
		if deviceModel.Name == "" {
			updates.Name = device.HostName
		}
		if err := tx.UpdateDevice(updates); err != nil {
			log.Errorf("update device error: %s", err)
			rsp.SetError(err)
			return
		}
	}

	// 收集上报的所有接口（使用 ip+netmask+name+mac 作为唯一标识）
	type interfaceKey struct {
		IP      string
		Netmask string
		Name    string
		MAC     string
	}
	reportedInterfaces := make(map[interfaceKey]bool)
	for _, iface := range device.Interfaces {
		for _, ipMask := range iface.IPMasks {
			key := interfaceKey{
				IP:      ipMask.IP,
				Netmask: ipMask.Netmask,
				Name:    iface.Name,
				MAC:     iface.MAC,
			}
			reportedInterfaces[key] = true
		}
	}

	// 创建或更新接口
	for _, iface := range device.Interfaces {
		for _, ipMask := range iface.IPMasks {
			// 查找是否已存在该接口（根据 device_id、ip、netmask、name、mac 唯一标识）
			existingIface, err := tx.GetEthernetInterface(deviceModel.ID, ipMask.IP, ipMask.Netmask, iface.Name, iface.MAC)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Errorf("get ethernet interface error: %s", err)
				rsp.SetError(err)
				return
			}

			interfaceModel := &model.EthernetInterface{
				DeviceID: deviceModel.ID,
				Name:     iface.Name,
				MAC:      iface.MAC,
				IP:       ipMask.IP,
				Netmask:  ipMask.Netmask,
			}

			if existingIface == nil {
				// 不存在，创建
				if err := tx.CreateEthernetInterface(interfaceModel); err != nil {
					log.Errorf("create ethernet interface error: %s", err)
					rsp.SetError(err)
					return
				}
			} else {
				// 存在，更新（虽然字段都相同，但可能其他字段需要更新）
				interfaceModel.ID = existingIface.ID
				if err := tx.UpdateEthernetInterface(interfaceModel); err != nil {
					log.Errorf("update ethernet interface error: %s", err)
					rsp.SetError(err)
					return
				}
			}
		}
	}

	// 删除不在上报列表中的接口
	existingInterfaces, err := tx.GetEthernetInterfacesByDeviceID(deviceModel.ID)
	if err != nil {
		log.Errorf("get ethernet interfaces error: %s", err)
		rsp.SetError(err)
		return
	}
	for _, existingIface := range existingInterfaces {
		key := interfaceKey{
			IP:      existingIface.IP,
			Netmask: existingIface.Netmask,
			Name:    existingIface.Name,
			MAC:     existingIface.MAC,
		}
		if !reportedInterfaces[key] {
			// 该接口不在上报列表中，删除
			if err := tx.DeleteEthernetInterface(existingIface.ID); err != nil {
				log.Errorf("delete ethernet interface error: %s", err)
				rsp.SetError(err)
				return
			}
		}
	}

	// 如果关联到edge，那么更新edge的device id
	if device.EdgeID > 0 {
		if err := tx.UpdateEdgeDeviceID(device.EdgeID, deviceModel.ID); err != nil {
			rsp.SetError(err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		rsp.SetError(err)
		return
	}
	committed = true
}
