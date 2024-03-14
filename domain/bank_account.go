package domain

type BankAccount struct {
	Id                string `json:"id"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
	UserId            string `json:"userId"`
}

type BankAccounts struct {
	Id                string `json:"id"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

type BankAccountsResponse struct {
	Message string         `json:"message"`
	Data    []BankAccounts `json:"data"`
}
