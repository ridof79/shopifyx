package util

import (
	"shopifyx/domain"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(c echo.Context, code int, message string) error {
	return c.JSON(code,
		map[string]string{
			"error": message,
		},
	)
}

func ResponseHandler(c echo.Context, code int, message string) error {
	return c.JSON(code,
		map[string]string{
			"message": message,
		},
	)
}

func UserSuccesResponseHandler(c echo.Context, code int, message, username, name, token string) error {
	return c.JSON(code, map[string]interface{}{
		"message": message,
		"data": map[string]interface{}{
			"username":    username,
			"name":        name,
			"accessToken": token,
		},
	})
}

func SerachProductPaginationResponseHandler(c echo.Context, code int, products []domain.ProductResponse, limit, offset, total int) error {
	return c.JSON(code, map[string]interface{}{
		"message": "ok",
		"data":    products,
		"meta": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}

func UploadImageResponseHandler(c echo.Context, code int, url string) error {
	return c.JSON(code, map[string]interface{}{
		"imageUrl": url,
	},
	)
}

func GetProductResponseHandler(c echo.Context, code int, product domain.ProductResponse, seller domain.SellerResponse) error {
	return c.JSON(code, map[string]interface{}{
		"message": "ok",
		"data": map[string]interface{}{
			"product": product,
			"seller":  seller,
		},
	})
}

func GetBankAccountsResposesHandler(c echo.Context, code int, bankAccounts []domain.BankAccounts) error {
	return c.JSON(code, map[string]interface{}{
		"message": "success",
		"data":    bankAccounts,
	})
}
