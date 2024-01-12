package server

import (
	"os"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/xprnio/go-serverless/internal/database"
	"github.com/xprnio/go-serverless/internal/runner"
)

type Server struct {
	App      *echo.Echo
	Docker   *client.Client
	Database *database.Database
	Manager  *runner.ContextManager
}

type Options struct {
	DatabaseName string
}

func New(opts Options) (*Server, error) {
	app := echo.New()

	db, err := database.New(opts.DatabaseName)
	if err != nil {
		return nil, err
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	manager, err := runner.NewManager(
		docker,
		os.Getenv("CONTEXT_NAME"),
		os.Getenv("CONTEXT_PATH"),
	)
	if err != nil {
		return nil, err
	}

	s := &Server{
		App:      app,
		Docker:   docker,
		Database: db,
		Manager:  manager,
	}

	app.GET("/v1/functions", s.handleFunctionsGetAll)
	app.POST("/v1/functions", s.handleFunctionsCreate)
	app.POST("/v1/functions/:id/pull", s.handleFunctionsPull)

	app.GET("/v1/routes", s.handleRoutesGetAll)
	app.GET("/v1/routes/:id", s.handleRoutesGet)
	app.POST("/v1/routes", s.handleRoutesCreate)

	app.POST("/r/*", s.handleRunner)

	return s, nil
}

func (s *Server) Start(addr string) error {
	return s.App.Start(addr)
}
