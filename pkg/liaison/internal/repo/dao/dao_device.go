package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

func (dao *dao) CreateDevice(device *model.Device) error {
	return dao.db.Create(device).Error
}

func (dao *dao) GetDeviceByID(id uint) (*model.Device, error) {
	var device model.Device
	if err := dao.db.Where("id = ?", id).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}
