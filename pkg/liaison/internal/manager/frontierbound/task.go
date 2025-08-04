package frontierbound

import (
	"context"
	"encoding/json"

	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/proto"
)

func (fb *frontierBound) reportTask(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	var task proto.Task
	if err := json.Unmarshal(req.Data(), &task); err != nil {
		rsp.SetError(err)
		return
	}

	err := fb.repo.UpdateTaskResult(task.ID, task.TaskResult)
	if err != nil {
		rsp.SetError(err)
		return
	}

	rsp.SetData(req.Data())
}
