package services

import (
	"context"
	"fmt"

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
