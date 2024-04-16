package client

import (
	"context"
	"encoding/json"
	"flaco/grpc_and_go/flaco_grpc"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

type DeviceOperation struct {
	Type         string `json:"type"`
	HasSucceeded bool   `json:"has_succeeded"`
}

type DeviceData struct {
	DeviceName string            `json:"device_name"`
	Operations []DeviceOperation `json:"operations"`
}

//localhost:8080
func Connect(addr string) {

	devices := GetDeviceData("./donnees/")

	conn, err := grpc.Dial(addr)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := flaco_grpc.NewDayServiceClient(&grpc.ClientConn{})

	client.SendDayInfoToServer(context.Background(), &flaco_grpc.Request{
		Device: devices,
	})
}

func ReadDeviceDataFromFiles(pathString string) ([]string, error) {
	var paths []string

	err := filepath.Walk(pathString, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			fmt.Printf("File found: %s\n", path)
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func GetDeviceData(path string) []DeviceData {
	var devices []DeviceData

	jsonData, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		return nil
	}

	err = json.Unmarshal(jsonData, &devices)
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return nil
	}

	return devices
}
