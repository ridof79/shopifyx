package domain

type ConditionEnum string

const (
	New    ConditionEnum = "new"
	Second ConditionEnum = "second"
)

type Product struct {
	Name           string        `json:"name"`
	Price          int           `json:"price"`
	ImageURL       string        `json:"imageUrl"`
	Stock          int           `json:"stock"`
	Condition      ConditionEnum `json:"condition"`
	Tags           []string      `json:"tags"`
	IsPurchaseable bool          `json:"isPurchaseable"`
	PurchaseCount  int           `json:"purchaseCount"`
}

type ProductResponse struct {
	Id             string        `json:"id"`
	Name           string        `json:"name"`
	Price          int           `json:"price"`
	ImageURL       string        `json:"imageUrl"`
	Stock          int           `json:"stock"`
	Condition      ConditionEnum `json:"condition"`
	Tags           []string      `json:"tags"`
	IsPurchaseable bool          `json:"isPurchaseable"`
	PurchaseCount  int           `json:"purchaseCount"`
}

type StockUpdate struct {
	Stock int `json:"stock"`
}
