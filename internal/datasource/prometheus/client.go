package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type Client struct {
	api v1.API
}

func NewClient(url string) (*Client, error) {
	client, err := api.NewClient(api.Config{Address: url})
	if err != nil {
		return nil, err
	}
	return &Client{api: v1.NewAPI(client)}, nil
}

func (c *Client) Query(ctx context.Context, query string, timeRange string) (interface{}, error) {
	duration, _ := time.ParseDuration(timeRange)
	endTime := time.Now()
	startTime := endTime.Add(-duration)

	result, _, err := c.api.QueryRange(ctx, query, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("query prometheus: %w", err)
	}

	return result, nil
}
