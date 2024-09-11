package sysinfo

import (
	"agent/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/yusufpapurcu/wmi"
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
	cpuInfo := &CpuInfo{}

	var cpus []processors
	query := "SELECT Name, NumberOfCores, NumberOfLogicalProcessors FROM Win32_Processor"
	err := wmi.Query(query, &cpus)
	if err != nil {
		return fmt.Errorf("error getting cpu: %v", err)
	}

	cpuInfo.Modelname = cpus[0].Name
	cpuInfo.Cores = fmt.Sprintf("%d", cpus[0].NumberOfLogicalProcessors)
	w.SystemInfo.CpuInfo = cpuInfo
	return nil
}

func (w *Windows) Ram() error {
	ram := &RamInfo{}

	var memories []memory

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

	var disks []disk
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
	var usbs []usb
	query := "SELECT Name FROM Win32_PnpEntity WHERE DeviceID LIKE '%USB%'"
	err := wmi.Query(query, &usbs)
	if err != nil {
		return fmt.Errorf("error getting usb: %v", err)
	}

	for _, usb := range usbs {
		w.SystemInfo.USBDevs = append(w.SystemInfo.USBDevs,
			USBDevs{
				Device: usb.Name,
			},
		)
	}
	return nil
}

func (w *Windows) Monitor() error {
	var monitors []monitor
	query := "SELECT Name, ScreenWidth, ScreenHeight FROM Win32_DesktopMonitor"
	err := wmi.Query(query, &monitors)
	if err != nil {
		return fmt.Errorf("error getting monitor: %v", err)
	}
	for _, mon := range monitors {
		if strings.Contains(mon.Name, "Default") {
			continue
		}
		w.SystemInfo.Monitor = append(w.SystemInfo.Monitor,
			Monitors{
				Device: mon.Name + " " + strconv.Itoa(int(mon.ScreenWidth)) + "x" + strconv.Itoa(int(mon.ScreenHeight)),
			},
		)
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
