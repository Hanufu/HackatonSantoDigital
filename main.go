package main

import (
	"log"
	"os"

	_ "github.com/Hanufu/HackatonSantoDigital/docs" // Importa a documentação gerada
	"github.com/Hanufu/HackatonSantoDigital/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	// Configuração do Logger
	file, err := os.Create("gin.log")
	if err != nil {
		log.Fatalf("could not create log file: %v", err)
	}
	gin.DefaultWriter = file

	// Middleware de Logging
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Rota para o Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas da API
	r.GET("/products/", handler.GetProducts)
	r.POST("/products/", handler.CreateProduct)
	r.GET("/products/:id", handler.GetProductByID)
	r.PUT("/products/:id", handler.UpdateProduct)
	r.DELETE("/products/:id", handler.DeleteProduct)

	// Obtém a porta do ambiente ou usa a padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
