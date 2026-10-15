package handler

import (
	// "encoding/json"
	// "errors"
	// "firstgoproject/internal/model"
	"firstgoproject/internal/service"
	// "io"
	"net/http"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product, err := decodeProduct(w, r)
	if err != nil {
		writeInputError(w, err)
		return
	}
	newProduct, err := h.service.Create(r.Context(), product)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newProduct)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	product, err := h.service.GetAll(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if len(product) == 0 {
		writeJSON(w, http.StatusOK, product)
	}
	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	product, err := decodeProduct(w, r)
	if err != nil {
		writeInputError(w, err)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	product.ID = id
	newproduct, err := h.service.Update(r.Context(), product)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, newproduct)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
