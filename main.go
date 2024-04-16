package main

import (
	"flaco/grpc_and_go/client"
	"fmt"
	"log"
)

func main() {
	paths, err := client.ReadDeviceDataFromFiles("./data/")

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(paths); i++ {
		devices := client.GetDeviceData(paths[i])
		fmt.Println("File : " + paths[i])
		fmt.Println(devices)
	}

	client.Connect("localhost:8080")

}
