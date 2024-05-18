package serveur

import (
	"context"
	"flaco/grpc_and_go/flaco_grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

// TestGetDeviceStatNoFailedOperations tests the calculation of device statistics when all operations have succeeded.
func TestGetDeviceStatNoFailedOperations(t *testing.T) {
	device := &flaco_grpc.Device{
		DeviceName: "device1",
		Operation:  []*flaco_grpc.Operation{{HasSucceeded: true}, {HasSucceeded: true}},
	}

	stat := GetDeviceStat(device)
	expectedTotal := int64(2)
	expectedFailed := int64(0)

	if stat.NbTotalOp != expectedTotal {
		t.Errorf("Expected total operations: %d, got: %d", expectedTotal, stat.NbTotalOp)
	}

	if stat.NbOpFailed != expectedFailed {
		t.Errorf("Expected failed operations: %d, got: %d", expectedFailed, stat.NbOpFailed)
	}
}

// TestGetDeviceStatWithFailedOperations tests the calculation of device statistics when some operations have failed.
func TestGetDeviceStatWithFailedOperations(t *testing.T) {
	device := &flaco_grpc.Device{
		DeviceName: "device1",
		Operation:  []*flaco_grpc.Operation{{HasSucceeded: true}, {HasSucceeded: false}},
	}

	stat := GetDeviceStat(device)
	expectedTotal := int64(2)
	expectedFailed := int64(1)

	if stat.NbTotalOp != expectedTotal {
		t.Errorf("Expected total operations: %d, got: %d", expectedTotal, stat.NbTotalOp)
	}

	if stat.NbOpFailed != expectedFailed {
		t.Errorf("Expected failed operations: %d, got: %d", expectedFailed, stat.NbOpFailed)
	}
}

// TestStoreToDatabaseSuccessful tests storing data to a MongoDB database with a successful connection.
func TestStoreToDatabaseSuccessful(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}
	defer client.Disconnect(context.Background())

	err = client.Database("flaco").Drop(context.Background())
	if err != nil {
		t.Errorf("Failed to drop database: %v", err)
	}

	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{Type: "type1", HasSucceeded: true}}},
		},
	}

	err = StoreToDatabase(req)
	if err != nil {
		t.Errorf("Error storing to database: %v", err)
	}

	collection := client.Database("flaco").Collection("test_device")
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		t.Errorf("Failed to count documents in collection: %v", err)
	}

	expectedCount := int64(1)
	if count != expectedCount {
		t.Errorf("Expected %d documents in collection, got %d", expectedCount, count)
	}

	err = client.Database("flaco").Drop(context.Background())
	if err != nil {
		t.Errorf("Failed to drop database: %v", err)
	}
}

// TestStoreToDatabaseFailedConnection tests storing data to a MongoDB database with a failed connection attempt.
func TestStoreToDatabaseFailedConnection(t *testing.T) {
	badURI := "mongodb://root:root@db:27017/"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(badURI))
	if err == nil {
		client.Disconnect(context.Background())
		t.Log("Expected connection to fail, but it succeeded")
		t.Failed()
	}

	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{HasSucceeded: true}}},
		},
	}

	err = StoreToDatabase(req)
	if err == nil {
		t.Failed()
	}
}

// TestStoreToDatabaseEmptyRequest tests storing data to a MongoDB database with an empty request.
func TestStoreToDatabaseEmptyRequest(t *testing.T) {
	req := &flaco_grpc.Request{}

	err := StoreToDatabase(req)
	if err != nil {
		t.Errorf("Expected no error for empty request, got: %v", err)
	}
}

// TestStoreToDatabaseNilRequest tests storing data to a MongoDB database with a nil request.
func TestStoreToDatabaseNilRequest(t *testing.T) {
	err := StoreToDatabase(nil)
	if err == nil {
		t.Log("Expected error for nil request, got nil")
		t.Failed()
	}
}

// TestStoreToDatabaseNoDevices tests storing data to a MongoDB database with a request containing no devices.
func TestStoreToDatabaseNoDevices(t *testing.T) {
	req := &flaco_grpc.Request{}

	err := StoreToDatabase(req)
	if err != nil {
		t.Errorf("Expected no error for request with no devices, got: %v", err)
	}
}

// TestStoreToDatabaseNoOperations tests storing data to a MongoDB database with a request containing devices but no operations.
func TestStoreToDatabaseNoOperations(t *testing.T) {
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{{DeviceName: "test_device"}},
	}

	err := StoreToDatabase(req)
	if err != nil {
		t.Errorf("Expected no error for request with no operations, got: %v", err)
	}
}

// TestStoreToDatabaseDuplicateDeviceName tests storing data to a MongoDB database with duplicate device names in the request.
func TestStoreToDatabaseDuplicateDeviceName(t *testing.T) {
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{HasSucceeded: true}}},
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{HasSucceeded: true}}},
		},
	}

	err := StoreToDatabase(req)
	if err != nil {
		t.Errorf("Expected no error for request with duplicate device names, got: %v", err)
	}
}

// TestStoreToDatabaseInvalidOperationType tests storing data to a MongoDB database with invalid operation types in the request.
func TestStoreToDatabaseInvalidOperationType(t *testing.T) {
	req := &flaco_grpc.Request{
		Device: []*flaco_grpc.Device{
			{DeviceName: "test_device", Operation: []*flaco_grpc.Operation{{Type: "", HasSucceeded: true}}},
		},
	}

	err := StoreToDatabase(req)
	if err != nil {
		t.Errorf("Expected no error for request with invalid operation type, got: %v", err)
	}
}
