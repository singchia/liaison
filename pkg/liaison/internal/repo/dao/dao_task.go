package dao

import "github.com/singchia/liaison/pkg/liaison/internal/repo/model"

func (d *dao) CreateTask(task *model.Task) error {
	return d.getDB().Create(task).Error
}

func (d *dao) UpdateTaskStatus(taskID uint, status model.TaskStatus) error {
	return d.getDB().Model(&model.Task{}).Where("id = ?", taskID).Update("task_status", status).Error
}

func (d *dao) UpdateTaskResult(taskID uint, result []byte) error {
	return d.getDB().Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"task_result": result,
		"task_status": model.TaskStatusCompleted,
	}).Error
}

func (d *dao) UpdateTaskError(taskID uint, error string) error {
	return d.getDB().Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"error":       error,
		"task_status": model.TaskStatusFailed,
	}).Error
}
