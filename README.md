# 🛒 Shopping Experience

This microservice handles shopping cart management and checkout processes for the Kiosko FIUBA e-commerce platform.

This project follows Domain-Driven Design principles with Hexagonal Architecture.

## 🏗️ Architecture

This project follows a Domain-Driven Design approach with Hexagonal Architecture (also known as Ports and Adapters) separating the domain core from infrastructure concerns:

```
shopping-experience/
├── cmd/                   # 🚀 Command entrypoints
│   └── shopping-experience/ # Main application
│       └── main.go        # Entry point
│
├── internal/              # 📦 Internal packages
│   ├── api/               # 🌐 API configuration
│   │   ├── router.go      # Router configuration
│   │   └── server.go      # Server setup
│   │
│   ├── cart/              # 🛒 Cart Management Bounded Context
│   │   ├── domain/        # 🧠 Domain Layer (Core)
│   │   │   ├── model/     # Domain entities and value objects
│   │   │   │   ├── cart.go      # Cart aggregate root
│   │   │   │   └── cart_item.go # CartItem value object
│   │   │   └── repository/ # Repository interfaces (ports)
│   │   │       └── cart_repository.go
│   │   │
│   │   ├── app/           # 📊 Application Layer
│   │   │   └── services/  # Application services
│   │   │       ├── cart_service.go # Cart service orchestrating domain objects
│   │   │       └── dto/   # Data Transfer Objects
│   │   │           └── cart_dto.go
│   │   │
│   │   └── infrastructure/ # 🔌 Infrastructure Layer (Adapters)
│   │       ├── http/      # HTTP API handlers
│   │       │   └── cart_handlers.go
│   │       └── postgresql/ # PostgreSQL repository implementation
│   │           └── cart_repository.go
│   │
│   ├── checkout/          # 📦 Checkout Process Bounded Context
│   │   ├── domain/        # 🧠 Domain Layer (Core)
│   │   │   ├── model/     # Domain entities and value objects
│   │   │   │   ├── checkout.go        # Checkout aggregate root
│   │   │   │   ├── delivery_option.go # DeliveryOption value object
│   │   │   │   ├── shipping_address.go # ShippingAddress entity
│   │   │   │   └── shipping_method.go # ShippingMethod entity
│   │   │   └── repository/ # Repository interfaces (ports)
│   │   │       ├── checkout_repository.go
│   │   │       └── shipping_repository.go
│   │   │
│   │   ├── app/           # 📊 Application Layer
│   │   │   └── services/  # Application services
│   │   │       ├── checkout_service.go # Checkout service
│   │   │       ├── shipping_service.go # Shipping service
│   │   │       └── dto/   # Data Transfer Objects
│   │   │           └── checkout_dto.go
│   │   │
│   │   └── infrastructure/ # 🔌 Infrastructure Layer (Adapters)
│   │       ├── http/      # HTTP API handlers
│   │       │   ├── checkout_handlers.go
│   │       │   └── shipping_handlers.go
│   │       ├── postgresql/ # PostgreSQL repository implementations
│   │       │   ├── checkout_repository.go
│   │       │   └── shipping_repository.go
│   │       └── clients/   # External service clients
│   │           ├── inventory_client.go
│   │           └── product_client.go
│   │
│   └── common/            # 🔄 Shared utilities
│       └── errors/        # Error handling
│
├── pkg/                   # 📚 Public packages
│   ├── logging/           # Logging utilities
│   └── validation/        # Input validation
│
├── docs/                  # 📝 API Documentation
│   ├── docs.go            # Generated Swagger docs
│   ├── swagger.json       # Swagger JSON specification
│   └── swagger.yaml       # Swagger YAML specification
│
└── migrations/            # 🔄 Database migrations
    ├── cart/              # Cart schema migrations
    └── checkout/          # Checkout schema migrations
```

## 🛒 Bounded Contexts

This microservice is organized into two main bounded contexts:

### Cart Management

The Cart Management bounded context handles operations related to shopping carts, including:

- Creating and retrieving carts
- Adding items to carts
- Updating quantities of items
- Removing items from carts

Key components:
- **Domain Models**: `Cart` (aggregate root), `CartItem` (value object)
- **Repository Interface**: `CartRepository`
- **Application Service**: `CartService`
- **Infrastructure**: PostgreSQL implementation, HTTP handlers

### Checkout Process

The Checkout Process bounded context handles operations related to the checkout process, including:

- Initiating a checkout from a cart
- Managing shipping addresses
- Selecting shipping methods
- Setting payment methods
- Completing the checkout

Key components:
- **Domain Models**: `Checkout` (aggregate root), `ShippingAddress` (entity), `ShippingMethod` (entity), `DeliveryOption` (value object)
- **Repository Interfaces**: `CheckoutRepository`, `ShippingRepository`
- **Application Services**: `CheckoutService`, `ShippingService`
- **Infrastructure**: PostgreSQL implementations, HTTP handlers

## 📝 API Documentation

The API is documented using Swagger (OpenAPI). The Swagger UI is available at:

```
http://localhost:8001/api/docs/
```

You can use this interface to:
- Explore all available API endpoints
- Test API calls directly from the browser
- View request/response models
- Understand the expected parameters and responses

To regenerate the Swagger documentation after making changes to API annotations, run:

```bash
swag init -g cmd/shopping-experience/main.go -o docs
```

## 🔌 API Endpoints

All API endpoints are available under the `/api` path prefix:

### Cart Management

- `POST /api/carts` - Create a new cart
- `GET /api/carts/{cartId}` - Get a cart by ID
- `DELETE /api/carts/{cartId}` - Delete a cart
- `POST /api/carts/{cartId}/items` - Add an item to a cart
- `PUT /api/carts/{cartId}/items/{itemId}` - Update a cart item
- `DELETE /api/carts/{cartId}/items/{itemId}` - Remove an item from a cart

### Checkout Process

- `POST /api/checkout/init` - Initialize a checkout from a cart
- `GET /api/checkout/{checkoutId}` - Get checkout details
- `PUT /api/checkout/{checkoutId}/shipping` - Update shipping details
- `PUT /api/checkout/{checkoutId}/payment-method` - Set payment method
- `POST /api/checkout/{checkoutId}/complete` - Complete the checkout process

### Shipping Management

- `POST /api/shipping/addresses` - Add a shipping address
- `GET /api/shipping/addresses` - Get all shipping addresses for a user
- `GET /api/shipping/addresses/{addressId}` - Get a shipping address by ID
- `PUT /api/shipping/addresses/{addressId}` - Update a shipping address
- `DELETE /api/shipping/addresses/{addressId}` - Delete a shipping address
- `GET /api/shipping/methods` - Get all available shipping methods

### Health Check

- `GET /api/health` - Check service health status

## 🐳 Docker

You can start the project with docker using this command:

```bash
docker compose up --build
```

This command exposes the microservice on port 8001 and enables auto-reload for faster development.

## 🔄 Migrations

### Running Migrations

Database migrations are automatically applied when using the migrations profile. Run migrations with:

```bash
# Apply all pending migrations
docker compose --profile migrations up
```

You can also run them manually:

```bash
# Apply all pending migrations
docker compose run --rm migrator

# Generate a new migration
docker compose run --rm api migrate create -ext sql -dir /migrations/cart -seq create_carts_table

# Revert last migration
docker compose run --rm api migrate down 1
```

## 🧪 Testing

To run tests:

```bash
go test ./...
```

To run tests with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🚀 Development

The codebase follows a modular structure based on Domain-Driven Design principles. When adding new features:

1. Start by defining domain models and repository interfaces in the domain layer
2. Implement application services that use these domain models
3. Create repository implementations in the infrastructure layer
4. Add HTTP handlers to expose the functionality via the API

## 🌐 Integration with Other Services

The Shopping Experience Microservice integrates with:

- **Product Catalog Microservice** - For retrieving product information
- **Inventory Microservice** - For checking product availability
- **Payment Hub** - For processing payments

## 📜 License

Copyright © 2025 Kiosko FIUBA. All rights reserved.
