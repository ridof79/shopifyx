package domain

import _ "database/sql"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SellerResponse struct {
	Name             string         `json:"name"`
	ProductSoldTotal int            `json:"productSoldTotal"`
	BankAccounts     []BankAccounts `json:"bankAccount"`
}
