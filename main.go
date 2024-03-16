package main

import (
	"log"
	"shopifyx/auth"
	"shopifyx/config"
	"shopifyx/delivery"

	middle "shopifyx/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/echoprometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e.Validator = middle.Validator

	// prometheus
	e.Use(echoprometheus.NewMiddleware("myapp"))   // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//	e.User(Prome)

	// Whitelist routes
	e.POST("/v1/user/register", delivery.RegisterUserHandler)
	e.POST("/v1/user/login", delivery.LoginUserHandler)
	e.GET("/v1/product", delivery.SearchProductHandler)
	e.GET("/v1/product/:productId", delivery.GetProductHandler)

	// Protected routes
	v1 := e.Group("/v1")
	v1.Use(echojwt.WithConfig(auth.ConfigJWT()))

	// Product routes
	product := v1.Group("/product")
	product.POST("", delivery.CreateProductHandler)
	product.PATCH("/:productId", delivery.UpdateProductHandler)
	product.DELETE("/:productId", delivery.DeleteProductHandler)
	product.POST("/:productId/stock", delivery.UpdateProductStockHandler)

	// Bank account routes
	bank := v1.Group("/bank")
	bank.POST("/account", delivery.AddBankAccountHandler)
	bank.GET("/account", delivery.GetBankAccountsHandler)
	bank.PATCH("/account/:bankAccountId", delivery.UpdateBankAccountHandler)

	// Payment route
	v1.POST("/product/:productId/buy", delivery.CreatePaymentHandler)

	// Image upload route
	v1.POST("/image", delivery.UploadImageHandler)

	e.Logger.Fatal(e.Start(":8000"))
}
