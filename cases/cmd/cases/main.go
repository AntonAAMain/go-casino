package main

import (
	"casino/cases/internal/cases/handlers"
	"casino/cases/internal/cases/repository"
	"casino/cases/internal/cases/service"
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
				Name: "cases_http_requests_total",
				Help: "Total HTTP requests to cases service",
			},
			[]string{"method", "path", "status"},
		)

		boxesOpened = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "cases_boxes_opened_total",
				Help: "Total number of boxes opened",
			},
		)
	)

	// Регистрируем метрики
	prometheus.MustRegister(httpRequests, boxesOpened)

	// Middleware для подсчета HTTP запросов
	prometheusMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			status := fmt.Sprintf("%d", c.Response().Status)
			httpRequests.WithLabelValues(c.Request().Method, c.Path(), status).Inc()
			return err
		}
	}

	casesRepo := repository.NewCasesRepository(db)
	casesService := service.NewBoxService(casesRepo)
	casesHandler := handlers.NewBoxCarHandler(casesService)

	e := echo.New()

	// Отдельный роутер для метрик БЕЗ аутентификации
	metrics := echo.New()
	metrics.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Запускаем метрики сервер на отдельном порту
	go func() {
		if err := metrics.Start(":8088"); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	e.Use(auth.Protected(false))
	e.Use(prometheusMiddleware)
	e.POST("/cases", func(c echo.Context) error {
		boxesOpened.Inc()
		return casesHandler.OpenBox(c)
	})

	fmt.Println("Server started on http://localhost:8084")
	if err := e.Start(":8084"); err != nil {
		log.Fatal(err)
	}
}
