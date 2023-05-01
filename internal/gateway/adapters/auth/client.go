package auth

import (
	pb "github.com/cardio-analyst/backend/api/proto/auth"
)

type Client struct {
	client pb.AuthServiceClient
}

func NewClient(client pb.AuthServiceClient) *Client {
	return &Client{
		client: client,
	}
}
