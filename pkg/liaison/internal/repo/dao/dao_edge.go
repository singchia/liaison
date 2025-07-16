package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

func (dao *dao) GetEdge(id uint64) (*model.Edge, error) {
	var edge model.Edge
	if err := dao.db.Where("id = ?", id).First(&edge).Error; err != nil {
		return nil, err
	}
	return &edge, nil
}

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

func (dao *dao) ListEdges(page, pageSize int) ([]*model.Edge, error) {
	var edges []*model.Edge
	if err := dao.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&edges).Error; err != nil {
		return nil, err
	}
	return edges, nil
}

func (dao *dao) UpdateEdge(edge *model.Edge) error {
	return dao.db.Save(edge).Error
}

func (dao *dao) DeleteEdge(id uint64) error {
	return dao.db.Delete(&model.Edge{}, id).Error
}
