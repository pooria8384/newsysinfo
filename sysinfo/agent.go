package sysinfo

type CpuInfo struct {
	Modelname string `json:"model"`
	Cores     int    `json:"cores"`
}

type DiskInfo struct {
	Device    string `json:"device"`
	TotalSize int    `json:"totalsize"`
	FreeSize  int    `json:"freesize"`
}

type OsInfo struct {
	OSType   string `json:"type"`
	OSArch   string `json:"arch"`
	Hostname string `json:"hostname"`
}

type RamInfo struct {
	Total       int     `json:"total"`
	Available   int     `json:"available"`
	Used        int     `json:"used"`
	UsedPercent float64 `json:"usedpercent"`
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
