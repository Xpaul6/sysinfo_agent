package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
	// "os"
)

var err error

const CHECK_INTERVAL = 500 * time.Millisecond
const BTOMB = 1024 * 1024

type CpuInfo struct {
	loadPercentage float64
	temperature    float64
}

type Meminfo struct {
	loadPercentage float64
	total          uint64 // bytes
	used           uint64 // bytes
}

func getCpuInfo() CpuInfo {
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

	var res CpuInfo = CpuInfo{loadPercentage[0], temperature}
	return res
}

func getMeminfo() Meminfo {
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err.Error())
	}

	var res Meminfo = Meminfo{vm.UsedPercent, vm.Total, vm.Used}
	return res
}

func main() {
	var cpuInfo CpuInfo = getCpuInfo()
	var memInfo Meminfo = getMeminfo()

	fmt.Println(cpuInfo)
	fmt.Println(memInfo)
}
