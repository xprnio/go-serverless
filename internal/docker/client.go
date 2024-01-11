package docker

import (
	"github.com/docker/docker/client"
)

type Client struct {
	client *client.Client
}

func New() (*Client, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}
