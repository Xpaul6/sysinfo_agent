package info

import (
	"log"
	"time"

	"github.com/Xpaul6/sysinfo_agent/models"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

var err error

const CHECK_INTERVAL = 500 * time.Millisecond

func GetCpuInfo() models.CpuInfo {
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

	var res models.CpuInfo = models.CpuInfo{
		LoadPercentage: loadPercentage[0],
		Temperature:    temperature,
	}
	return res
}

func GetMeminfo() models.Meminfo {
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err.Error())
	}

	var res models.Meminfo = models.Meminfo{
		LoadPercentage: vm.UsedPercent,
		Total:          vm.Total,
		Used:           vm.Used,
	}
	return res
}
