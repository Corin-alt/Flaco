package serveur

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net"

	"flaco/grpc_and_go/flaco_grpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// Server struct represents the gRPC server
type Server struct {
	flaco_grpc.UnimplementedDayServiceServer // Embedding the unimplemented server for forward compatibility
}

// CompletedDeviceInfo holds the information about devices and operations, including totals and failures
type CompletedDeviceInfo struct {
	Devices    []*flaco_grpc.Device // List of device data
	NbTotalOp  int64                // Total number of operations
	NbOpFailed int64                // Number of failed operations
}

// SendDayInfoToServer processes the request from the client, stores data in the database, and returns a response
func (s *Server) SendDayInfoToServer(ctx context.Context, req *flaco_grpc.Request) (*flaco_grpc.Response, error) {
	err := StoreToDatabase(req) // Store the request data in the database
	if err != nil {
		return nil, err // Return an error if storage fails
	}
	return &flaco_grpc.Response{}, nil // Return an empty response on
}

// StoreToDatabase connects to the database and stores the device data and calculated values
func StoreToDatabase(req *flaco_grpc.Request) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	completedDevice := CalculateValues(req) // Calculate the total and failed operations

	// Store each device's information in the database
	for _, deviceInfo := range req.GetDevice() {
		for _, operation := range deviceInfo.Operation {
			state := "FAILED"
			if operation.HasSucceeded {
				state = "SUCCESS"
			}

			coll := client.Database("flaco").Collection(deviceInfo.DeviceName)
			_, err = coll.InsertOne(context.Background(), bson.M{
				"type":  operation.Type,
				"state": state,
			})
			if err != nil {
				return err // Return an error if insertion fails
			}
		}
	}

	// Store the calculated total and failed operations in the database
	// Insert aggregate operation information into the database
	collection := client.Database("flaco").Collection("OperationsInformation")
	_, err = collection.InsertOne(context.Background(), bson.M{
		"total":      completedDevice.NbTotalOp,
		"successful": completedDevice.NbTotalOp - completedDevice.NbOpFailed,
		"failed":     completedDevice.NbOpFailed,
	})
	if err != nil {
		return err
	}

	return nil
}

// CalculateValues calculates the total and failed operations from the request
func CalculateValues(req *flaco_grpc.Request) *CompletedDeviceInfo {
	nbTotal := 0
	nbFailed := 0
	// Iterate over each device and its operations to count total and failed operations
	for _, deviceInfo := range req.GetDevice() {
		for _, operation := range deviceInfo.Operation {
			if !operation.HasSucceeded {
				nbFailed++ // Increment the failed operations count
			}
			nbTotal++ // Increment the total operations count
		}
	}
	return &CompletedDeviceInfo{
		Devices:    req.GetDevice(), // List of devices from the request
		NbOpFailed: int64(nbFailed), // Total number of failed operations
		NbTotalOp:  int64(nbTotal),  // Total number of operations
	}
}

// Connect initializes the gRPC server and listens for incoming connections
func Connect() {
	listener, err := net.Listen("tcp", ":8082") // Create a TCP listener on port 8082
	if err != nil {
		panic(err) // Terminate if listener creation fails
	}

	s := grpc.NewServer()                             // Create a new gRPC server
	flaco_grpc.RegisterDayServiceServer(s, &Server{}) // Register the DayService server

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err) // Log and terminate if server fails to start
	}
}
