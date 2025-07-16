package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

func (d *dao) CreateDevice(device *model.Device) error {
	return d.getDB().Create(device).Error
}

func (d *dao) GetDeviceByID(id uint) (*model.Device, error) {
	var device model.Device
	if err := d.getDB().Where("id = ?", id).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}
