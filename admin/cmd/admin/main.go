// package main

// import (
// 	"casino/admin/internal/admin/handlers"
// 	"casino/admin/internal/admin/repository"
// 	"casino/admin/internal/admin/service"
// 	"casino/pkg/middleware/auth"
// 	"fmt"
// 	"log"

// 	"github.com/labstack/echo/v4"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func main() {
// 	dsn := "postgres://postgres:0611@localhost:5432/go-casino?sslmode=disable"

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("failed to connect database:", err)
// 	}

// 	fmt.Println("Database connected!")

// 	carsRepo := repository.NewCarsRepository(db)
// 	carsService := service.NewCarsService(carsRepo)
// 	carsHandler := handlers.NewCarHandler(carsService)

// 	boxRepo := repository.NewBoxRepository(db)
// 	boxService := service.NewBoxService(boxRepo)
// 	boxHandler := handlers.NewBoxHandler(boxService)

// 	boxCarRepo := repository.NewBoxCarsRepository(db)
// 	boxCarService := service.NewBoxCarService(boxCarRepo)
// 	boxCarHandler := handlers.NewBoxCarHandler(boxCarService)

// 	e := echo.New()

// 	e.Use(auth.Protected(true))
// 	e.POST("/admin/car/create", carsHandler.CreateCar)
// 	e.GET("/admin/car/all", carsHandler.GetAllCars)

// 	e.POST("/admin/box/create", boxHandler.CreateBox)
// 	e.POST("/admin/box/car/create", boxCarHandler.CreateBoxCar)

// 	fmt.Println("Server started on http://localhost:8082")
// 	if err := e.Start(":8082"); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"casino/admin/internal/admin/handlers"
	"casino/admin/internal/admin/repository"
	"casino/admin/internal/admin/service"
	"casino/pkg/middleware/auth"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ------------------ Prometheus метрики ------------------
var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "admin_http_requests_total",
			Help: "Total HTTP requests to admin service",
		},
		[]string{"method", "path", "status"},
	)

	carsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "admin_cars_created_total",
			Help: "Total number of cars created",
		},
	)

	boxesCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "admin_boxes_created_total",
			Help: "Total number of boxes created",
		},
	)
)

func initMetrics() {
	prometheus.MustRegister(httpRequests, carsCreated, boxesCreated)
}

func prometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		status := fmt.Sprintf("%d", c.Response().Status)
		httpRequests.WithLabelValues(c.Request().Method, c.Path(), status).Inc()
		return err
	}
}

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

	carsRepo := repository.NewCarsRepository(db)
	carsService := service.NewCarsService(carsRepo)
	carsHandler := handlers.NewCarHandler(carsService)

	boxRepo := repository.NewBoxRepository(db)
	boxService := service.NewBoxService(boxRepo)
	boxHandler := handlers.NewBoxHandler(boxService)

	boxCarRepo := repository.NewBoxCarsRepository(db)
	boxCarService := service.NewBoxCarService(boxCarRepo)
	boxCarHandler := handlers.NewBoxCarHandler(boxCarService)

	initMetrics()

	e := echo.New()
	e.Use(auth.Protected(true))
	e.Use(prometheusMiddleware)

	// Отдельный роутер для метрик БЕЗ аутентификации
	metrics := echo.New()
	metrics.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Запускаем метрики сервер на отдельном порту
	go func() {
		if err := metrics.Start(":8083"); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	e.POST("/admin/car/create", func(c echo.Context) error {
		carsCreated.Inc()
		return carsHandler.CreateCar(c)
	})
	e.GET("/admin/car/all", carsHandler.GetAllCars)

	e.POST("/admin/box/create", func(c echo.Context) error {
		boxesCreated.Inc()
		return boxHandler.CreateBox(c)
	})
	e.POST("/admin/box/car/create", boxCarHandler.CreateBoxCar)

	fmt.Println("Server started on http://localhost:8082")
	if err := e.Start(":8082"); err != nil {
		log.Fatal(err)
	}
}
