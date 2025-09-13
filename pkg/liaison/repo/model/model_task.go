package model

import "gorm.io/gorm"

type TaskType int

const (
	TaskTypeScan TaskType = iota + 1
)

type TaskSubType int

const (
	TaskSubTypeScanApplication TaskSubType = iota + 1
)

type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota + 1
	TaskStatusRunning
	TaskStatusCompleted
	TaskStatusFailed
)

type Task struct {
	gorm.Model
	EdgeID      uint64      `gorm:"column:edge_id;type:bigint;not null"` // edge的任务
	TaskType    TaskType    `gorm:"column:task_type;type:int;not null"`
	TaskSubType TaskSubType `gorm:"column:task_sub_type;type:int;not null"`
	TaskStatus  TaskStatus  `gorm:"column:task_status;type:int;not null"`
	TaskParams  []byte      `gorm:"column:task_params;type:blob;not null"`
	TaskResult  []byte      `gorm:"column:task_result;type:blob;not null"`
	Error       string      `gorm:"column:error;type:text;not null"`
}

func (Task) TableName() string {
	return "task"
}

// task scan application
type TaskScanApplicationParams struct {
	Network  string `json:"network"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type ScannedApplication struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type TaskScanApplicationResult struct {
	ScannedApplications []ScannedApplication `json:"scanned_applications"`
}
