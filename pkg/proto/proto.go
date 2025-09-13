package proto

import "context"

// manager <-> entry
// 一个Proxy是三元组
type Proxy struct {
	ID        int
	Name      string
	ProxyPort int
	IP        string
	Port      int
}

type ProxyManager interface {
	CreateProxy(ctx context.Context, proxy *Proxy) error
	DeleteProxy(ctx context.Context, proxy *Proxy) error
}

// manager <-> edge
type Meta struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// device
type Device struct {
	Fingerprint string                    `json:"fingerprint"`
	HostName    string                    `json:"host_name"`
	CPU         int                       `json:"cpu"`
	Memory      int                       `json:"memory"`
	OS          string                    `json:"os"`
	OSVersion   string                    `json:"os_version"`
	DeviceUsage DeviceUsage               `json:"device_usage"`
	Interfaces  []DeviceEthernetInterface `json:"interfaces"`
	EdgeID      uint64                    `json:"edge_id"` // 如果包含edge id，则表示是edge所在的机器
}

type DeviceEthernetInterface struct {
	Name    string `json:"name"`
	MAC     string `json:"mac"`
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

type DeviceUsage struct {
	Fingerprint string  `json:"fingerprint"`
	CPUUsage    float32 `json:"cpu_usage"`
	MemoryUsage float32 `json:"memory_usage"`
	DiskUsage   float32 `json:"disk_usage"`
}

// task
type ScanApplicationTaskRequest struct {
	TaskID   uint     `json:"task_id"`
	Nets     []string `json:"nets"`
	Port     int      `json:"port"`
	Protocol string   `json:"protocol"`
}

type ScannedApplication struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type ScanApplicationTaskResult struct {
	TaskID              uint                 `json:"task_id"`
	ScannedApplications []ScannedApplication `json:"scanned_applications"`
	Error               string               `json:"error"`
	Status              string               `json:"status"` // running, completed, failed
}

type Dst struct {
	Addr string `json:"addr"`
}
