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

// @Summary CreateEdge
// @Tags 1.0
// @Param params query v1.CreateEdgeRequest true "queriies"
// @Success 200 {object} v1.CreateEdgeResponse
// @Router /api/v1/edge [post]
func (cp *controlPlane) CreateEdge(ctx context.Context, req *v1.CreateEdgeRequest) (*v1.CreateEdgeResponse, error) {
	// 在事务中创建edge和ak/sk
	tx := cp.repo.Begin()

	edge := &model.Edge{
		Name:        req.Name,
		Status:      model.EdgeStatusOffline,
		Description: req.Description,
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

// @Summary GetEdge
// @Tags 1.0
// @Param id query int true "edge id"
// @Success 200 {object} v1.GetEdgeResponse
// @Router /api/v1/edge/{id} [get]
func (cp *controlPlane) GetEdge(ctx context.Context, req *v1.GetEdgeRequest) (*v1.GetEdgeResponse, error) {
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
			Status:      int32(edge.Status),
			CreatedAt:   edge.CreatedAt.Format(time.DateTime),
			UpdatedAt:   edge.UpdatedAt.Format(time.DateTime),
		},
	}, nil
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
