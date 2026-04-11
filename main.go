package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Xpaul6/sysinfo_agent/info"
	. "github.com/Xpaul6/sysinfo_agent/models"
)

func main() {
	var port string = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot load .env file")
	} else {
		port = os.Getenv("PORT")
	}

	var currSysInfo = SysInfo{
		CPU:   CpuInfo{},
		Mem:   MemInfo{},
		Disks: []DiskInfo{},
		Net:   []NetInfo{},
	}

	// Info update loop
	go func() {
		for {
			currSysInfo = getSystemInfo()
			time.Sleep(7 * time.Second)
		}
	}()

	// Gin router setup
	router := gin.Default()
	router.GET("/sysinfo", func(c *gin.Context) { c.IndentedJSON(http.StatusOK, currSysInfo) })

	router.Run("localhost:" + string(port))
}

func getSystemInfo() SysInfo {
	const gatherCallsNumber = 4
	var wg = sync.WaitGroup{}

	var cpuInfo CpuInfo
	var memInfo MemInfo
	var diskInfo []DiskInfo
	var netInfo []NetInfo

	errChan := make(chan error, gatherCallsNumber)
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
		Net:   netInfo,
	}

	// Error logging
	close(errChan)
	for v := range errChan {
		log.Println(v)
	}

	return sysInfo
}
