package main

import (
	"log"
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
	errChan := make(chan error, 3)
	var wg = sync.WaitGroup{}

	const gatherCallsNumber = 4
	var cpuInfo CpuInfo
	var memInfo MemInfo
	var diskInfo []DiskInfo
	var netInfo []NetInfo

	wg.Add(gatherCallsNumber)

	// Goroutines function calls
	go func() {
		defer wg.Done()
		var err error
		cpuInfo, err = info.GetCpuInfo()
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		memInfo, err = info.GetMemInfo()
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		diskInfo, err = info.GetDiskInfo()
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		netInfo, err = info.GetNetInfo()
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	var sysInfo = SysInfo{
		CPU:   cpuInfo,
		Mem:   memInfo,
		Disks: diskInfo,
		Net: netInfo,
	}

	c.IndentedJSON(http.StatusOK, sysInfo)

	close(errChan)
	// Error logging
	for v := range errChan {
		log.Println(v)
	}
}
