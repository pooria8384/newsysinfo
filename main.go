package main

import (
	"agent/services"
	"agent/sysinfo"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {

	currentInfo, err := hashFile("info.json")
	if err != nil {
		fmt.Println("Creating info.json file for first time")
	}

	scanner := sysinfo.NewScanner()
	scanner.Do()

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

	newInfo, err := hashFile("info.json")
	if err != nil {
		fmt.Printf("Failed to get file hashing info: %v", err)
	}

	if currentInfo != newInfo {
		sender := services.NewSender()
		sender.Send()
	}

}

func hashFile(filePath string) ([16]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return [16]byte{}, err
	}
	return md5.Sum(data), nil
}
