#!/bin/bash

BASE_URL="http://localhost:8080"

echo -e "\n=== 1. Login ==="
LOGIN_RESP=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin"}')
TOKEN=$(echo $LOGIN_RESP | grep -o '"access_token":"[^"]*' | grep -o '[^"]*$')

echo -e "\n=== 2. Create Ingredient ==="
ING_RESP=$(curl -s -X POST $BASE_URL/api/v1/ingredients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Sugar","unit_of_measure":"kg","type":"dry","stock":100}')
echo "Ingredient: $ING_RESP"
ING_ID=$(echo $ING_RESP | grep -o '"id":[^,]*' | head -1 | grep -o '[0-9]*')

echo -e "\n=== 3. Create Basic Product & Menu ==="
CAT_RESP=$(curl -s -X POST $BASE_URL/api/v1/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Beverages","description":"Drinks","is_active":true}')
echo "Category: $CAT_RESP"

PROD_RESP=$(curl -s -X POST $BASE_URL/api/v1/menu \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Coca Cola",
    "sales_price":2.5,
    "description":"Soda",
    "category_id":1,
    "ingredients":[{"ingredient_id":'$ING_ID',"quantity":0.1}]
  }')
echo "Menu: $PROD_RESP"
PROD_ID=$(echo $PROD_RESP | grep -o '"id":[^,]*' | head -1 | grep -o '[0-9]*')

echo -e "\n=== 4. Get Menu ==="
curl -s -X GET $BASE_URL/api/v1/menu \
  -H "Authorization: Bearer $TOKEN" | jq '.data[0]'

echo -e "\n=== 5. Create Order ==="
CREATE_ORDER_RESP=$(curl -s -X POST $BASE_URL/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "table_id": 1,
    "items": [
      {
        "product_id": '$PROD_ID',
        "quantity": 2,
        "unit_price": 2.5
      }
    ]
  }')
echo "Create Order: $CREATE_ORDER_RESP"
ORDER_ID=$(echo $CREATE_ORDER_RESP | head -1 | grep -o '"id":[^,]*' | head -1 | grep -o '[0-9]*')

echo -e "\n=== 6. Get Order By Table ==="
curl -s -X GET $BASE_URL/api/v1/orders?table_id=1 \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\n=== 7. Checkout Order ==="
if [ -n "$ORDER_ID" ]; then
    CHECKOUT_RESP=$(curl -s -X POST $BASE_URL/api/v1/orders/$ORDER_ID/checkout \
      -H "Authorization: Bearer $TOKEN")
    echo "Checkout Order: $CHECKOUT_RESP"
else
    echo "❌ No Order ID found to checkout."
fi

echo -e "\nAPI tests completed!"
