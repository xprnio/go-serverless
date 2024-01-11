package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (server *Server) handleRoutesGetAll(c echo.Context) error {
	routes, err := server.db.GetRoutes()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"success": true,
			"data":    routes,
		},
	)
}

func (server *Server) handleRoutesGet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	result, err := server.db.GetRoute(id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"success": false,
			"data":    result,
		},
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
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	functionId, err := uuid.Parse(body.FunctionId)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	result, err := server.db.CreateRoute(
		body.Path,
		functionId,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}
	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"success": true,
			"data":    result,
		},
	)
}
