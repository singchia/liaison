package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

// FingerprintComponents 指纹组成信息
type FingerprintComponents struct {
	MACs    []string
	CPUInfo string
	DiskIDs []string
	Raw     string
}

// GetFingerprint 获取设备指纹
// 基于 MAC、CPU、磁盘序列号生成指纹
// 所有组件都进行排序，确保指纹稳定性
// 返回指纹字符串、组成信息和错误
func GetFingerprint() (string, *FingerprintComponents, error) {
	// 获取所有非 loopback 网卡的 MAC 地址，并排序
	// 在 Mac 上，需要过滤掉不稳定的虚拟接口（如 VPN、虚拟网卡等）
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}
	macs := make([]string, 0)
	for _, iface := range interfaces {
		// 使用 isVirtualInterface 判断，如果是虚拟接口则跳过
		if isVirtualInterface(&iface) {
			continue
		}
		macs = append(macs, iface.HardwareAddr.String())
	}
	if len(macs) == 0 {
		// 如果过滤后没有 MAC 地址，回退到只做基本检查（不检查 UP 标志和 PointToPoint）
		for _, iface := range interfaces {
			// 跳过 loopback 和没有 MAC 的接口
			if iface.Flags&net.FlagLoopback != 0 || len(iface.HardwareAddr) == 0 {
				continue
			}
			// 根据平台判断是否是物理接口（不检查 UP 和 PointToPoint）
			isPhysical := false
			if runtime.GOOS == "darwin" {
				isPhysical = IsMacPhysicalInterface(&iface)
			} else if runtime.GOOS == "windows" {
				isPhysical = IsWindowsPhysicalInterface(&iface)
			} else if runtime.GOOS == "linux" {
				isPhysical = IsLinuxPhysicalInterface(&iface)
			} else {
				// 其他系统，只做基本名称过滤
				name := strings.ToLower(iface.Name)
				isPhysical = !strings.Contains(name, "utun") &&
					!strings.Contains(name, "bridge") &&
					!strings.Contains(name, "vmnet") &&
					!strings.Contains(name, "vboxnet") &&
					!strings.Contains(name, "awdl") &&
					!strings.Contains(name, "anpi")
			}
			if isPhysical {
				macs = append(macs, iface.HardwareAddr.String())
			}
		}
	}
	if len(macs) == 0 {
		return "", nil, fmt.Errorf("failed to get MAC address")
	}
	// 对 MAC 地址进行排序，确保顺序稳定
	sort.Strings(macs)
	mac := strings.Join(macs, ",") // 使用所有 MAC 地址，用逗号分隔

	// CPU 信息
	cpuID := ""
	cpuInfo, err := cpu.Info()
	if err != nil {
		// 如果无法获取 CPU 信息，使用 hostname 作为备用
		hostname, _ := os.Hostname()
		cpuID = hostname
	} else if len(cpuInfo) > 0 {
		// 如果有多个 CPU，对 CPU 信息进行排序
		cpuStrings := make([]string, 0, len(cpuInfo))
		for _, info := range cpuInfo {
			cpuStr := info.ModelName + info.VendorID + info.Family
			cpuStrings = append(cpuStrings, cpuStr)
		}
		// 去重并排序
		cpuMap := make(map[string]bool)
		uniqueCpus := make([]string, 0)
		for _, cpuStr := range cpuStrings {
			if !cpuMap[cpuStr] {
				cpuMap[cpuStr] = true
				uniqueCpus = append(uniqueCpus, cpuStr)
			}
		}
		sort.Strings(uniqueCpus)
		cpuID = strings.Join(uniqueCpus, ",")
	}

	// 磁盘序列号（收集所有磁盘的序列号，并排序）
	// 在 Mac 上，只使用物理磁盘的序列号，过滤掉虚拟磁盘和挂载点
	diskIDs := make([]string, 0)
	counts, err := disk.IOCounters()
	if err == nil && len(counts) > 0 {
		for _, stats := range counts {
			// 只使用有序列号的磁盘
			if stats.SerialNumber != "" {
				// 在 Mac 上，过滤掉可能的虚拟磁盘（名称包含特定关键词）
				name := strings.ToLower(stats.Name)
				// 跳过明显的虚拟磁盘（可以根据实际情况调整）
				if !strings.Contains(name, "disk image") &&
					!strings.Contains(name, "virtual") {
					diskIDs = append(diskIDs, stats.SerialNumber)
				}
			}
		}
	}
	// 如果无法获取磁盘序列号，使用 hostname 作为备用
	if len(diskIDs) == 0 {
		hostname, _ := os.Hostname()
		diskIDs = []string{hostname}
	}
	// 对磁盘序列号进行排序，确保顺序稳定
	sort.Strings(diskIDs)
	diskID := strings.Join(diskIDs, ",") // 使用所有磁盘序列号，用逗号分隔

	raw := fmt.Sprintf("%s|%s|%s", mac, cpuID, diskID)
	sum := sha256.Sum256([]byte(raw))
	fingerprint := hex.EncodeToString(sum[:])

	components := &FingerprintComponents{
		MACs:    macs,
		CPUInfo: cpuID,
		DiskIDs: diskIDs,
		Raw:     raw,
	}

	return fingerprint, components, nil
}
