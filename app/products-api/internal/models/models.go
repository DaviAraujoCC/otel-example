package models

type Product struct {
	ProductID   int                `json:"productId"`
	Title       string             `json:"title"`
	Sku         string             `json:"sku"`
	Barcodes    []string           `json:"barcodes"`
	Description string             `json:"description"`
	Attributes  []ProductAttribute `json:"attributes"`
	Price       string             `json:"price"`
	Created     string             `json:"created"`
	LastUpdated string             `json:"lastUpdated"`
}

type ProductAttribute struct {
	ProductID int    `json:"productId"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type ProductGetResponse struct {
	TotalCount int       `json:"totalCount"`
	Items      []Product `json:"items"`
}