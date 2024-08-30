package model

import (
	"github.com/go-playground/validator/v10"
)

type Product struct {
	ProductKey            string  `json:"productKey" validate:"required"`
	ProductSubcategoryKey string  `json:"productSubcategoryKey" validate:"required"`
	ProductSKU            string  `json:"productSKU" validate:"required"`
	ProductName           string  `json:"productName" validate:"required,min=1"`
	ModelName             string  `json:"modelName" validate:"required,min=1"`
	ProductDescription    string  `json:"productDescription" validate:"required,min=1"`
	ProductColor          string  `json:"productColor" validate:"required,min=1"`
	ProductSize           string  `json:"productSize" validate:"required,min=1"`
	ProductStyle          string  `json:"productStyle" validate:"required,min=1"`
	ProductCost           float64 `json:"productCost" validate:"required,gt=0"`
	ProductPrice          float64 `json:"productPrice" validate:"required,gt=0"`
}

var validate = validator.New()

// ValidateProduct valida um objeto Product e retorna um erro se a validação falhar.
func ValidateProduct(product Product) error {
	err := validate.Struct(product)
	if err != nil {
		// Adicionar mensagens de erro personalizadas se necessário
		return err
	}
	return nil
}
