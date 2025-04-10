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

## 🔌 API Endpoints

The microservice exposes the following API endpoints:

### Cart Management

- `POST /carts` - Create a new cart
- `GET /carts/{cartId}` - Get a cart by ID
- `DELETE /carts/{cartId}` - Delete a cart
- `POST /carts/{cartId}/items` - Add an item to a cart
- `PUT /carts/{cartId}/items/{itemId}` - Update a cart item
- `DELETE /carts/{cartId}/items/{itemId}` - Remove an item from a cart

### Checkout Process

- `POST /checkout/init` - Initialize a checkout from a cart
- `GET /checkout/{checkoutId}` - Get checkout details
- `PUT /checkout/{checkoutId}/shipping` - Update shipping details
- `PUT /checkout/{checkoutId}/payment-method` - Set payment method
- `POST /checkout/{checkoutId}/complete` - Complete the checkout process

### Shipping Management

- `POST /shipping/addresses` - Add a shipping address
- `GET /shipping/addresses` - Get all shipping addresses for a user
- `GET /shipping/addresses/{addressId}` - Get a shipping address by ID
- `PUT /shipping/addresses/{addressId}` - Update a shipping address
- `DELETE /shipping/addresses/{addressId}` - Delete a shipping address
- `GET /shipping/methods` - Get all available shipping methods

## 🐳 Docker

You can start the project with docker using this command:

```bash
docker-compose up --build
```

If you want to develop in docker with auto-reload add `-f docker-compose.dev.yml` to your docker command:

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

This command exposes the microservice on port 8001 and enables auto-reload for faster development.

## 🔄 Migrations

### Running Migrations

Database migrations are automatically applied when starting the application with docker-compose. However, you can also run them manually:

```bash
# Apply all pending migrations
docker-compose run --rm api migrate up

# Generate a new migration
docker-compose run --rm api migrate create -ext sql -dir /migrations/cart -seq create_carts_table

# Revert last migration
docker-compose run --rm api migrate down 1
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
