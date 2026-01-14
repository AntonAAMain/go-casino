package handlers

import (
	"casino/admin/internal/admin/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BoxHandler struct {
	service *service.BoxService
}

func NewBoxHandler(service *service.BoxService) *BoxHandler {
	return &BoxHandler{service: service}
}

func (a *BoxHandler) CreateBox(c echo.Context) error {
	req := struct {
		Name  string  `json:"name"`
		Price float32 `json:"price"`
		Mode  string  `json:"mode"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	box, err := a.service.CreateBox(req.Name, req.Mode, req.Price)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    box,
	})
}
