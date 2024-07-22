// components/cpuinfo.go
package sysinfo

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

type CPUInfo struct {
	Modelname string `json:"cpu_ModelName"`
	Cores     int32  `json:"cpu_cores"`
}

func (c *CPUInfo) FetchInfo() error {
	cpuInfos, err := cpu.Info()
	if err != nil {
		return fmt.Errorf("error getting CPU info: %v", err)
	}

	if len(cpuInfos) > 0 {
		c.Modelname = cpuInfos[0].ModelName
		c.Cores = cpuInfos[0].Cores
	}
	return nil
}

func (c *CPUInfo) GetInfo() error {
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling CPUInfo to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
	return nil
}
