package sysinfo

import (
	"agent/utils"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type UnixLike struct {
	*SystemInfo
}

func NewUnixLike() Iagent {
	return &UnixLike{
		&SystemInfo{
			OsInfo:    &OsInfo{},
			DiskInfos: []DiskInfo{},
			CpuInfo:   &CpuInfo{},
			RamInfo:   &RamInfo{},
			USBDevs:   []USBDevs{},
		},
	}
}

func (u *UnixLike) Get() *SystemInfo {
	return u.SystemInfo
}

func (u *UnixLike) Cpu() error {
	cpuInfo := &CpuInfo{}

	data, err := cpu.Info()
	if err != nil {
		return fmt.Errorf("error getting CPU info: %v", err)
	}
	if len(data) > 0 {
		cpuInfo.Modelname = data[0].ModelName
		cpuInfo.Cores = uint32(len(data))
	}
	u.SystemInfo.CpuInfo = cpuInfo
	return nil
}

func (u *UnixLike) Ram() error {
	ram := &RamInfo{}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("error getting memory info: %v", err)
	}
	ram.Total = utils.ToHuman(float32(vmStat.VMallocTotal/1024), 0)
	ram.Available = utils.ToHuman(float32(vmStat.Available), 0)
	ram.Used = utils.ToHuman(float32(vmStat.Used), 0)
	ram.UsedPercent = fmt.Sprintf("%.2f %s", vmStat.UsedPercent, "%")

	u.SystemInfo.RamInfo = ram
	return nil
}

func (u *UnixLike) Disk() error {
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
		u.SystemInfo.DiskInfos = append(u.SystemInfo.DiskInfos,
			DiskInfo{
				Device:    dev,
				TotalSize: utils.ToHuman(float32(dsk.total), 0),
				FreeSize:  utils.ToHuman(float32(dsk.free), 0),
			},
		)
	}
	return nil
}

func (u *UnixLike) Os() error {
	o := &OsInfo{}
	o.OSType = runtime.GOOS
	o.OSArch = runtime.GOARCH

	hostName, err := os.Hostname()

	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}
	o.Hostname = hostName

	u.SystemInfo.OsInfo = o

	return nil
}

func (u *UnixLike) USB() error {
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return fmt.Errorf("error getting USB devices: %v", err)
	}

	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			u.SystemInfo.USBDevs = append(u.SystemInfo.USBDevs,
				USBDevs{
					Device: line,
				},
			)
		}
	}
	return nil
}

func (u *UnixLike) Monitor() error {
	cmd := exec.Command("xrandr")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error getting monitors: %v", err)
	}

	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, " connected") {
			u.SystemInfo.Monitor = append(u.SystemInfo.Monitor,
				Monitors{
					Device: line,
				},
			)
		}
	}
	return nil
}

func (u *UnixLike) Do() {
	if err := u.Cpu(); err != nil {
		log.Println("Failed to fetch CPU info")
	}
	if err := u.Os(); err != nil {
		log.Println("Failed to fetch OS info")
	}
	if err := u.Disk(); err != nil {
		log.Println("Failed to fetch Disk info")
	}
	if err := u.Ram(); err != nil {
		log.Println("Failed to fetch RAM info")
	}
	if err := u.USB(); err != nil {
		log.Println("Failed to fetch USB info")
	}
	if err := u.Monitor(); err != nil {
		log.Println("Failed to fetch Monitors info")
	}
}
