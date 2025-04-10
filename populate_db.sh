#!/bin/bash

# Script to populate the database with sample data for the shopping-experience microservice

API_URL="http://localhost:8001/api"
CARTS_URL="$API_URL/carts"

echo "Populating the shopping-experience database..."

# Define UUIDs manually since some endpoints may not be implemented yet
STANDARD_ID="11111111-1111-1111-1111-111111111111"
EXPRESS_ID="22222222-2222-2222-2222-222222222222"
SAME_DAY_ID="33333333-3333-3333-3333-333333333333"
ADDRESS1_ID="44444444-4444-4444-4444-444444444444"
ADDRESS2_ID="55555555-5555-5555-5555-555555555555"

echo "Using predefined shipping method IDs:"
echo "Standard Shipping ID: $STANDARD_ID"
echo "Express Shipping ID: $EXPRESS_ID"
echo "Same Day Delivery ID: $SAME_DAY_ID"
echo "Shipping Address 1 ID: $ADDRESS1_ID"
echo "Shipping Address 2 ID: $ADDRESS2_ID"

# Create sample user carts
echo "Creating sample carts..."

USER1_CART_ID=$(curl -s -X POST "$CARTS_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "b7f060bc-aaba-4861-b71b-74c3c85badd3"
  }' | jq -r '.id')

if [ -z "$USER1_CART_ID" ] || [ "$USER1_CART_ID" == "null" ]; then
  echo "❌ Failed to create cart for User 1"
else
  echo "Cart created for User 1 with ID: $USER1_CART_ID"

  # Add items to the cart
  echo "Adding items to User 1's cart..."

  ITEM_RESPONSE=$(curl -s -X POST "$CARTS_URL/$USER1_CART_ID/items" \
    -H "Content-Type: application/json" \
    -d '{
      "productId": "4abe25d0-c7f3-4d98-9a90-e21587ada874",
      "name": "Laptop HP Pavilion 15",
      "price": 12999.99,
      "quantity": 1,
      "imageUrl": "https://example.com/hp-pavilion-15.jpg"
    }')

  if [[ "$ITEM_RESPONSE" == *"error"* ]] || [[ "$ITEM_RESPONSE" == *"status"*":"*"500"* ]]; then
    echo "❌ Failed to add laptop to User 1's cart: $ITEM_RESPONSE"
  else
    echo "✅ Laptop added to User 1's cart"
  fi

  ITEM_RESPONSE=$(curl -s -X POST "$CARTS_URL/$USER1_CART_ID/items" \
    -H "Content-Type: application/json" \
    -d '{
      "productId": "8c6e7315-95b0-4f94-b7ac-1e95f738ce7b",
      "name": "Mouse Logitech G305",
      "price": 3999.99,
      "quantity": 2,
      "imageUrl": "https://example.com/logitech-g305.jpg"
    }')

  if [[ "$ITEM_RESPONSE" == *"error"* ]] || [[ "$ITEM_RESPONSE" == *"status"*":"*"500"* ]]; then
    echo "❌ Failed to add mouse to User 1's cart: $ITEM_RESPONSE"
  else
    echo "✅ Mouse added to User 1's cart"
  fi
fi

# Create another cart for a different user
USER2_CART_ID=$(curl -s -X POST "$CARTS_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "dce9e6b8-f84d-4ab8-8fcb-f5225a0fc213"
  }' | jq -r '.id')

if [ -z "$USER2_CART_ID" ] || [ "$USER2_CART_ID" == "null" ]; then
  echo "❌ Failed to create cart for User 2"
else
  echo "Cart created for User 2 with ID: $USER2_CART_ID"

  # Add an item to the second cart
  echo "Adding items to User 2's cart..."

  ITEM_RESPONSE=$(curl -s -X POST "$CARTS_URL/$USER2_CART_ID/items" \
    -H "Content-Type: application/json" \
    -d '{
      "productId": "b47a547b-2da3-4b9e-a68f-c81b9c6cd2db",
      "name": "Termo Stanley Clásico 1.4lts Negro",
      "price": 4999.99,
      "quantity": 1,
      "imageUrl": "https://example.com/stanley-black.jpg"
    }')

  if [[ "$ITEM_RESPONSE" == *"error"* ]] || [[ "$ITEM_RESPONSE" == *"status"*":"*"500"* ]]; then
    echo "❌ Failed to add thermos to User 2's cart: $ITEM_RESPONSE"
  else
    echo "✅ Thermos added to User 2's cart"
  fi
fi

echo "✅ Database seeded with basic cart data successfully"
echo ""
echo "Note: Shipping and checkout endpoints were skipped as they may not be fully implemented yet."
echo "The following predefined IDs can be used for development:"
echo "- Standard Shipping: $STANDARD_ID"
echo "- Express Shipping: $EXPRESS_ID"
echo "- Same Day Delivery: $SAME_DAY_ID"
echo "- Shipping Address 1: $ADDRESS1_ID"
echo "- Shipping Address 2: $ADDRESS2_ID" 