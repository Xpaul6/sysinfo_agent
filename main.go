package main

import (
	"sync"
	"net/http"
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

	const gatherCallsNumber = 2
	var cpuInfo CpuInfo
	var memInfo MemInfo

	wg.Add(gatherCallsNumber)

	go func(){
		defer wg.Done()
		cpuInfo = info.GetCpuInfo()
	}()

	go func(){
		defer wg.Done()
		memInfo = info.GetMemInfo()
	}()

	wg.Wait()

	sysInfo := struct {
		CPU CpuInfo `json:"cpu"`
		Mem MemInfo `json:"mem"`
	}{
		CPU: cpuInfo,
		Mem: memInfo,
	}

	c.IndentedJSON(http.StatusOK, sysInfo)
}
