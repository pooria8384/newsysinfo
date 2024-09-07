package sysinfo

import "runtime"

type CpuInfo struct {
	Modelname string `json:"model"`
	Cores     uint32 `json:"cores"`
}

type DiskInfo struct {
	Device    string `json:"device"`
	TotalSize string `json:"totalsize"`
	FreeSize  string `json:"freesize"`
}

type OsInfo struct {
	OSType   string `json:"type"`
	OSArch   string `json:"arch"`
	Hostname string `json:"hostname"`
}

type RamInfo struct {
	Total       string `json:"total"`
	Available   string `json:"available"`
	Used        string `json:"used"`
	UsedPercent string `json:"usedpercent"`
}

type USBDevs struct {
	Device string `json:"device"`
}

type SystemInfo struct {
	*OsInfo   `json:"os"`
	*RamInfo  `json:"ram"`
	DiskInfos []DiskInfo `json:"disk"`
	*CpuInfo  `json:"cpu"`
	USBDevs   []USBDevs
}

type Iagent interface {
	Cpu() error
	Ram() error
	Disk() error
	Os() error
	USB() error
	Get() *SystemInfo
	Do()
}

func NewScanner() Iagent {
	osType := runtime.GOOS
	switch osType {
	case "windows":
		return nil
		// 	return NewWindows()
	default:
		return NewUnixLike()
	}
}
