package main

import (
	"fmt"
	"time"

	"github.com/Xpaul6/sysinfo_agent/info"
	"github.com/Xpaul6/sysinfo_agent/models"
)

var err error

const CHECK_INTERVAL = 500 * time.Millisecond
const BTOMB = 1024 * 1024

func main() {
	var cpuInfo models.CpuInfo = info.GetCpuInfo()
	var memInfo models.Meminfo = info.GetMeminfo()

	fmt.Println(cpuInfo)
	fmt.Println(memInfo)
}
