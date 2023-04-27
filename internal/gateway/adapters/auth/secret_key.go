package auth

import (
	"context"

	pb "github.com/cardio-analyst/backend/pkg/api/proto/auth"
)

func (c *Client) GenerateSecretKey(ctx context.Context, userLogin, userEmail string) (string, error) {
	request := &pb.GenerateSecretKeyRequest{
		UserLogin: userLogin,
		UserEmail: userEmail,
	}

	response, err := c.client.GenerateSecretKey(ctx, request)
	if err != nil {
		return "", err
	}

	return response.GetSecretKey(), nil
}
