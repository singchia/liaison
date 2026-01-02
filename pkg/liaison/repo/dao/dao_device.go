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

func (d *dao) GetEthernetInterface(deviceID uint, ip, netmask, name, mac string) (*model.EthernetInterface, error) {
	var iface model.EthernetInterface
	if err := d.getDB().Where("device_id = ? AND ip = ? AND netmask = ? AND name = ? AND mac = ?",
		deviceID, ip, netmask, name, mac).First(&iface).Error; err != nil {
		return nil, err
	}
	return &iface, nil
}

func (d *dao) GetEthernetInterfacesByDeviceID(deviceID uint) ([]*model.EthernetInterface, error) {
	var interfaces []*model.EthernetInterface
	if err := d.getDB().Where("device_id = ?", deviceID).Find(&interfaces).Error; err != nil {
		return nil, err
	}
	return interfaces, nil
}

func (d *dao) UpdateEthernetInterface(iface *model.EthernetInterface) error {
	return d.getDB().Model(&model.EthernetInterface{}).Where("id = ?", iface.ID).Updates(map[string]interface{}{
		"name":    iface.Name,
		"mac":     iface.MAC,
		"ip":      iface.IP,
		"netmask": iface.Netmask,
	}).Error
}

func (d *dao) DeleteEthernetInterface(id uint) error {
	return d.getDB().Delete(&model.EthernetInterface{}, id).Error
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

func (d *dao) GetDeviceByFingerprint(fingerprint string) (*model.Device, error) {
	var device model.Device
	if err := d.getDB().Where("fingerprint = ?", fingerprint).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (d *dao) CountDevices(query *ListDevicesQuery) (int64, error) {
	var count int64
	db := d.getDB()
	if len(query.IDs) > 0 {
		db = db.Where("id IN ?", query.IDs)
	}
	if err := db.Model(&model.Device{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *dao) ListDevices(query *ListDevicesQuery) ([]*model.Device, error) {
	db := d.getDB()
	if query.Page > 0 && query.PageSize > 0 {
		db = db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
	}
	if len(query.IDs) > 0 {
		db = db.Where("id IN ?", query.IDs)
	}
	var devices []*model.Device
	if err := db.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (d *dao) UpdateDevice(device *model.Device) error {
	updates := map[string]interface{}{}
	if device.Name != "" {
		updates["name"] = device.Name
	}
	if device.Description != "" {
		updates["description"] = device.Description
	}
	if device.CPU > 0 {
		updates["cpu"] = device.CPU
	}
	if device.Memory > 0 {
		updates["memory"] = device.Memory
	}
	if device.Disk > 0 {
		updates["disk"] = device.Disk
	}
	if device.OS != "" {
		updates["os"] = device.OS
	}
	if device.OSVersion != "" {
		updates["os_version"] = device.OSVersion
	}
	if device.HostName != "" {
		updates["host_name"] = device.HostName
	}
	if device.CPUUsage > 0 {
		updates["cpu_usage"] = device.CPUUsage
	}
	if device.MemoryUsage > 0 {
		updates["memory_usage"] = device.MemoryUsage
	}
	if device.DiskUsage > 0 {
		updates["disk_usage"] = device.DiskUsage
	}
	if len(updates) == 0 {
		return nil
	}
	return d.getDB().Model(&model.Device{}).Where("id = ?", device.ID).Updates(updates).Error
}

func (d *dao) UpdateDeviceUsage(deviceID uint, cpuUsage, memoryUsage, diskUsage float32) error {
	return d.getDB().Model(&model.Device{}).Where("id = ?", deviceID).Omit("updated_at").Updates(map[string]interface{}{
		"cpu_usage":    cpuUsage,
		"memory_usage": memoryUsage,
		"disk_usage":   diskUsage,
		"heartbeat_at": time.Now(),
	}).Error
}
