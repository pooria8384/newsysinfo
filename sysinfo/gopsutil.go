package sysinfo

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type Gopsutil struct {
	*SystemInfo
}

func NewGopsutil() Iagent {
	return &Gopsutil{
		&SystemInfo{
			OsInfo:    &OsInfo{},
			DiskInfos: []DiskInfo{},
			CpuInfo:   &CpuInfo{},
			RamInfo:   &RamInfo{},
		},
	}
}

func (g *Gopsutil) Get() *SystemInfo {
	return g.SystemInfo
}

func (g *Gopsutil) Cpu() error {
	cpuInfo := &CpuInfo{}

	data, err := cpu.Info()
	if err != nil {
		return fmt.Errorf("error getting CPU info: %v", err)
	}

	if len(data) > 0 {
		cpuInfo.Modelname = data[0].ModelName
		cpuInfo.Cores = uint32(data[0].Cores)
	}
	g.SystemInfo.CpuInfo = cpuInfo
	return nil
}

func (g *Gopsutil) Ram() error {
	ram := &RamInfo{}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("error getting memory info: %v", err)
	}

	ram.Total = uint32(vmStat.Total)
	ram.Available = uint32(vmStat.Available)
	ram.Used = uint32(vmStat.Used)
	ram.UsedPercent = vmStat.UsedPercent

	g.SystemInfo.RamInfo = ram
	return nil
}

func (g *Gopsutil) Disk() error {
	disks := []DiskInfo{}

	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("error getting disk partitions: %v", err)
	}

	for _, partition := range partitions {
		diskInfo, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return fmt.Errorf("error getting disk usage info for %s: %v", partition.Mountpoint, err)
		}

		disks = append(disks, DiskInfo{
			Device:    partition.Device,
			TotalSize: uint32(diskInfo.Total),
			FreeSize:  uint32(diskInfo.Free),
		})
	}

	g.SystemInfo.DiskInfos = disks
	return nil
}

func (g *Gopsutil) Os() error {
	o := &OsInfo{}
	o.OSType = runtime.GOOS
	o.OSArch = runtime.GOARCH

	hostName, err := os.Hostname()

	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}
	o.Hostname = hostName

	g.SystemInfo.OsInfo = o

	return nil
}
