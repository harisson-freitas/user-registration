package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// fmt.Println("******ADD USER******")
	// AddUser(client)

	// fmt.Println("******ADD USER VERBOSE******")
	// AddUserVerbose(client)

	// fmt.Println("******ADD USERS******")
	// AddUsers(client)

	fmt.Println("******ADD USER STREAM BOTH******")
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := createUser()
	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := createUser()
	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status, "\n", "User: ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := createUsers()

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 2)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := createUsers()
	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.FirstName)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Recebendo user %v com status: %v\n", res.GetUser().GetFirstName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}

func createUser() *pb.User {
	return &pb.User{
		Id:             1,
		FirstName:      "Leon",
		LastName:       "Kennedy",
		Email:          "leon.skennedy@rpd.com",
		DocumentNumber: "3434343434-900",
		CellPhone:      "555 334343434",
	}
}

func createUsers() []*pb.User {
	return []*pb.User{
		{
			Id:             1,
			FirstName:      "Leon",
			LastName:       "Kennedy",
			Email:          "leon.skennedy@rpd.com",
			DocumentNumber: "3434343434-900",
			CellPhone:      "555 334343434",
		},
		{
			Id:             2,
			FirstName:      "Chris",
			LastName:       "Redfield",
			Email:          "chris.redfield@stars.com",
			DocumentNumber: "4545454545-800",
			CellPhone:      "555 0903940909",
		},
		{
			Id:             3,
			FirstName:      "Jill",
			LastName:       "Valentine",
			Email:          "jill.valentine@stars.com",
			DocumentNumber: "545454545-908",
			CellPhone:      "555 7879897978",
		},
		{
			Id:             4,
			FirstName:      "Claire",
			LastName:       "Redfield",
			Email:          "claire.redfield@fox.com",
			DocumentNumber: "45433223232-700",
			CellPhone:      "555 9988787878",
		},
	}
}
