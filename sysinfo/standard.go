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
			OsInfo:   &OsInfo{},
			DiskInfo: &DiskInfo{},
			CpuInfo:  &CpuInfo{},
			RamInfo:  &RamInfo{},
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

	c.Cores = int(runtime.NumCPU())
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
		r.Total = int(memTotal * 1024)

		output, err = exec.Command("sh", "-c", "grep MemAvailable /proc/meminfo").Output()
		if err != nil {
			return nil, err
		}
		memAvailableStr := strings.Fields(string(output))[1]
		memAvailable, err := strconv.ParseUint(memAvailableStr, 10, 64)
		if err != nil {
			return nil, err
		}
		r.Available = int(memAvailable * 1024)

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
		r.Total = int(memTotal)

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
		r.Available = int(pageFree * pageSizeInt)

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
				r.Total = int(memTotal * 1024)
			}
			if strings.Contains(line, "FreePhysicalMemory") {
				memAvailableStr := strings.Fields(line)[2]
				memAvailable, err := strconv.ParseUint(memAvailableStr, 10, 64)
				if err != nil {
					return nil, err
				}
				r.Available = int(memAvailable * 1024)
			}
		}
	}

	r.Used = r.Total - r.Available
	r.UsedPercent = (float64(r.Used) / float64(r.Total)) * 100

	s.SystemInfo.RamInfo = r
	return s, nil
}

func (s *Standard) Disk() (Iagent, error) {
	d := &DiskInfo{}

	switch runtime.GOOS {
	case "linux", "darwin":
		output, err := exec.Command("sh", "-c", "df -k --output=size,avail / | tail -1").Output()
		if err != nil {
			return nil, err
		}
		fields := strings.Fields(string(output))
		if len(fields) >= 2 {
			totalSize, err := strconv.ParseUint(fields[0], 10, 64)
			if err != nil {
				return nil, err
			}
			d.TotalSize = int(totalSize * 1024)

			freeSize, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return nil, err
			}
			d.FreeSize = int(freeSize * 1024)
		}
	case "windows":
		output, err := exec.Command("powershell", "Get-PSDrive", "-PSProvider", "FileSystem").Output()
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "C") {
				fields := strings.Fields(line)
				totalSizeStr := fields[1]
				freeSizeStr := fields[2]
				totalSize, err := strconv.ParseUint(totalSizeStr, 10, 64)
				if err != nil {
					return nil, err
				}
				freeSize, err := strconv.ParseUint(freeSizeStr, 10, 64)
				if err != nil {
					return nil, err
				}
				d.TotalSize = int(totalSize * 1024)
				d.FreeSize = int(freeSize * 1024)
			}
		}
	}

	s.SystemInfo.DiskInfo = d
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
