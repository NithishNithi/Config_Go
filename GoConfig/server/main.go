package main

import (
	pb "GoConfig/proto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

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
	collection  *mongo.Collection
}

func NewServer() (*server, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database("kishore").Collection("nithish")
	return &server{
		mongoClient: client,
		collection:  collection,
	}, nil
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *server) AddApplication(ctx context.Context, request *pb.AddApplicationRequest) (*pb.AddApplicationResponse, error) {
	time := time.Now()
	filter := bson.M{"applicationName": request.ApplicationName, "createdAt": time.Format("2006-01-02 15:04:05")}
	response, err := s.collection.InsertOne(ctx, filter)
	if err != nil {
		log.Print("err:", err)
		return nil, err
	}
	insertedID, ok := response.InsertedID.(primitive.ObjectID)
	if !ok {
		// Handle the case where the insertedID is not an ObjectId, if needed.
		fmt.Println("InsertedID is not an ObjectId")
		return nil, err
	}

	insertedIDStr := insertedID.Hex()

	result := &pb.AddApplicationResponse{
		ApplicationId: insertedIDStr,
	}
	return result, nil
}

func (s *server) WatchConfig(req *pb.WatchDataRequest, stream pb.MyService_WatchConfigServer) error {
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"fullDocument.key", bson.D{{"$regex", "^" + req.Key + `(\.|$)`}}}}}},
	}
	changeStream, err := s.collection.Watch(context.Background(), pipeline)
	if err != nil {
		return err
	}
	defer changeStream.Close(context.Background())

	for changeStream.Next(context.Background()) {
		var change bson.M
		if err := changeStream.Decode(&change); err != nil {
			return err
		}

		// Convert the change to a GetDataResponse and send it to the client
		response := &pb.GetDataResponse{
			Data: convertMapToStructValue(change["fullDocument"].(bson.M)),
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) InsertConfig(ctx context.Context, req *pb.Request) (*emptypb.Empty, error) {
	// Find the application by its ID
	var application bson.M
	err := s.collection.FindOne(context.Background(), bson.M{"applicationName": req.ApplicationName}).Decode(&application)
	if err != nil {
		return nil, err
	}
	var value interface{}
	// Create the document to be inserted
	err = json.Unmarshal([]byte(req.Value), &value)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	time := time.Now()
	// Insert the document into the application
	update := bson.M{
		"$set": bson.M{
			"key":       req.Key,
			"value":     value,
			"updatedAt": time.Format("2006-01-02 15:04:05"),
		},
	}
	_, err = s.collection.UpdateOne(context.Background(), bson.M{"_id": application["_id"]}, update)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) GetConfig(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {

	filter := bson.M{"key": bson.M{"$regex": "^" + req.Key + `(\.|$)`}}

	cur, err := s.collection.Find(context.TODO(), filter)
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
	formattedOutput := formatData(convertedData, "")
	fmt.Println(formattedOutput)

	response := &pb.GetDataResponse{
		Data: convertedData,
	}

	return response, nil
}

func convertMapToStructValue(inputMap map[string]interface{}) map[string]*structpb.Value {
	convertedMap := make(map[string]*structpb.Value)
	for key, value := range inputMap {
		switch v := value.(type) {
		case map[string]interface{}:
			convertedMap[key] = structpb.NewStructValue(&structpb.Struct{
				Fields: convertMapToStructValue(v),
			})
		default:
			sv, err := structpb.NewValue(value)
			if err != nil {
				// Handle the error if needed
			}
			convertedMap[key] = sv
		}
	}
	return convertedMap
}
func formatData(inputMap map[string]*structpb.Value, indent string) string {
	var output string

	for key, value := range inputMap {
		output += indent + "key: " + key + "\n"
		if value.Kind != nil {
			switch value.Kind.(type) {
			case *structpb.Value_NullValue:
				output += indent + "value: null\n"
			case *structpb.Value_BoolValue:
				output += indent + "value: " + fmt.Sprintf("%t", value.GetBoolValue()) + "\n"
			case *structpb.Value_NumberValue:
				output += indent + "value: " + fmt.Sprintf("%f", value.GetNumberValue()) + "\n"
			case *structpb.Value_StringValue:
				output += indent + "value: " + value.GetStringValue() + "\n"
			case *structpb.Value_StructValue:
				nestedMap := value.GetStructValue().GetFields()
				output += formatData(nestedMap, indent+"  ")
			default:
				output += indent + "value: unsupported type\n"
			}
		}
	}

	return output
}

func (s *server) AddConfig(ctx context.Context, req *pb.AddConfigRequest) (*emptypb.Empty, error) {
	var value interface{}
	err := json.Unmarshal(req.Value.Value, &value)
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$push": bson.M{"Config": bson.M{"Key": req.Key, "Value": value}},
	}

	res, err := s.collection.UpdateOne(context.Background(), filter, update)
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
