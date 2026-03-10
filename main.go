package main

import (
	"fmt"

	"github.com/Xpaul6/sysinfo_agent/info"
	. "github.com/Xpaul6/sysinfo_agent/models"
)

func main() {
	var cpuInfo CpuInfo = info.GetCpuInfo()
	var memInfo Meminfo = info.GetMeminfo()

	fmt.Println(cpuInfo)
	fmt.Println(memInfo)
}
