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

type UnixLike struct {
	*SystemInfo
}

func NewUnixLike() Iagent {
	return &UnixLike{
		&SystemInfo{
			OsInfo:    &OsInfo{},
			DiskInfos: []DiskInfo{},
			CpuInfo:   &CpuInfo{},
			RamInfo:   &RamInfo{},
			USBDevs:   []USBDevs{},
		},
	}
}

func (u *UnixLike) Get() *SystemInfo {
	return u.SystemInfo
}

func (u *UnixLike) Cpu() error {
	cpuInfo := &CpuInfo{}

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return fmt.Errorf("error getting cpu info: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if u.SystemInfo.CpuInfo.Cores != "" && u.SystemInfo.CpuInfo.Modelname != "" {
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
	u.SystemInfo.CpuInfo = cpuInfo
	return nil
}

func (u *UnixLike) Ram() error {
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
	u.SystemInfo.RamInfo = ram
	return nil
}

func (u *UnixLike) Disk() error {
	cmd := exec.Command("lsblk", "-o", "MODEL,SERIAL,SIZE", "-d", "-n")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Join(strings.Fields(line), " ")
		u.SystemInfo.DiskInfos = append(u.SystemInfo.DiskInfos, DiskInfo{Device: line})
	}

	return nil
}

func (u *UnixLike) Os() error {
	o := &OsInfo{}
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error getting uname: %v", err)
	}
	o.OSArch = strings.TrimSpace(string(output))
	o.OSType = runtime.GOOS

	hostName, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}
	o.Hostname = hostName

	u.SystemInfo.OsInfo = o

	return nil
}

func (u *UnixLike) USB() error {
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return fmt.Errorf("error getting USB devices: %v", err)
	}

	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			u.SystemInfo.USBDevs = append(u.SystemInfo.USBDevs,
				USBDevs{
					Device: line,
				},
			)
		}
	}
	return nil
}

func (u *UnixLike) Monitor() error {
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
			u.SystemInfo.Monitor = append(u.SystemInfo.Monitor,
				Monitors{
					Device: line,
				},
			)
		}
	}
	return nil
}

func (u *UnixLike) Do() {
	if err := u.Cpu(); err != nil {
		log.Println("Failed to fetch CPU info")
	}
	if err := u.Os(); err != nil {
		log.Println("Failed to fetch OS info")
	}
	if err := u.Disk(); err != nil {
		log.Println("Failed to fetch Disk info")
	}
	if err := u.Ram(); err != nil {
		log.Println("Failed to fetch RAM info")
	}
	if err := u.USB(); err != nil {
		log.Println("Failed to fetch USB info")
	}
	if err := u.Monitor(); err != nil {
		log.Println("Failed to fetch Monitors info")
	}
}
