package model

type Car struct {
	Id           string `json:"id"`
	Year         string `json:"year"`
	Make         string `json:"make"`
	Model        string `json:"model"`
	Trim         string `json:"trim"`
	Body         string `json:"body"`
	Transmission string `json:"transmission"`
	State        string `json:"state"`
	Condition    string `json:"condition"`
	Odometer     string `json:"odometer"`
	Color        string `json:"color"`
	Interior     string `json:"interior"`
	Mmr          string `json:"mmr"`
	Seller       string `json:"seller"`
	SellingPrice string `json:"selling-price"`
	SaleDate     string `json:"sale-date"`
}
