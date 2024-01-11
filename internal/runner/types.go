package runner

import (
	"fmt"
	"os"

	"github.com/xprnio/go-serverless/internal/context"
	"github.com/xprnio/go-serverless/internal/database"
	"github.com/xprnio/go-serverless/internal/docker"
)

type Runner struct {
	Function *database.Function
	client   *docker.Client
}

func NewRunner(
	client *docker.Client,
	function *database.Function,
) *Runner {
	return &Runner{
		Function: function,
		client:   client,
	}
}

func (r *Runner) RunRequest(ctx *context.Context) ([]byte, error) {
	resp, err := r.client.RunImage(docker.RunImageOptions{
		Image:       r.Function.Image,
		RequestId:   ctx.RequestId,
		RequestBody: ctx.RequestBody,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return resp, nil
}
