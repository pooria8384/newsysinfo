// components/types.go
package sysinfo

type Printable interface {
	FetchInfo() error
	GetInfo() error
}

type SystemInfo struct {
	OSInfo
	RAMInfo
	DiskInfo
	CPUInfo
}
