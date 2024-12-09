package api

import (
	"net/http"
	"product-management/cache"  // Import Redis cache functions
	"product-management/database"
	"product-management/models"
	"strconv" // For converting numbers to strings
	"github.com/gin-gonic/gin"
)

// CreateProduct handles creating a new product
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save product to the database
	if result := database.DB.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Cache the product in Redis after saving it
	cache.CacheProduct(strconv.FormatUint(uint64(product.ID), 10), product)

	c.JSON(http.StatusCreated, product)
}

// GetProductByID handles fetching a product by ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// Check if the product exists in Redis cache
	cachedProduct, err := cache.GetFromCache(id)
	if err == nil {
		// If product is found in cache, return it
		c.JSON(http.StatusOK, cachedProduct)
		return
	}

	// If not found in cache, fetch the product from the database
	if result := database.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Cache the product in Redis
	cache.CacheProduct(id, product)

	c.JSON(http.StatusOK, product)
}

// GetAllProducts handles fetching all products with optional filtering
func GetAllProducts(c *gin.Context) {
	var products []models.Product
	userID := c.DefaultQuery("user_id", "")
	priceMin := c.DefaultQuery("price_min", "")
	priceMax := c.DefaultQuery("price_max", "")
	name := c.DefaultQuery("name", "")

	// Build query with optional filters
	query := database.DB.Model(&models.Product{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if priceMin != "" && priceMax != "" {
		query = query.Where("product_price BETWEEN ? AND ?", priceMin, priceMax)
	}
	if name != "" {
		query = query.Where("product_name LIKE ?", "%"+name+"%")
	}

	query = query.Where("deleted_at IS NULL")
	// Execute query
	result := query.Find(&products)

	// Check if any products were found
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No products found"})
		return
	}

	// Return the products
	c.JSON(http.StatusOK, products)
}
