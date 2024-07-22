package main

import (
	"agent/sysinfo"
	"flag"
	"fmt"
	"log"
)

func main() {

	lib := flag.String("lib", "gopsutil", "Choose a library")

	flag.Parse()

	scanner := sysinfo.NewScanner(lib)
	if _, err := scanner.Cpu(); err != nil {
		log.Println("Failed to fetch CPU info")
	}
	if _, err := scanner.Os(); err != nil {
		log.Println("Failed to fetch OS info")
	}
	if _, err := scanner.Disk(); err != nil {
		log.Println("Failed to fetch Disk info")
	}
	if _, err := scanner.Ram(); err != nil {
		log.Println("Failed to fetch RAM info")
	}

	result := scanner.Get()

	fmt.Println(*result)

	// systemInfoType := []string{"OS", "Disk", "RAM", "Cpu"}

	// for _, sysType := range systemInfoType {
	// 	sysType, err := sysinfo.NewSysInfo(sysType)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	sysType.GetInfo()
	// 	sysType.FetchInfo()

	// }

	// osInfo, err := sysinfo.FactorySysInfo("OS")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// osInfo.GatherInfo()
	// osInfo.PrintInfo()

	// for _, componentType := range componetsType {
	// 	err, component := sysinfo.FactorySysInfo(componentType)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	err,_ = sysinfo.FactorySysInfo(componentType)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	component.Printable()
	// }

}
