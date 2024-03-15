package main

import (
	"fmt"
	"net/http"
	"shopifyx/auth"
	"shopifyx/config"
	"shopifyx/delivery"
	"time"

	"log"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "shopifyx_request",
		Help:    "Histogram of the /shopifyx request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inisialisasi koneksi database
	config.InitDB()
	defer config.CloseDB()

	// Inisialisasi Echo framework
	e := echo.New()

	// Custom logger
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Middleware

	e.Use(echojwt.WithConfig(auth.ConfigJWT()))

	//auth
	//e.POST("/v1/user/register", delivery.RegisterUserHandler)
	NewRoute(e, "/v1/user/register", "POST", delivery.RegisterUserHandler)

	//e.POST("/v1/user/login", delivery.LoginUserHandler)
	NewRoute(e, "/v1/user/login", "POST", delivery.LoginUserHandler)

	//product
	e.POST("/v1/product", delivery.CreateProductHandler)
	e.PATCH("/v1/product/:productId", delivery.UpdateProductHandler)
	e.DELETE("/v1/product/:productId", delivery.DeleteProductHandler)

	//stock managemenet
	e.POST("/v1/product/:productId/stock", delivery.UpdateProductStockHandler)

	//bank account
	e.POST("/v1/bank/account", delivery.AddBankAccountHandler)
	e.GET("/v1/bank/account", delivery.GetBankAccountsHandler)
	e.PATCH("/v1/bank/account/:bankAccountId", delivery.UpdateBankAccountHandler)

	//payment
	e.POST("/v1/product/:productId/buy", delivery.CreatePaymentHandler)

	//seach
	//e.GET("/v1/product", delivery.SearchProductHandler)
	NewRoute(e, "/v1/product", "GET", delivery.SearchProductHandler)

	//get product
	//e.GET("/v1/product/:productId", delivery.GetProductHandler)
	NewRoute(e, "/v1/product/:productId", "GET", delivery.GetProductHandler)

	//image upload
	e.POST("/v1/image", delivery.UploadImageHandler)

	e.Logger.Fatal(e.Start(":8000"))
}

func NewRoute(c *echo.Echo, path string, method string, handler echo.HandlerFunc) {
	c.Add(method, path, wrapHandlerWithMetrics(path, method, handler))
}

func wrapHandlerWithMetrics(path string, method string, handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()

		// Execute the actual handler and catch any errors
		err := handler(c)

		// Regardless of whether an error occurred, record the metrics
		duration := time.Since(startTime).Seconds()
		statusCode := fmt.Sprintf("%d", c.Response().Status)

		if err != nil {
			if c.Response().Status == http.StatusOK { // Default status code
				c.Response().Status = http.StatusInternalServerError // Assume internal server error if not set
			}
			c.String(http.StatusInternalServerError, err.Error()) // Ensure the response reflects the error
		}

		requestHistogram.WithLabelValues(path, method, statusCode).Observe(duration)
		return err
	}
}
