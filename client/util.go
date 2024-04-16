package client

import (
	"flaco/grpc_and_go/flaco_grpc"
)

func ConvertDeviceDataToGRPCDevice(in DeviceData) *flaco_grpc.Device {
	var operations []*flaco_grpc.Operation
	for _, op := range in.Operations {
		operation := &flaco_grpc.Operation{
			Type:         op.Type,
			HasSucceeded: op.HasSucceeded,
		}
		operations = append(operations, operation)
	}

	return &flaco_grpc.Device{
		DeviceName: in.DeviceName,
		Operation:  operations,
	}
}
