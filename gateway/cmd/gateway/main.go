package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"casino/pkg/middleware/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type GatewayHandler struct {
	httpClient *http.Client
}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *GatewayHandler) proxyRequest(c echo.Context, serviceURL string) error {
	// Получаем оригинальный запрос
	req := c.Request()

	// Преобразуем путь для целевого сервиса
	// Например: "/api/v1/auth/register" + "/auth" = "/auth/register"
	targetPath := strings.Replace(req.URL.Path, "/api/v1", "", 1)
	targetURL := serviceURL + targetPath
	if req.URL.RawQuery != "" {
		targetURL += "?" + req.URL.RawQuery
	}

	// Копируем тело запроса
	var body io.Reader
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read request body"})
		}
		body = bytes.NewReader(bodyBytes)
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Восстанавливаем тело для Echo
	}

	// Создаем прокси запрос
	proxyReq, err := http.NewRequest(req.Method, targetURL, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create proxy request"})
	}

	// Копируем заголовки
	for key, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Выполняем запрос
	resp, err := g.httpClient.Do(proxyReq)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "service unavailable"})
	}
	defer resp.Body.Close()

	// Читаем ответ
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read response"})
	}

	// Устанавливаем статус код
	c.Response().WriteHeader(resp.StatusCode)

	// Копируем заголовки ответа (кроме тех, которые управляет Echo)
	for key, values := range resp.Header {
		if strings.ToLower(key) != "content-length" {
			for _, value := range values {
				c.Response().Header().Add(key, value)
			}
		}
	}

	// Возвращаем тело ответа
	return c.String(resp.StatusCode, string(respBody))
}

func main() {
	// Получаем URL сервисов из переменных окружения или используем дефолтные
	authURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8080")
	paymentURL := getEnv("PAYMENT_SERVICE_URL", "http://localhost:8081")
	adminURL := getEnv("ADMIN_SERVICE_URL", "http://localhost:8082")
	casesURL := getEnv("CASES_SERVICE_URL", "http://localhost:8084")
	inventoryURL := getEnv("INVENTORY_SERVICE_URL", "http://localhost:8085")

	// Prometheus метрики
	var (
		httpRequests = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_total",
				Help: "Total HTTP requests to gateway service",
			},
			[]string{"method", "path", "status", "service"},
		)
	)

	// Регистрируем метрики
	prometheus.MustRegister(httpRequests)

	// Middleware для подсчета HTTP запросов
	prometheusMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			status := fmt.Sprintf("%d", c.Response().Status)
			service := extractServiceFromPath(c.Path())
			httpRequests.WithLabelValues(c.Request().Method, c.Path(), status, service).Inc()
			return err
		}
	}

	gatewayHandler := NewGatewayHandler()

	e := echo.New()

	// CORS middleware
	e.Use(middleware.CORS())

	// Логирование
	e.Use(middleware.Logger())

	// Prometheus middleware
	e.Use(prometheusMiddleware)

	// Отдельный роутер для метрик БЕЗ аутентификации
	metrics := echo.New()
	metrics.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Запускаем метрики сервер на отдельном порту
	go func() {
		if err := metrics.Start(":8086"); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// Auth endpoints (без аутентификации)
	e.POST("/api/v1/auth/register", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, authURL)
	})
	e.POST("/api/v1/auth/login", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, authURL)
	})

	// Payment endpoints (с аутентификацией)
	e.POST("/api/v1/payment/transaction", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, paymentURL)
	}, auth.Protected(false))

	// Admin endpoints (только для админов)
	e.POST("/api/v1/admin/cars", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, adminURL)
	}, auth.Protected(true))
	e.GET("/api/v1/admin/cars", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, adminURL)
	}, auth.Protected(true))
	e.POST("/api/v1/admin/boxes", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, adminURL)
	}, auth.Protected(true))
	e.POST("/api/v1/admin/box-cars", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, adminURL)
	}, auth.Protected(true))

	// Cases endpoints (с аутентификацией)
	e.POST("/api/v1/cases/open", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, casesURL)
	}, auth.Protected(false))

	// Inventory endpoints (с аутентификацией)
	e.POST("/api/v1/profile/inventory/sell/all", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, inventoryURL)
	}, auth.Protected(false))
	e.POST("/api/v1/profile/inventory/sell/car", func(c echo.Context) error {
		return gatewayHandler.proxyRequest(c, inventoryURL)
	}, auth.Protected(false))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	fmt.Println("Gateway server started on http://localhost:4000")
	if err := e.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func extractServiceFromPath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/"), "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}
