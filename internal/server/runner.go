package server

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/xprnio/go-serverless/internal/context"
	"github.com/xprnio/go-serverless/internal/runner"
	"github.com/labstack/echo/v4"
	"github.com/xprnio/go-serverless/internal/server/responses"
)

func (server *Server) handleRunner(c echo.Context) error {
	req := c.Request()
	path := strings.TrimPrefix(req.URL.Path, "/r")
	route, err := server.db.GetRouteByPath(path)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			responses.NewErrorResponse(err.Error()),
		)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(req.Body); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	ctx, err := context.New(buf.Bytes())
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	runner := runner.NewRunner(server.docker, route.Function)
	resp, err := runner.RunRequest(ctx)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}
	return c.JSONBlob(http.StatusOK, resp)
}
