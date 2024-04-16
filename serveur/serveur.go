package serveur

import (
	"context"
	"flaco/grpc_and_go/flaco_grpc"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type Server struct {
	flaco_grpc.UnimplementedDayServiceServer
}
type CompletedDeviceInfo struct {
	devices    []*flaco_grpc.Device
	nbTotalOp  int64
	nbOpFailed int64
}

func StoreToDatabase(DeviceInfo *flaco_grpc.Request) error {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		return err
	}

	CompletedDevice := CalculVal(DeviceInfo)

	for _, DeviceInfoComp := range CompletedDevice.devices {
		coll := client.Database("flaco").Collection(DeviceInfoComp.DeviceName)
		_, err = coll.InsertOne(context.Background(), DeviceInfoComp, nil)
		if err != nil {
			return err
		}
	}

	coll := client.Database("flaco").Collection("Information_Calculs")
	_, err = coll.InsertOne(context.Background(), CompletedDevice.nbTotalOp, nil)
	if err != nil {
		return err
	}
	_, err = coll.InsertOne(context.Background(), CompletedDevice.nbOpFailed, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s Server) SendDayInfoToServer(ctx context.Context, req *flaco_grpc.Request) (*flaco_grpc.Response, error) {
	DeviceInfo := flaco_grpc.Request{
		Device: req.GetDevice(),
	}

	err := StoreToDatabase(&DeviceInfo)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func CalculVal(DeviceInfo *flaco_grpc.Request) *CompletedDeviceInfo {
	nbTotal := 0
	nbFailed := 0

	for _, DeviceInfo := range DeviceInfo.Device {
		for _, Operation := range DeviceInfo.Operation {
			if !Operation.HasSucceeded {
				nbFailed++
			}
			nbTotal++
		}
	}

	return &CompletedDeviceInfo{
		devices:    DeviceInfo.Device,
		nbOpFailed: int64(nbFailed),
		nbTotalOp:  int64(nbTotal),
	}
}

func ServeurListen() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	flaco_grpc.RegisterDayServiceServer(s, &Server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
