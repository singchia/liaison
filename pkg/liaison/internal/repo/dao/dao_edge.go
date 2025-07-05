package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

func (dao *dao) CreateEdge(edge *model.Edge) error {
	return dao.db.Create(edge).Error
}

func (dao *dao) GetEdgeByDeviceID(deviceID uint) (*model.Edge, error) {
	var edge model.Edge
	if err := dao.db.Where("device_id = ?", deviceID).First(&edge).Error; err != nil {
		return nil, err
	}
	return &edge, nil
}
