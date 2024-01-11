package server

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/xprnio/go-serverless/internal/context"
	"github.com/xprnio/go-serverless/internal/runner"
	"github.com/labstack/echo/v4"
)

func (server *Server) handleRunner(c echo.Context) error {
	req := c.Request()
	path := strings.TrimPrefix(req.URL.Path, "/r")
	route, err := server.db.GetRouteByPath(path)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(req.Body); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	ctx, err := context.New(buf.Bytes())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	runner := runner.NewRunner(server.docker, route.Function)
	resp, err := runner.RunRequest(ctx)
	return c.JSONBlob(200, resp)
}
