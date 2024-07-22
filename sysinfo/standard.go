package sysinfo

type Standard struct {
	*SystemInfo
}

func (s *Standard) Scan() error {
	return nil
}

func (s *Standard) Get() *SystemInfo {
	return s.SystemInfo
}
