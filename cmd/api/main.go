package main

import (
	"firstgoproject/internal/handler"
	"firstgoproject/internal/repository"
	"firstgoproject/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	productRepo := repository.NewInMemoryProductRepository()
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	router := mux.NewRouter()

	router.HandleFunc("/health/live", GetStatus).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/products", productHandler.CreateProduct).
		Methods(http.MethodPost)

	router.HandleFunc("/api/v1/products", productHandler.GetProducts).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/products/{id}", productHandler.GetProductByID).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/products/{id}", productHandler.UpdateProduct).
		Methods(http.MethodPut)

	router.HandleFunc("/api/v1/products/{id}", productHandler.DeleteProduct).
		Methods(http.MethodDelete)

	log.Println("API server running at :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
