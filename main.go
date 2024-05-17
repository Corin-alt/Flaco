package main

import (
	"flaco/grpc_and_go/client"
	"flaco/grpc_and_go/serveur"
	"time"
)

func main() {
	//run server
	go serveur.Connect()

	//wait 10 seconds
	time.Sleep(5 * time.Second)

	//run client
	client.NewClient("0.0.0.0:8082")
}
