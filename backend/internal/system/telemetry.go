package system

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type NetworkInterfaceStats struct {
	RxBytes   uint64 `json:"rx_bytes"`
	TxBytes   uint64 `json:"tx_bytes"`
	OperState string `json:"operstate"`
	IPv4      string `json:"ipv4"`
	IPv4GW    string `json:"ipv4_gw"`
	IPv6      string `json:"ipv6"`
	IPv6GW    string `json:"ipv6_gw"`
}

type SystemUsage struct {
	CPU  float64 `json:"cpu"`
	Ram  float64 `json:"ram"`
	Swap float64 `json:"swap"`
}

var (
	prevIdle  uint64
	prevTotal uint64
)

type AddrInfo struct {
	Ifname    string `json:"ifname"`
	Operstate string `json:"operstate"`
	AddrInfo  []struct {
		Family    string `json:"family"`
		Local     string `json:"local"`
		Prefixlen int    `json:"prefixlen"`
	} `json:"addr_info"`
}

type RouteInfo struct {
	Dst     string `json:"dst"`
	Gateway string `json:"gateway"`
	Dev     string `json:"dev"`
}

// GetNetworkInterfaces parses /proc/net/dev to get interface stats and ip to get IPs/routes.
func GetNetworkInterfaces() map[string]NetworkInterfaceStats {
	stats := make(map[string]NetworkInterfaceStats)
	data, err := os.ReadFile("/proc/net/dev")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines[2:] { // skip header
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			iface := strings.TrimSpace(parts[0])
			fields := strings.Fields(parts[1])
			if len(fields) >= 9 {
				rx, _ := strconv.ParseUint(fields[0], 10, 64)
				tx, _ := strconv.ParseUint(fields[8], 10, 64)
				stats[iface] = NetworkInterfaceStats{
					RxBytes: rx,
					TxBytes: tx,
				}
			}
		}
	} else {
		log.Printf("Failed to read /proc/net/dev, using mock data: %v", err)
		stats["eth0"] = NetworkInterfaceStats{RxBytes: 1000, TxBytes: 2000}
		stats["wlan0"] = NetworkInterfaceStats{RxBytes: 500, TxBytes: 100}
	}

	var addrs []AddrInfo
	outAddr, errAddr := exec.Command("ip", "-j", "addr", "show").Output()
	if errAddr == nil {
		json.Unmarshal(outAddr, &addrs)
	}

	var routes4 []RouteInfo
	outRoute4, errRoute4 := exec.Command("ip", "-j", "route", "show").Output()
	if errRoute4 == nil {
		json.Unmarshal(outRoute4, &routes4)
	}

	var routes6 []RouteInfo
	outRoute6, errRoute6 := exec.Command("ip", "-j", "-6", "route", "show").Output()
	if errRoute6 == nil {
		json.Unmarshal(outRoute6, &routes6)
	}

	for _, addr := range addrs {
		iface := addr.Ifname
		stat := stats[iface]

		stat.OperState = addr.Operstate

		var ipv4, ipv6 string
		for _, a := range addr.AddrInfo {
			if a.Family == "inet" && ipv4 == "" {
				ipv4 = a.Local + "/" + strconv.Itoa(a.Prefixlen)
			} else if a.Family == "inet6" && ipv6 == "" && !strings.HasPrefix(a.Local, "fe80") {
				ipv6 = a.Local + "/" + strconv.Itoa(a.Prefixlen)
			}
		}

		var gw4, gw6 string
		for _, r := range routes4 {
			if r.Dev == iface && r.Dst == "default" {
				gw4 = r.Gateway
			}
		}
		for _, r := range routes6 {
			if r.Dev == iface && r.Dst == "default" {
				gw6 = r.Gateway
			}
		}

		stat.IPv4 = ipv4
		stat.IPv4GW = gw4
		stat.IPv6 = ipv6
		stat.IPv6GW = gw6

		stats[iface] = stat
	}

	return stats
}

// GetSystemUsage parses /proc/stat and /proc/meminfo to get CPU, RAM, and Swap usage percentages.
func GetSystemUsage() SystemUsage {
	// CPU usage
	var cpuUsage float64
	dataStat, err := os.ReadFile("/proc/stat")
	if err == nil {
		lines := strings.Split(string(dataStat), "\n")
		if len(lines) > 0 {
			fields := strings.Fields(lines[0])
			if len(fields) >= 5 && fields[0] == "cpu" {
				var total uint64
				var idle uint64

				for i := 1; i < len(fields); i++ {
					val, _ := strconv.ParseUint(fields[i], 10, 64)
					total += val
					if i == 4 || i == 5 { // idle and iowait
						idle += val
					}
				}

				deltaTotal := total - prevTotal
				deltaIdle := idle - prevIdle

				prevTotal = total
				prevIdle = idle

				if deltaTotal > 0 {
					cpuUsage = float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100.0
				}
			}
		}
	} else {
		log.Printf("Failed to read /proc/stat, using mock data: %v", err)
		cpuUsage = 25.5
	}

	// Mem usage
	var ramUsedPct float64
	var swapUsedPct float64

	dataMem, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		lines := strings.Split(string(dataMem), "\n")
		mem := make(map[string]float64)
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				key := strings.TrimSuffix(fields[0], ":")
				val, _ := strconv.ParseFloat(fields[1], 64)
				mem[key] = val
			}
		}

		if mem["MemTotal"] > 0 {
			if memAvailable, ok := mem["MemAvailable"]; ok {
				ramUsedPct = (mem["MemTotal"] - memAvailable) / mem["MemTotal"] * 100.0
			} else {
				ramUsedPct = (mem["MemTotal"] - mem["MemFree"] - mem["Buffers"] - mem["Cached"]) / mem["MemTotal"] * 100.0
			}
		}

		if mem["SwapTotal"] > 0 {
			swapUsedPct = (mem["SwapTotal"] - mem["SwapFree"]) / mem["SwapTotal"] * 100.0
		}
	} else {
		log.Printf("Failed to read /proc/meminfo, using mock data: %v", err)
		ramUsedPct = 40.2
		swapUsedPct = 10.5
	}

	return SystemUsage{
		CPU:  cpuUsage,
		Ram:  ramUsedPct,
		Swap: swapUsedPct,
	}
}
