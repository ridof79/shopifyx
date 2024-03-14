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
	buyerId := claims.Id

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
				"message": err.Error(),
			})
	}
	defer tx.Rollback()

	// bank account id user == product id user
	// get user id dari product id
	validBankId, sellerId, _ := repository.ProductAndBankAccountValid(tx, payment.BankAccountId, productId)
	if !validBankId {
		tx.Rollback()
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"message": "payment details invalid",
			},
		)
	}

	if err := repository.CreatePayment(tx, &payment, productId, buyerId, sellerId); err != nil {
		tx.Rollback()
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err.Error(),
			})
	}

	productStock, err := repository.GetProductStockTx(tx, productId)
	if err != nil {
		tx.Rollback()
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err.Error(),
			})
	}

	if productStock < payment.Quantity {
		tx.Rollback()
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"message": "Insufficient stock",
			},
		)
	}

	if err := repository.UpdateProductStockTx(tx, productId, productStock-payment.Quantity); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err.Error(),
			})
	}

	return c.JSON(
		http.StatusCreated,
		map[string]interface{}{
			"message": "Payment added successfully!",
		})
}
