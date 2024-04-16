package serveur

import (
	"context"
	"flaco/grpc_and_go/flaco_grpc"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	flaco_grpc.UnimplementedDayServiceServer
}

func StoreToDatabase(DeviceInfo *flaco_grpc.Request) error {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		return err
	}
	for _, DeviceInfo := range DeviceInfo.Device {
		coll := client.Database("flaco").Collection(DeviceInfo.DeviceName)
		_, err = coll.InsertOne(context.Background(), DeviceInfo, nil)
		if err != nil {
			return err
		}
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

func main() {

}
