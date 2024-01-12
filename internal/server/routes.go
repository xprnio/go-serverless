package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xprnio/go-serverless/internal/server/responses"
)

func (server *Server) handleRoutesGetAll(c echo.Context) error {
	routes, err := server.Database.GetRoutes()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		responses.NewResourceResponse(routes),
	)
}

func (server *Server) handleRoutesGet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	result, err := server.Database.GetRoute(id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		responses.NewResourceResponse(result),
	)
}

type CreateRouteRequest struct {
	Path       string `json:"path"`
	FunctionId string `json:"function_id"`
}

func (server *Server) handleRoutesCreate(c echo.Context) error {
	var body CreateRouteRequest
	req := c.Request()

	d := json.NewDecoder(req.Body)
	if err := d.Decode(&body); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	functionId, err := uuid.Parse(body.FunctionId)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	result, err := server.Database.CreateRoute(
		body.Path,
		functionId,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}
	return c.JSON(
		http.StatusOK,
		responses.NewResourceResponse(result),
	)
}
