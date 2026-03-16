#!/bin/bash
# Quick API smoke test
# Uso: bash test_api.sh [BASE_URL]

BASE_URL="${1:-http://72.61.73.95:8080}"

echo -e "\n=== 1. Login ==="
LOGIN_RESP=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"admin"}')
echo "Login: $LOGIN_RESP"
TOKEN=$(echo $LOGIN_RESP | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])" 2>/dev/null || echo "")

if [ -z "$TOKEN" ]; then
    echo "ERROR: No se obtuvo token. Abortando."
    exit 1
fi
echo "Token obtenido OK"

echo -e "\n=== 2. Create Ingredient ==="
ING_RESP=$(curl -s -X POST $BASE_URL/ingredientes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Sugar","unit_of_measure":"kg","type":"dry","stock":100}')
echo "Ingredient: $ING_RESP"
ING_ID=$(echo $ING_RESP | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

echo -e "\n=== 3. Create Category & Menu ==="
CAT_RESP=$(curl -s -X POST $BASE_URL/categorias \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Beverages"}')
echo "Category: $CAT_RESP"

PROD_RESP=$(curl -s -X POST $BASE_URL/menu \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\":\"Coca Cola\",
    \"sales_price\":2.5,
    \"ingredients\":[{\"ingredient_id\":$ING_ID,\"quantity\":0.1}]
  }")
echo "Menu: $PROD_RESP"
PROD_ID=$(echo $PROD_RESP | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

echo -e "\n=== 4. Get Menu ==="
curl -s -X GET $BASE_URL/menu \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool 2>/dev/null || echo "(no jq/python3)"

echo -e "\n=== 5. Create Mesa ==="
MESA_NUM=$((RANDOM % 9000 + 1000))
MESA_RESP=$(curl -s -X POST $BASE_URL/mesas \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"numero\":$MESA_NUM,\"capacidad\":4}")
echo "Mesa: $MESA_RESP"
MESA_ID=$(echo $MESA_RESP | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

echo -e "\n=== 6. Create Order ==="
CREATE_ORDER_RESP=$(curl -s -X POST $BASE_URL/ordenes \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"mesa_id\":\"$MESA_ID\",\"mesero_id\":\"1\"}")
echo "Create Order: $CREATE_ORDER_RESP"
ORDER_ID=$(echo $CREATE_ORDER_RESP | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

echo -e "\n=== 7. Add Items to Order ==="
if [ "$ORDER_ID" != "0" ] && [ "$PROD_ID" != "0" ]; then
    ADD_RESP=$(curl -s -X POST $BASE_URL/ordenes/$ORDER_ID/items \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"items\":[{\"menu_item_id\":\"$PROD_ID\",\"cantidad\":2,\"notas\":\"Sin hielo\"}]}")
    echo "Add Items: $ADD_RESP"
fi

echo -e "\n=== 8. Get Order By Table ==="
curl -s -X GET "$BASE_URL/ordenes?table_id=$MESA_ID" \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool 2>/dev/null || echo "(no python3)"

echo -e "\n=== 9. Checkout Order ==="
if [ "$ORDER_ID" != "0" ]; then
    CHECKOUT_RESP=$(curl -s -X POST $BASE_URL/ordenes/$ORDER_ID/checkout \
      -H "Authorization: Bearer $TOKEN")
    echo "Checkout Order: $CHECKOUT_RESP"
fi

echo -e "\nAPI smoke tests completed!"
