# 🏪 Golang Store API

A RESTful API for store management built with **Go**, following **Clean Architecture** principles. This API provides full CRUD operations for Categories and Products, along with server-side pagination, filtering, sorting, Redis caching, and a dashboard report endpoint.

---
## 📑 API Documentation

You can import [this collection ](https://github.com/ZaidanNur/golang_store_api/blob/main/test-elabram.postman_collection.json) into Postman to test the API endpoints.

---
## 📚 Tech Stack

| Technology | Version | Description |
|------------|---------|-------------|
| **Go** | 1.25+ | Primary programming language |
| **Gin** | v1.11 | Fast and lightweight HTTP web framework |
| **GORM** | v1.31 | ORM (Object-Relational Mapping) for Go |
| **PostgreSQL** | 15+ | Primary relational database |
| **Redis** | - | In-memory cache for frequently accessed data |
| **Atlas** | - | Declarative schema-based database migration tool |
| **Air** | - | Live-reload tool for development |

---

## 📦 Key Libraries

### Core Dependencies

| Library | Purpose |
|---------|---------|
| [`github.com/gin-gonic/gin`](https://github.com/gin-gonic/gin) | Main HTTP framework — routing, middleware, and request handling |
| [`gorm.io/gorm`](https://gorm.io) | ORM for database interaction (query builder, associations, migrations) |
| [`gorm.io/driver/postgres`](https://gorm.io/docs/connecting_to_the_database.html) | PostgreSQL driver for GORM |
| [`github.com/redis/go-redis/v9`](https://github.com/redis/go-redis) | Redis client for data caching |
| [`github.com/go-playground/validator/v10`](https://github.com/go-playground/validator) | Struct-level validation for request bodies |
| [`github.com/joho/godotenv`](https://github.com/joho/godotenv) | Loads environment variables from `.env` file |
| [`ariga.io/atlas-provider-gorm`](https://github.com/ariga/atlas-provider-gorm) | Atlas migration integration with GORM schema definitions |

### Development Tools

| Tool | Purpose |
|------|---------|
| [Air](https://github.com/air-verse/air) | Automatic hot-reload during development |
| [Atlas](https://atlasgo.io/) | Declarative database migration management |

---

## 🏗️ Project Structure (Clean Architecture)

```
├── cmd/
│   ├── app/
│   │   └── main.go                  # Application entry point
│   └── loader/
│       └── main.go                  # Schema loader for Atlas migrations
├── internal/
│   ├── cache/
│   │   └── redis_cache.go           # Redis cache wrapper (Get, Set, Delete)
│   ├── delivery/
│   │   ├── helper/
│   │   │   └── validator_helper.go  # Custom validation error messages
│   │   └── http/
│   │       ├── category_handler.go  # HTTP handlers for Category endpoints
│   │       └── product_handler.go   # HTTP handlers for Product endpoints
│   ├── domain/
│   │   ├── category.go              # Category entity & interfaces
│   │   └── product.go               # Product entity & interfaces
│   ├── dto/
│   │   ├── category_dto.go          # Request/Response DTOs for Category
│   │   └── product_dto.go           # Request/Response DTOs for Product
│   ├── repository/
│   │   ├── category_repository.go   # Category data access layer
│   │   └── product_repository.go    # Product data access layer
│   └── usecase/
│       ├── category_usecase.go      # Category business logic
│       └── product_usecase.go       # Product business logic
├── migrations/                      # Atlas database migration files
├── .air.toml                        # Air configuration (hot-reload)
├── atlas.hcl                        # Atlas migration configuration
├── .env.example                     # Environment variable template
├── go.mod                           # Go module dependencies
└── go.sum                           # Dependency checksums
```

---

## ⚙️ Setup & Installation

### Prerequisites

- **Go** >= 1.25 → [Download Go](https://go.dev/dl/)
- **PostgreSQL** >= 15 → [Download PostgreSQL](https://www.postgresql.org/download/)
- **Redis** (optional, for caching) → [Download Redis](https://redis.io/download)
- **Atlas CLI** (for migrations) → [Install Atlas](https://atlasgo.io/getting-started#installation)
- **Air** (for dev hot-reload) → `go install github.com/air-verse/air@latest`

### Step-by-step Setup

#### 1. Clone the Repository

```bash
git clone https://github.com/ZaidanNur/golang_store_api.git
cd golang_store_api
```

#### 2. Install Dependencies

```bash
go mod download
```

#### 3. Configure Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

`.env` file contents:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=test_elabram
DB_PORT=5432
SERVER_PORT=8080
REDIS_URL=localhost:6379
```

#### 4. Create the PostgreSQL Database

```sql
CREATE DATABASE test_elabram;
```

#### 5. Run Database Migrations (Atlas)

```bash
atlas migrate apply --env local --url "postgres://postgres:postgres@localhost:5432/test_elabram?sslmode=disable"
```

#### 6. Start the Server

**Development (with hot-reload):**

```bash
air
```

**Production:**

```bash
go build -o app.exe ./cmd/app
./app.exe
```

The server will start at `http://localhost:8080` (or the port configured in `.env`).

---

## 📖 API Documentation

**Base URL:** `http://localhost:8080`

All responses are in JSON format. Successful responses generally follow this structure:

```json
{
  "status": 200,
  "message": "operation success",
  "data": { ... }
}
```

---

### 📂 Categories

#### Get All Categories

```
GET /category
```

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "get categories success",
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices and gadgets",
      "created_at": "2026-02-15T10:00:00+07:00",
      "updated_at": "2026-02-15T10:00:00+07:00"
    }
  ]
}
```

---

#### Get Category by ID

```
GET /category/:id
```

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| `id` | `int` | Path | Category ID |

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "get category success",
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "created_at": "2026-02-15T10:00:00+07:00",
    "updated_at": "2026-02-15T10:00:00+07:00"
  }
}
```

**Error** `404 Not Found`:

```json
{ "error": "category not found" }
```

---

#### Create Category

```
POST /category
```

**Request Body:**

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `name` | `string` | ✅ | `required` | Category name |
| `description` | `string` | ✅ | `required` | Category description |

**Example Request:**

```json
{
  "name": "Electronics",
  "description": "Electronic devices and gadgets"
}
```

**Response** `201 Created`:

```json
{
  "status": 201,
  "message": "category created successfully",
  "data": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "created_at": "2026-02-15T10:00:00+07:00",
    "updated_at": "2026-02-15T10:00:00+07:00"
  }
}
```

**Error** `400 Bad Request` (Validation failed):

```json
{
  "status": 400,
  "message": "Validation failed",
  "errors": {
    "Name": "This field is required",
    "Description": "This field is required"
  }
}
```

---

#### Update Category

```
PUT /category/:id
```

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| `id` | `int` | Path | Category ID |

**Request Body** (partial update):

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | `string` | ❌ | New category name |
| `description` | `string` | ❌ | New category description |

**Example Request:**

```json
{
  "name": "Updated Electronics"
}
```

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "category updated successfully",
  "data": {
    "id": 1,
    "name": "Updated Electronics",
    "description": "Electronic devices and gadgets",
    "created_at": "2026-02-15T10:00:00+07:00",
    "updated_at": "2026-02-15T12:00:00+07:00"
  }
}
```

---

#### Delete Category

```
DELETE /category/:id
```

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| `id` | `int` | Path | Category ID |

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "category deleted successfully"
}
```

---

### 📦 Products

#### Get All Products (Paginated)

```
GET /products
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | `int` | `1` | Page number (min: 1) |
| `limit` | `int` | `10` | Items per page (min: 1, max: 100) |
| `name` | `string` | - | Filter by product name (partial match) |
| `category_id` | `int` | - | Filter by category ID |
| `price_min` | `int` | - | Minimum price filter |
| `price_max` | `int` | - | Maximum price filter |
| `stock_min` | `int` | - | Minimum stock filter |
| `stock_max` | `int` | - | Maximum stock filter |
| `sort_by` | `string` | `created_at` | Column to sort by |
| `sort_order` | `string` | `desc` | Sort direction (`asc` / `desc`) |

**Example Request:**

```
GET /products?page=1&limit=5&name=laptop&sort_by=price&sort_order=asc
```

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "get products success",
  "data": [
    {
      "id": 1,
      "name": "Laptop Pro",
      "description": "High-end laptop",
      "price": 15000000,
      "stock_quantity": 50,
      "is_active": true,
      "category_id": 1,
      "category": {
        "id": 1,
        "name": "Electronics",
        "description": "Electronic devices",
        "created_at": "2026-02-15T10:00:00+07:00",
        "updated_at": "2026-02-15T10:00:00+07:00"
      },
      "created_at": "2026-02-15T10:00:00+07:00",
      "updated_at": "2026-02-15T10:00:00+07:00"
    }
  ],
  "page": 1,
  "limit": 5,
  "total_items": 25,
  "total_pages": 5
}
```

---

#### Get Product by ID

```
GET /products/:id
```

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| `id` | `int` | Path | Product ID |

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "get product success",
  "data": {
    "id": 1,
    "name": "Laptop Pro",
    "description": "High-end laptop",
    "price": 15000000,
    "stock_quantity": 50,
    "is_active": true,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices",
      "created_at": "2026-02-15T10:00:00+07:00",
      "updated_at": "2026-02-15T10:00:00+07:00"
    },
    "created_at": "2026-02-15T10:00:00+07:00",
    "updated_at": "2026-02-15T10:00:00+07:00"
  }
}
```

---

#### Create Product

```
POST /products
```

**Request Body:**

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `name` | `string` | ✅ | `required` | Product name |
| `description` | `string` | ✅ | `required` | Product description |
| `price` | `int` | ✅ | `required, gt=0` | Price (must be > 0) |
| `stock_quantity` | `int` | ✅ | `required, gte=0` | Stock quantity (must be >= 0) |
| `is_active` | `bool` | ❌ | - | Whether product is active |
| `category_id` | `int` | ✅ | `required, gt=0` | Associated category ID |

**Example Request:**

```json
{
  "name": "Laptop Pro",
  "description": "High-end laptop for professionals",
  "price": 15000000,
  "stock_quantity": 50,
  "is_active": true,
  "category_id": 1
}
```

**Response** `201 Created`:

```json
{
  "status": 201,
  "message": "product created successfully",
  "data": {
    "id": 1,
    "name": "Laptop Pro",
    "description": "High-end laptop for professionals",
    "price": 15000000,
    "stock_quantity": 50,
    "is_active": true,
    "category_id": 1,
    "category": { ... },
    "created_at": "2026-02-15T10:00:00+07:00",
    "updated_at": "2026-02-15T10:00:00+07:00"
  }
}
```

**Error** `400 Bad Request` (Validation failed):

```json
{
  "status": 400,
  "message": "Validation failed",
  "errors": {
    "Name": "This field is required",
    "Price": "Must be greater than 0",
    "StockQuantity": "Must be greater than or equal to 0",
    "CategoryID": "Must be greater than 0"
  }
}
```

---

#### Update Product

```
PUT /products/:id
```

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| `id` | `int` | Path | Product ID |

**Request Body** (partial update — all fields optional):

| Field | Type | Validation | Description |
|-------|------|------------|-------------|
| `name` | `string` | - | New product name |
| `description` | `string` | - | New product description |
| `price` | `int` | `gt=0` | New price |
| `stock_quantity` | `int` | `gte=0` | New stock quantity |
| `is_active` | `bool` | - | New active status |
| `category_id` | `int` | `gt=0` | New category ID |

**Example Request:**

```json
{
  "price": 14000000,
  "stock_quantity": 45
}
```

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "name": "Laptop Pro",
    "description": "High-end laptop for professionals",
    "price": 14000000,
    "stock_quantity": 45,
    "is_active": true,
    "category_id": 1,
    "category": { ... },
    "created_at": "2026-02-15T10:00:00+07:00",
    "updated_at": "2026-02-15T12:00:00+07:00"
  }
}
```

---

#### Delete Product

```
DELETE /products/:id
```

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| `id` | `int` | Path | Product ID |

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "Product deleted successfully"
}
```

---

#### Get Product Report (Dashboard)

```
GET /products/report
```

Returns a dashboard-style summary report of all products. This data is **cached with Redis** for improved performance.

**Response** `200 OK`:

```json
{
  "status": 200,
  "message": "get product report success",
  "data": {
    "total_products": 25,
    "total_stock": 1250,
    "average_price": 5000000.00,
    "products": [
      {
        "id": 1,
        "name": "Laptop Pro",
        "category_name": "Electronics",
        "price": 15000000,
        "stock_quantity": 50
      }
    ]
  }
}
```

---

## 🔧 Key Features

### ✅ Clean Architecture
The project follows Clean Architecture with clearly separated layers:
- **Domain** — Entities and interface contracts for repositories & usecases
- **Repository** — Data access layer (GORM)
- **Usecase** — Business logic layer
- **Delivery/HTTP** — HTTP handlers for receiving and responding to requests
- **DTO** — Data Transfer Objects for request validation & response formatting

### ✅ Server-side Pagination, Filtering & Sorting
The `GET /products` endpoint supports:
- Pagination (`page`, `limit`)
- Filtering by name, category, price range, stock range
- Sorting by any column (`sort_by`, `sort_order`)

### ✅ Redis Caching
Product report data is cached in Redis to reduce database query load. The cache is automatically invalidated when data changes occur.

### ✅ Request Validation
Automatic request body validation using `go-playground/validator` with informative per-field error messages.

### ✅ Database Migrations (Atlas)
Database schema is managed declaratively from GORM models using Atlas, ensuring the schema always stays in sync with the code.

### ✅ Hot-Reload Development
Uses Air for a comfortable development experience — the server automatically restarts when code changes are detected.

---

## 📝 Notes

- Redis is **optional**. If Redis is unavailable, the application runs normally without caching.
- All endpoints return JSON responses.
- Validation errors return per-field details for easy debugging.
