package handlers

import (
	"casino/admin/internal/admin/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CarsHandler struct {
	service *service.CarsService
}

func NewCarHandler(service *service.CarsService) *CarsHandler {
	return &CarsHandler{service: service}
}

func (a *CarsHandler) CreateCar(c echo.Context) error {
	req := struct {
		Name  string  `json:"name"`
		Price float32 `json:"price"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	car, err := a.service.CreateCar(req.Name, req.Price)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    car,
	})
}

func (a *CarsHandler) GetAllCars(c echo.Context) error {

	cars, err := a.service.GetAllCars()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    cars,
	})
}
