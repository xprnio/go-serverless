package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

func splitImageName(name string) (string, string) {
	parts := strings.SplitN(name, ":", 2)

	if len(parts) == 1 {
		return parts[0], "latest"
	}

	return parts[0], parts[1]
}

func ImageExists(client *client.Client, name string) (bool, error) {
	ctx := context.Background()
	opts := types.ImageListOptions{}
	images, err := client.ImageList(ctx, opts)
	if err != nil {
		return false, nil
	}

	imageName, imageTag := splitImageName(name)
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == fmt.Sprintf("%s:%s", imageName, imageTag) {
				return true, nil
			}
		}
	}

	return false, nil
}

func PullImageAuthenticated(
	client *client.Client,
	name string,
	auth registry.AuthConfig,
) error {
	enc, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	opts := types.ImagePullOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(enc),
	}
	return pullImage(client, name, opts)
}

func PullImage(client *client.Client, name string) error {
	return pullImage(client, name, types.ImagePullOptions{})
}

func pullImage(
	client *client.Client,
	name string,
	opts types.ImagePullOptions,
) error {
	r, err := client.ImagePull(context.Background(), name, opts)
	if err != nil {
		return err
	}
	defer r.Close()

	s := bufio.NewScanner(r)
	for s.Scan() {
		log.Println(s.Text())
	}

	return nil
}
