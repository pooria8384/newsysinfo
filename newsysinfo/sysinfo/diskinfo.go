// components/diskinfo.go
package sysinfo

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/disk"
)

type DiskInfo struct {
	Device    string `json:"diskinfo_device"`
	TotalSize uint64 `json:"diskinfo_totalsize"`
	FreeSize  uint64 `json:"diskinfo_freesize"`
}

func (d *DiskInfo) FetchInfo() error {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("error getting disk partitions: %v", err)
	}

	for _, partition := range partitions {
		diskInfo, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Printf("error getting disk usage info for %s: %v", partition.Mountpoint, err)
		}
		d.Device = partition.Device
		d.TotalSize = diskInfo.Total / 1024 / 1024 / 1024
		d.FreeSize = diskInfo.Free / 1024 / 1024 / 1024
	}
	return nil
}

func (d *DiskInfo) GetInfo() error {
	jsonData, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling DiskInfo to JSON: +%v", err)
	}
	fmt.Println(string(jsonData), "GB")
	return nil
}
