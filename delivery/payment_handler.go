package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"shopifyx/auth"
	"shopifyx/db"
	"shopifyx/domain"
	"shopifyx/repository"
	"shopifyx/util"

	"github.com/labstack/echo/v4"
)

const (
	PaymentDetailsInvalid = "payment details invalid or product not purchaseable"
	InsufficientStock     = "Insufficient stock"
	FailedToMakePayment   = "failed to make payment"

	PaymentAddedSuccessfully = "payment added successfully"
)

func CreatePaymentHandler(c echo.Context) error {
	buyerId := auth.GetUserIdFromToken(c)
	log.Println("buyerId", buyerId)

	var payment domain.Payment
	productId := c.Param("productId")

	if err := json.NewDecoder(c.Request().Body).Decode(&payment); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	tx, _ := db.GetDB().Begin()
	defer tx.Rollback()

	validBankId, productStock, sellerId, err := repository.CheckStockProductAndBankAccountValid(tx, payment.BankAccountId, productId)
	if !validBankId || err != nil {
		log.Println("Failed to commit transaction:", err)
		tx.Rollback()
		return util.ErrorHandler(c, http.StatusBadRequest, PaymentDetailsInvalid)
	}

	if productStock < payment.Quantity {
		log.Println("Failed to commit transaction:", err)
		tx.Rollback()
		return util.ErrorHandler(c, http.StatusBadRequest, InsufficientStock)
	}

	if err := repository.CreatePayment(tx, &payment, productId, buyerId, sellerId); err != nil {
		log.Println("Failed to commit transaction:", err)
		tx.Rollback()
	}

	if err := repository.UpdateProductStockTx(tx, productId, productStock-payment.Quantity); err != nil {
		log.Println("Failed to commit transaction:", err)
		tx.Rollback()
	}

	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit transaction:", err)
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToMakePayment)
	}

	return util.ResponseHandler(c, http.StatusCreated, PaymentAddedSuccessfully)
}
