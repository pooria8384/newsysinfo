// components/osinfo.go
package sysinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

type OSInfo struct {
	OSType   string `json:"os_type"`
	OSArch   string `json:"os_arch"`
	Hostname string `json:"hostname"`
}

func (o *OSInfo) FetchInfo() error {

	o.OSType = runtime.GOOS
	o.OSArch = runtime.GOARCH

	hostName, err := os.Hostname()

	if err != nil {
		return fmt.Errorf("error getting hostname: %v", err)
	}
	o.Hostname = hostName
	return nil
}

func (o *OSInfo) GetInfo() error {
	jsonData, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling OSInfo to JSON: %v", err)
	}
	fmt.Printf("%v", string(jsonData))
	return nil
}
