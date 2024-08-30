package handler

import (
	"net/http"
	"strconv"

	"github.com/Hanufu/HackatonSantoDigital/service"
	"github.com/gin-gonic/gin"
)

const filePath = "archives/AdventureWorks_Products.csv"

// CreateProduct godoc
// @Summary Create a new product
// @Description Adds a new product to the CSV file
// @Tags products
// @Accept json
// @Produce json
// @Param product body service.Product true "Product data"
// @Success 201 {object} service.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/ [post]
func CreateProduct(c *gin.Context) {
	var newProduct service.Product
	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	if err := service.AddProduct(filePath, newProduct); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, newProduct)
}

// GetProducts godoc
// @Summary Get a list of products
// @Description Returns a list of products with pagination, filtering, and sorting
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param filter query string false "Filter criteria"
// @Param sort query string false "Sorting criteria"
// @Success 200 {array} service.Product
// @Failure 500 {object} ErrorResponse
// @Router /products/ [get]
func GetProducts(c *gin.Context) {
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSizeStr := c.Query("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	filter := c.Query("filter")
	sort := c.Query("sort")

	products, err := service.GetProducts(filePath, page, pageSize, filter, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Returns details of a single product based on ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} service.Product
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := service.GetProductByID(filePath, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Updates a product based on ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body service.Product true "Updated product data"
// @Success 200 {object} service.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct service.Product
	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		return
	}

	if err := service.UpdateProduct(filePath, id, updatedProduct); err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Removes a product based on ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := service.DeleteProduct(filePath, id); err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
