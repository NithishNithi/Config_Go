// package grpc

// import (
// 	pb "github.com/NithishNithi/GoShop/proto"
// 	"google.golang.org/grpc"
// )

// func NewGoShopClient() (pb.GoShopServiceClient, *grpc.ClientConn, error) {
// 	// Set up a connection to the gRPC server
// 	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Create a gRPC client
// 	client := pb.NewGoShopServiceClient(conn)

// 	return client, conn, nil
// }

package grpcclient

import (
	"log"
	"sync"

	pb "Config/proto"
	"google.golang.org/grpc"
)

var once sync.Once

type GrpcClient pb.ApplicationServiceClient

var (
	instance GrpcClient
)

func GetGrpcClientInstance() (GrpcClient,*grpc.ClientConn) {
	var conn *grpc.ClientConn
	once.Do(func() { // <-- atomic, does not allow repeating
		conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		//defer conn.Close()

		instance = pb.NewApplicationServiceClient(conn)
	})

	return instance,conn
}

