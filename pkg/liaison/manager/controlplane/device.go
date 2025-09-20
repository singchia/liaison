package controlplane

import (
	"context"
	"time"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
)

func (cp *controlPlane) ListDevices(_ context.Context, req *v1.ListDevicesRequest) (*v1.ListDevicesResponse, error) {
	devices, err := cp.repo.ListDevices(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	count, err := cp.repo.CountDevices()
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
		devicesV1[i] = transformDevice(device)
	}
	return devicesV1
}

func transformDevice(device *model.Device) *v1.Device {
	return &v1.Device{
		Id:          uint64(device.ID),
		Name:        device.Name,
		Description: device.Description,
		CreatedAt:   device.CreatedAt.Format(time.DateTime),
		UpdatedAt:   device.UpdatedAt.Format(time.DateTime),
	}
}
