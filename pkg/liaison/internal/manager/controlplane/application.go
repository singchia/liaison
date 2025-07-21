package controlplane

import (
	"context"
	"time"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/internal/repo/dao"
	"github.com/singchia/liaison/pkg/liaison/internal/repo/model"
)

func (cp *controlPlane) ListApplications(_ context.Context, req *v1.ListApplicationsRequest) (*v1.ListApplicationsResponse, error) {
	applications, err := cp.repo.ListApplications(&dao.ListApplicationsQuery{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		DeviceID: uint(req.DeviceId),
	})
	if err != nil {
		return nil, err
	}
	return &v1.ListApplicationsResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Applications{
			Applications: transformApplications(applications),
		},
	}, nil
}

func (cp *controlPlane) UpdateApplication(_ context.Context, req *v1.UpdateApplicationRequest) (*v1.UpdateApplicationResponse, error) {
	application, err := cp.repo.GetApplicationByID(uint(req.Id))
	if err != nil {
		return nil, err
	}
	application.Name = req.Name
	err = cp.repo.UpdateApplication(application)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateApplicationResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Application{
			Id:        uint64(application.ID),
			Name:      application.Name,
			CreatedAt: application.CreatedAt.Format(time.DateTime),
			UpdatedAt: application.UpdatedAt.Format(time.DateTime),
		},
	}, nil
}

func transformApplications(applications []*model.Application) []*v1.Application {
	applicationsV1 := make([]*v1.Application, len(applications))
	for i, application := range applications {
		applicationsV1[i] = transformApplication(application)
	}
	return applicationsV1
}

func transformApplication(application *model.Application) *v1.Application {
	return &v1.Application{
		Id:        uint64(application.ID),
		Name:      application.Name,
		CreatedAt: application.CreatedAt.Format(time.DateTime),
		UpdatedAt: application.UpdatedAt.Format(time.DateTime),
	}
}
