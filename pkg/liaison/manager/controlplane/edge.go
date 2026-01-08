package controlplane

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jumboframes/armorigo/log"
	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/manager/frontierbound"
	"github.com/singchia/liaison/pkg/liaison/repo/dao"
	"github.com/singchia/liaison/pkg/liaison/repo/model"
)

func (cp *controlPlane) CreateEdge(_ context.Context, req *v1.CreateEdgeRequest) (*v1.CreateEdgeResponse, error) {
	// 在事务中创建edge和ak/sk
	tx := cp.repo.Begin()

	edge := &model.Edge{
		Name:        req.Name,
		Description: req.Description,
		Status:      model.EdgeStatusRunning, // 默认状态为运行中
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

	// 生成安装命令
	serverURL := cp.conf.Manager.ServerURL
	serverAddr := ""
	if serverURL == "" {
		// 如果没有配置，从 Listen 地址生成
		listen := cp.conf.Manager.Listen
		if listen.TLS.Enable {
			serverURL = fmt.Sprintf("https://%s", listen.Addr)
		} else {
			serverURL = fmt.Sprintf("http://%s", listen.Addr)
		}
		serverAddr = listen.Addr
	} else {
		// 从 serverURL 中提取地址（移除 http:// 或 https:// 前缀）
		if strings.HasPrefix(serverURL, "https://") {
			serverAddr = strings.TrimPrefix(serverURL, "https://")
		} else if strings.HasPrefix(serverURL, "http://") {
			serverAddr = strings.TrimPrefix(serverURL, "http://")
		} else {
			serverAddr = serverURL
		}
	}
	installCommand := fmt.Sprintf("curl -sSL %s/install.sh | bash -s -- --access-key=%s --secret-key=%s --server-addr=%s",
		serverURL, accessKey, secretKey, serverAddr)

	// 返回响应
	return &v1.CreateEdgeResponse{
		Code:    200,
		Message: "success",
		Data: &v1.AccessKey{
			AccessKey: accessKey,
			SecretKey: secretKey,
			Command:   installCommand,
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
			Status:      int32(edge.Status),
			Online:      int32(edge.Online),
			CreatedAt:   edge.CreatedAt.Format(time.DateTime),
			UpdatedAt:   edge.UpdatedAt.Format(time.DateTime),
		},
	}, nil
}

func (cp *controlPlane) ListEdges(_ context.Context, req *v1.ListEdgesRequest) (*v1.ListEdgesResponse, error) {
	var (
		deviceIDs       []uint
		devices         []*model.Device
		err             error
		preDeviceSearch bool
	)
	if req.DeviceName != "" {
		devices, err = cp.repo.ListDevices(&dao.ListDevicesQuery{
			Query: dao.Query{
				Order: "id",
				Desc:  true,
			},
			Name: req.DeviceName,
		})
		if err != nil {
			return nil, err
		}
		for _, device := range devices {
			deviceIDs = append(deviceIDs, device.ID)
		}
		preDeviceSearch = true
	}

	query := &dao.ListEdgesQuery{
		Query: dao.Query{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
			Order:    "id",
			Desc:     true,
		},
	}
	if len(deviceIDs) > 0 {
		query.DeviceIDs = deviceIDs
	}
	// 搜索Edges
	edges, err := cp.repo.ListEdges(query)
	if err != nil {
		return nil, err
	}
	if !preDeviceSearch {
		// 如果没有提前搜索设备， 则后置关联devices
		deviceIDs := []uint{}
		for _, edge := range edges {
			deviceIDs = append(deviceIDs, edge.DeviceID)
		}
		devices, err = cp.repo.ListDevices(&dao.ListDevicesQuery{
			Query: dao.Query{
				Order: "id",
				Desc:  true,
			},
			IDs: deviceIDs,
		})
		if err != nil {
			return nil, err
		}
	}
	// 创建设备映射，通过 ID 匹配
	deviceMap := make(map[uint]*model.Device)
	for _, device := range devices {
		deviceMap[device.ID] = device
	}
	// 关联设备到 edge
	for _, edge := range edges {
		if device, ok := deviceMap[edge.DeviceID]; ok {
			edge.Device = device
		}
	}

	count, err := cp.repo.CountEdges(query)
	if err != nil {
		return nil, err
	}
	return &v1.ListEdgesResponse{
		Code:    200,
		Message: "success",
		Data: &v1.Edges{
			Total: int32(count),
			Edges: transformEdges(edges),
		},
	}, nil
}

func (cp *controlPlane) UpdateEdge(_ context.Context, req *v1.UpdateEdgeRequest) (*v1.UpdateEdgeResponse, error) {
	edge, err := cp.repo.GetEdge(req.Id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		edge.Name = req.Name
	}
	if req.Description != "" {
		edge.Description = req.Description
	}
	if req.Status != 0 {
		edge.Status = model.EdgeStatus(req.Status)
	}
	err = cp.repo.UpdateEdge(edge)
	if err != nil {
		return nil, err
	}
	// 重新获取更新后的 edge 以返回完整数据
	updatedEdge, err := cp.repo.GetEdge(req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateEdgeResponse{
		Code:    200,
		Message: "success",
		Data:    transformEdge(updatedEdge),
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
	// 获取edge
	edge, err := cp.repo.GetEdge(req.EdgeId)
	if err != nil {
		log.Errorf("get edge error: %s", err)
		return nil, err
	}
	if edge.Online != model.EdgeOnlineStatusOnline {
		log.Errorf("edge is not online")
		return nil, errors.New("edge is not online")
	}

	// 获取设备
	device, err := cp.repo.GetDeviceByID(uint(edge.DeviceID))
	if err != nil {
		log.Errorf("get device error: %s", err)
		return nil, err
	}
	nets := []string{}
	for _, iface := range device.Interfaces {
		// 获取IP地址所在网段
		//ip := net.ParseIP(iface.IP)
		//if ip == nil {
		//	continue
		//}
		//nets = append(nets, fmt.Sprintf("%s/%s", ip.Mask(net.IPMask(iface.Netmask)), iface.Netmask))
		nets = append(nets, iface.IP)
	}

	// 创建任务
	params := model.TaskScanApplicationParams{
		Nets:     nets,
		Port:     int(req.Port),
		Protocol: req.Protocol,
	}
	data, err := json.Marshal(params)
	if err != nil {
		log.Errorf("marshal params error: %s", err)
		return nil, err
	}
	expiration := 10 * time.Minute
	task := &model.Task{
		EdgeID:      req.EdgeId,
		TaskType:    model.TaskTypeScan,
		TaskSubType: model.TaskSubTypeScanApplication,
		TaskStatus:  model.TaskStatusPending,
		ExpiredAt:   time.Now().Add(expiration), // 10分钟过期
		TaskParams:  data,
		TaskResult:  []byte(`{"scanned_applications":[]}`), // 初始化为空结果
		Error:       "",                                    // 初始化为空错误
	}
	err = cp.repo.CreateTask(task)
	if err != nil {
		log.Errorf("create task error: %s", err)
		return nil, err
	}
	go func() {
		time.Sleep(expiration)
		task, err := cp.repo.GetTask(task.ID)
		if err != nil {
			log.Errorf("get task error: %s", err)
			return
		}
		switch task.TaskStatus {
		case model.TaskStatusPending:
			err = cp.repo.UpdateTaskError(task.ID, "task expired")
			if err != nil {
				log.Errorf("update task error: %s", err)
			}
		case model.TaskStatusRunning:
			err = cp.repo.UpdateTaskError(task.ID, "task expired")
			if err != nil {
				log.Errorf("update task error: %s", err)
			}
		}
	}()

	// 下发扫描任务
	err = cp.frontierBound.EmitScanApplications(context.Background(), task.ID, req.EdgeId, &frontierbound.Net{
		Nets:     nets,
		Protocol: req.Protocol,
		Port:     int(req.Port),
	})
	if err != nil {
		err = cp.repo.UpdateTaskError(task.ID, err.Error())
		if err != nil {
			log.Errorf("update task error: %s", err)
		}
		return nil, err
	}
	return &v1.CreateEdgeScanApplicationTaskResponse{
		Code:    200,
		Message: "success",
	}, nil
}

func (cp *controlPlane) GetEdgeScanApplicationTask(_ context.Context, req *v1.GetEdgeScanApplicationTaskRequest) (*v1.GetEdgeScanApplicationTaskResponse, error) {
	tasks, err := cp.repo.ListTasks(&dao.ListTasksQuery{
		EdgeID:      uint(req.EdgeId),
		TaskType:    model.TaskTypeScan,
		TaskSubType: model.TaskSubTypeScanApplication,
		//Status:      []model.TaskStatus{model.TaskStatusPending, model.TaskStatusRunning},
		Query: dao.Query{
			Page:     1,
			PageSize: 1,
			Order:    "id",
			Desc:     true,
		},
	})
	if err != nil {
		log.Errorf("list tasks error: %s", err)
		return nil, err
	}
	if len(tasks) == 0 {
		return &v1.GetEdgeScanApplicationTaskResponse{
			Code:    200,
			Message: "success",
			Data:    nil,
		}, nil
	}

	// result
	result := model.TaskScanApplicationResult{}
	err = json.Unmarshal(tasks[0].TaskResult, &result)
	if err != nil {
		log.Errorf("unmarshal task result error: %s", err)
		return nil, err
	}
	applications := []string{}
	for _, application := range result.ScannedApplications {
		applications = append(applications, fmt.Sprintf("%s:%d:%s", application.IP, application.Port, application.Protocol))
	}

	// 返回
	return &v1.GetEdgeScanApplicationTaskResponse{
		Code:    200,
		Message: "success",
		Data: &v1.EdgeScanApplicationTask{
			Id:           uint64(tasks[0].ID),
			EdgeId:       uint64(tasks[0].EdgeID),
			TaskStatus:   tasks[0].TaskStatus.String(),
			CreatedAt:    tasks[0].CreatedAt.Format(time.DateTime),
			UpdatedAt:    tasks[0].UpdatedAt.Format(time.DateTime),
			Applications: applications,
			Error:        tasks[0].Error,
		},
	}, nil
}

func transformEdges(edges []*model.Edge) []*v1.Edge {
	edgesV1 := make([]*v1.Edge, len(edges))
	for i, edge := range edges {
		edgesV1[i] = transformEdge(edge)
	}
	return edgesV1
}

func transformEdge(edge *model.Edge) *v1.Edge {
	edgeV1 := &v1.Edge{
		Id:          uint64(edge.ID),
		Name:        edge.Name,
		Description: edge.Description,
		Status:      int32(edge.Status),
		Online:      int32(edge.Online),
		CreatedAt:   edge.CreatedAt.Format(time.DateTime),
		UpdatedAt:   edge.UpdatedAt.Format(time.DateTime),
	}
	// 填充设备信息
	if edge.Device != nil {
		edgeV1.Device = transformDevice(edge.Device)
	}
	return edgeV1
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
