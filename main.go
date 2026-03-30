package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/Xpaul6/sysinfo_agent/info"
	. "github.com/Xpaul6/sysinfo_agent/models"
)

func main() {
	router := gin.Default()
	router.GET("/sysinfo", getSystemInfo)

	router.Run("localhost:8080")
}

func getSystemInfo(c *gin.Context) {
	var wg = sync.WaitGroup{}

	const gatherCallsNumber = 3
	var cpuInfo CpuInfo
	var memInfo MemInfo
	var diskInfo []DiskInfo

	wg.Add(gatherCallsNumber)

	// Goroutines function calls
	go func() {
		defer wg.Done()
		cpuInfo = info.GetCpuInfo()
	}()

	go func() {
		defer wg.Done()
		memInfo = info.GetMemInfo()
	}()

	go func() {
		defer wg.Done()
		diskInfo = info.GetDiskInfo()
	}()

	wg.Wait()

	var sysInfo = SysInfo{
		CPU:   cpuInfo,
		Mem:   memInfo,
		Disks: diskInfo,
	}

	c.IndentedJSON(http.StatusOK, sysInfo)
}
