package main

import (
	"fmt"
	"log"
	"shopifyx/auth"
	"shopifyx/config"
	"shopifyx/db"
	"shopifyx/delivery"

	middle "shopifyx/middleware"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inisialisasi koneksi database
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	fmt.Println(config)

	db.InitDB(config)
	defer db.CloseDB()

	err = db.GetDB().Ping()
	if err != nil {
		fmt.Println("Failed to ping database:", err)
		return
	}

	fmt.Println("Database connection OK")

	// Inisialisasi Echo framework
	e := echo.New()
	// Validator
	e.Validator = middle.Validator

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
