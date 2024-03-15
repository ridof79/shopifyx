package delivery

import (
	"encoding/json"
	"net/http"
	"shopifyx/auth"
	"shopifyx/config"
	"shopifyx/domain"
	"shopifyx/repository"
	"shopifyx/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	PaymentDetailsInvalid = "payment details invalid"
	InsufficientStock     = "Insufficient stock"
	FailedToMakePayment   = "failed to make payment"

	PaymentAddedSuccessfully = "payment added successfully!"
)

func CreatePaymentHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	buyerId := claims.Id

	var payment domain.Payment
	productId := c.Param("productId")

	if err := json.NewDecoder(c.Request().Body).Decode(&payment); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	tx, err := config.GetDB().Begin()
	defer tx.Rollback()

	validBankId, sellerId, _ := repository.ProductAndBankAccountValid(tx, payment.BankAccountId, productId)
	if !validBankId {
		tx.Rollback()
		return util.ErrorHandler(c, http.StatusBadRequest, PaymentDetailsInvalid)
	}

	if err := repository.CreatePayment(tx, &payment, productId, buyerId, sellerId); err != nil {
		tx.Rollback()
	}

	productStock, err := repository.GetProductStockTx(tx, productId)
	if err != nil {
		tx.Rollback()
	}

	if productStock < payment.Quantity {
		tx.Rollback()
		return util.ErrorHandler(c, http.StatusBadRequest, InsufficientStock)
	}

	if err := repository.UpdateProductStockTx(tx, productId, productStock-payment.Quantity); err != nil {
		tx.Rollback()
	}

	if err := tx.Commit(); err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToMakePayment)
	}

	return util.ResponseHandler(c, http.StatusCreated, PaymentAddedSuccessfully)
}
