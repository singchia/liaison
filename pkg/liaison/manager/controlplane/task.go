package controlplane

import "github.com/singchia/liaison/pkg/liaison/repo/model"

func (cp *controlPlane) checkTask() {
	// 查询所有pending和running状态的任务
	tasks, err := cp.repo.ListTasks(model.TaskStatusPending, model.TaskStatusRunning)
	if err != nil {
		return
	}
	for _, task := range tasks {
		cp.repo.UpdateTaskError(task.ID, "task expired")
	}
}
