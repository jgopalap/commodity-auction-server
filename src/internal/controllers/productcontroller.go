package controllers

import (
	"encoding/json"
	"internal/models"
	"internal/utils"
	"net/http"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {
	product := &models.Product{}
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal product " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	isNewRecord := database.Table("auction_db.public.Products").NewRecord(product)
	if isNewRecord {
		err = database.Table("auction_db.public.Products").Create(&product).Error
		if err != nil {
			resp := utils.Message(false, "Unable to add product " + err.Error())
			utils.Respond(w, resp)
			return
		}
		resp := utils.Message(true, "success")
		utils.Respond(w, resp)
	} else {
		resp := utils.Message(false, "Unable to add product, already exits.")
		utils.Respond(w, resp)
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

func RemoveProduct(w http.ResponseWriter, r *http.Request) {
	product := &models.ProductDetailed{}
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal product " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Products").Where("product_id = ?", product.ProductId).
		Delete(models.ProductDetailed{}).Error
	if err != nil {
		resp := utils.Message(false, "Unable to remove product " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "success")
	utils.Respond(w, resp)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	productId := r.URL.Query().Get("product_id")
	product := &models.ProductDetailed{}
	err := database.Table("auction_db.public.Products").Where("product_id = ?", productId).
		First(&product).Error
	if err != nil {
		resp := utils.Message(false, "Unable to retrieve product " + err.Error() + "." )
		utils.Respond(w, resp)
		return
	}
	resp := utils.Message(true, "success")
	resp["product"] = product
	utils.Respond(w, resp)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	product := &models.ProductDetailed{}
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		resp := utils.Message(false,  "Unable to unmarshal product " + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	err = database.Table("auction_db.public.Products").Where("product_id = ?", product.ProductId).
		Update("producer_id", product.UserId).
		Update("category_id", product.CategoryId).
		Update("description", product.Description).Error
	if err != nil {
		resp := utils.Message(false, "Unable to update product" + err.Error() + ".")
		utils.Respond(w, resp)
		return
	}

	resp := utils.Message(true, "Success")
	utils.Respond(w, resp)
}

