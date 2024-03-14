package domain

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SellerResponse struct {
	Name          string         `json:"name"`
	ProductTotal  int            `json:"productTotal"`
	PurchaseTotal int            `json:"purchaseTotal"`
	BankAccounts  []BankAccounts `json:"bankAccount"`
}
