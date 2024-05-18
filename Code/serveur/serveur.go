package serveur

import (
	"context"
	"flaco/grpc_and_go/flaco_grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Server struct represents the gRPC server
type Server struct {
	flaco_grpc.UnimplementedDayServiceServer // Embedding the unimplemented server for forward compatibility
}

// DeviceStat struct holds statistics about device operations
type DeviceStat struct {
	DeviceName  string // Device name
	NbTotalOp   int64  // Total number of operations
	NbOpSuccess int64  // Number of successful operations
	NbOpFailed  int64  // Number of failed operations
}

// SendDayInfoToServer processes the request from the client, stores data in the database, and returns a response
func (s *Server) SendDayInfoToServer(ctx context.Context, req *flaco_grpc.Request) (*flaco_grpc.Response, error) {
	err := StoreToDatabase(req) // Store the request data in the database
	if err != nil {
		return nil, err // Return an error if storage fails
	}
	return &flaco_grpc.Response{}, nil // Return an empty response
}

// StoreToDatabase connects to the database and stores the device data and calculated values
func StoreToDatabase(req *flaco_grpc.Request) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatal(err) // Log and exit if database connection fails
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err) // Log and exit if database disconnection fails
		}
	}()

	println("[LOGS] => Storing information into database...")

	statCollection := client.Database("flaco").Collection("StatByDevice")

	// Store each device's information in the database
	for _, deviceInfo := range req.GetDevice() {
		statDevice := GetDeviceStat(deviceInfo)
		for _, operation := range deviceInfo.Operation {
			state := "FAILED"
			if operation.HasSucceeded {
				state = "SUCCESS"
			}

			// Insert operation details into a collection named after the device
			coll := client.Database("flaco").Collection(deviceInfo.DeviceName)
			_, err = coll.InsertOne(context.Background(), bson.M{
				"type":  operation.Type,
				"state": state,
			})
			if err != nil {
				return err // Return an error if insertion fails
			}
		}

		// Update the statistics collection with the device's operations count
		filter := bson.M{"name": statDevice.DeviceName}
		update := bson.M{
			"$inc": bson.M{
				"total":      statDevice.NbTotalOp,
				"successful": statDevice.NbOpSuccess,
				"failed":     statDevice.NbOpFailed,
			},
			"$setOnInsert": bson.M{
				"device": statDevice.DeviceName,
			},
		}
		opts := options.Update().SetUpsert(true)

		_, err = statCollection.UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			return err // Return an error if update fails
		}
	}

	return nil
}

// GetDeviceStat calculates the statistics for a given device
func GetDeviceStat(device *flaco_grpc.Device) *DeviceStat {
	nbTotal := 0
	nbFailed := 0
	nbSuccess := 0

	// Iterate over each operation of the device to count total, successful, and failed operations
	for _, operation := range device.Operation {
		if !operation.HasSucceeded {
			nbFailed++ // Increment the failed operations count
		} else {
			nbSuccess++ // Increment the successful operations count
		}
		nbTotal++ // Increment the total operations count
	}

	// Return the calculated statistics for the device
	return &DeviceStat{
		DeviceName:  device.DeviceName,
		NbOpSuccess: int64(nbSuccess),
		NbOpFailed:  int64(nbFailed),
		NbTotalOp:   int64(nbTotal),
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
