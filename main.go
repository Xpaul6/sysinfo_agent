package main

import (
	"fmt"
	"sync"

	"github.com/Xpaul6/sysinfo_agent/info"
	. "github.com/Xpaul6/sysinfo_agent/models"
)

func main() {
	var wg = sync.WaitGroup{}

	var cpuInfo CpuInfo
	var memInfo Meminfo

	wg.Add(2)

	go func(){
		defer wg.Done()
		cpuInfo = info.GetCpuInfo()
	}()

	go func(){
		defer wg.Done()
		memInfo = info.GetMemInfo()
	}()

	wg.Wait()

	fmt.Println(cpuInfo)
	fmt.Println(memInfo)
}
