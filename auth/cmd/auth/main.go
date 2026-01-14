package main

import (
	"casino/auth/internal/auth/handlers"
	"casino/auth/internal/auth/repository"
	"casino/auth/internal/auth/service"
	"fmt"
	"log"
	"os"

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
				Name: "auth_http_requests_total",
				Help: "Total HTTP requests to auth service",
			},
			[]string{"method", "path", "status"},
		)

		usersRegistered = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "auth_users_registered_total",
				Help: "Total number of users registered",
			},
		)

		usersLoggedIn = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "auth_users_logged_in_total",
				Help: "Total number of user logins",
			},
		)
	)

	// Регистрируем метрики
	prometheus.MustRegister(httpRequests, usersRegistered, usersLoggedIn)

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
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	e := echo.New()

	// Отдельный роутер для метрик БЕЗ аутентификации
	metrics := echo.New()
	metrics.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Запускаем метрики сервер на отдельном порту
	go func() {
		if err := metrics.Start(":8087"); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// Добавляем middleware для подсчета запросов
	e.Use(prometheusMiddleware)

	e.POST("/auth/register", func(c echo.Context) error {
		usersRegistered.Inc()
		return authHandler.Register(c)
	})
	e.POST("/auth/login", func(c echo.Context) error {
		usersLoggedIn.Inc()
		return authHandler.Login(c)
	})

	fmt.Println("Server started on http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
