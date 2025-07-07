package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

func (dao *dao) CreateApplication(application *model.Application) error {
	return dao.db.Create(application).Error
}

func (dao *dao) GetApplicationByID(id uint) (*model.Application, error) {
	var application model.Application
	if err := dao.db.Where("id = ?", id).First(&application).Error; err != nil {
		return nil, err
	}
	return &application, nil
}
