package sysinfo

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

func (s *Standard) Get() *SystemInfo {
	return s.SystemInfo
}

func (s *Standard) Cpu() (Iagent, error) {
	c := &CpuInfo{}
	c.Modelname = "..."
	c.Cores = 0
	s.SystemInfo.CpuInfo = c
	return s, nil
}

func (s *Standard) Ram() (Iagent, error) {
	r := &RamInfo{}
	r.Available = 0
	r.Total = 0
	r.Used = 0
	r.UsedPercent = 0
	return s, nil
}

func (s *Standard) Os() (Iagent, error) {
	o := &OsInfo{}
	o.Hostname = "..."
	o.OSArch = "..."
	o.OSType = "..."
	return s, nil
}

func (s *Standard) Disk() (Iagent, error) {
	d := &DiskInfo{}
	d.Device = "..."
	d.FreeSize = 0
	d.FreeSize = 0
	return s, nil
}
