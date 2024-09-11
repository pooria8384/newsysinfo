package sysinfo

import "runtime"

type CpuInfo struct {
	Modelname string `json:"model"`
	Cores     string `json:"cores"`
}

type DiskInfo struct {
	Device string `json:"device"`
}

type OsInfo struct {
	OSType   string `json:"type"`
	OSArch   string `json:"arch"`
	Hostname string `json:"hostname"`
}

type RamInfo struct {
	Total string `json:"total"`
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
	USBDevs   []USBDevs  `json:"usbs"`
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
		return NewWindows()
	default:
		return NewUnixLike()
	}
}
