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

func normalizeDeviceName(device string) string {
	if strings.HasPrefix(device, "/dev/") {
		for i := len(device) - 1; i >= 0; i-- {
			if device[i] < '0' || device[i] > '9' {
				return device[:i+1]
			} else {
				return device
			}
		}
	}

	return ""
}

func GetDiskInfo() []DiskInfo {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatal(err.Error())
	}

	var diskMap map[string]*DiskInfo = make(map[string]*DiskInfo)

	for _, v := range partitions {
		usage, err := disk.Usage(v.Mountpoint)
		if err != nil {
			log.Fatal(err.Error())
		}

		diskName := normalizeDeviceName(v.Device)
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
	return res
}
