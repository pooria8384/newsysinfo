package main

import (
	"agent/sysinfo"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {

	scanner := sysinfo.NewScanner()
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

	file, err := os.Create("info.json")
	if err != nil {
		log.Println("There is an error: " + err.Error())
	}
	defer file.Close()

	_, err = file.Write(jsonresult)
	if err != nil {
		log.Println("There is an error: " + err.Error())
	}
}
