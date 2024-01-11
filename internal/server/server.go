package server

import (
	"fmt"
	"os"

	"github.com/xprnio/go-serverless/internal/database"
	"github.com/xprnio/go-serverless/internal/docker"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e      *echo.Echo
	db     *database.Database
	docker *docker.Client
}

type Options struct {
	DatabaseName string
	Port         uint16
}

func New(opts Options) (*Server, error) {
	e := echo.New()
	db, err := database.New(opts.DatabaseName)
	docker, err := docker.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err != nil {
		return nil, err
	}

	s := &Server{e, db, docker}

	e.GET("/v1/functions", s.handleFunctionsGetAll)
	e.POST("/v1/functions", s.handleFunctionsCreate)
	e.POST("/v1/functions/:id/pull", s.handleFunctionsPull)

	e.GET("/v1/routes", s.handleRoutesGetAll)
	e.GET("/v1/routes/:id", s.handleRoutesGet)
	e.POST("/v1/routes", s.handleRoutesCreate)

	e.POST("/r/*", s.handleRunner)

	addr := fmt.Sprintf(":%d", opts.Port)
	if err := e.Start(addr); err != nil {
		return nil, err
	}

	return s, nil
}
