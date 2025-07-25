package dao

import (
	"time"

	"github.com/singchia/liaison/pkg/liaison/internal/repo/model"
)

func (d *dao) GetEdge(id uint64) (*model.Edge, error) {
	var edge model.Edge
	if err := d.getDB().Where("id = ?", id).First(&edge).Error; err != nil {
		return nil, err
	}
	return &edge, nil
}

func (d *dao) GetEdgeByAccessKey(accessKey string) (*model.AccessKey, *model.Edge, error) {
	var ak model.AccessKey
	if err := d.getDB().Where("access_key = ?", accessKey).First(&ak).Error; err != nil {
		return nil, nil, err
	}
	edge, err := d.GetEdge(uint64(ak.EdgeID))
	if err != nil {
		return nil, nil, err
	}
	return &ak, edge, nil
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
	if edge.Name != "" {
		if err := d.getDB().Model(&model.Edge{}).Where("id = ?", edge.ID).Update("name", edge.Name).Error; err != nil {
			return err
		}
	}
	if edge.Description != "" {
		if err := d.getDB().Model(&model.Edge{}).Where("id = ?", edge.ID).Update("description", edge.Description).Error; err != nil {
			return err
		}
	}
	if edge.Status != 0 {
		if err := d.getDB().Model(&model.Edge{}).Where("id = ?", edge.ID).Update("status", edge.Status).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *dao) UpdateEdgeOnlineStatus(edgeID uint64, onlineStatus model.EdgeOnlineStatus) error {
	return d.getDB().Model(&model.Edge{}).Where("id = ?", edgeID).Update("online", onlineStatus).Error
}

func (d *dao) UpdateEdgeHeartbeatAt(edgeID uint64, heartbeatAt time.Time) error {
	return d.getDB().Model(&model.Edge{}).Where("id = ?", edgeID).Update("heartbeat_at", heartbeatAt).Error
}

func (d *dao) DeleteEdge(id uint64) error {
	return d.getDB().Delete(&model.Edge{}, id).Error
}
