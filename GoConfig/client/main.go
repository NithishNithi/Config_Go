package main

import (
	pb "GoConfig/proto"
	"context"
	"log"

	"google.golang.org/grpc"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMyServiceClient(conn)

	// insert Data ----->

	// name := "Kyc-App2"
	// value := KeyValuePair{
	// 	Key:   "11",
	// 	Value: "54321",
	// }

	// valueJSON, err := json.Marshal(value)
	// if err != nil {
	// 	log.Fatalf("Failed to marshal value to JSON: %v", err)
	// }

	// req := &pb.Request{
	// 	Name: name,
	// 	Value: &anypb.Any{
	// 		Value: valueJSON,
	// 	},
	// }

	// _, err = client.InsertData(context.Background(), req)
	// if err != nil {
	// 	log.Fatalf("Failed to insert data: %v", err)
	// }

	// <-----------AddData

	// ---------> GetData

	id := "65391e475de65bfc69dd549a"
	key := ""

	req := &pb.GetDataRequest{
		Id:  id,
		Key: key,
	}

	_, err = client.GetData(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}

	// <---------- GetData

	// -------->addConfig

	// id := "65391e475de65bfc69dd549a"
	// key := "3"
	// value := true
	// valueJSON, err := json.Marshal(value)
	// if err != nil {
	// 	log.Fatalf("Failed to marshal value to JSON: %v", err)
	// }

	// req := &pb.AddConfigRequest{
	// 	Id:  id,
	// 	Key: key,
	// 	Value: &anypb.Any{
	// 		Value: valueJSON,
	// 	},
	// }
	// fmt.Println("111111")
	// fmt.Println(req)
	// _, err = client.AddConfig(context.Background(), req)
	// if err != nil {
	// 	fmt.Println("2222")
	// 	log.Fatalf("Failed to insert data: %v", err)
	// }

	// <--------------- AddConfig

}
