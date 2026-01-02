package dao

import (
	"time"

	"github.com/singchia/liaison/pkg/liaison/repo/model"
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

func (d *dao) CountEdges() (int64, error) {
	var count int64
	if err := d.getDB().Model(&model.Edge{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *dao) ListEdges(page, pageSize int) ([]*model.Edge, error) {
	db := d.getDB()
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize)
	}
	var edges []*model.Edge
	if err := db.Find(&edges).Error; err != nil {
		return nil, err
	}
	return edges, nil
}

// 更新Name Description Status
func (d *dao) UpdateEdge(edge *model.Edge) error {
	updates := make(map[string]interface{})
	if edge.Name != "" {
		updates["name"] = edge.Name
	}
	if edge.Description != "" {
		updates["description"] = edge.Description
	}
	if edge.Status != 0 {
		updates["status"] = edge.Status
	}
	if len(updates) > 0 {
		if err := d.getDB().Model(&model.Edge{}).Where("id = ?", edge.ID).Updates(updates).Error; err != nil {
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

func (d *dao) UpdateEdgeDeviceID(edgeID uint64, deviceID uint) error {
	return d.getDB().Model(&model.Edge{}).Where("id = ?", edgeID).Update("device_id", deviceID).Error
}

func (d *dao) DeleteEdge(id uint64) error {
	return d.getDB().Delete(&model.Edge{}, id).Error
}
