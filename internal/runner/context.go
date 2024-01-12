package runner

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/docker/docker/api/types/mount"
	"github.com/xprnio/go-serverless/internal/utils"
)

const RunnerContextPath = "/context"

type Context struct {
	Image       string
	ContextName string
	ContextPath string

	RequestId   string
	RequestBody []byte
}

func NewContext(
	image string,
	contextName, contextPath string,
	requestId string,
	requestBody []byte,
) (*Context, error) {
	ctx := &Context{
		Image:       image,
		ContextName: contextName,
		ContextPath: contextPath,
		RequestId:   requestId,
		RequestBody: requestBody,
	}

	return ctx, ctx.prepare()
}

func (c Context) Mounts() []mount.Mount {
	if utils.IsDocker() {
		return []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: c.ContextName,
				Target: c.ContextPath,
			},
		}
	}

	return []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: path.Join(c.ContextPath, c.RequestId),
			Target: path.Join(c.ContextPath, c.RequestId),
		},
	}
}

func (c *Context) AddContext(filename string, data []byte) error {
	p := path.Join(c.ContextPath, c.RequestId, filename)
	log.Println("adding context:", p)
	return os.WriteFile(p, data, 0600)
}

func (c *Context) GetContext(filename string) ([]byte, error) {
	p := path.Join(c.ContextPath, c.RequestId, filename)
	return os.ReadFile(p)
}

func (c *Context) prepare() error {
	if err := c.prepareDirectory(); err != nil {
		return err
	}

	if err := c.AddContext("request.json", c.RequestBody); err != nil {
		return err
	}

	emptyResponse := []byte("{}")
	if err := c.AddContext("response.json", emptyResponse); err != nil {
		return err
	}

	return nil
}

func (c *Context) prepareDirectory() error {
	p := path.Join(c.ContextPath, c.RequestId)
	info, err := os.Stat(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		log.Println("creating context:", p)
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}

		return nil
	}

	if !info.IsDir() {
		return fmt.Errorf("%s exists, but is not a directory", p)
	}

	return nil
}
