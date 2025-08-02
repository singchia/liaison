package controlplane

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/internal/repo/model"
)

func (cp *controlPlane) CreateEdge(_ context.Context, req *v1.CreateEdgeRequest) (*v1.CreateEdgeResponse, error) {
	// 在事务中创建edge和ak/sk
	tx := cp.repo.Begin()

	edge := &model.Edge{
		Name:        req.Name,
		Description: req.Description,
		Online:      model.EdgeOnlineStatusOffline,
	}

	err := tx.CreateEdge(edge)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 生成 AK/SK
	accessKey, secretKey := generateAccessKeyPair()

	err = tx.CreateAccessKey(&model.AccessKey{
		EdgeID:    edge.ID,
		AccessKey: accessKey,
		SecretKey: secretKey,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// 返回响应
	return &v1.CreateEdgeResponse{
		Code:    200,
		Message: "success",
		Data: &v1.AccessKey{
			AccessKey: accessKey,
			SecretKey: secretKey,
		},
	}, nil
}

func (cp *controlPlane) GetEdge(_ context.Context, req *v1.GetEdgeRequest) (*v1.GetEdgeResponse, error) {
	edge, err := cp.repo.GetEdge(req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetEdgeResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Edge{
			Id:          uint64(edge.ID),
			Name:        edge.Name,
			Description: edge.Description,
			Online:      int32(edge.Online),
			CreatedAt:   edge.CreatedAt.Format(time.DateTime),
			UpdatedAt:   edge.UpdatedAt.Format(time.DateTime),
		},
	}, nil
}

func (cp *controlPlane) ListEdges(_ context.Context, req *v1.ListEdgesRequest) (*v1.ListEdgesResponse, error) {
	edges, err := cp.repo.ListEdges(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &v1.ListEdgesResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Edges{
			Edges: transformEdges(edges),
		},
	}, nil
}

func (cp *controlPlane) UpdateEdge(_ context.Context, req *v1.UpdateEdgeRequest) (*v1.UpdateEdgeResponse, error) {
	edge, err := cp.repo.GetEdge(req.Id)
	if err != nil {
		return nil, err
	}
	edge.Name = req.Name
	edge.Description = req.Description
	err = cp.repo.UpdateEdge(edge)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateEdgeResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (cp *controlPlane) DeleteEdge(_ context.Context, req *v1.DeleteEdgeRequest) (*v1.DeleteEdgeResponse, error) {
	err := cp.repo.DeleteEdge(req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteEdgeResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (cp *controlPlane) CreateEdgeScanApplicationTask(_ context.Context, req *v1.CreateEdgeScanApplicationTaskRequest) (*v1.CreateEdgeScanApplicationTaskResponse, error) {

}

func (cp *controlPlane) GetEdgeScanApplicationTask(_ context.Context, req *v1.GetEdgeScanApplicationTaskRequest) (*v1.GetEdgeScanApplicationTaskResponse, error) {

}

func transformEdges(edges []*model.Edge) []*v1.Edge {
	edgesV1 := make([]*v1.Edge, len(edges))
	for i, edge := range edges {
		edgesV1[i] = transformEdge(edge)
	}
	return edgesV1
}

func transformEdge(edge *model.Edge) *v1.Edge {
	return &v1.Edge{
		Id:          uint64(edge.ID),
		Name:        edge.Name,
		Description: edge.Description,
		Online:      int32(edge.Online),
		CreatedAt:   edge.CreatedAt.Format(time.DateTime),
		UpdatedAt:   edge.UpdatedAt.Format(time.DateTime),
	}
}

// generateAccessKey 生成 Access Key
// 格式: 时间戳 + 随机字符串
func generateAccessKey() string {
	// 获取时间戳
	timestamp := time.Now().UnixNano()

	// 生成随机字节
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	// 组合并编码
	data := fmt.Sprintf("%d%s", timestamp, hex.EncodeToString(randomBytes))
	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(data))

	// 限制长度
	if len(encoded) > 20 {
		encoded = encoded[:20]
	}

	return encoded
}

// generateSecretKey 生成 Secret Key
// 格式: 32字节随机数据，Base64编码
func generateSecretKey() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)
}

// generateAccessKeyPair 生成 Access Key 和 Secret Key 对
func generateAccessKeyPair() (accessKey, secretKey string) {
	return generateAccessKey(), generateSecretKey()
}
