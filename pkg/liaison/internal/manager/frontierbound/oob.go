package frontierbound

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/singchia/liaison/pkg/liaison/internal/repo/model"
	"github.com/singchia/liaison/pkg/proto"
)

// 获取设备ID
func (fb *frontierBound) getID(meta []byte) (uint64, error) {
	var m proto.Meta
	if err := json.Unmarshal(meta, &m); err != nil {
		return 0, err
	}
	ak, edge, err := fb.repo.GetEdgeByAccessKey(m.AccessKey)
	if err != nil {
		return 0, err
	}
	if ak.SecretKey != m.SecretKey {
		return 0, errors.New("invalid secret key")
	}

	return uint64(edge.ID), nil
}

func (fb *frontierBound) online(edgeID uint64, meta []byte, addr net.Addr) error {
	err := fb.repo.UpdateEdgeOnlineStatus(edgeID, model.EdgeOnlineStatusOnline)
	if err != nil {
		return err
	}
	return nil
}

func (fb *frontierBound) offline(edgeID uint64, meta []byte, addr net.Addr) error {
	err := fb.repo.UpdateEdgeOnlineStatus(edgeID, model.EdgeOnlineStatusOffline)
	if err != nil {
		return err
	}
	return nil
}
