package models

type Product struct {
	 UserId int `json:"user_id"`
	 CategoryId int `json:"category_id"`
	 Description string `json:"description"`
}

type ProductDetailed struct {
	Product
	ProductId int `json:"product_id"`
}

type Category struct {
	CategoryId int `json:"category_id"`
	Category string `json:"category"`
	Subcategory string `json:"sub_category"`
}

