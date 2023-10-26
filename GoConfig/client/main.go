package main

import (
	pb "GoConfig/proto"
	"context"
	"encoding/json"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var client pb.MyServiceClient

func main() {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewMyServiceClient(conn)

	// InsertData()
	// AddConfig()
	GetData()
}

func InsertData() {

	Key := "my.server.bool"
	Value := false

	valueJSON, err := json.Marshal(Value)
	if err != nil {
		log.Fatalf("Failed to marshal value to JSON: %v", err)
	}
	req := &pb.Request{
		Key: Key,
		Value: string(valueJSON),
	}
	_, err = client.InsertData(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}
}

func GetData() {

	key := "my.client.integer"
	req := &pb.GetDataRequest{
		Key: key,
	}
	_, err := client.GetData(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
}

func AddConfig() {
	id := "653a0683e9807ff9306101f1"
	key := ""
	value := 123
	valueJSON, err := json.Marshal(value)
	if err != nil {
		log.Fatalf("Failed to marshal value to JSON: %v", err)
	}
	req := &pb.AddConfigRequest{
		Id:  id,
		Key: key,
		Value: &anypb.Any{
			Value: valueJSON,
		},
	}
	_, err = client.AddConfig(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}
}
