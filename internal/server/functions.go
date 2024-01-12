package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xprnio/go-serverless/internal/database"
	"github.com/xprnio/go-serverless/internal/docker"
	"github.com/xprnio/go-serverless/internal/server/responses"
)

func (server *Server) handleFunctionsGetAll(c echo.Context) error {
	result, err := server.Database.GetFunctions()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		responses.NewResourceResponse(result),
	)
}

func (server *Server) handleFunctionsGet(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	result, err := server.Database.GetFunction(id)
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

func (server *Server) handleFunctionsCreate(c echo.Context) error {
	var body database.Function
	req := c.Request()

	d := json.NewDecoder(req.Body)
	if err := d.Decode(&body); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	exists, err := docker.ImageExists(server.Docker, body.Image)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	if !exists {
		log.Println("pulling image", body.Image)
		err := docker.PullImage(server.Docker, body.Image)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				responses.NewErrorResponse(err.Error()),
			)
		}
	} else {
		log.Println("image exists")
	}

	result, err := server.Database.SaveFunction(&body)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	return c.JSON(
		http.StatusCreated,
		responses.NewResourceResponse(result),
	)
}

func (server *Server) handleFunctionsPull(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	result, err := server.Database.GetFunction(id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			responses.NewErrorResponse(err.Error()),
		)
	}

	err = docker.PullImage(server.Docker, result.Image)
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
