package dao

import "github.com/singchia/liaison/pkg/liaison/repo/model"

func (d *dao) CreateApplication(application *model.Application) error {
	return d.getDB().Create(application).Error
}

func (d *dao) GetApplicationByID(id uint) (*model.Application, error) {
	var application model.Application
	if err := d.getDB().Where("id = ?", id).First(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (d *dao) DeleteApplication(id uint) error {
	return d.getDB().Delete(&model.Application{}, id).Error
}

func (d *dao) UpdateApplication(application *model.Application) error {
	return d.getDB().Save(application).Error
}

func (d *dao) ListApplications(query *ListApplicationsQuery) ([]*model.Application, error) {
	var applications []*model.Application
	db := d.getDB()
	// device_id
	if query.DeviceID != 0 {
		db = db.Where("device_id = ?", query.DeviceID)
	}
	// page & page_size
	if query.Page > 0 && query.PageSize > 0 {
		db = db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
	}

	if err := db.Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

type ListApplicationsQuery struct {
	Page     int
	PageSize int
	DeviceID uint
}
