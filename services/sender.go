package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Sender interface {
	Send()
}

type Agent struct {
}

func NewSender() Sender {
	return &Agent{}
}

func (a *Agent) Send() {
	api := os.Getenv("ASSET_AGENT_API")
	token := os.Getenv("ASSET_AGENT_TOKEN")

	requestBody, _ := json.Marshal(map[string]string{
		"name":        "New Asset",
		"description": "Asset description",
	})
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error making request: ", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request: ", err.Error())
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error in reading response body: ", err.Error())
	}

	logFile, err := os.Create("agent.log")
	if err != nil {
		log.Println("Failed to create log file: ", err.Error())
	}
	defer logFile.Close()

	logFile.Write(body)

	fmt.Println(res.Status)
	// fmt.Println(string(body))
}
