package server

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xprnio/go-serverless/internal/runner"
	"github.com/xprnio/go-serverless/internal/server/responses"
)

func (server *Server) handleRunner(c echo.Context) error {
	req := c.Request()
	path := strings.TrimPrefix(req.URL.Path, "/r")

	route, err := server.Database.GetRouteByPath(path)
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

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	ctx, err := server.Manager.NewContext(
		route.Function.Image,
		buf.Bytes(),
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	runner, err := runner.New(server.Docker, route.Function)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	resp, err := runner.Run(ctx)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	// TODO: Handle function-level errors and status codes
	return c.JSONBlob(http.StatusOK, resp)
}
