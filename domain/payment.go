package domain

type Payment struct {
	Id                   string `json:"id"`
	BankAccountId        string `json:"bankAccountId"`
	PaymentProofImageURL string `json:"paymentProofImageUrl"`
	Quantity             int    `json:"quantity"`
}
