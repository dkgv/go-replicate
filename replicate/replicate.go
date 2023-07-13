package replicate

import (
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
	token  string
	common service

	Predictions *PredictionsService
	Models      *ModelsService
}

func NewClient(token string) *Client {
	c := &Client{
		token: token,
	}
	c.common.client = c

	return &Client{
		token:       token,
		Predictions: (*PredictionsService)(&c.common),
		Models:      (*ModelsService)(&c.common),
	}
}

func (c *Client) baseRequest(method string, endpoint string) (*http.Request, error) {
	url := fmt.Sprintf(baseURL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	authorizationHeader := fmt.Sprintf("Token %s", c.token)
	req.Header.Set("Authorization", authorizationHeader)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
