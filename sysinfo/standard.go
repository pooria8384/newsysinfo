package sysinfo

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Standard struct {
	*SystemInfo
}

func NewStandard() Iagent {
	return &Standard{
		&SystemInfo{
			OsInfo:    &OsInfo{},
			DiskInfos: []DiskInfo{},
			CpuInfo:   &CpuInfo{},
			RamInfo:   &RamInfo{},
		},
	}
}

func (s *Standard) Cpu() (Iagent, error) {
	c := &CpuInfo{}

	switch runtime.GOOS {
	case "linux":
		output, err := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep 'model name' | head -1").Output()
		if err != nil {
			return nil, err
		}
		c.Modelname = strings.Split(strings.TrimSpace(string(output)), ": ")[1]
	case "darwin":
		output, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
		if err != nil {
			return nil, err
		}
		c.Modelname = strings.TrimSpace(string(output))
	case "windows":
		output, err := exec.Command("wmic", "cpu", "get", "name").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			c.Modelname = strings.TrimSpace(lines[1])
		}
	default:
		c.Modelname = "Unknown CPU Model"
	}

	c.Cores = uint32(runtime.NumCPU())
	s.SystemInfo.CpuInfo = c
	return s, nil
}

func (s *Standard) Ram() (Iagent, error) {
	r := &RamInfo{}

	switch runtime.GOOS {
	case "linux":
		output, err := exec.Command("sh", "-c", "grep MemTotal /proc/meminfo").Output()
		if err != nil {
			return nil, err
		}
		memTotalStr := strings.Fields(string(output))[1]
		memTotal, err := strconv.ParseUint(memTotalStr, 10, 64)
		if err != nil {
			return nil, err
		}
		r.Total = uint32(memTotal)

		output, err = exec.Command("sh", "-c", "grep MemAvailable /proc/meminfo").Output()
		if err != nil {
			return nil, err
		}
		memAvailableStr := strings.Fields(string(output))[1]
		memAvailable, err := strconv.ParseUint(memAvailableStr, 10, 64)
		if err != nil {
			return nil, err
		}
		r.Available = uint32(memAvailable)

	case "darwin":
		output, err := exec.Command("sysctl", "hw.memsize").Output()
		if err != nil {
			return nil, err
		}
		memTotalStr := strings.TrimSpace(strings.Split(string(output), ": ")[1])
		memTotal, err := strconv.ParseUint(memTotalStr, 10, 64)
		if err != nil {
			return nil, err
		}
		r.Total = uint32(memTotal)

		output, err = exec.Command("vm_stat").Output()
		if err != nil {
			return nil, err
		}
		vmStat := strings.Split(string(output), "\n")
		pageFreeStr := strings.Fields(vmStat[1])[2]
		pageFree, err := strconv.ParseUint(pageFreeStr[:len(pageFreeStr)-1], 10, 64)
		if err != nil {
			return nil, err
		}
		pageSize, err := exec.Command("sysctl", "-n", "hw.pagesize").Output()
		if err != nil {
			return nil, err
		}
		pageSizeInt, err := strconv.ParseUint(strings.TrimSpace(string(pageSize)), 10, 64)
		if err != nil {
			return nil, err
		}
		r.Available = uint32(pageFree * pageSizeInt)

	case "windows":
		output, err := exec.Command("powershell", "Get-WmiObject", "Win32_OperatingSystem").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "TotalVisibleMemorySize") {
				memTotalStr := strings.Fields(line)[2]
				memTotal, err := strconv.ParseUint(memTotalStr, 10, 64)
				if err != nil {
					return nil, err
				}
				r.Total = uint32(memTotal)
			}
			if strings.Contains(line, "FreePhysicalMemory") {
				memAvailableStr := strings.Fields(line)[2]
				memAvailable, err := strconv.ParseUint(memAvailableStr, 10, 64)
				if err != nil {
					return nil, err
				}
				r.Available = uint32(memAvailable)
			}
		}
	}

	r.Used = r.Total - r.Available
	r.UsedPercent = (float64(r.Used) / float64(r.Total)) * 100

	s.SystemInfo.RamInfo = r
	return s, nil
}
func (s *Standard) Disk() (Iagent, error) {
	disks := []DiskInfo{}

	switch runtime.GOOS {
	case "linux":
		output, err := exec.Command("sh", "-c", "df -B1").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(output), "\n")

		for _, line := range lines[1:] { // Skip header line
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				device := fields[0]
				totalSize, err := strconv.Atoi(fields[1])
				if err != nil {
					return nil, err
				}
				freeSize, err := strconv.Atoi(fields[3])
				if err != nil {
					return nil, err
				}
				disks = append(disks, DiskInfo{
					Device:    device,
					TotalSize: uint32(totalSize),
					FreeSize:  uint32(freeSize),
				})
			}
		}
	case "darwin":
		output, err := exec.Command("sh", "-c", "df -B1").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(output), "\n")

		for _, line := range lines[1:] { // Skip header line
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				device := fields[0]
				totalSize, err := strconv.Atoi(fields[1])
				if err != nil {
					return nil, err
				}
				freeSize, err := strconv.Atoi(fields[3])
				if err != nil {
					return nil, err
				}
				disks = append(disks, DiskInfo{
					Device:    device,
					TotalSize: uint32(totalSize),
					FreeSize:  uint32(freeSize),
				})
			}
		}
	case "windows":
		output, err := exec.Command("powershell", "Get-PSDrive", "-PSProvider", "FileSystem").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(output), "\n")

		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				device := fields[0]
				totalSizeStr := fields[2]
				freeSizeStr := fields[3]
				totalSize, err := strconv.Atoi(totalSizeStr)
				if err != nil {
					return nil, err
				}
				freeSize, err := strconv.Atoi(freeSizeStr)
				if err != nil {
					return nil, err
				}
				disks = append(disks, DiskInfo{
					Device:    device,
					TotalSize: uint32(totalSize),
					FreeSize:  uint32(freeSize),
				})
			}
		}
	}

	s.SystemInfo.DiskInfos = disks
	return s, nil
}
func (s *Standard) Os() (Iagent, error) {
	o := &OsInfo{}
	o.OSType = runtime.GOOS
	o.OSArch = runtime.GOARCH

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	o.Hostname = hostName

	s.SystemInfo.OsInfo = o
	return s, nil
}

func (s *Standard) Get() *SystemInfo {
	return s.SystemInfo
}
