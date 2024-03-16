package auth

import (
	"net/http"
	"os"
	"shopifyx/domain"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type JwtCustomClaims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *domain.UserLogin) (string, error) {

	var jwtSecretKey = os.Getenv("JWT_SECRET")
	var jwtExpiredMinutes = os.Getenv("JWT_EXPIRED_MINUTES")

	var tokenExpirationTime, err = strconv.Atoi(jwtExpiredMinutes)
	if err != nil {
		panic(err)
	}

	claims := &JwtCustomClaims{
		Id: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenExpirationTime) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return t, err
	}

	return t, nil
}

func ConfigJWT() echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			if err == echojwt.ErrJWTMissing {
				return echo.NewHTTPError(http.StatusForbidden, "you dont have access")
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized1")
		},
	}
}

func GetUserIdFromHeader(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token, _ := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_SECRET"), nil
	})
	claims, _ := token.Claims.(*JwtCustomClaims)
	return claims.Id
}

func HashPassword(password string) (string, error) {
	cost, _ := strconv.Atoi(os.Getenv("BCRYPT_SALT"))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserIdFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims.Id
}
