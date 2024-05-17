package client

import (
	"context"
	"flaco/grpc_and_go/flaco_grpc"
	"flaco/grpc_and_go/serveur"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

// TestCalculateValuesWithNoFailedOperations tests the calculation of values when all operations have succeeded.
func TestCalculateValuesWithNoFailedOperations(t *testing.T) {
	// Creating a test request with successful operations
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{Operation: []*flaco_grpc.Operation{{HasSucceeded: true}, {HasSucceeded: true}}},
		},
	}

	// Calling CalculateValues function with the test request
	result := serveur.CalculateValues(req)
	expectedTotal := int64(2)
	expectedFailed := int64(0)

	// Checking if the calculated values match the expected values
	if result.NbTotalOp != expectedTotal {
		t.Errorf("Expected total operations: %d, got: %d", expectedTotal, result.NbTotalOp)
	}

	if result.NbOpFailed != expectedFailed {
		t.Errorf("Expected failed operations: %d, got: %d", expectedFailed, result.NbOpFailed)
	}
}

// TestCalculateValuesWithFailedOperations tests the calculation of values when some operations have failed.
func TestCalculateValuesWithFailedOperations(t *testing.T) {
	// Creating a test request with one successful and one failed operation
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{Operation: []*flaco_grpc.Operation{{HasSucceeded: true}, {HasSucceeded: false}}},
		},
	}

	// Calling CalculateValues function with the test request
	result := serveur.CalculateValues(req)
	expectedTotal := int64(2)
	expectedFailed := int64(1)

	// Checking if the calculated values match the expected values
	if result.NbTotalOp != expectedTotal {
		t.Errorf("Expected total operations: %d, got: %d", expectedTotal, result.NbTotalOp)
	}

	if result.NbOpFailed != expectedFailed {
		t.Errorf("Expected failed operations: %d, got: %d", expectedFailed, result.NbOpFailed)
	}
}

// TestStoreToDatabaseSuccessful tests storing data to a MongoDB database with a successful connection.
func TestStoreToDatabaseSuccessful(t *testing.T) {
	// Connecting to MongoDB database
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Dropping the existing 'flaco' database
	err = client.Database("flaco").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Creating a test request with a successful operation
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{HasSucceeded: true}}},
		},
	}

	// Storing data to the database
	err = serveur.StoreToDatabase(req)
	if err != nil {
		t.Errorf("Error storing to database: %v", err)
	}

	// Checking the number of documents in the collection
	collection := client.Database("flaco").Collection("test_device")
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		t.Errorf("Failed to count documents in collection: %v", err)
	}

	expectedCount := int64(1)
	if count != expectedCount {
		t.Errorf("Expected %d documents in collection, got %d", expectedCount, count)
	}

	// Dropping the existing 'flaco' database
	err = client.Database("flaco").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// TestStoreToDatabaseFailed tests storing data to a MongoDB database with a failed connection attempt.
func TestStoreToDatabaseFailed(t *testing.T) {
	// Attempting to connect to MongoDB database with a bad URI
	badURI := "mongodb://root:root@db:27017/"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(badURI))
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Creating a test request with a successful operation
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{HasSucceeded: true}}},
		},
	}

	// Attempting to store data to the database (which should fail)
	err = serveur.StoreToDatabase(req)
	if err == nil {
		t.Log("Expected error storing to database, got nil")
	} else {
		t.Failed()
	}
}

// TestStoreToDatabaseEmptyRequest tests storing data to a MongoDB database with an empty request.
func TestStoreToDatabaseEmptyRequest(t *testing.T) {
	req := &flaco_grpc.Request{} // Creating an empty request

	// Attempting to store data to the database
	err := serveur.StoreToDatabase(req)
	if err == nil {
		t.Log("Expected error for empty request, got nil")
	} else {
		t.Failed()
	}
}

// TestStoreToDatabaseNilRequest tests storing data to a MongoDB database with a nil request.
func TestStoreToDatabaseNilRequest(t *testing.T) {
	// Attempting to store data to the database with a nil request
	err := serveur.StoreToDatabase(nil)
	if err == nil {
		t.Log("Expected error for nil request, got nil")
	} else {
		t.Failed()
	}
}

// TestStoreToDatabaseNoDevices tests storing data to a MongoDB database with a request containing no devices.
func TestStoreToDatabaseNoDevices(t *testing.T) {
	req := &flaco_grpc.Request{} // Creating a request with no devices

	// Attempting to store data to the database
	err := serveur.StoreToDatabase(req)
	if err == nil {
		t.Log("Expected error for request with no devices, got nil")
	} else {
		t.Failed()
	}
}

// TestStoreToDatabaseNoOperations tests storing data to a MongoDB database with a request containing devices but no operations.
func TestStoreToDatabaseNoOperations(t *testing.T) {
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{{}}, // Creating a request with devices but no operations
	}

	// Attempting to store data to the database
	err := serveur.StoreToDatabase(req)
	if err == nil {
		t.Log("Expected error for request with no operations, got nil")
	} else {
		t.Failed()
	}
}

// TestStoreToDatabaseDuplicateDeviceName tests storing data to a MongoDB database with duplicate device names in the request.
func TestStoreToDatabaseDuplicateDeviceName(t *testing.T) {
	// Creating a request with duplicate device names and successful operations
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{HasSucceeded: true}}},
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{HasSucceeded: true}}},
		},
	}

	// Attempting to store data to the database
	err := serveur.StoreToDatabase(req)
	if err == nil {
		t.Log("Expected error for request with duplicate device names, got nil")
	} else {
		t.Failed()
	}
}

// TestStoreToDatabaseInvalidOperationType tests storing data to a MongoDB database with invalid operation types in the request.
func TestStoreToDatabaseInvalidOperationType(t *testing.T) {
	// Creating a request with invalid operation type
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{Operation: []*flaco_grpc.Operation{{Type: "invalid_type"}}},
		},
	}

	// Attempting to store data to the database
	err := serveur.StoreToDatabase(req)
	if err == nil {
		t.Log("Expected error for request with invalid operation type, got nil")
	} else {
		t.Failed()
	}
}
