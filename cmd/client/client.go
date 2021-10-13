package main

import (
	"context"
	"fmt"
	"log"

	"github.com/harisson-freitas/user-registration/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50052", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	AddUser(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:             1,
		FirstName:      "Leon",
		LastName:       "Kennedy",
		Email:          "leon.skennedy@rpd.com",
		DocumentNumber: "3434343434-900",
		CellPhone:      "555 334343434",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}
