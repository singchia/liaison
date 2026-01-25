package utils

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	// Windows 物理适配器缓存（MAC 地址 -> 适配器信息）
	windowsPhysicalAdapters map[string]*WindowsNetworkAdapter
	windowsAdapterOnce      sync.Once
	windowsAdapterErr        error
)

// WindowsNetworkAdapter Windows 网络适配器信息（通过 WMI 获取）
type WindowsNetworkAdapter struct {
	Name            string
	Description     string
	MACAddress      *string
	PhysicalAdapter bool
	NetEnabled      *bool
	Manufacturer    *string
}

// getWindowsPhysicalAdapters 通过 WMI 获取物理网络适配器（Windows 专用）
func getWindowsPhysicalAdapters() (map[string]*WindowsNetworkAdapter, error) {
	if runtime.GOOS != "windows" {
		return make(map[string]*WindowsNetworkAdapter), nil
	}
	return getWindowsPhysicalAdaptersWMI()
}

// isWindowsPhysicalInterface Windows 专用的物理接口判断（使用 WMI 查询结果）
func isWindowsPhysicalInterface(name, mac string) bool {
	// 确保 WMI 查询已执行
	windowsAdapterOnce.Do(func() {
		windowsPhysicalAdapters, windowsAdapterErr = getWindowsPhysicalAdapters()
		if windowsAdapterErr != nil {
			// 如果 WMI 查询失败，回退到名称过滤
			windowsPhysicalAdapters = make(map[string]*WindowsNetworkAdapter)
		}
	})

	// 如果 WMI 查询失败，回退到名称过滤
	if windowsAdapterErr != nil {
		nameLower := strings.ToLower(name)
		virtualKeywords := []string{
			"virtual", "vmware", "hyper-v",
			"vbox", "tap", "tun", "vpn", "wsl",
			"loopback", "isatap", "teredo",
		}
		for _, k := range virtualKeywords {
			if strings.Contains(nameLower, k) {
				return false
			}
		}
		return true
	}

	// 标准化 MAC 地址
	macNormalized := NormalizeMAC(mac)

	// 检查是否在物理适配器映射中
	_, exists := windowsPhysicalAdapters[macNormalized]
	return exists
}

// IsPhysicalInterfaceLinux 判断是否是物理网络接口（Linux 专用）
// 通过 sysfs 文件系统判断
func IsPhysicalInterfaceLinux(iface net.Interface) bool {
	if runtime.GOOS != "linux" {
		return false
	}

	// 1. 排除 loopback
	if iface.Flags&net.FlagLoopback != 0 {
		return false
	}

	// 2. 接口名快速排除
	name := iface.Name
	virtualPrefixes := []string{
		"lo", "docker", "veth", "cni",
		"flannel", "virbr", "br-",
		"tun", "tap",
	}
	for _, p := range virtualPrefixes {
		if strings.HasPrefix(name, p) {
			return false
		}
	}

	// 3. 必须有 MAC
	if len(iface.HardwareAddr) == 0 {
		return false
	}

	// 4. MAC 不能是本地生成
	if iface.HardwareAddr[0]&2 != 0 {
		return false
	}

	// 5. sysfs device 判断（⭐️最关键）
	devicePath := filepath.Join("/sys/class/net", name, "device")
	if _, err := os.Stat(devicePath); err != nil {
		return false
	}

	// 6. type 判断（兜底）
	typePath := filepath.Join("/sys/class/net", name, "type")
	if data, err := os.ReadFile(typePath); err == nil {
		if strings.TrimSpace(string(data)) != "1" {
			return false
		}
	}

	return true
}

// IsPhysicalInterface 判断是否是物理网络接口
// macOS: 1. 接口名称必须以 "en" 开头 2. MAC 地址不能是本地管理的
// Windows: 通过 WMI 查询物理适配器
// Linux: 通过 sysfs 判断（需要 net.Interface，在 isVirtualInterface 中处理）
func IsPhysicalInterface(name, mac string) bool {
	if runtime.GOOS == "darwin" {
		// macOS: enX 才可能是物理
		if !strings.HasPrefix(name, "en") {
			return false
		}
		// MAC 地址不能是本地管理的
		if IsLocallyAdministeredMAC(mac) {
			return false
		}
		return true
	}

	if runtime.GOOS == "windows" {
		// Windows: 必须有 MAC 地址
		if mac == "" {
			return false
		}
		// 使用 WMI 查询物理适配器
		return isWindowsPhysicalInterface(name, mac)
	}

	if runtime.GOOS == "linux" {
		// Linux: 使用 sysfs 判断物理接口
		// 注意：这里需要 net.Interface，但函数签名只有 name 和 mac
		// 所以 Linux 的判断会在 isVirtualInterface 中处理
		// 这里只做基本的 MAC 地址检查
		if mac == "" {
			return false
		}
		// MAC 不能是本地生成
		if IsLocallyAdministeredMAC(mac) {
			return false
		}
		return true
	}

	// 其他系统，默认返回 true（由其他逻辑过滤）
	return true
}
