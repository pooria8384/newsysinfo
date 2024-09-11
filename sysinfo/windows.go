package sysinfo

import (
	"agent/utils"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

type Processors struct {
	Name                      string
	NumberOfLogicalProcessors uint16
}

type Memory struct {
	Capacity uint64
}

type Disk struct {
	Model string
	Size  uint64
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
	cpuInfo := &CpuInfo{}

	var cpus []Processors
	query := "SELECT Name, NumberOfCores, NumberOfLogicalProcessors FROM Win32_Processor"
	err := wmi.Query(query, &cpus)
	if err != nil {
		return fmt.Errorf("error getting cpu: %v", err)
	}

	fmt.Println(cpus)
	cpuInfo.Modelname = cpus[0].Name
	cpuInfo.Cores = fmt.Sprintf("%d", cpus[0].NumberOfLogicalProcessors)
	w.SystemInfo.CpuInfo = cpuInfo
	return nil
}

func (w *Windows) Ram() error {
	ram := &RamInfo{}

	var memories []Memory

	query := "SELECT Capacity FROM Win32_PhysicalMemory"
	err := wmi.Query(query, &memories)
	if err != nil {
		return fmt.Errorf("error getting memories: %v", err)
	}
	var total uint64
	for _, m := range memories {
		total += m.Capacity
	}
	ram.Total = utils.ToHuman(float64(total), 0)
	w.SystemInfo.RamInfo = ram
	return nil
}

func (w *Windows) Disk() error {

	var disks []Disk
	query := "SELECT Model, Size FROM Win32_DiskDrive"
	err := wmi.Query(query, &disks)
	if err != nil {
		return fmt.Errorf("error getting disks: %v", err)
	}

	for _, d := range disks {
		dd := fmt.Sprintf("%s %s", d.Model, utils.ToHuman(float64(d.Size), 0))
		w.SystemInfo.DiskInfos = append(w.SystemInfo.DiskInfos, DiskInfo{Device: dd})
	}
	return nil
}

func (w *Windows) Os() error {
	o := &OsInfo{}
	cmd := exec.Command("wmic", "os", "get", "OSArchitecture")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error getting wmic: %v", err)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		o.OSArch = strings.TrimSpace(lines[1])
	}
	o.OSType = runtime.GOOS

	hostName, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}
	o.Hostname = hostName

	w.SystemInfo.OsInfo = o

	return nil
}

func (w *Windows) USB() error {
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return fmt.Errorf("error getting USB devices: %v", err)
	}

	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			w.SystemInfo.USBDevs = append(w.SystemInfo.USBDevs,
				USBDevs{
					Device: line,
				},
			)
		}
	}
	return nil
}

func (w *Windows) Monitor() error {
	cmd := exec.Command("xrandr")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error getting monitors: %v", err)
	}

	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, " connected") {
			w.SystemInfo.Monitor = append(w.SystemInfo.Monitor,
				Monitors{
					Device: line,
				},
			)
		}
	}
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
