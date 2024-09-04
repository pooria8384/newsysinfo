package main

import (
	"agent/sysinfo"
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

func main() {

	lib := flag.String("lib", "gopsutil", "Choose a library")

	flag.Parse()

	scanner := sysinfo.NewScanner(lib)
	if err := scanner.Cpu(); err != nil {
		log.Println("Failed to fetch CPU info")
	}
	if err := scanner.Os(); err != nil {
		log.Println("Failed to fetch OS info")
	}
	if err := scanner.Disk(); err != nil {
		log.Println("Failed to fetch Disk info")
	}
	if err := scanner.Ram(); err != nil {
		log.Println("Failed to fetch RAM info")
	}

	result := scanner.Get()

	jsonresult, err := json.MarshalIndent(result, "", " ")

	if err != nil {
		fmt.Printf("Failed to marshal system info: %v", err)
	}
	_ = jsonresult
	// fmt.Println(string(jsonresult))
}
