package sysinfo

import (
	"agent/utils"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

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

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return fmt.Errorf("error getting cpu info: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if w.SystemInfo.CpuInfo.Cores != "" && w.SystemInfo.CpuInfo.Modelname != "" {
			break
		}
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			fields := strings.Split(line, ":")
			if len(fields) > 1 {
				cpuInfo.Modelname = strings.TrimSpace(fields[1])
			}
		}
		if strings.HasPrefix(line, "siblings") {
			fields := strings.Split(line, ":")
			if len(fields) > 1 {
				cpuInfo.Cores = strings.TrimSpace(fields[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error getting cpu info: %v", err)
	}
	w.SystemInfo.CpuInfo = cpuInfo
	return nil
}

func (w *Windows) Ram() error {
	ram := &RamInfo{}

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return fmt.Errorf("error getting meminfo: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				var totalRAM uint64
				fmt.Sscanf(fields[1], "%d", &totalRAM)
				ram.Total = utils.ToHuman(float64(totalRAM), 1)
			}
		}
	}
	w.SystemInfo.RamInfo = ram
	return nil
}

func (w *Windows) Disk() error {
	cmd := exec.Command("lsblk", "-o", "MODEL,SERIAL,SIZE", "-d", "-n")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Join(strings.Fields(line), " ")
		w.SystemInfo.DiskInfos = append(w.SystemInfo.DiskInfos, DiskInfo{Device: line})
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
