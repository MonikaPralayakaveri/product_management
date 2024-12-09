package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

// Initialize Redis connection
func ConnectRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis default port
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis successfully")
}

// Cache product data in Redis
func CacheProduct(key string, product interface{}) {
	// Serialize the product into a JSON string
	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Println("Failed to serialize product:", err)
		return
	}

	// Store the serialized product in Redis with a 10-minute expiration time
	err = rdb.Set(ctx, key, productJSON, 10*time.Minute).Err()
	if err != nil {
		log.Println("Failed to cache product:", err)
	}
}

// Retrieve product data from Redis
func GetFromCache(key string) (string, error) {
	return rdb.Get(ctx, key).Result()  // Retrieve the data using the key
}