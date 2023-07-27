package replicate

import (
	"context"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://api.replicate.com/v1/%s"
)

type service struct {
	client *Client
}

type Client struct {
	token   string
	common  service
	baseURL string

	Predictions *PredictionsService
	Models      *ModelsService
}

type ClientOption func(c *Client)

func NewClient(token string, options ...ClientOption) *Client {
	c := &Client{
		token:   token,
		baseURL: baseURL,
	}

	for _, option := range options {
		option(c)
	}

	c.common.client = c

	return &Client{
		token:       token,
		Predictions: (*PredictionsService)(&c.common),
		Models:      (*ModelsService)(&c.common),
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func (c *Client) baseRequest(ctx context.Context, method string, endpoint string) (*http.Request, error) {
	url := fmt.Sprintf(c.baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	authorizationHeader := fmt.Sprintf("Token %s", c.token)
	req.Header.Set("Authorization", authorizationHeader)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
