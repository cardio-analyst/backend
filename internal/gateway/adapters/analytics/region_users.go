package analytics

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Client) UsersByRegions(ctx context.Context) (map[string]int64, error) {
	request := &emptypb.Empty{}

	response, err := c.client.UsersByRegions(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.GetUsersByRegions(), nil
}
