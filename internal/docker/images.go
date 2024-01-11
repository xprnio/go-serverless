package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	pkg "github.com/docker/docker/pkg/archive"
)

type BuildImageOptions struct {
	Tags    []string
	Context string
}

type DockerStream struct {
	Content string `json:"stream"`
}

func (c *Client) splitImageName(name string) (string, string) {
	parts := strings.SplitN(name, ":", 2)

	if len(parts) == 1 {
		return parts[0], "latest"
	}

	return parts[0], parts[1]
}

func (c *Client) ImageExists(name string) (bool, error) {
	ctx := context.Background()
	opts := types.ImageListOptions{}
	images, err := c.client.ImageList(ctx, opts)
	if err != nil {
		return false, nil
	}

	imageName, imageTag := c.splitImageName(name)
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == fmt.Sprintf("%s:%s", imageName, imageTag) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *Client) PullImageAuthenticated(
	name string,
	auth registry.AuthConfig,
) error {
	enc, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	r, err := c.client.ImagePull(
		context.Background(), name,
		types.ImagePullOptions{
			RegistryAuth: base64.URLEncoding.EncodeToString(enc),
		},
	)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(r)
	defer r.Close()
	for s.Scan() {
		log.Println(s.Text())
	}
	return nil
}

func (c *Client) PullImage(name string) error {
	r, err := c.client.ImagePull(
		context.Background(), name,
		types.ImagePullOptions{
			RegistryAuth: "",
		},
	)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(r)
	defer r.Close()
	for s.Scan() {
		log.Println(s.Text())
	}
	return nil
}

func (c *Client) BuildImage(opts BuildImageOptions) (io.Reader, error) {
	tar, err := pkg.TarWithOptions(
		opts.Context,
		&pkg.TarOptions{},
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.ImageBuild(
		context.Background(), tar,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       opts.Tags,
		},
	)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var stream DockerStream
		json.Unmarshal(scanner.Bytes(), &stream)
		buf.WriteString(stream.Content)
	}

	return bufio.NewReader(buf), nil
}
