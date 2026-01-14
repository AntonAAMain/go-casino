package main

import (
	"casino/payment/internal/payment/handlers"
	"casino/payment/internal/payment/repository"
	"casino/payment/internal/payment/service"
	"fmt"
	"log"
	"os"

	"casino/pkg/middleware/auth"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// По умолчанию используем localhost для локальной разработки
	dsn := "postgres://postgres:0611@localhost:5432/go-casino?sslmode=disable"

	// Если установлена переменная окружения DATABASE_URL (из docker-compose), используем её
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
				Name: "payment_http_requests_total",
				Help: "Total HTTP requests to payment service",
			},
			[]string{"method", "path", "status"},
		)

		paymentsProcessed = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "payment_payments_processed_total",
				Help: "Total number of payments processed",
			},
		)
	)

	// Регистрируем метрики
	prometheus.MustRegister(httpRequests, paymentsProcessed)

	// Middleware для подсчета HTTP запросов
	prometheusMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			status := fmt.Sprintf("%d", c.Response().Status)
			httpRequests.WithLabelValues(c.Request().Method, c.Path(), status).Inc()
			return err
		}
	}

	userRepo := repository.NewUserRepository(db)
	paymentService := service.NewPaymentService(userRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	e := echo.New()

	// Отдельный роутер для метрик БЕЗ аутентификации
	metrics := echo.New()
	metrics.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Запускаем метрики сервер на отдельном порту
	go func() {
		if err := metrics.Start(":8090"); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	e.Use(auth.Protected(false))
	e.Use(prometheusMiddleware)
	e.POST("/payment/create", func(c echo.Context) error {
		paymentsProcessed.Inc()
		return paymentHandler.CreateTransaction(c)
	})

	fmt.Println(`Server started on http://localhost:8081`)
	if err := e.Start(":8081"); err != nil {
		log.Fatal(err)
	}
}
