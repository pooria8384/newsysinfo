package sysinfo

type CpuInfo struct {
	Modelname string `json:"model"`
	Cores     int32  `json:"cores"`
}

type DiskInfo struct {
	Device    string `json:"diskinfo_device"`
	TotalSize uint64 `json:"diskinfo_totalsize"`
	FreeSize  uint64 `json:"diskinfo_freesize"`
}

type OsInfo struct {
	OSType   string `json:"os_type"`
	OSArch   string `json:"os_arch"`
	Hostname string `json:"hostname"`
}

type RamInfo struct {
	Total       uint64  `json:"total_ram"`
	Available   uint64  `json:"available_ram"`
	Used        uint64  `json:"used_ram"`
	UsedPercent float64 `json:"usedpercent_ram"`
}

type SystemInfo struct {
	*OsInfo   `json:"os"`
	*RamInfo  `json:"ram"`
	*DiskInfo `json:"disk"`
	*CpuInfo  `json:"cpu"`
}

type Iagent interface {
	Cpu() (Iagent, error)
	Ram() (Iagent, error)
	Disk() (Iagent, error)
	Os() (Iagent, error)
	Get() *SystemInfo
}

func NewScanner(lib *string) Iagent {
	switch *lib {
	case "gopsutil":
		return NewGopsutil()
	case "standard":
		return NewStandard()
	default:
		return NewGopsutil()
	}
}
