package frontierbound

import (
	"context"

	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/liaison/internal/repo/model"
)

// 上报edge，更改初始化到Running状态
func (fb *frontierBound) reportEdge(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	tx := fb.repo.Begin()
	defer tx.Rollback()

	edge, err := tx.GetEdge(uint64(req.ClientID()))
	if err != nil {
		rsp.SetError(err)
		return
	}
	if edge.Status == model.EdgeStatusIniting {
		err := fb.repo.UpdateEdgeStatus(uint64(req.ClientID()), model.EdgeStatusRunning)
		if err != nil {
			rsp.SetError(err)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		rsp.SetError(err)
		return
	}
	return
}
