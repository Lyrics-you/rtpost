package main

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/shirou/gopsutil/v3/cpu"
)

func GetCpuPercent() float64 {
	usage, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Panicf("cpu.Percent error, %v", err)
	}
	return usage[0]
}

func PrintCupInformation() (cores int32) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Printf("cpu.Info error, %v", err)
	}
	infoStat := cpuInfos[0]
	fmt.Println("ModelName : ", infoStat.ModelName)
	fmt.Println("Cores : ", infoStat.Cores)
	return infoStat.Cores
}

func main() {
	cores := PrintCupInformation()
	for {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("%.2f%%", GetCpuPercent()*100/float64(cores)))

		// time.Sleep(time.Second)
	}
}
