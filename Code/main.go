package main

import (
	"flaco/grpc_and_go/client"
	"flaco/grpc_and_go/serveur"
	"time"
)

func main() {
	println("[LOGS] => Server launch...")

	go serveur.Connect() // Start the gRPC server in a goroutine

	println("[LOGS] => Wait 5 seconds before the client connection...")

	time.Sleep(5 * time.Second) // Wait for 5 seconds to ensure the server is up and running

	client.NewClient("0.0.0.0:8082") // Connect the gRPC client to the server

	println("[LOGS] => Client disconnected.")
}
