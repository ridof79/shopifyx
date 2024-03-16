package domain

import _ "database/sql"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username" validate:"required,min=5,max=15"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type UserLogin struct {
	Id       string `json:"id"`
	Username string `json:"username" validate:"required,min=5,max=15"`
	Name     string `json:"name"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}

type SellerResponse struct {
	Name             string         `json:"name"`
	ProductSoldTotal int            `json:"productSoldTotal"`
	BankAccounts     []BankAccounts `json:"bankAccount"`
}
