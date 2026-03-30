package info

import (
	"strings"
	"time"

	. "github.com/Xpaul6/sysinfo_agent/models"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/sensors"
)

const CHECK_INTERVAL = 3000 * time.Millisecond

// CPU info
func GetCpuInfo() (CpuInfo, error) {
	loadPercentage, err := cpu.Percent(CHECK_INTERVAL, false)
	if err != nil {
		return CpuInfo{}, nil
	}

	tempSensors, err := sensors.SensorsTemperatures()
	var temperature float64
	if err != nil {
		return CpuInfo{}, nil
	}

	for _, v := range tempSensors {
		if v.SensorKey == "coretemp_package_id_0" || v.SensorKey == "PMU tcal" {
			temperature = v.Temperature
			break
		}
	}

	var res CpuInfo = CpuInfo{
		LoadPercentage: loadPercentage[0],
		Temperature:    temperature,
	}
	return res, nil
}

// RAM info
func GetMemInfo() (MemInfo, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return MemInfo{}, err
	}

	var res MemInfo = MemInfo{
		LoadPercentage: vm.UsedPercent,
		Total:          vm.Total,
		Used:           vm.Used,
	}
	return res, nil
}

// Disks info
func normalizeDiskDeviceName(device string) string {
	// TODO: needs more optimal way for detectng physical drives, but works so far
	if strings.HasPrefix(device, "/dev/s") {
		for i := len(device) - 1; i >= 0; i-- {
			if device[i] < '0' || device[i] > '9' {
				return device[:i+1]
			}
		}
	}

	return ""
}

func GetDiskInfo() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskMap map[string]*DiskInfo = make(map[string]*DiskInfo)

	for _, v := range partitions {
		usage, err := disk.Usage(v.Mountpoint)
		if err != nil {
			return nil, err
		}

		diskName := normalizeDiskDeviceName(v.Device)
		if diskName == "" {
			continue
		}

		if _, exists := diskMap[diskName]; !exists {
			diskMap[diskName] = &DiskInfo{
				MountPoint: diskName,
			}
		}

		diskMap[diskName].Total += usage.Total
		diskMap[diskName].Used += usage.Used
	}
	var res []DiskInfo
	for _, v := range diskMap {
		res = append(res, *v)
	}
	return res, nil
}

// Net devices info
func isOkNetDeviceName(name string) bool {
	var filter []string = []string{"lo", "docker", "veth", "br-", "bridge", "utun"}
	for _, v := range filter {
		if strings.HasPrefix(name, v) {
			return false
		}
	}
	return true
}

func GetNetInfo() ([]NetInfo, error) {
	before, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	time.Sleep(CHECK_INTERVAL)

	after, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	beforeMap := make(map[string]net.IOCountersStat)
	for _, v := range before {
		if !isOkNetDeviceName(v.Name) {
			continue
		}
		beforeMap[v.Name] = v
	}

	var res []NetInfo
	for _, a := range after {
		b, ok := beforeMap[a.Name]
		if !ok {
			continue
		}

		var rBytes = a.BytesRecv - b.BytesRecv
		var sBytes = a.BytesSent - b.BytesSent
		res = append(res, NetInfo{
			Name: a.Name,
			RBpS: float64(rBytes) / CHECK_INTERVAL.Seconds(),
			SBpS: float64(sBytes) / CHECK_INTERVAL.Seconds(),
		})
	}

	return res, nil
}
