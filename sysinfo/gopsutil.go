package sysinfo

import (
	"agent/utils"
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
		cpuInfo.Cores = uint32(len(data))
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
	ram.Total = utils.KbToHumanReadable(uint(vmStat.VMallocTotal))
	ram.Available = utils.KbToHumanReadable(uint(vmStat.Available))
	ram.Used = utils.KbToHumanReadable(uint(vmStat.Used))
	ram.UsedPercent = fmt.Sprintf("%.2f %s", vmStat.UsedPercent, "%")

	g.SystemInfo.RamInfo = ram
	return nil
}

func (g *Gopsutil) Disk() error {
	disks := map[string]struct {
		total uint
		free  uint
	}{}

	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Errorf("error getting disk partitions: %v", err)
	}

	for _, partition := range partitions {
		diskInfo, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return fmt.Errorf("error getting disk usage info for %s: %v", partition.Mountpoint, err)
		}
		deviceName := disk.GetDiskSerialNumber(partition.Device)
		if existsDevice, ok := disks[deviceName]; ok {
			existsDevice.total += uint(diskInfo.Total)
			existsDevice.free += uint(diskInfo.Free)
			disks[deviceName] = existsDevice
		} else {
			disks[deviceName] = struct {
				total uint
				free  uint
			}{
				total: uint(diskInfo.Total),
				free:  uint(diskInfo.Free),
			}
		}
	}

	for dev, dsk := range disks {
		g.SystemInfo.DiskInfos = append(g.SystemInfo.DiskInfos,
			DiskInfo{
				Device:    dev,
				TotalSize: utils.KbToHumanReadable(dsk.total),
				FreeSize:  utils.KbToHumanReadable(dsk.free),
			},
		)
	}
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
