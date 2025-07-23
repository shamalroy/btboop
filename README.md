# btboop

A lightweight Bluetooth switch HTTP server for macOS

## Overview

This Go program provides a minimal HTTP server to manage Bluetooth connections for multiple devices (e.g., Apple Magic Keyboard and Trackpad).

## Features

- `GET /status` → Returns current connection status of devices
- `PUT /connect` → Attempts to connect devices
- `PUT /disconnect` → Attempts to disconnect devices

## Requirements

- **blueutil CLI tool** (https://github.com/toy/blueutil)
    - Install via Homebrew: `brew install blueutil`
- Devices must be previously paired with the host Mac
- Touch ID devices (e.g., Magic Keyboard with Touch ID) must be paired once via USB for full functionality
- MAC addresses must be configured in the `devices` map in `main.go`

## Setup

### 1. Find MAC Addresses

Run `blueutil --paired` in Terminal to list known Bluetooth devices and their MAC addresses. Copy the MAC address for each device and paste them into the `devices` map in `main.go`.

### 2. Build and Run

```bash
# Clone or download the project
cd btboop

# Initialize Go module (if not already done)
go mod init btboop

# Build the binary
go build -o btboop

# Run the server
./btboop
```
The server will start at [http://localhost:5151](http://localhost:5151)


## API Usage
### Check Status
``` bash
curl http://localhost:5151/status
```
### Connect Devices
``` bash
curl -X PUT http://localhost:5151/connect
```
### Disconnect Devices
``` bash
curl -X PUT http://localhost:5151/disconnect
```
## Configuration
Edit the `devices` map in with your device MAC addresses: `main.go`
``` go
var devices = map[string]string{
    "keyboard": "XX-XX-XX-XX-XX-XX",
    "trackpad": "YY-YY-YY-YY-YY-YY",
}
```
