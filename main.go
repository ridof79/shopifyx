package main

import (
	"os"
	"shopifyx/auth"
	"shopifyx/config"
	"shopifyx/delivery"

	"strings"

	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

	// Custom logger
	// Configure Echo to use LogrusAdapter

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/v1/user/")
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	}))

	//auth
	e.POST("/v1/user/register", delivery.RegisterUserHandler)
	e.POST("/v1/user/login", delivery.LoginUserHandler)

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
	e.GET("/v1/product", delivery.SearchProductHandler)

	//get product
	e.GET("/v1/product/:productId", delivery.GetProductHandler)

	//image upload
	e.POST("/v1/image", delivery.UploadImageHandler)

	e.Logger.Fatal(e.Start(":8000"))
}
