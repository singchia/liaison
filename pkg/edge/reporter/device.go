package reporter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jumboframes/armorigo/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/singchia/liaison/pkg/proto"
)

func (r *reporter) loopReportDevice(ctx context.Context) {
	for {
		device, err := getDevice()
		if err != nil {
			log.Errorf("get device error: %v", err)
			continue
		}
		r.reportDevice(ctx, device)
		time.Sleep(time.Hour)
	}
}

func (r *reporter) loopReportDeviceUsage(ctx context.Context) {
	for {
		deviceUsage, err := getDeviceUsage()
		if err != nil {
			log.Errorf("get device usage error: %v", err)
			continue
		}
		if err != nil {
			log.Errorf("get device usage error: %v", err)
			continue
		}
		r.reportDeviceUsage(ctx, deviceUsage)
		time.Sleep(time.Minute)
	}
}

func (r *reporter) reportDevice(ctx context.Context, device *proto.Device) error {
	data, err := json.Marshal(device)
	if err != nil {
		return err
	}

	req := r.frontierBound.NewRequest(data)
	rsp, err := r.frontierBound.Call(ctx, "report_device", req)
	if err != nil {
		return err
	}
	if rsp.Error() != nil {
		return rsp.Error()
	}
	return nil
}

func (r *reporter) reportDeviceUsage(ctx context.Context, deviceUsage *proto.DeviceUsage) error {
	data, err := json.Marshal(deviceUsage)
	if err != nil {
		return err
	}
	req := r.frontierBound.NewRequest(data)
	rsp, err := r.frontierBound.Call(ctx, "report_device_usage", req)
	if err != nil {
		return err
	}
	if rsp.Error() != nil {
		return rsp.Error()
	}
	return nil
}

// 注意：指纹才是唯一标识一台主机，而不是边缘和AKSK
// 如果当前edge扫描的指纹发生改变，则会重新建立edge和device的关联
func getDevice() (*proto.Device, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// 基础信息
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	memMB := memory.Total / 1024 / 1024

	// OS
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	os := info.OS
	osVersion := info.KernelVersion

	// 网络接口
	interfaces, err := getDeviceEthernetInterface()
	if err != nil {
		return nil, err
	}
	// 指纹
	fingerprint, err := getFingerprint()
	if err != nil {
		return nil, err
	}

	return &proto.Device{
		Fingerprint: fingerprint,
		HostName:    hostname,
		CPU:         cpuCount,
		Memory:      int(memMB),
		OS:          os,
		OSVersion:   osVersion,
		Interfaces:  interfaces,
	}, nil
}

func getDeviceEthernetInterface() ([]*proto.DeviceEthernetInterface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ethernetInterfaces := make([]*proto.DeviceEthernetInterface, 0)

	for _, iface := range interfaces {
		// 跳过没有MAC或未启用的接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		ipMasks := make([]*proto.IPMask, 0)

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet: // 带掩码的情况
				ip := v.IP.String()
				mask := ""
				if v.Mask != nil {
					// IPv4 -> 点分十进制，IPv6 -> 前缀长度
					if v.IP.To4() != nil {
						mask = net.IP(v.Mask).String()
					} else {
						ones, _ := v.Mask.Size()
						mask = fmt.Sprintf("%d", ones) // IPv6 用前缀长度
					}
				}
				ipMasks = append(ipMasks, &proto.IPMask{
					IP:      ip,
					Netmask: mask,
				})

			case *net.IPAddr: // 只有IP没有掩码
				ipMasks = append(ipMasks, &proto.IPMask{
					IP:      v.IP.String(),
					Netmask: "",
				})
			}
		}

		ethernetInterfaces = append(ethernetInterfaces, &proto.DeviceEthernetInterface{
			Name:    iface.Name,
			MAC:     iface.HardwareAddr.String(),
			IPMasks: ipMasks,
		})
	}
	return ethernetInterfaces, nil
}

func getDeviceUsage() (*proto.DeviceUsage, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	// 资源使用
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}
	deviceUsage := &proto.DeviceUsage{
		CPUUsage:    float32(cpuUsage[0]),
		MemoryUsage: float32(memory.UsedPercent),
		DiskUsage:   float32(diskUsage.UsedPercent),
	}
	return deviceUsage, nil
}

// 基于 MAC、CPU、磁盘序列号生成指纹
func getFingerprint() (string, error) {
	// 获取主网卡 MAC
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	mac := ""
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 && iface.Flags&net.FlagLoopback == 0 {
			mac = iface.HardwareAddr.String()
			break
		}
	}

	// CPU 信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		return "", err
	}
	cpuID := ""
	if len(cpuInfo) > 0 {
		cpuID = cpuInfo[0].ModelName + cpuInfo[0].VendorID + cpuInfo[0].Family
	}

	// 磁盘序列号（取根分区对应的磁盘）
	diskID := ""
	counts, _ := disk.IOCounters()
	if len(counts) > 0 {
		for _, stats := range counts {
			if stats.SerialNumber != "" {
				diskID = stats.SerialNumber
				break
			}
		}
	}

	raw := fmt.Sprintf("%s|%s|%s", mac, cpuID, diskID)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:]), nil
}
