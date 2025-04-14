FROM golang:1.23-alpine3.19 AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o shopping-experience ./cmd/shopping-experience/main.go

# Build the migrator
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o migrator ./cmd/migrator/main.go

# Development stage
FROM golang:1.23-alpine AS dev

# Set working directory
WORKDIR /app

# Install development dependencies
RUN apk add --no-cache gcc musl-dev postgresql-client git

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy air configuration from project
COPY .air.toml .

# Command to run air for hot reloading
CMD ["air", "-c", ".air.toml"]

# Runtime stage - This is now the final and default stage
FROM alpine:3.19 AS prod

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata postgresql-client

# Set working directory
WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/shopping-experience .
COPY --from=builder /app/migrator .

# Create migrations directory (will be populated at runtime if needed)
RUN mkdir -p /app/migrations/

# Set executable permissions
RUN chmod +x /app/shopping-experience
RUN chmod +x /app/migrator

# Expose the application port
EXPOSE 8001

# Command to run the application
CMD ["/app/shopping-experience"]
