#!/bin/bash

# Script to populate the database with sample data for the shopping-experience microservice

API_URL="http://localhost:8001/api"
SHIPPING_URL="$API_URL/shipping/methods"
CARTS_URL="$API_URL/carts"
CHECKOUT_URL="$API_URL/checkout"
SHIPPING_ADDRESS_URL="$API_URL/shipping/addresses"

echo "Populating the shopping-experience database..."

# Create shipping methods
echo "Creating shipping methods..."

STANDARD_ID=$(curl -s -X POST "$SHIPPING_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "11111111-1111-1111-1111-111111111111",
    "name": "Standard Shipping",
    "description": "Delivery in 3-5 business days",
    "price": 5.99,
    "estimatedDays": 5
  }' | jq -r '.id')

echo "Standard Shipping created with ID: $STANDARD_ID"

EXPRESS_ID=$(curl -s -X POST "$SHIPPING_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "22222222-2222-2222-2222-222222222222",
    "name": "Express Shipping",
    "description": "Delivery in 1-2 business days",
    "price": 12.99,
    "estimatedDays": 2
  }' | jq -r '.id')

echo "Express Shipping created with ID: $EXPRESS_ID"

SAME_DAY_ID=$(curl -s -X POST "$SHIPPING_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "33333333-3333-3333-3333-333333333333",
    "name": "Same Day Delivery",
    "description": "Delivery within 24 hours (select areas only)",
    "price": 19.99,
    "estimatedDays": 1
  }' | jq -r '.id')

echo "Same Day Delivery created with ID: $SAME_DAY_ID"

# Create sample user carts
echo "Creating sample carts..."

USER1_CART_ID=$(curl -s -X POST "$CARTS_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "b7f060bc-aaba-4861-b71b-74c3c85badd3"
  }' | jq -r '.id')

echo "Cart created for User 1 with ID: $USER1_CART_ID"

# Add items to the cart
echo "Adding items to User 1's cart..."

curl -s -X POST "$CARTS_URL/$USER1_CART_ID/items" \
  -H "Content-Type: application/json" \
  -d '{
    "productId": "4abe25d0-c7f3-4d98-9a90-e21587ada874",
    "name": "Laptop HP Pavilion 15",
    "price": 12999.99,
    "quantity": 1,
    "imageUrl": "https://example.com/hp-pavilion-15.jpg"
  }'

curl -s -X POST "$CARTS_URL/$USER1_CART_ID/items" \
  -H "Content-Type: application/json" \
  -d '{
    "productId": "8c6e7315-95b0-4f94-b7ac-1e95f738ce7b",
    "name": "Mouse Logitech G305",
    "price": 3999.99,
    "quantity": 2,
    "imageUrl": "https://example.com/logitech-g305.jpg"
  }'

echo "Items added to User 1's cart"

# Create another cart for a different user
USER2_CART_ID=$(curl -s -X POST "$CARTS_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "dce9e6b8-f84d-4ab8-8fcb-f5225a0fc213"
  }' | jq -r '.id')

echo "Cart created for User 2 with ID: $USER2_CART_ID"

# Add an item to the second cart
echo "Adding items to User 2's cart..."

curl -s -X POST "$CARTS_URL/$USER2_CART_ID/items" \
  -H "Content-Type: application/json" \
  -d '{
    "productId": "b47a547b-2da3-4b9e-a68f-c81b9c6cd2db",
    "name": "Termo Stanley Clásico 1.4lts Negro",
    "price": 4999.99,
    "quantity": 1,
    "imageUrl": "https://example.com/stanley-black.jpg"
  }'

echo "Items added to User 2's cart"

# Create shipping addresses
echo "Creating shipping addresses..."

ADDRESS1_ID=$(curl -s -X POST "$SHIPPING_ADDRESS_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "b7f060bc-aaba-4861-b71b-74c3c85badd3",
    "name": "Juan Pérez",
    "streetAddress": "Av. Paseo Colón 850",
    "city": "Ciudad Autónoma de Buenos Aires",
    "state": "Buenos Aires",
    "postalCode": "C1063ACV",
    "country": "Argentina",
    "phone": "+5491122334455",
    "isDefault": true
  }' | jq -r '.id')

echo "Shipping address created for User 1 with ID: $ADDRESS1_ID"

ADDRESS2_ID=$(curl -s -X POST "$SHIPPING_ADDRESS_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "dce9e6b8-f84d-4ab8-8fcb-f5225a0fc213",
    "name": "María González",
    "streetAddress": "Av. Las Heras 2214",
    "city": "Ciudad Autónoma de Buenos Aires",
    "state": "Buenos Aires",
    "postalCode": "C1127AAL",
    "country": "Argentina",
    "phone": "+5491155667788",
    "isDefault": true
  }' | jq -r '.id')

echo "Shipping address created for User 2 with ID: $ADDRESS2_ID"

# Create a checkout for User 1
echo "Creating a checkout for User 1..."

CHECKOUT1_ID=$(curl -s -X POST "$CHECKOUT_URL/init" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "b7f060bc-aaba-4861-b71b-74c3c85badd3",
    "cartId": "'"$USER1_CART_ID"'"
  }' | jq -r '.id')

echo "Checkout created for User 1 with ID: $CHECKOUT1_ID"

# Update the checkout with shipping information
echo "Updating checkout with shipping information..."

curl -s -X PUT "$CHECKOUT_URL/$CHECKOUT1_ID/shipping" \
  -H "Content-Type: application/json" \
  -d '{
    "shippingAddressId": "'"$ADDRESS1_ID"'",
    "shippingMethodId": "'"$EXPRESS_ID"'"
  }'

echo "Checkout updated with shipping information"

# Set payment method for the checkout
echo "Setting payment method for the checkout..."

curl -s -X PUT "$CHECKOUT_URL/$CHECKOUT1_ID/payment-method" \
  -H "Content-Type: application/json" \
  -d '{
    "paymentMethod": "CREDIT_CARD"
  }'

echo "Payment method set for checkout"

echo "✅ Database populated successfully" 