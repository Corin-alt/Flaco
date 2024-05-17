package main

import (
	"flaco/grpc_and_go/client"
	"flaco/grpc_and_go/serveur"
	"time"
)

func main() {
	go serveur.Connect()             // Start the gRPC server in a goroutine
	time.Sleep(5 * time.Second)      // Wait for 5 seconds to ensure the server is up and running
	client.NewClient("0.0.0.0:8082") // Connect the gRPC client to the server
}
