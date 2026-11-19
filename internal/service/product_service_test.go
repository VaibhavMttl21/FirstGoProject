package service

import (
	"firstgoproject/internal/model"
	"testing"
)

// func TestExample(t *testing.T) {
//     expected := 10
//     actual := 5 + 5

//     if actual != expected {
//         t.Errorf("expected %d, got %d", expected, actual)
//     }
// }

func TestValidateProduct(t *testing.T) {
	tests := []struct {
		name        string
		product     model.Product
		expectError bool
	}{
		{
			name: "valid product",
			product: model.Product{
				Name:          "Laptop",
				PriceCents:    75000,
				Currency:      "INR",
				StockQuantity: 10,
			},
			expectError: false,
		},
		{
			name: "empty product name",
			product: model.Product{
				Name:          "",
				PriceCents:    75000,
				Currency:      "INR",
				StockQuantity: 10,
			},
			expectError: true,
		},
		{
			name: "negative price",
			product: model.Product{
				Name:          "Laptop",
				PriceCents:    -100,
				Currency:      "INR",
				StockQuantity: 10,
			},
			expectError: true,
		},
		{
			name: "negative stock",
			product: model.Product{
				Name:          "Laptop",
				PriceCents:    75000,
				Currency:      "INR",
				StockQuantity: -5,
			},
			expectError: true,
		},
		{
			name: "invalid currency",
			product: model.Product{
				Name:          "Laptop",
				PriceCents:    75000,
				Currency:      "inr",
				StockQuantity: 10,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProduct(tt.product)

			if tt.expectError && err == nil {
				t.Errorf("expected an error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}