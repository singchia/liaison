package proto

type Meta struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

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
