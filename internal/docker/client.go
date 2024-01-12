package docker

import (
	"github.com/docker/docker/client"
)

type Client struct {
	client *client.Client
}
