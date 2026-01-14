package main

import (
	"casino/inventory/internal/inventory/handlers"
	"casino/inventory/internal/inventory/repository"
	"casino/inventory/internal/inventory/service"
	"os"

	"casino/pkg/middleware/auth"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgres://postgres:0611@localhost:5432/go-casino?sslmode=disable"

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	fmt.Println("Database connected!")

	// Prometheus метрики
	var (
		httpRequests = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "inventory_http_requests_total",
				Help: "Total HTTP requests to inventory service",
			},
			[]string{"method", "path", "status"},
		)

		itemsSold = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "inventory_items_sold_total",
				Help: "Total number of items sold",
			},
		)
	)

	// Регистрируем метрики
	prometheus.MustRegister(httpRequests, itemsSold)

	// Middleware для подсчета HTTP запросов
	prometheusMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			status := fmt.Sprintf("%d", c.Response().Status)
			httpRequests.WithLabelValues(c.Request().Method, c.Path(), status).Inc()
			return err
		}
	}

	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handlers.NewInventoryCarHandler(inventoryService)

	e := echo.New()

	// Отдельный роутер для метрик БЕЗ аутентификации
	metrics := echo.New()
	metrics.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Запускаем метрики сервер на отдельном порту
	go func() {
		if err := metrics.Start(":8089"); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	e.Use(auth.Protected(false))
	e.Use(prometheusMiddleware)
	e.POST("/profile/inventory/sell/all", func(c echo.Context) error {
		itemsSold.Inc()
		return inventoryHandler.SellUserInventory(c)
	})
	e.POST("/profile/inventory/sell/car", func(c echo.Context) error {
		itemsSold.Inc()
		return inventoryHandler.SellUserCar(c)
	})

	fmt.Println("Server started on http://localhost:8085")
	if err := e.Start(":8085"); err != nil {
		log.Fatal(err)
	}
}
