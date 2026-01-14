package handlers

import (
	"casino/payment/internal/payment/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (a *PaymentHandler) CreateTransaction(c echo.Context) error {

	req := struct {
		Amount int `json:"amount"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input",
		})
	}

	userID := c.Get("user_id").(int)

	transaction, err := a.service.CreateTransaction(req.Amount, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, transaction)

}
