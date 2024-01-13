package runner

import (
	"bufio"
	"log"
	"path"

	"github.com/docker/docker/client"
	"github.com/xprnio/go-serverless/internal/database"
	"github.com/xprnio/go-serverless/internal/docker"
)

type Runner struct {
	Function *database.Function
	Docker   *client.Client
}

func (r *Runner) Run(ctx *Context) ([]byte, error) {
	env := map[string]string{
		"REQUEST_ID":   ctx.RequestId,
		"CONTEXT_PATH": path.Join(ctx.ContextPath, ctx.RequestId),
	}
	mounts := ctx.Mounts()
	container, err := docker.CreateContainer(r.Docker, ctx.Image, env, mounts)
	if err != nil {
		return nil, err
	}

	if err := docker.StartContainer(r.Docker, container.ID); err != nil {
		return nil, err
	}

	reader, err := docker.ReadLogs(r.Docker, container.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Implement proper function logging
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	return ctx.GetContext("response.json")
}

func New(
	client *client.Client,
	function *database.Function,
) (*Runner, error) {
	r := &Runner{
		Function: function,
		Docker:   client,
	}
	return r, nil
}
