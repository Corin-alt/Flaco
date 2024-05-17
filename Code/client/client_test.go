package client

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDeviceDataFromFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Remove the temporary directory at the end of the test

	// Create two temporary JSON files in the temporary directory
	file1, err := os.Create(filepath.Join(tempDir, "file1.json"))
	if err != nil {
		t.Errorf("Unable to create temporary file: %v", err)
	}
	file1.Close()

	file2, err := os.Create(filepath.Join(tempDir, "file2.json"))
	if err != nil {
		t.Errorf("Unable to create temporary file: %v", err)
	}
	file2.Close()

	// Read data from files in the temporary directory
	paths, err := ReadDeviceDataFromFiles(tempDir)
	if err != nil {
		t.Errorf("Error reading files: %v", err)
	}

	// Check if the number of paths returned is correct
	if len(paths) != 2 {
		t.Fail()
		t.Logf("Expected number of files: 2, obtained: %d", len(paths))
	}
}

func TestGetDeviceData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Remove the temporary directory at the end of the test

	// Create a valid JSON file containing device data
	validJSON := `[{"device_name":"device1","operations":[{"type":"CREATE","has_succeeded":true}]}]`
	err = os.WriteFile(filepath.Join(tempDir, "valid.json"), []byte(validJSON), 0644)
	if err != nil {
		t.Errorf("Unable to write valid JSON file: %v", err)
	}

	// Create an invalid (empty) JSON file to test behavior with incorrect data
	invalidJSON := ""
	err = os.WriteFile(filepath.Join(tempDir, "invalid.json"), []byte(invalidJSON), 0644)
	if err != nil {
		t.Errorf("Unable to write invalid JSON file: %v", err)
	}

	// Retrieve device data from files in the temporary directory
	devices, err := GetDeviceData(tempDir)
	if err != nil {
		t.Failed()
	}

	// Check if the number of devices returned is correct
	if len(devices) != 1 {
		t.Fail()
		t.Logf("Expected number of devices: 1, obtained: %d", len(devices))
	}
}

func TestConvertDeviceDataToGRPCDevice(t *testing.T) {
	// Create device data for testing
	deviceData := DeviceData{
		DeviceName: "device1",
		Operations: []DeviceOperation{
			{Type: "CREATE", HasSucceeded: true},
			{Type: "DELETE", HasSucceeded: false},
		},
	}

	// Convert device data to GRPC format
	grpcDevice := ConvertDeviceDataToGRPCDevice(deviceData)

	// Check if the device name is correct
	if grpcDevice.DeviceName != "device1" {
		t.Fail()
		t.Logf("Expected device name: device1, obtained: %s", grpcDevice.DeviceName)
	}

	// Check if the number of operations is correct
	if len(grpcDevice.Operation) != 2 {
		t.Fail()
		t.Logf("Expected number of operations: 2, obtained: %d", len(grpcDevice.Operation))
	}

	// Check if the operations are correct
	if grpcDevice.Operation[0].Type != "CREATE" || !grpcDevice.Operation[0].HasSucceeded {
		t.Fail()
		t.Logf("Expected operation: {Type: \"CREATE\", HasSucceeded: true}, obtained: %+v", grpcDevice.Operation[0])
	}

	if grpcDevice.Operation[1].Type != "DELETE" || grpcDevice.Operation[1].HasSucceeded {
		t.Fail()
		t.Logf("Expected operation: {Type: \"DELETE\", HasSucceeded: false}, obtained: %+v", grpcDevice.Operation[1])
	}
}
