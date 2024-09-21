package model

type Item struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Price  string `json:"price"`
	Detail string `json:"detail"`
	Open   string `json:"open"`
}
