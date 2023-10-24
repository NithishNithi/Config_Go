package main

import (
	"Config/constants"
	"Config/controllers"
	"Config/database"
	pro "Config/proto"
	"Config/services"
	"context"
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initDatabase(client *mongo.Client) {
	ConfigCollection := database.GetCollection(client, "GoConfig", "config")
	controllers.Configservice = services.InitCustomerService(ConfigCollection,context.Background())
}

func main() {
	mongoClient, err := database.ConnectDatabase()
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(context.TODO())
	initDatabase(mongoClient)
	lis, err := net.Listen("tcp", constants.Port)
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pro.RegisterApplicationServiceServer(s, &controllers.RPCServer{})
	fmt.Println("server listening on ", constants.Port)
	err = s.Serve(lis)
	if err != nil {
		fmt.Println("failed to serve", err)
	}

}
