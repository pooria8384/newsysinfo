package sysinfo

import (
	"fmt"
)

func NewSysInfo(systeminfo string) (Printable, error) {
	switch systeminfo {
	case "OS":
		return &OSInfo{}, nil
	case "Disk":
		return &DiskInfo{}, nil
	case "RAM":
		return &RAMInfo{}, nil
	case "CPU":
		return &CPUInfo{}, nil
	default:
		return nil, fmt.Errorf("unknown component type")
	}
}

func CreateSystemInfo() (, error) {
	osInfo, err := NewSysInfo("OS")
	if err != nil {
		return nil, err
	}

	ramInfo, err := NewSysInfo("RAM")
	if err != nil {
		return nil, err
	}

	diskInfo, err := NewSysInfo("Disk")
	if err != nil {
		return nil, err
	}

	cpuInfo, err := NewSysInfo("CPU")
	if err != nil {
		return nil, err
	}

	return 
}

