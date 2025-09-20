package dao

import (
	"time"

	"github.com/singchia/liaison/pkg/liaison/repo/model"
)

func (d *dao) CreateDevice(device *model.Device) error {
	return d.getDB().Create(device).Error
}

func (d *dao) CreateEthernetInterface(iface *model.EthernetInterface) error {
	return d.getDB().Create(iface).Error
}

func (d *dao) GetDeviceByID(id uint) (*model.Device, error) {
	var device model.Device
	if err := d.getDB().Where("id = ?", id).First(&device).Error; err != nil {
		return nil, err
	}
	if err := d.getDB().Where("device_id = ?", id).Find(&device.Interfaces).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (d *dao) CountDevices() (int64, error) {
	var count int64
	if err := d.getDB().Model(&model.Device{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *dao) ListDevices(page, pageSize int) ([]*model.Device, error) {
	db := d.getDB()
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	var devices []*model.Device
	if err := db.Find(&devices).Error; err != nil {
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

func (d *dao) UpdateDeviceUsage(deviceID uint, cpuUsage, memoryUsage, diskUsage float32) error {
	return d.getDB().Model(&model.Device{}).Where("id = ?", deviceID).Omit("updated_at").Updates(map[string]interface{}{
		"cpu_usage":    cpuUsage,
		"memory_usage": memoryUsage,
		"disk_usage":   diskUsage,
		"heartbeat_at": time.Now(),
	}).Error
}
