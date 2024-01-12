package docker

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func VolumeExists(client *client.Client, name string) (bool, error) {
	filters := filters.NewArgs(
		filters.KeyValuePair{Key: "name", Value: name},
	)
	resp, err := client.VolumeList(
		context.Background(),
		volume.ListOptions{Filters: filters},
	)
	if err != nil {
		return false, err
	}

	for _, v := range resp.Volumes {
		if v.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func CreateVolume(client *client.Client, name string) (volume.Volume, error) {
	return client.VolumeCreate(
		context.Background(),
		volume.CreateOptions{Name: name},
	)
}
