package service

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Product struct {
	ProductKey            string  `json:"ProductKey"`
	ProductSubcategoryKey string  `json:"ProductSubcategoryKey"`
	ProductSKU            string  `json:"ProductSKU"`
	ProductName           string  `json:"ProductName"`
	ModelName             string  `json:"ModelName"`
	ProductDescription    string  `json:"ProductDescription"`
	ProductColor          string  `json:"ProductColor"`
	ProductSize           string  `json:"ProductSize"`
	ProductStyle          string  `json:"ProductStyle"`
	ProductCost           float64 `json:"ProductCost"`
	ProductPrice          float64 `json:"ProductPrice"`
}

// ReadProducts reads products from a CSV file.
func ReadProducts(filePath string) ([]Product, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header
	if err != nil {
		log.Printf("Error reading header: %v", err)
		return nil, err
	}

	var products []Product
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Error reading record: %v", err)
			return nil, err
		}
		if len(record) < 11 {
			continue // Skip incomplete records
		}

		cost, err := strconv.ParseFloat(record[9], 64)
		if err != nil {
			log.Printf("Error parsing cost: %v", err)
			return nil, err
		}
		price, err := strconv.ParseFloat(record[10], 64)
		if err != nil {
			log.Printf("Error parsing price: %v", err)
			return nil, err
		}

		products = append(products, Product{
			ProductKey:            record[0],
			ProductSubcategoryKey: record[1],
			ProductSKU:            record[2],
			ProductName:           record[3],
			ModelName:             record[4],
			ProductDescription:    record[5],
			ProductColor:          record[6],
			ProductSize:           record[7],
			ProductStyle:          record[8],
			ProductCost:           cost,
			ProductPrice:          price,
		})
	}

	return products, nil
}

// WriteProducts writes products to a CSV file.
func WriteProducts(filePath string, products []Product) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{
		"ProductKey", "ProductSubcategoryKey", "ProductSKU", "ProductName",
		"ModelName", "ProductDescription", "ProductColor", "ProductSize",
		"ProductStyle", "ProductCost", "ProductPrice",
	})
	if err != nil {
		log.Printf("Error writing header: %v", err)
		return err
	}

	for _, product := range products {
		err = writer.Write([]string{
			product.ProductKey, product.ProductSubcategoryKey, product.ProductSKU, product.ProductName,
			product.ModelName, product.ProductDescription, product.ProductColor, product.ProductSize,
			product.ProductStyle, strconv.FormatFloat(product.ProductCost, 'f', 4, 64),
			strconv.FormatFloat(product.ProductPrice, 'f', 4, 64),
		})
		if err != nil {
			log.Printf("Error writing product: %v", err)
			return err
		}
	}

	return nil
}

// AddProduct adds a new product to the CSV file.
func AddProduct(filePath string, newProduct Product) error {
	products, err := ReadProducts(filePath)
	if err != nil {
		return err
	}

	products = append(products, newProduct)
	return WriteProducts(filePath, products)
}

// GetProductByID retrieves a product by ID.
func GetProductByID(filePath, id string) (Product, error) {
	products, err := ReadProducts(filePath)
	if err != nil {
		return Product{}, err
	}

	for _, product := range products {
		if product.ProductKey == id {
			return product, nil
		}
	}

	return Product{}, errors.New("product not found")
}

// UpdateProduct updates an existing product by ID.
func UpdateProduct(filePath, id string, updatedProduct Product) error {
	products, err := ReadProducts(filePath)
	if err != nil {
		return err
	}

	for i, product := range products {
		if product.ProductKey == id {
			products[i] = updatedProduct
			return WriteProducts(filePath, products)
		}
	}

	return errors.New("product not found")
}

// DeleteProduct removes a product by ID.
func DeleteProduct(filePath, id string) error {
	products, err := ReadProducts(filePath)
	if err != nil {
		return err
	}

	var updatedProducts []Product
	for _, product := range products {
		if product.ProductKey != id {
			updatedProducts = append(updatedProducts, product)
		}
	}

	if len(updatedProducts) == len(products) {
		return errors.New("product not found")
	}

	return WriteProducts(filePath, updatedProducts)
}

// GetProducts retrieves products with pagination, filtering, and sorting.
func GetProducts(filePath string, page, pageSize int, filter, sort string) ([]Product, error) {
	products, err := ReadProducts(filePath)
	if err != nil {
		return nil, err
	}

	if filter != "" {
		var filteredProducts []Product
		for _, product := range products {
			if matchesFilter(product, filter) {
				filteredProducts = append(filteredProducts, product)
			}
		}
		products = filteredProducts
	}

	sortProducts(products, sort)

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(products) {
		return []Product{}, nil
	}
	if end > len(products) {
		end = len(products)
	}

	return products[start:end], nil
}

func matchesFilter(product Product, filter string) bool {
	filter = strings.ToLower(filter)
	return strings.Contains(strings.ToLower(product.ProductName), filter) ||
		strings.Contains(strings.ToLower(product.ProductDescription), filter) ||
		strings.Contains(strings.ToLower(product.ProductColor), filter) ||
		strings.Contains(strings.ToLower(product.ProductStyle), filter)
}

func sortProducts(products []Product, sortBy string) {
	switch sortBy {
	case "price":
		sort.Slice(products, func(i, j int) bool {
			return products[i].ProductPrice < products[j].ProductPrice
		})
	case "name":
		sort.Slice(products, func(i, j int) bool {
			return products[i].ProductName < products[j].ProductName
		})
	default:
		sort.Slice(products, func(i, j int) bool {
			return products[i].ProductKey < products[j].ProductKey
		})
	}
}
