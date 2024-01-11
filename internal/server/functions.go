package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xprnio/go-serverless/internal/database"
)

func (server *Server) handleFunctionsGetAll(c echo.Context) error {
	result, err := server.db.GetFunctions()
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
			"data":    result,
		},
	)
}

func (server *Server) handleFunctionsGet(c echo.Context) error {
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

	result, err := server.db.GetFunction(id)
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

func (server *Server) handleFunctionsCreate(c echo.Context) error {
	var body database.Function
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

	exists, err := server.docker.ImageExists(body.Image)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	if !exists {
		log.Println("pulling image", body.Image)
		err := server.docker.PullImage(body.Image)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"success": false,
					"message": err.Error(),
				},
			)
		}
	} else {
		log.Println("image exists")
	}

	result, err := server.db.SaveFunction(&body)
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
		http.StatusCreated,
		map[string]interface{}{
			"success": false,
			"data":    result,
		},
	)
}

func (server *Server) handleFunctionsPull(c echo.Context) error {
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

	result, err := server.db.GetFunction(id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"success": false,
				"message": err.Error(),
			},
		)
	}

	err = server.docker.PullImage(result.Image)
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
