package delivery

import (
	"encoding/json"
	"net/http"
	"shopifyx/auth"
	"shopifyx/domain"
	"shopifyx/repository"
	"shopifyx/util"

	"github.com/labstack/echo/v4"
)

const (
	InvalidRequestBody        = "invalid request body"
	RequredFieldsMissing      = "required fields are missing or invalid"
	FailedToAddBankAccount    = "failed to add bank account"
	FailedToUpdateBankAccount = "failed to update bank account"
	FailedToDeleteBankAccount = "failed to delete bank account"

	BankAccountNotFound = "bank account not found"
	DontHavePermission  = "you don't have permission to perform this action"

	AccountAddedSuccessfully   = "account added successfully"
	AccountUpdateSuccessfully  = "account updated successfully"
	AccountDeletedSuccessfully = "account deleted successfully"
)

func AddBankAccountHandler(c echo.Context) error {
	userId := auth.GetUserIdFromToken(c)
	var bankAccount domain.BankAccount

	if err := json.NewDecoder(c.Request().Body).Decode(&bankAccount); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	err := repository.AddBankAccount(&bankAccount, userId)

	if err != nil {
		if repository.IsConstrainViolations(err) {
			return util.ErrorHandler(c, http.StatusBadRequest, RequredFieldsMissing)
		}
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToAddBankAccount)
	}
	return util.ResponseHandler(c, http.StatusOK, AccountAddedSuccessfully)
}

func GetBankAccountsHandler(c echo.Context) error {
	userId := auth.GetUserIdFromToken(c)

	bankAccounts, err := repository.GetBankAccounts(userId)
	if err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToAddBankAccount)
	}

	var bankAccountsResponse []domain.BankAccounts
	for _, acc := range bankAccounts {
		bankAccountResponse := domain.BankAccounts{
			BankAccountId:     acc.Id,
			BankName:          acc.BankName,
			BankAccountName:   acc.BankAccountName,
			BankAccountNumber: acc.BankAccountNumber,
		}
		bankAccountsResponse = append(bankAccountsResponse, bankAccountResponse)
	}

	return util.GetBankAccountsResposesHandler(c, http.StatusOK, bankAccountsResponse)
}

func UpdateBankAccountHandler(c echo.Context) error {
	userId := auth.GetUserIdFromToken(c)

	bankAccountId := c.Param("bankAccountId")
	if len(bankAccountId) != 36 {
		return util.ErrorHandler(c, http.StatusNotFound, BankAccountNotFound)
	}

	var updatedBankAccount domain.BankAccountUpdate

	if err := c.Bind(&updatedBankAccount); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	if err := c.Validate(updatedBankAccount); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, err.Error())
	}

	result, err := repository.UpdateBankAccount(&updatedBankAccount, bankAccountId, userId)

	switch result {
	case 1:
		return util.ResponseHandler(c, http.StatusOK, AccountUpdateSuccessfully)
	case 2:
		return util.ErrorHandler(c, http.StatusNotFound, BankAccountNotFound)
	case 3:
		return util.ErrorHandler(c, http.StatusForbidden, DontHavePermission)
	}

	if err != nil {
		if repository.IsConstrainViolations(err) {
			return util.ErrorHandler(c, http.StatusBadRequest, RequredFieldsMissing)
		}

		if repository.IdNotFound(err) {
			return util.ErrorHandler(c, http.StatusNotFound, BankAccountNotFound)
		}
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToUpdateBankAccount)
	}
	return nil
}

func DeleteBankAccountHandler(c echo.Context) error {
	userId := auth.GetUserIdFromToken(c)

	bankAccountId := c.Param("bankAccountId")

	err := repository.DeleteBankAccount(bankAccountId, userId)

	if err != nil {
		if repository.IdNotFound(err) {
			return util.ErrorHandler(c, http.StatusNotFound, BankAccountNotFound)
		}

		if repository.DontHavePermission(err) {
			return util.ErrorHandler(c, http.StatusForbidden, DontHavePermission)
		}
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToDeleteBankAccount)
	}

	return util.ResponseHandler(c, http.StatusOK, AccountDeletedSuccessfully)
}
