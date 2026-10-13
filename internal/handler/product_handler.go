package handler

import (
	"encoding/json"
	"firstgoproject/internal/model"
	"firstgoproject/internal/service"
	"net/http"
)

type ProductHandler struct{
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler{
	return &ProductHandler{
		service:service,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r*http.Request){
	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil{
		
	}
	// newproduct,err := h.service.Create(r.Context(),product)
	// if err != nil{

	// }
	// return 
}