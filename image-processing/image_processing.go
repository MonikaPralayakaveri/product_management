package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // PostgreSQL driver
	"product-management/models" // Assuming you have a models package with Product struct
	"encoding/json"  // For JSON marshalling/unmarshalling
	"strings"        // For splitting the image URLs
	"image"
)

// Initialize database connection
var db *gorm.DB

func init() {
	var err error
	// Connect to your PostgreSQL database (update the dsn with your database credentials)
	dsn := "host=localhost user=postgres password=yourpassword dbname=product_management port=5432 sslmode=disable"
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully")
}

// Struct to handle the message from the queue
type ImageProcessingTask struct {
	ProductID     uint   `json:"product_id"`
	ProductImages string `json:"product_images"`
}

// Process and compress the image
func processImage(imageURL string) ([]byte, error) {
	// Download the image from the URL
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Resize the image (for example, resize to 800px wide while maintaining aspect ratio)
	img = resize.Resize(800, 0, img, resize.Lanczos3)

	// Compress to JPEG format (you can also choose PNG or another format)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	return buf.Bytes(), nil
}

// Save the image to disk (or to an S3 bucket, depending on your use case)
func saveImage(compressedImage []byte, productID uint) string {
	// Save the image locally (you can use AWS S3 or other cloud storage)
	fileName := fmt.Sprintf("compressed_%d.jpg", productID)
	err := ioutil.WriteFile(fileName, compressedImage, 0644)
	if err != nil {
		log.Printf("Failed to save compressed image: %v", err)
		return ""
	}
	return fileName
}

// Update the product in the database with the compressed image URL
func updateProductCompressedImages(productID uint, compressedImageURL string) {
	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		log.Printf("Product not found with ID %d: %v", productID, err)
		return
	}

	// Update the compressed images field
	product.CompressedImages = compressedImageURL
	if err := db.Save(&product).Error; err != nil {
		log.Printf("Failed to update product %d: %v", productID, err)
		return
	}

	log.Printf("Updated product %d with compressed image URL %s", productID, compressedImageURL)
}

// Process the messages from RabbitMQ
func processQueue() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue to receive messages
	q, err := ch.QueueDeclare(
		"image_queue", // Queue name
		false,         // Durable
		false,         // Delete when unused
		false,         // Exclusive
		false,         // No-wait
		nil,           // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Start consuming messages from the queue
	msgs, err := ch.Consume(
		q.Name,       // Queue name
		"",           // Consumer name
		true,         // Auto-acknowledge
		false,        // Exclusive
		false,        // No-wait
		false,        // Arguments
		amqp.Table{}, // Provide an empty table for additional arguments
	)
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	// Process each message
	for msg := range msgs {
		var task ImageProcessingTask
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Process and compress each image
		compressedImages := []string{}
		for _, imageURL := range splitImageURLs(task.ProductImages) {
			compressedImage, err := processImage(imageURL)
			if err != nil {
				log.Printf("Failed to process image %s: %v", imageURL, err)
				continue
			}

			// Save the compressed image
			compressedImageURL := saveImage(compressedImage, task.ProductID)
			if compressedImageURL != "" {
				compressedImages = append(compressedImages, compressedImageURL)
			}
		}

		// Update the product's compressed images in the database
		if len(compressedImages) > 0 {
			updateProductCompressedImages(task.ProductID, compressedImages[0]) // Only store the first compressed image for simplicity
		}
	}
}

// Split the product image URLs into an array
func splitImageURLs(imageURLs string) []string {
	return strings.Split(imageURLs, ",")
}

// Main entry point for the image processing service
func main() {
	log.Println("Image processing service started...")
	processQueue()
}