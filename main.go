package main

import (
	"log"
	"product-management/database"
	"product-management/api"  // Import the API handlers from the `api` package
	"product-management/cache"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.Connect()

	// Connect to Redis for caching
	cache.ConnectRedis()

	// Create a new Gin router
	r := gin.Default()

	// Define API routes
	r.POST("/products", api.CreateProduct)       // Correctly route to CreateProduct handler
	r.GET("/products/:id", api.GetProductByID)   // Correctly route to GetProductByID handler
	r.GET("/products", api.GetAllProducts)       // Correctly route to GetAllProducts handler

	// Start the server
	log.Println("Server is running on port 8080")
	r.Run(":8080") // Run on port 8080
}