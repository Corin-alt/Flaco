package client

type DeviceOperation struct {
	Type         string `json:"type"`
	HasSucceeded bool   `json:"has_succeeded"`
}

type DeviceData struct {
	DeviceName string            `json:"device_name"`
	Operations []DeviceOperation `json:"operations"`
}
