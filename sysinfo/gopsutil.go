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
			OsInfo:   &OsInfo{},
			DiskInfo: &DiskInfo{},
			CpuInfo:  &CpuInfo{},
			RamInfo:  &RamInfo{},
		},
	}
}

func (g *Gopsutil) Get() *SystemInfo {
	return g.SystemInfo
}

func (g *Gopsutil) Cpu() (Iagent, error) {
	cpuInfo := &CpuInfo{}

	data, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting CPU info: %v", err)
	}

	if len(data) > 0 {
		cpuInfo.Modelname = data[0].ModelName
		cpuInfo.Cores = int(data[0].Cores)
	}
	g.SystemInfo.CpuInfo = cpuInfo
	return g, nil
}

func (g *Gopsutil) Ram() (Iagent, error) {
	ram := &RamInfo{}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting memory info: %v", err)
	}

	ram.Total = int(vmStat.Total)
	ram.Available = int(vmStat.Available)
	ram.Used = int(vmStat.Used)
	ram.UsedPercent = vmStat.UsedPercent

	g.SystemInfo.RamInfo = ram
	return g, nil
}

func (g *Gopsutil) Disk() (Iagent, error) {
	d := &DiskInfo{}
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("error getting disk partitions: %v", err)
	}

	for _, partition := range partitions {
		diskInfo, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Printf("error getting disk usage info for %s: %v", partition.Mountpoint, err)
		}
		d.Device = partition.Device
		d.TotalSize = int(diskInfo.Total)
		d.FreeSize = int(diskInfo.Free)
	}
	g.SystemInfo.DiskInfo = d
	return g, nil

}

func (g *Gopsutil) Os() (Iagent, error) {
	o := &OsInfo{}
	o.OSType = runtime.GOOS
	o.OSArch = runtime.GOARCH

	hostName, err := os.Hostname()

	if err != nil {
		return nil, fmt.Errorf("error getting hostname: %v", err)
	}
	o.Hostname = hostName

	g.SystemInfo.OsInfo = o

	return g, nil
}
