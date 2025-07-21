package main

import (
	"context"
	"log"

	"github.com/krishanu7/grpc/api/generated/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":50051",grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	//Create User
	createResp, err := client.CreateUser(context.Background(), &user.CreateUserRequest{
		Name:  "Krishanu1",
		Email: "krishanu1@gmail.com",
	})

	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Printf("Created user : %+v", createResp)
	// Get User
	getResp, err := client.GetUser(context.Background(), &user.GetUserRequest{
		Id: createResp.Id,
	})

	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	log.Printf("Get user: %+v", getResp)
}
