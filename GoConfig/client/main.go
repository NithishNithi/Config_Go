package main

import (
	pb "GoConfig/proto"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var configClient pb.MyServiceClient

func main() {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	configClient = pb.NewMyServiceClient(conn)
	
	// InsertConfig()
	// AddConfig()
	// GetConfig()
	Insert("Myapp", "Myapp.client.integer", 123)
	// Init("Myapp")

}

func Init(ApplicationName string) {
	result, err := configClient.AddApplication(context.Background(), &pb.AddApplicationRequest{
		ApplicationName: ApplicationName,
	})
	if err != nil {
		log.Print("err:", err)
		return
	}
	fmt.Println("ApplicationId:", result.ApplicationId)
}

// func InsertConfig() {
// 	Key := "you.client.float"
// 	Value := 123.123

// 	valueJSON, err := json.Marshal(Value)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal value to JSON: %v", err)
// 	}
// 	req := &pb.Request{
// 		Key:   Key,
// 		Value: string(valueJSON),
// 	}
// 	_, err = configClient.InsertConfig(context.Background(), req)
// 	if err != nil {
// 		log.Fatalf("Failed to insert data: %v", err)
// 	}
// }

func Insert[T any](ApplicationName, key string, val T) (err error) {
	valueJSON, err := json.Marshal(val)
	if err != nil {
		log.Fatalf("Failed to marshal value to JSON: %v", err)
	}

	req := &pb.Request{
		ApplicationName: ApplicationName,
		Key:             key,
		Value:           string(valueJSON),
	}
	_, err = configClient.InsertConfig(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}
	return err
}

func GetConfig() {
	key := "you.collection.string"
	req := &pb.GetDataRequest{
		Key: key,
	}
	response, err := configClient.GetConfig(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
	fmt.Println(response)
}

// func AddConfig() {
// 	id := "653a0683e9807ff9306101f1"
// 	key := ""
// 	value := 123
// 	valueJSON, err := json.Marshal(value)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal value to JSON: %v", err)
// 	}
// 	req := &pb.AddConfigRequest{
// 		Id:  id,
// 		Key: key,
// 		Value: &anypb.Any{
// 			Value: valueJSON,
// 		},
// 	}
// 	_, err = client.AddConfig(context.Background(), req)
// 	if err != nil {
// 		log.Fatalf("Failed to insert data: %v", err)
// 	}
// }

func WatchConfig() {
	req := &pb.WatchDataRequest{
		Key: "you.client",
	}
	stream, err := configClient.WatchConfig(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to watch data: %v", err)
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive data: %v", err)
		}
		fmt.Println(response.Data)
	}
}
