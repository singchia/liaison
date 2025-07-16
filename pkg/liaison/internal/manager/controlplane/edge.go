package controlplane

import (
	"context"

	v1 "github.com/singchia/liaison/api/v1"
)

func (cp *controlPlane) CreateEdge(ctx context.Context, req *v1.CreateEdgeRequest) (*v1.CreateEdgeResponse, error) {
	// 在事务中创建edge和ak/sk
	return nil, nil
}
