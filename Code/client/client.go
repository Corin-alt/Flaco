package client

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"path/filepath"

	"flaco/grpc_and_go/flaco_grpc"
	"google.golang.org/grpc"
)

// DeviceOperation represents a device operation with its type and success status
type DeviceOperation struct {
	Type         string `json:"type"`          // Type of the operation
	HasSucceeded bool   `json:"has_succeeded"` // Success status of the operation
}

// DeviceData represents data for a device including its name and operations performed
type DeviceData struct {
	DeviceName string            `json:"device_name"` // Name of the device
	Operations []DeviceOperation `json:"operations"`  // List of operations performed on the device
}

// NewClient initializes a new gRPC client, reads device data, converts it, and sends it to the server
func NewClient(addr string) {
	// Get device data from the specified directory
	devices, err := GetDeviceData("./donnees/")
	if err != nil {
		log.Fatalf("Error getting device data: %v", err)
	}

	// Establish a connection to the gRPC server using insecure credentials
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Fail to dial: %v", err) // Log and terminate if connection fails
	}
	defer conn.Close() // Ensure the connection is closed when done

	// Create a new DayService client from the connection
	client := flaco_grpc.NewDayServiceClient(conn)

	// Convert the device data to the format expected by the gRPC service
	var devicesConverted []*flaco_grpc.Device
	for _, device := range devices {
		devicesConverted = append(devicesConverted, ConvertDeviceDataToGRPCDevice(device))
	}

	// Send the converted device data to the server
	_, err = client.SendDayInfoToServer(context.Background(), &flaco_grpc.Request{Device: devicesConverted})
	if err != nil {
		fmt.Println("[LOGS] => Error sending data to server:", err) // Print error if the request fails
	}
}

// ReadDeviceDataFromFiles reads all file paths in the given directory
func ReadDeviceDataFromFiles(pathString string) ([]string, error) {
	var paths []string
	// Walk through the directory and collect all file paths
	err := filepath.Walk(pathString, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("[LOGS] => Error during file walk:", err)
			return err // Return the error encountered during file traversal
		}
		if !info.IsDir() {
			fmt.Printf("[LOGS] => File found: %s\n", path)
			paths = append(paths, path) // Add file path to the list if it's not a directory
		}
		return nil
	})
	if err != nil {
		return nil, err // Return paths and any error encountered
	}
	return paths, nil
}

// GetDeviceData reads device data from a specified directory and unmarshals the JSON content into a slice of DeviceData
func GetDeviceData(dirPath string) ([]DeviceData, error) {
	paths, err := ReadDeviceDataFromFiles(dirPath)
	if err != nil {
		return nil, err
	}

	var devices []DeviceData
	for _, path := range paths {
		// Read the JSON file from the given path
		jsonData, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("[LOGS] => Error reading file:", err)
			return nil, err
		}
		// Unmarshal the JSON data into a slice of DeviceData
		var deviceData []DeviceData
		err = json.Unmarshal(jsonData, &deviceData)
		// Append the device data to the devices slice
		devices = append(devices, deviceData...)
	}
	return devices, nil
}

// ConvertDeviceDataToGRPCDevice converts DeviceData to the gRPC Device type
func ConvertDeviceDataToGRPCDevice(deviceData DeviceData) *flaco_grpc.Device {
	var operations []*flaco_grpc.Operation
	// Convert each operation to the gRPC Operation type
	for _, op := range deviceData.Operations {
		operations = append(operations, &flaco_grpc.Operation{
			Type:         op.Type,
			HasSucceeded: op.HasSucceeded,
		})
	}
	// Return a new gRPC Device with the converted operations
	return &flaco_grpc.Device{
		DeviceName: deviceData.DeviceName,
		Operation:  operations,
	}
}
