package main

import (
	"shopifyx/auth"
	"shopifyx/config"
	"shopifyx/delivery"

	"log"

	prometheus "shopifyx/middleware"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.WithConfig(auth.ConfigJWT()))

	//auth
	//e.POST("/v1/user/register", delivery.RegisterUserHandler)
	prometheus.NewRoute(e, "/v1/user/register", "POST", delivery.RegisterUserHandler)

	//e.POST("/v1/user/login", delivery.LoginUserHandler)
	prometheus.NewRoute(e, "/v1/user/login", "POST", delivery.LoginUserHandler)

	//product
	//e.POST("/v1/product", delivery.CreateProductHandler)
	prometheus.NewRoute(e, "/v1/product", "POST", delivery.CreateProductHandler)
	//e.PATCH("/v1/product/:productId", delivery.UpdateProductHandler)
	prometheus.NewRoute(e, "/v1/product/:productId", "PATCH", delivery.UpdateProductHandler)
	//e.DELETE("/v1/product/:productId", delivery.DeleteProductHandler)
	prometheus.NewRoute(e, "/v1/product/:productId", "DELETE", delivery.DeleteProductHandler)

	//stock managemenet
	//e.POST("/v1/product/:productId/stock", delivery.UpdateProductStockHandler)
	prometheus.NewRoute(e, "/v1/product/:productId/stock", "POST", delivery.UpdateProductStockHandler)

	//bank account
	//e.POST("/v1/bank/account", delivery.AddBankAccountHandler)
	prometheus.NewRoute(e, "/v1/bank/account", "POST", delivery.AddBankAccountHandler)
	//e.GET("/v1/bank/account", delivery.GetBankAccountsHandler)
	prometheus.NewRoute(e, "/v1/bank/account", "GET", delivery.GetBankAccountsHandler)
	//e.PATCH("/v1/bank/account/:bankAccountId", delivery.UpdateBankAccountHandler)
	prometheus.NewRoute(e, "/v1/bank/account/:bankAccountId", "PATCH", delivery.UpdateBankAccountHandler)

	//payment
	//e.POST("/v1/product/:productId/buy", delivery.CreatePaymentHandler)
	prometheus.NewRoute(e, "/v1/product/:productId/buy", "POST", delivery.CreatePaymentHandler)

	//seach
	//e.GET("/v1/product", delivery.SearchProductHandler)
	prometheus.NewRoute(e, "/v1/product", "GET", delivery.SearchProductHandler)

	//get product
	//e.GET("/v1/product/:productId", delivery.GetProductHandler)
	prometheus.NewRoute(e, "/v1/product/:productId", "GET", delivery.GetProductHandler)

	//image upload
	//e.POST("/v1/image", delivery.UploadImageHandler)
	prometheus.NewRoute(e, "/v1/image", "POST", delivery.UploadImageHandler)

	e.Logger.Fatal(e.Start(":8000"))
}
