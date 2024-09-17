package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Sender interface {
	GetAccessToken() error
	Send()
}

type Agent struct {
}

func NewSender() Sender {
	return &Agent{}
}

func (a *Agent) GetAccessToken() error {
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     os.Getenv("ASSET_AUTH_ID"),
		"client_secret": os.Getenv("ASSET_AUTH_SECRET"),
		"scope":         "",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Error to marshal data: ", err.Error())
		return err
	}
	authURL := os.Getenv("ASSET_AUTH_URL")

	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln("Error to get token: ", err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error to read body: ", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Failed to read access token")
		err := errors.New("Failed to read access token")
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalln("Failed to unmarshal access token: ", err.Error())
		return err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		log.Fatalln("Access Token not found in respose: ", err.Error())
		return err
	}

	// a.AccessToken = accessToken
	fmt.Println(accessToken)
	return nil
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
