package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/harisson-freitas/user-registration/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	fmt.Printf("User add: %v", req)
	return &pb.User{
		Id:             1,
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		Email:          req.GetEmail(),
		DocumentNumber: req.GetDocumentNumber(),
		CellPhone:      req.GetCellPhone(),
	}, nil
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	startRegistration(stream)
	insertRegistration(req, stream)
	finishRegistration(req, stream)
	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		users = append(users, createUser(req))
		fmt.Println("Adding", req.GetFirstName())
	}
}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving stream from the client: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   req,
		})

		if err != nil {
			log.Fatalf("Error sending stream to the client: %v", err)
		}
	}
}

func startRegistration(stream pb.UserService_AddUserVerboseServer) error {
	init, err := retrieveEmptyUser("Init")
	if err != nil {
		log.Fatalf("Error starting user creation status update: %v", err)
	}

	stream.Send(init)
	time.Sleep(time.Second * 2)

	return nil
}

func insertRegistration(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	insert, err := retrieveEmptyUser("Inserting")
	if err != nil {
		log.Fatalf("Error inserting user creation status update: %v", err)
	}

	stream.Send(insert)
	time.Sleep(time.Second * 2)

	update, err := retrieveFilledUser(req, "User has been inserted")
	if err != nil {
		log.Fatalf("Error updating user creation status: %v", err)
	}

	stream.Send(update)
	time.Sleep(time.Second * 2)

	return nil
}

func finishRegistration(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	complete, err := retrieveFilledUser(req, "Completed")
	if err != nil {
		log.Fatalf("Error updating user creation status: %v", err)
	}

	stream.Send(complete)
	time.Sleep(time.Second * 2)

	return nil
}

func retrieveEmptyUser(status string) (*pb.UserResultStream, error) {
	return &pb.UserResultStream{
		Status: status,
		User:   &pb.User{},
	}, nil
}

func retrieveFilledUser(req *pb.User, status string) (*pb.UserResultStream, error) {
	return &pb.UserResultStream{
		Status: status,
		User:   createUser(req),
	}, nil
}

func createUser(req *pb.User) *pb.User {
	return &pb.User{
		Id:             req.GetId(),
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		Email:          req.GetEmail(),
		DocumentNumber: req.GetDocumentNumber(),
		CellPhone:      req.GetCellPhone(),
	}
}
