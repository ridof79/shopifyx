package delivery

import (
	"encoding/json"
	"net/http"
	"shopifyx/auth"
	"shopifyx/config"
	"shopifyx/domain"
	"shopifyx/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreatePaymentHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	var payment domain.Payment
	productId := c.Param("productId")

	if err := json.NewDecoder(c.Request().Body).Decode(&payment); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	tx, err := config.GetDB().Begin()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err,
			})
	}
	defer tx.Rollback()

	if err := repository.CreatePayment(tx, &payment, productId, userId); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err,
			})
	}

	productStock, err := repository.GetProductStockTx(tx, productId)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err,
			})
	}

	if productStock < payment.Quantity {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"message": "Insufficient stock",
			},
		)
	}

	if err := repository.UpdateProductStockTx(tx, productId, productStock-payment.Quantity); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err,
			})
	}

	return c.JSON(
		http.StatusCreated,
		map[string]interface{}{
			"message": "Payment added successfully!",
		})
}
