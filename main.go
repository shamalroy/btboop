package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

// DeviceStatus represents the Bluetooth connection status for a device.
type DeviceStatus struct {
	Device    string `json:"device"`
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

// StatusResponse holds an array of device statuses for JSON output.
type StatusResponse struct {
	Devices []DeviceStatus `json:"devices"`
}

// Replace with your actual MAC addresses for Bluetooth devices.
var devices = map[string]string{
	"keyboard": "XX-XX-XX-XX-XX-XX",
	"trackpad": "YY-YY-YY-YY-YY-YY",
}

// checkConnection returns whether the given MAC address is connected.
func checkConnection(mac string) (bool, error) {
	cmd := exec.Command("blueutil", "--is-connected", mac)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return string(output) == "1\n", nil
}

// connectDevice attempts to connect the device at the given MAC address.
func connectDevice(mac string) error {
	cmd := exec.Command("blueutil", "--connect", mac)
	return cmd.Run()
}

// disconnectDevice attempts to disconnect the device at the given MAC address.
func disconnectDevice(mac string) error {
	cmd := exec.Command("blueutil", "--disconnect", mac)
	return cmd.Run()
}

// getStatus handles GET /status to report connection state for all tracked devices.
func getStatus(w http.ResponseWriter, r *http.Request) {
	var response StatusResponse
	for name, mac := range devices {
		connected, err := checkConnection(mac)
		status := DeviceStatus{Device: name, Connected: connected}
		if err != nil {
			status.Error = err.Error()
		}
		response.Devices = append(response.Devices, status)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// connectDevices handles PUT /connect to connect all devices.
func connectDevices(w http.ResponseWriter, r *http.Request) {
	var response StatusResponse
	for name, mac := range devices {
		err := connectDevice(mac)
		connected, _ := checkConnection(mac)
		status := DeviceStatus{Device: name, Connected: connected}
		if err != nil {
			status.Error = fmt.Sprintf("Failed to connect: %v", err)
		}
		response.Devices = append(response.Devices, status)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// disconnectDevices handles PUT /disconnect to disconnect all devices.
func disconnectDevices(w http.ResponseWriter, r *http.Request) {
	var response StatusResponse
	for name, mac := range devices {
		err := disconnectDevice(mac)
		connected, _ := checkConnection(mac)
		status := DeviceStatus{Device: name, Connected: connected}
		if err != nil {
			status.Error = fmt.Sprintf("Failed to disconnect: %v", err)
		}
		response.Devices = append(response.Devices, status)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// main sets up the HTTP server and route handlers.
func main() {
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			connectDevices(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
		}
	})

	http.HandleFunc("/disconnect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			disconnectDevices(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
		}
	})

	log.Println("Bluetooth switch server started on :5151")
	log.Fatal(http.ListenAndServe(":5151", nil))
}
