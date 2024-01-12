package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/xprnio/go-serverless/internal/utils"
)

func CreateContainer(
	client *client.Client,
	image string,
	env map[string]string,
	mnt []mount.Mount,
) (container.CreateResponse, error) {
	ctx := context.Background()
	conf := &container.Config{
		Env:          utils.EnvMapToList(env),
		Image:        image,
		AttachStdin:  true,
		AttachStderr: true,
	}
	hostConf := &container.HostConfig{Mounts: mnt}
	return client.ContainerCreate(ctx, conf, hostConf, nil, nil, "")
}

func StartContainer(
	client *client.Client,
	containerId string,
) error {
	ctx := context.Background()
	err := client.ContainerStart(
		ctx, containerId,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}

	return nil
}

func AttachContainer(
	client *client.Client,
	containerId string,
) (types.HijackedResponse, error) {
	ctx := context.Background()
	return client.ContainerAttach(
		ctx, containerId,
		types.ContainerAttachOptions{
			Stderr: true,
			Stream: true,
		},
	)
}
