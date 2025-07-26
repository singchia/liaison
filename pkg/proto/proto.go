package proto

type Meta struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type ReportDeviceUsage struct {
	Fingerprint string  `json:"fingerprint"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}
