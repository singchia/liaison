package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

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
