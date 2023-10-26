package main

import (
	pb "GoConfig/proto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type server struct {
	pb.UnimplementedMyServiceServer
	mongoClient *mongo.Client
}

func NewServer() (*server, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	return &server{
		mongoClient: client,
	}, nil
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *server) InsertData(ctx context.Context, req *pb.Request) (*emptypb.Empty, error) {
	// Unmarshal the value from the request
	var value interface{}
	err := json.Unmarshal([]byte(req.Value), &value)
	if err != nil {
		return nil, err
	}

	collection := s.mongoClient.Database("kishore").Collection("nithish")

	// Create the document to be inserted
	document := bson.M{
		"key":   req.Key,
		"value": value,
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
	collection := s.mongoClient.Database("kishore").Collection("nithish")
	filter := bson.M{"key": bson.M{"$regex": "^" + req.Key + `(\.|$)`}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	result := make(map[string]interface{})
	for cur.Next(context.TODO()) {
		var entry bson.M
		if err := cur.Decode(&entry); err != nil {
			return nil, err
		}

		key := entry["key"].(string)
		parts := strings.Split(key, ".")
		current := result

		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = entry["value"]
			} else {
				if _, ok := current[part].(map[string]interface{}); !ok {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}
	}
	convertedData := convertMapToStructValue(result)

	response := &pb.GetDataResponse{
		Data: convertedData,
	}
	fmt.Println("result:",response)

	return response, nil
}

func convertMapToStructValue(inputMap map[string]interface{}) map[string]*structpb.Value {
	convertedMap := make(map[string]*structpb.Value)
	for key, value := range inputMap {
		sv, err := structpb.NewValue(value)
		if err != nil {
			// Handle the error if needed
		}
		convertedMap[key] = sv
	}
	return convertedMap
}


func (s *server) AddConfig(ctx context.Context, req *pb.AddConfigRequest) (*emptypb.Empty, error) {
	var value interface{}
	err := json.Unmarshal(req.Value.Value, &value)
	if err != nil {
		return nil, err
	}

	collection := s.mongoClient.Database("kishore").Collection("nithish")
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
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
		return
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterMyServiceServer(s, server)
	if err2 := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err2)
	}
	defer server.mongoClient.Disconnect(context.Background())
}
