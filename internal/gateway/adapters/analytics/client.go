package analytics

import pb "github.com/cardio-analyst/backend/api/proto/analytics"

type Client struct {
	client pb.AnalyticsServiceClient
}

func NewClient(client pb.AnalyticsServiceClient) *Client {
	return &Client{
		client: client,
	}
}
