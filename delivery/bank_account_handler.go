package delivery

import (
	"encoding/json"
	"net/http"
	"shopifyx/auth"
	"shopifyx/domain"
	"shopifyx/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AddBankAccountHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	var bankAccount domain.BankAccount

	if err := json.NewDecoder(c.Request().Body).Decode(&bankAccount); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	err := repository.AddBankAccount(&bankAccount, userId)

	if err != nil {
		if repository.IsConstrainViolations(err) {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "required fields are missing or invalid",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err,
			})
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"message": "account added successfully/",
		})
}

func GetBankAccountsHandler(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	bankAccounts, err := repository.GetBankAccounts(userId)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to get bank accounts",
		})
	}

	var bankAccountsResponse []domain.BankAccounts
	for _, acc := range bankAccounts {
		bankAccountResponse := domain.BankAccounts{
			Id:                acc.Id,
			BankName:          acc.BankName,
			BankAccountName:   acc.BankAccountName,
			BankAccountNumber: acc.BankAccountNumber,
		}
		bankAccountsResponse = append(bankAccountsResponse, bankAccountResponse)
	}

	response := domain.BankAccountsResponse{
		Message: "success",
		Data:    bankAccountsResponse,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateBankAccountHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	bankAccountId := c.Param("bankAccountId")

	var updatedBankAccount domain.BankAccount

	if err := json.NewDecoder(c.Request().Body).Decode(&updatedBankAccount); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	err := repository.UpdateBankAccount(&updatedBankAccount, bankAccountId, userId)

	if err != nil {
		if repository.IdNotFound(err) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "bank account not found",
				},
			)
		}

		if repository.DontHavePermission(err) {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "you don't have permission to delete this account",
				},
			)
		}

		if repository.IsConstrainViolations(err) {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "required fields are missing or invalid",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err,
			})
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"message": "account updated successfully",
		})
}

func DeleteBankAccountHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	bankAccountId := c.Param("bankAccountId")

	err := repository.DeleteBankAccount(bankAccountId, userId)

	if err != nil {
		if repository.IdNotFound(err) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "bank account not found",
				},
			)
		}

		if repository.DontHavePermission(err) {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "you don't have permission to delete this account",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err,
			})
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"message": "account deleted successfully",
		})
}
