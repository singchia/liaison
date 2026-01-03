package controlplane

import (
	"context"
	"time"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/repo/dao"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
)

func (cp *controlPlane) ListDevices(_ context.Context, req *v1.ListDevicesRequest) (*v1.ListDevicesResponse, error) {
	query := dao.ListDevicesQuery{
		Query: dao.Query{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	}
	devices, err := cp.repo.ListDevices(&query)
	if err != nil {
		return nil, err
	}
	// 获取网卡
	for _, device := range devices {
		interfaces, err := cp.repo.GetEthernetInterfacesByDeviceID(uint(device.ID))
		if err != nil {
			return nil, err
		}
		device.Interfaces = interfaces
	}
	count, err := cp.repo.CountDevices(&query)
	if err != nil {
		return nil, err
	}
	return &v1.ListDevicesResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Devices{
			Total:   int32(count),
			Devices: transformDevices(devices),
		},
	}, nil
}

func (cp *controlPlane) GetDevice(_ context.Context, req *v1.GetDeviceRequest) (*v1.GetDeviceResponse, error) {
	device, err := cp.repo.GetDeviceByID(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &v1.GetDeviceResponse{
		Code:    200,
		Message: "success",
		Data:    transformDevice(device),
	}, nil
}

func (cp *controlPlane) UpdateDevice(_ context.Context, req *v1.UpdateDeviceRequest) (*v1.UpdateDeviceResponse, error) {
	device, err := cp.repo.GetDeviceByID(uint(req.Id))
	if err != nil {
		return nil, err
	}
	device.Name = req.Name
	device.Description = req.Description
	err = cp.repo.UpdateDevice(device)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateDeviceResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func transformDevices(devices []*model.Device) []*v1.Device {
	devicesV1 := make([]*v1.Device, len(devices))
	for i, device := range devices {
		deviceV1 := transformDevice(device)
		deviceV1.Interfaces = transformEthernetInterfaces(device.Interfaces)
		devicesV1[i] = deviceV1
	}
	return devicesV1
}

func transformEthernetInterfaces(interfaces []*model.EthernetInterface) []*v1.EthernetInterface {
	// 每个网卡可能有多个IP地址
	v1ifaces := map[string]*v1.EthernetInterface{}
	for _, iface := range interfaces {
		v1iface, ok := v1ifaces[iface.Name+iface.MAC]
		if !ok {
			v1ifaces[iface.Name+iface.MAC] = &v1.EthernetInterface{
				Name: iface.Name,
				Mac:  iface.MAC,
				Ip:   []string{iface.IP},
			}
		} else {
			v1iface.Ip = append(v1iface.Ip, iface.IP)
		}
	}
	v1ifaceslice := make([]*v1.EthernetInterface, 0, len(v1ifaces))
	for _, v1iface := range v1ifaces {
		v1ifaceslice = append(v1ifaceslice, v1iface)
	}
	return v1ifaceslice
}

func transformDevice(device *model.Device) *v1.Device {
	return &v1.Device{
		Id:          uint64(device.ID),
		Name:        device.Name,
		Description: device.Description,
		Cpu:         int32(device.CPU),
		Memory:      int32(device.Memory),
		Disk:        int32(device.Disk),
		Os:          device.OS,
		Version:     device.OSVersion,
		CreatedAt:   device.CreatedAt.Format(time.DateTime),
		UpdatedAt:   device.UpdatedAt.Format(time.DateTime),
	}
}
