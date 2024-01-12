package runner

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/xprnio/go-serverless/internal/docker"
	"github.com/xprnio/go-serverless/internal/utils"
)

type ContextManager struct {
	Docker      *client.Client
	ContextName string
	ContextPath string
}

const (
	defaultContextName = "serverless_context"
	defaultContextPath = "/tmp/context"
)

func NewManager(
	client *client.Client,
	contextName, contextPath string,
) (*ContextManager, error) {
	if contextName == "" {
		contextName = defaultContextName
	}

	if contextPath == "" {
		contextPath = defaultContextPath
	}

	m := &ContextManager{
		Docker:      client,
		ContextName: contextName,
		ContextPath: contextPath,
	}

	if err := m.prepareVolume(); err != nil {
		return nil, err
	}

	if err := m.prepareContextPath(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *ContextManager) NewContext(
	image string,
	body []byte,
) (*Context, error) {
	requestId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return m.NewContextWithId(requestId, image, body)
}

func (m *ContextManager) NewContextWithId(
	id uuid.UUID,
	image string,
	body []byte,
) (*Context, error) {
	return NewContext(
		image,
		m.ContextName,
		m.ContextPath,
		id.String(),
		body,
	)
}

func (m *ContextManager) prepareContextPath() error {
	if utils.IsDocker() {
		// TODO
		return nil
	}

	stat, err := os.Stat(m.ContextPath)
	if err == nil {
		if !stat.IsDir() {
			return fmt.Errorf(
				"%s exists, but is not a directory",
				m.ContextPath,
			)
		}

		log.Println("context path exists:", m.ContextPath)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	log.Println("creating context path:", m.ContextPath)
	return os.MkdirAll(m.ContextPath, 0755)
}

func (m *ContextManager) prepareVolume() error {
	if utils.IsDocker() {
		// TODO Verify that ContextPath is mounted to ContextName
		return nil
	}

	exists, err := docker.VolumeExists(m.Docker, m.ContextName)
	if err != nil {
		return err
	}

	if !exists {
		_, err := docker.CreateVolume(m.Docker, m.ContextName)
		if err != nil {
			return err
		}

		log.Println("volume created:", m.ContextName)
	} else {
		log.Println("volume exists:", m.ContextName)
	}

	return nil
}
