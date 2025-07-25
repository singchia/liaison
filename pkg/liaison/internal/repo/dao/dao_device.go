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

func (d *dao) ListDevices(page, pageSize int) ([]*model.Device, error) {
	var devices []*model.Device
	if err := d.getDB().Offset((page - 1) * pageSize).Limit(pageSize).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (d *dao) UpdateDevice(device *model.Device) error {
	if device.Name != "" {
		if err := d.getDB().Model(&model.Device{}).Where("id = ?", device.ID).Update("name", device.Name).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *dao) UpdateDeviceUsage(deviceID uint, cpuUsage, memoryUsage, diskUsage int) error {
	return d.getDB().Model(&model.Device{}).Where("id = ?", deviceID).Omit("updated_at").Updates(map[string]interface{}{
		"cpu_usage":    cpuUsage,
		"memory_usage": memoryUsage,
		"disk_usage":   diskUsage,
	}).Error
}
