package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

func Connect(addr string) error {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	return nil
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
