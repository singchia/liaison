package frontierbound

import (
	"context"
	"encoding/json"

	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
	"github.com/singchia/liaison/pkg/proto"
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
	defer tx.Rollback()

	deviceModel := &model.Device{
		Fingerprint: device.Fingerprint,
		HostName:    device.HostName,
		CPU:         device.CPU,
		Memory:      device.Memory,
		OS:          device.OS,
		OSVersion:   device.OSVersion,
		CPUUsage:    device.DeviceUsage.CPUUsage,
		MemoryUsage: device.DeviceUsage.MemoryUsage,
		DiskUsage:   device.DeviceUsage.DiskUsage,
	}
	if err := tx.CreateDevice(deviceModel); err != nil {
		rsp.SetError(err)
		return
	}

	for _, iface := range device.Interfaces {
		interfaceModel := &model.EthernetInterface{
			DeviceID: deviceModel.ID,
			Name:     iface.Name,
			MAC:      iface.MAC,
			IP:       iface.IP,
			Netmask:  iface.Netmask,
			Gateway:  iface.Gateway,
		}
		if err := tx.CreateEthernetInterface(interfaceModel); err != nil {
			rsp.SetError(err)
			return
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
}
