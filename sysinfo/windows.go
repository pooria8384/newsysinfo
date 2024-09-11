package sysinfo

import (
	"log"
)

type processors struct {
	Name                      string
	NumberOfLogicalProcessors uint16
}

type memory struct {
	Capacity uint64
}

type disk struct {
	Model string
	Size  uint64
}

type usb struct {
	Name string
}

type monitor struct {
	Name         string
	ScreenWidth  uint32
	ScreenHeight uint32
}

type Windows struct {
	*SystemInfo
}

func NewWindows() Iagent {
	return &Windows{
		&SystemInfo{
			OsInfo:    &OsInfo{},
			DiskInfos: []DiskInfo{},
			CpuInfo:   &CpuInfo{},
			RamInfo:   &RamInfo{},
			USBDevs:   []USBDevs{},
		},
	}
}

func (w *Windows) Get() *SystemInfo {
	return w.SystemInfo
}

func (w *Windows) Cpu() error {
	return nil
}

func (w *Windows) Ram() error {
	return nil
}

func (w *Windows) Disk() error {
	return nil
}

func (w *Windows) Os() error {
	return nil
}

func (w *Windows) USB() error {
	return nil
}

func (w *Windows) Monitor() error {
	return nil
}

func (w *Windows) Do() {
	if err := w.Cpu(); err != nil {
		log.Println("Failed to fetch CPU info")
	}
	if err := w.Os(); err != nil {
		log.Println("Failed to fetch OS info")
	}
	if err := w.Disk(); err != nil {
		log.Println("Failed to fetch Disk info")
	}
	if err := w.Ram(); err != nil {
		log.Println("Failed to fetch RAM info")
	}
	if err := w.USB(); err != nil {
		log.Println("Failed to fetch USB info")
	}
	if err := w.Monitor(); err != nil {
		log.Println("Failed to fetch Monitors info")
	}
}
