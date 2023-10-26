package main

import (
	pb "GoConfig/proto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedMyServiceServer
}
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}


func (s *server) InsertData(ctx context.Context, req *pb.Request) (*emptypb.Empty, error) {
	// Unmarshal the value from the request
	var value interface{}
	err := json.Unmarshal(req.Value.Value, &value)
	if err != nil {
		return nil, err
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	// Get the collection
	collection := client.Database("kishore").Collection("nithish")

	// Create the document to be inserted
	document := bson.M{
		"name":   req.Name,
		"Config": []interface{}{value}, // Store the value in an array
	}

	// Insert the document into the collection
	result, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return &emptypb.Empty{}, nil
}

func (s *server) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("kishore").Collection("nithish")

	// Create a filter based on _id
	filter := bson.M{"_id": id}

	// If req.Key is not empty, add the Config.Key condition to the filter
	if req.Key != "" {
		filter["Config"] = bson.M{
			"$elemMatch": bson.M{"Key": req.Key},
		}
	}

	// Define a variable to store the result
	var result bson.M

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Result:")
	fmt.Println(result)

	// Convert the result to a gRPC response
	response := &pb.GetDataResponse{}
	response.GDRA = append(response.GDRA, &pb.Application{
		Id:   id.Hex(), // Convert the ObjectID to a string
		Name: result["name"].(string),
		// Add other fields as needed
	})
	return response, nil
}


func (s *server) AddConfig(ctx context.Context, req *pb.AddConfigRequest) (*emptypb.Empty, error) {
	var value interface{}
	err := json.Unmarshal(req.Value.Value, &value)
	if err != nil {
		return nil, err
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer client.Disconnect(context.Background())
	collection := client.Database("kishore").Collection("nithish")
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$push": bson.M{"Config": bson.M{"Key": req.Key, "Value": value}},
	}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return nil, err // Return the error to the client
	}
	fmt.Println(res.UpsertedID)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Listening")
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterMyServiceServer(s, &server{})
	if err2 := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err2)
	}
}
