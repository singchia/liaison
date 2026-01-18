package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/jumboframes/armorigo/log"
	"github.com/jumboframes/armorigo/sigaction"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/singchia/liaison/pkg/edge"
	"github.com/singchia/liaison/pkg/edge/config"
)

func main() {
	// 检查是否需要显示指纹
	if config.ShouldShowFingerprint() {
		// 显示详细指纹信息
		showFingerprintDetails()
		os.Exit(0)
	}

	edge, err := edge.NewEdge()
	if err != nil {
		log.Errorf("new edge err: %s", err)
		return
	}

	sig := sigaction.NewSignal()
	sig.Wait(context.TODO())

	edge.Close()
}

func showFingerprintDetails() {
	// 获取所有非 loopback 网卡的 MAC 地址，并排序
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting network interfaces: %v\n", err)
		os.Exit(1)
	}
	macs := make([]string, 0)
	macList := []string{}
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 && iface.Flags&net.FlagLoopback == 0 {
			macStr := iface.HardwareAddr.String()
			macs = append(macs, macStr)
			macList = append(macList, fmt.Sprintf("%s (%s)", macStr, iface.Name))
		}
	}
	if len(macs) == 0 {
		fmt.Fprintf(os.Stderr, "Error: failed to get MAC address\n")
		os.Exit(1)
	}
	// 对 MAC 地址进行排序
	sort.Strings(macs)
	mac := strings.Join(macs, ",")

	// CPU 信息
	cpuID := ""
	cpuInfo, err := cpu.Info()
	if err != nil {
		hostname, _ := os.Hostname()
		cpuID = hostname
		fmt.Printf("CPU Info: Error getting CPU info, using hostname as fallback: %s\n", hostname)
	} else if len(cpuInfo) > 0 {
		// 收集所有 CPU 信息并排序
		cpuStrings := make([]string, 0, len(cpuInfo))
		for _, info := range cpuInfo {
			cpuStr := info.ModelName + info.VendorID + info.Family
			cpuStrings = append(cpuStrings, cpuStr)
			fmt.Printf("CPU Info[%d]: ModelName=%s, VendorID=%s, Family=%s\n", len(cpuStrings)-1, info.ModelName, info.VendorID, info.Family)
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
	diskIDs := make([]string, 0)
	counts, err := disk.IOCounters()
	if err == nil && len(counts) > 0 {
		for _, stats := range counts {
			if stats.SerialNumber != "" {
				diskIDs = append(diskIDs, stats.SerialNumber)
				fmt.Printf("Disk Serial: %s (from %s)\n", stats.SerialNumber, stats.Name)
			}
		}
	}
	if len(diskIDs) == 0 {
		hostname, _ := os.Hostname()
		diskIDs = []string{hostname}
		fmt.Printf("Disk Serial: Error getting disk serial, using hostname as fallback: %s\n", hostname)
	}
	// 对磁盘序列号进行排序
	sort.Strings(diskIDs)
	diskID := strings.Join(diskIDs, ",")

	// 显示详细信息
	fmt.Println("\n=== Fingerprint Details ===")
	fmt.Printf("All MAC Addresses:\n")
	for _, m := range macList {
		fmt.Printf("  - %s\n", m)
	}
	fmt.Printf("Sorted MACs (used): %s\n", mac)
	fmt.Printf("CPU ID (sorted): %s\n", cpuID)
	fmt.Printf("Disk ID (sorted): %s\n", diskID)

	// 计算指纹
	raw := fmt.Sprintf("%s|%s|%s", mac, cpuID, diskID)
	fmt.Printf("\nRaw String: %s\n", raw)
	sum := sha256.Sum256([]byte(raw))
	fingerprint := hex.EncodeToString(sum[:])
	fmt.Printf("Fingerprint: %s\n", fingerprint)
}
