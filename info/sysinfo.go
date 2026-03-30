package info

import (
	"log"
	"time"
	"strings"

	. "github.com/Xpaul6/sysinfo_agent/models"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

const CHECK_INTERVAL = 500 * time.Millisecond

func GetCpuInfo() CpuInfo {
	loadPercentage, err := cpu.Percent(CHECK_INTERVAL, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	tempSensors, err := sensors.SensorsTemperatures()
	var temperature float64
	if err != nil {
		log.Fatal(err.Error())
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
	return res
}

func GetMemInfo() MemInfo {
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err.Error())
	}

	var res MemInfo = MemInfo{
		LoadPercentage: vm.UsedPercent,
		Total:          vm.Total,
		Used:           vm.Used,
	}
	return res
}

func GetDiskInfo() []DiskInfo {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatal(err.Error())
	}
	var res []DiskInfo
	for _, v := range partitions {
		usage, err := disk.Usage(v.Mountpoint)
		if err != nil {
			log.Fatal(err.Error())
		}
		var curr DiskInfo = DiskInfo{
			MountPoint: v.Mountpoint,
			Total:      usage.Total,
			Used:       usage.Used,
		}
		res = append(res, curr)
	}
	return res
}

func normalizeDeviceName(device string) string {
	if strings.HasPrefix(device, "/dev/") {
		for i := len(device) - 1; i >= 0; i-- {
			if device[i] < '0' || device[i] > '9' {
				return device[:i+1]
			}
		}
	}

	if strings.HasPrefix(device, "disk") {
		if idx := strings.Index(device, "s"); idx != -1 {
			return device[:idx]
		}
	}

	return device
}

func GetPhysicalDisks() []DiskInfo {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil
	}

	diskMap := make(map[string]*DiskInfo)

	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		diskName := normalizeDeviceName(p.Device)

		if _, exists := diskMap[diskName]; !exists {
			diskMap[diskName] = &DiskInfo{
				MountPoint: diskName,
			}
		}

		diskMap[diskName].Total += usage.Total
		diskMap[diskName].Used += usage.Used
	}

	var result []DiskInfo
	for _, d := range diskMap {
		result = append(result, *d)
	}

	return result
}
