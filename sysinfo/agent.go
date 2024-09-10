package sysinfo

import "runtime"

type CpuInfo struct {
	Modelname string `json:"model"`
	Cores     string `json:"cores"`
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
	Total string `json:"total"`
	Free  string `json:"free"`
}

type USBDevs struct {
	Device string `json:"device"`
}

type Monitors struct {
	Device string `json:"device"`
}

type SystemInfo struct {
	*OsInfo   `json:"os"`
	*RamInfo  `json:"ram"`
	DiskInfos []DiskInfo `json:"disk"`
	*CpuInfo  `json:"cpu"`
	USBDevs   []USBDevs
	Monitor   []Monitors `json:"monitor"`
}

type Iagent interface {
	Cpu() error
	Ram() error
	Disk() error
	Os() error
	USB() error
	Monitor() error
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
