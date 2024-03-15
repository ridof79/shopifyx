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
)

type JwtCustomClaims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *domain.User) (string, error) {

	var jwtSecretKey = os.Getenv("JWT_SECRET")
	var jwtExpiredMinutes = os.Getenv("JWT_EXPIRED_MINUTES")

	var tokenExpirationTime, err = strconv.Atoi(jwtExpiredMinutes)
	if err != nil {
		panic(err)
	}

	claims := &JwtCustomClaims{
		Id:   user.Id,
		Name: user.Name,
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
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/v1/user/") || strings.HasPrefix(c.Path(), "/metrics")
		},
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
