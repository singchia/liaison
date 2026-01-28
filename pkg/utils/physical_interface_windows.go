//go:build windows
// +build windows

package utils

import (
	"strings"

	"github.com/yusufpapurcu/wmi"
)

// getWindowsPhysicalAdaptersWMI 通过 WMI 获取物理网络适配器（Windows 专用实现）
func getWindowsPhysicalAdaptersWMI() (map[string]*WindowsNetworkAdapter, error) {
	var adapters []WindowsNetworkAdapter
	query := `
		SELECT Name, Description, MACAddress,
		       PhysicalAdapter, NetEnabled, Manufacturer
		FROM Win32_NetworkAdapter
	`

	err := wmi.Query(query, &adapters)
	if err != nil {
		return nil, err
	}

	// 创建 MAC 地址到适配器的映射
	adapterMap := make(map[string]*WindowsNetworkAdapter)
	for i := range adapters {
		adapter := &adapters[i]
		// 只处理物理适配器
		if !adapter.PhysicalAdapter {
			continue
		}
		// 必须有 MAC 地址
		if adapter.MACAddress == nil || *adapter.MACAddress == "" {
			continue
		}
		// 必须启用（如果 NetEnabled 不为 nil）
		if adapter.NetEnabled != nil && !*adapter.NetEnabled {
			continue
		}
		// 排除明显虚拟厂商（通过描述）
		desc := strings.ToLower(adapter.Description)
		virtualKeywords := []string{
			"virtual", "vmware", "hyper-v",
			"vbox", "tap", "tun", "vpn", "wsl",
		}
		isVirtual := false
		for _, k := range virtualKeywords {
			if strings.Contains(desc, k) {
				isVirtual = true
				break
			}
		}
		if isVirtual {
			continue
		}

		// 将 MAC 地址标准化（去除分隔符，转为小写）
		mac := NormalizeMAC(*adapter.MACAddress)
		adapterMap[mac] = adapter
	}

	return adapterMap, nil
}
