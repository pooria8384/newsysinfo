// components/raminfo.go
package sysinfo

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/mem"
)

type RAMInfo struct {
	Total       uint64  `json:"total_ram"`
	Available   uint64  `json:"available_ram"`
	Used        uint64  `json:"used_ram"`
	UsedPercent float64 `json:"usedpercent_ram"`
}

func (r *RAMInfo) FetchInfo() error {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("error getting memory info: %v", err)
	}

	r.Total = vmStat.Total / 1024 / 1024
	r.Available = vmStat.Available / 1024 / 1024
	r.Used = vmStat.Used / 1024 / 1024
	r.UsedPercent = vmStat.UsedPercent
	return nil
}

func (r *RAMInfo) GetInfo() error {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling RAMInfo to JSON: +%v", err)
	}
	fmt.Println(string(jsonData), "MB")
	return nil
}
