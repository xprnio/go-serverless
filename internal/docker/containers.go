package docker

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

type RunImageOptions struct {
	Image       string
	RequestId   string
	RequestBody []byte
}

func (c *Client) RunImage(opts RunImageOptions) ([]byte, error) {
	ctx := context.Background()
	runContext, err := CreateContext(opts)
	if err != nil {
		return nil, err
	}

	createResp, err := c.client.ContainerCreate(
		ctx,
		&container.Config{
			Image:        opts.Image,
			AttachStdin:  true,
			AttachStderr: true,
		},
		&container.HostConfig{
			Mounts: runContext.Mounts,
		},
		nil, nil, "",
	)
	if err != nil {
		return nil, err
	}

	err = c.client.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	attachResp, err := c.client.ContainerAttach(ctx, createResp.ID, types.ContainerAttachOptions{
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(attachResp.Reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return readResponse(runContext)
}

func readResponse(ctx RunContext) ([]byte, error) {
	content, err := os.ReadFile(path.Join(ctx.Dir, "response.json"))
	if err != nil {
		return nil, err
	}

	return content, nil
}

type RunContext struct {
	Dir    string
	Mounts []mount.Mount
}

func CreateContext(opts RunImageOptions) (RunContext, error) {
	var ctx RunContext
	dir, err := CreateContextDir(opts.RequestId)
	if err != nil {
		return ctx, err
	}

	_, err = addContext(dir, "request.json", string(opts.RequestBody))
	if err != nil {
		return ctx, err
	}

	_, err = addContext(dir, "response.json", "{}")
	if err != nil {
		return ctx, err
	}

	ctx.Dir = dir
	ctx.Mounts = []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: dir,
			Target: "/context",
		},
	}

	return ctx, nil
}

func addContext(dir, filename string, json string) (mount.Mount, error) {
	err := os.WriteFile(path.Join(dir, filename), []byte(json), 0600)
	if err != nil {
		return mount.Mount{}, err
	}

	return mount.Mount{
		Type:   mount.TypeBind,
		Source: path.Join(dir, filename),
		Target: path.Join("/context", filename),
	}, nil
}

func CreateContextDir(requestId string) (string, error) {
	tmp := path.Join(os.TempDir(), requestId)
	info, err := os.Stat(tmp)
	if err == nil {
		if !info.IsDir() {
			return "",
				fmt.Errorf("%s exists, but is not a directory", tmp)
		}

		fmt.Printf("%s exists\n", tmp)
		return tmp, nil
	}

	if os.IsNotExist(err) {
		fmt.Printf("creating directory: %s\n", tmp)
		if err := os.Mkdir(tmp, 0755); err != nil {
			return tmp, err
		}

		return tmp, nil
	}

	return tmp, err
}
