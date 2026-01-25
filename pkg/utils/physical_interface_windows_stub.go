//go:build !windows
// +build !windows

package utils

// getWindowsPhysicalAdaptersWMI 非 Windows 平台返回空映射
func getWindowsPhysicalAdaptersWMI() (map[string]*WindowsNetworkAdapter, error) {
	return make(map[string]*WindowsNetworkAdapter), nil
}
