package dao

import "github.com/singchia/liaison/pkg/liaison/repo/model"

type Query struct {
	// Pagination
	Page, PageSize int
	// Time range
	StartTime, EndTime int64
	// Order
	Order string
	Desc  bool
}

type ListTasksQuery struct {
	Query
	EdgeID      uint
	Status      []model.TaskStatus
	TaskType    model.TaskType
	TaskSubType model.TaskSubType
}

type ListApplicationsQuery struct {
	Query
	DeviceIDs []uint
	IDs       []uint
}

type ListDevicesQuery struct {
	Query
	IDs  []uint
	Name string
	IP   string
}

type ListProxiesQuery struct {
	Query
	IDs            []uint
	ApplicationIDs []uint
}

type ListEdgesQuery struct {
	Query
	DeviceIDs []uint
}
