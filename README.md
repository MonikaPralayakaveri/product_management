# product_management

# Product Management System

This is a product management system built with Go, Gin, PostgreSQL, RabbitMQ, Redis, and a simple image processing microservice. The system allows users to add products, fetch products by ID, and manage product data with caching and asynchronous image processing.

## **Architecture Overview**

- **Modular Design**: The application follows a modular architecture with clear separation between different components:
  - **API Layer**: Handles HTTP requests and routes them to the appropriate handlers.
  - **Database Layer**: Manages database interactions using GORM ORM.
  - **Image Processing**: Asynchronous microservice for compressing product images.
  - **Caching**: Redis is used to cache frequently accessed product data.
  - **Message Queue**: RabbitMQ is used to send image processing tasks asynchronously.

## **Setup Instructions**

### **Prerequisites**
Before running the application, ensure you have the following installed:
- **Go 1.18+**
- **PostgreSQL**
- **Redis**
- **RabbitMQ**

### **Step-by-Step Setup**

1. **Clone the Repository**

   ```bash
   git clone https://github.com/MonikaPralayakaveri/product_management.git
   cd product_management
2.**Install Dependencies**

Ensure you have the necessary Go dependencies:

go mod tidy

3.**Set Up PostgreSQL**

Install PostgreSQL (if not already installed).

Create a database for the product management system:
sql code
CREATE DATABASE product_management;
Ensure the PostgreSQL server is running and accessible.
4.**Set Up Redis**

Install Redis (if not already installed).
Start Redis server on localhost:6379.
5.**Set Up RabbitMQ**

Install RabbitMQ (if not already installed).
Start RabbitMQ server on localhost:5672.
6.**Environment Configuration**

You can configure the necessary environment variables in a .env file or directly in your codebase. Example configuration:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=product_management
REDIS_HOST=localhost
REDIS_PORT=6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

7.**Run the Application**

Start the main application:

go run main.go
This will start the API server on http://localhost:8080.

**Database Schema**
The application uses PostgreSQL with the following table schema for products and users:

*Product Table*

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_price NUMERIC,
    product_images TEXT,
    compressed_images TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
*User Table*
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);


**Testing Coverage**
*Unit and Integration Tests*
The application includes unit and integration tests with over 90% code coverage. These tests cover the core components such as:

API routes
Image processing
Database interactions
To run tests, use the following command:

go test -v ./...
Note: Make sure your PostgreSQL, Redis, and RabbitMQ services are running before testing.

**Caching Strategy**
Redis is used to cache product data retrieved via GET /products/:id to reduce load on the database.
Cache invalidation is handled manually, ensuring that updates to product data are reflected in the cache.
**Image Processing (Asynchronous)**
After a product is created, its images are processed asynchronously.
RabbitMQ is used to queue image processing tasks.
Image Compression: The image processing service downloads images, resizes them, and stores them in the appropriate location (locally or on cloud storage).
Redis is used to store compressed image URLs for quick retrieval.
**Assumptions and Limitations**
The application assumes that the images provided for products are publicly accessible via URLs.
The image processing service currently only compresses and resizes images to a fixed width (800px), with further customizations possible.
**Usage**
POST /products
Creates a new product with the provided data.

Request Body Example:
json code
{
  "user_id": 1,
  "product_name": "Test Product",
  "product_description": "Sample description",
  "product_price": 99.99,
  "product_images": "image1.jpg,image2.jpg"
}
GET /products/:id
Fetches a product by its ID. If the product is cached, the response is returned from Redis for faster access.

Response Example:
json
Copy code
{
  "id": 1,
  "user_id": 1,
  "product_name": "Test Product",
  "product_description": "Sample description",
  "product_price": 99.99,
  "product_images": "image1.jpg,image2.jpg",
  "compressed_images": "compressed_image1.jpg,compressed_image2.jpg",
  "created_at": "2024-12-09T11:16:24+05:30",
  "updated_at": "2024-12-09T11:16:24+05:30"
}
