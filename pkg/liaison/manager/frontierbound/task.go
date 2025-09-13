package frontierbound

import (
	"context"
	"encoding/json"

	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
	"github.com/singchia/liaison/pkg/proto"
)

type Net struct {
	Nets     []string
	Protocol string
	Port     int
}

func (fb *frontierBound) EmitScanApplications(ctx context.Context, taskID uint, edgeID uint64, net *Net) error {
	request := &proto.ScanApplicationTaskRequest{
		TaskID:   taskID,
		Nets:     net.Nets,
		Protocol: net.Protocol,
		Port:     net.Port,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req := fb.svc.NewRequest(data)
	fb.svc.Call(ctx, edgeID, "scan_application", req)
	return nil
}

func (fb *frontierBound) reportTaskScanApplication(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	var task proto.ScanApplicationTaskResult
	err := json.Unmarshal(req.Data(), &task)
	if err != nil {
		rsp.SetError(err)
		return
	}

	switch task.Status {
	case "running":
		data, err := json.Marshal(task.ScannedApplications)
		if err != nil {
			rsp.SetError(err)
			return
		}
		err = fb.repo.UpdateTaskResult(task.TaskID, model.TaskStatusRunning, data)
		if err != nil {
			rsp.SetError(err)
			return
		}
	case "completed":
		data, err := json.Marshal(task.ScannedApplications)
		if err != nil {
			rsp.SetError(err)
			return
		}
		err = fb.repo.UpdateTaskResult(task.TaskID, model.TaskStatusCompleted, data)
		if err != nil {
			rsp.SetError(err)
			return
		}
	case "failed":
		err = fb.repo.UpdateTaskError(task.TaskID, task.Error)
		if err != nil {
			rsp.SetError(err)
			return
		}
	}
}
