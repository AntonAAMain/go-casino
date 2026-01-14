package handlers

import (
	"casino/admin/internal/admin/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BoxCarHandler struct {
	service *service.BoxCarService
}

func NewBoxCarHandler(service *service.BoxCarService) *BoxCarHandler {
	return &BoxCarHandler{service: service}
}

func (a *BoxCarHandler) CreateBoxCar(c echo.Context) error {
	req := struct {
		CarIds []int `json:"car_id"`
		BoxId  int   `json:"box_id"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	boxCar, err := a.service.CreateBoxCar(req.CarIds, req.BoxId)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    boxCar,
	})
}
