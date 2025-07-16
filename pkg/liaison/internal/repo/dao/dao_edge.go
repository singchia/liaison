package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

func (d *dao) GetEdge(id uint64) (*model.Edge, error) {
	var edge model.Edge
	if err := d.getDB().Where("id = ?", id).First(&edge).Error; err != nil {
		return nil, err
	}
	return &edge, nil
}

func (d *dao) CreateEdge(edge *model.Edge) error {
	return d.getDB().Create(edge).Error
}

func (d *dao) GetEdgeByDeviceID(deviceID uint) (*model.Edge, error) {
	var edge model.Edge
	if err := d.getDB().Where("device_id = ?", deviceID).First(&edge).Error; err != nil {
		return nil, err
	}
	return &edge, nil
}

func (d *dao) ListEdges(page, pageSize int) ([]*model.Edge, error) {
	var edges []*model.Edge
	if err := d.getDB().Offset((page - 1) * pageSize).Limit(pageSize).Find(&edges).Error; err != nil {
		return nil, err
	}
	return edges, nil
}

func (d *dao) UpdateEdge(edge *model.Edge) error {
	return d.getDB().Save(edge).Error
}

func (d *dao) DeleteEdge(id uint64) error {
	return d.getDB().Delete(&model.Edge{}, id).Error
}
