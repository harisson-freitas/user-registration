package main

import (
	"log"
	"net"

	"github.com/harisson-freitas/user-registration/pb"
	"github.com/harisson-freitas/user-registration/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Could not server: %v", err)
	}
}
