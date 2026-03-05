#!/bin/bash
# =============================================================================
# POS API - Script de Pruebas Automatizadas
# Prueba todos los endpoints del sistema POS contra un servidor en ejecucion.
# Uso: bash scripts/test_api.sh [BASE_URL]
# =============================================================================

set -euo pipefail

BASE="${1:-http://72.61.73.95:8080}"
H="Content-Type: application/json"
PASS=0
FAIL=0
TOTAL=0
TOKEN=""
ORDER_ID=""
PAGO_ID=""
MESA_ID=""
ING_ID=""

# --- Utilidades ---

red()   { echo -e "\033[0;31m$1\033[0m"; }
green() { echo -e "\033[0;32m$1\033[0m"; }
bold()  { echo -e "\033[1m$1\033[0m"; }

assert_status() {
    local label="$1"
    local expected="$2"
    local actual="$3"
    local body="$4"
    TOTAL=$((TOTAL + 1))

    if [ "$actual" -eq "$expected" ]; then
        PASS=$((PASS + 1))
        green "  [PASS] $label (HTTP $actual)"
    else
        FAIL=$((FAIL + 1))
        red   "  [FAIL] $label (esperado $expected, obtuvo $actual)"
        red   "         $body"
    fi
}

assert_json_field() {
    local label="$1"
    local json="$2"
    local field="$3"
    local expected="$4"
    TOTAL=$((TOTAL + 1))

    local actual
    actual=$(echo "$json" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d${field})" 2>/dev/null || echo "__ERROR__")

    if [ "$actual" = "$expected" ]; then
        PASS=$((PASS + 1))
        green "  [PASS] $label ($field = $expected)"
    else
        FAIL=$((FAIL + 1))
        red   "  [FAIL] $label ($field esperado '$expected', obtuvo '$actual')"
    fi
}

api() {
    local method="$1"
    local path="$2"
    shift 2
    local args=(-s -w "\n%{http_code}" -X "$method" "${BASE}${path}" -H "$H")
    if [ -n "$TOKEN" ]; then
        args+=(-H "Authorization: Bearer $TOKEN")
    fi
    args+=("$@")
    curl "${args[@]}" 2>/dev/null
}

parse_response() {
    local raw="$1"
    BODY=$(echo "$raw" | sed '$d')
    STATUS=$(echo "$raw" | tail -1)
}

# =============================================================================
bold "============================================"
bold "  POS API Test Suite"
bold "  Servidor: $BASE"
bold "============================================"
echo ""

# --- 1. HEALTH ---
bold "--- 1. Salud y Monitoreo ---"

RAW=$(api GET /health)
parse_response "$RAW"
assert_status "GET /health" 200 "$STATUS" "$BODY"
assert_json_field "GET /health body" "$BODY" "['status']" "ok"

RAW=$(api GET /ping)
parse_response "$RAW"
assert_status "GET /ping" 200 "$STATUS" "$BODY"

echo ""

# --- 2. AUTH ---
bold "--- 2. Autenticacion ---"

# Login con credenciales invalidas
RAW=$(api POST /auth/login -d '{"email":"bad@test.com","password":"wrong"}')
parse_response "$RAW"
assert_status "POST /auth/login (invalido)" 401 "$STATUS" "$BODY"

# Login correcto
RAW=$(api POST /auth/login -d '{"email":"admin@test.com","password":"admin"}')
parse_response "$RAW"
assert_status "POST /auth/login (valido)" 200 "$STATUS" "$BODY"
assert_json_field "Login success field" "$BODY" "['success']" "True"
TOKEN=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])" 2>/dev/null || echo "")

if [ -z "$TOKEN" ]; then
    red "ERROR CRITICO: No se obtuvo token. Abortando."
    exit 1
fi
green "  Token obtenido correctamente"

# Auth/me
RAW=$(api GET /auth/me)
parse_response "$RAW"
assert_status "GET /auth/me" 200 "$STATUS" "$BODY"
assert_json_field "Me success" "$BODY" "['success']" "True"
assert_json_field "Me tiene rol" "$BODY" "['data']['rol']" "PROPIETARIO"

# Acceso sin token
OLD_TOKEN="$TOKEN"
TOKEN=""
RAW=$(api GET /auth/me)
parse_response "$RAW"
assert_status "GET /auth/me (sin token)" 401 "$STATUS" "$BODY"
TOKEN="$OLD_TOKEN"

echo ""

# --- 3. MESAS ---
bold "--- 3. Mesas ---"

RAW=$(api GET /mesas)
parse_response "$RAW"
assert_status "GET /mesas" 200 "$STATUS" "$BODY"
assert_json_field "Mesas success" "$BODY" "['success']" "True"

# Crear mesa con numero unico
MESA_NUM=$((RANDOM % 9000 + 1000))
RAW=$(api POST /mesas -d "{\"numero\":$MESA_NUM,\"capacidad\":4}")
parse_response "$RAW"
assert_status "POST /mesas (crear)" 201 "$STATUS" "$BODY"
MESA_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

if [ "$MESA_ID" != "0" ]; then
    # Get by ID
    RAW=$(api GET /mesas/$MESA_ID)
    parse_response "$RAW"
    assert_status "GET /mesas/$MESA_ID" 200 "$STATUS" "$BODY"

    # Update estado
    RAW=$(api PATCH /mesas/$MESA_ID/estado -d '{"estado":"occupied"}')
    parse_response "$RAW"
    assert_status "PATCH /mesas/$MESA_ID/estado" 200 "$STATUS" "$BODY"



    # Delete
    RAW=$(api DELETE /mesas/$MESA_ID)
    parse_response "$RAW"
    assert_status "DELETE /mesas/$MESA_ID" 200 "$STATUS" "$BODY"
fi

echo ""

# --- 4. USUARIOS ---
bold "--- 4. Usuarios ---"

RAW=$(api GET /usuarios)
parse_response "$RAW"
assert_status "GET /usuarios" 200 "$STATUS" "$BODY"
assert_json_field "Usuarios success" "$BODY" "['success']" "True"

# Get usuario por ID
RAW=$(api GET /usuarios/1)
parse_response "$RAW"
assert_status "GET /usuarios/1" 200 "$STATUS" "$BODY"
assert_json_field "Usuario tiene nombre" "$BODY" "['data']['nombre']" "Admin"

# Get usuario inexistente
RAW=$(api GET /usuarios/99999)
parse_response "$RAW"
assert_status "GET /usuarios/99999 (no existe)" 404 "$STATUS" "$BODY"

echo ""

# --- 5. INGREDIENTES ---
bold "--- 5. Ingredientes ---"

RAW=$(api GET /ingredientes)
parse_response "$RAW"
assert_status "GET /ingredientes" 200 "$STATUS" "$BODY"

# Crear ingrediente
RAW=$(api POST /ingredientes -d '{"name":"TestIngredient","unit_of_measure":"kg","type":"dry","stock":100}')
parse_response "$RAW"
assert_status "POST /ingredientes" 201 "$STATUS" "$BODY"
ING_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

if [ "$ING_ID" != "0" ]; then
    # Get by ID
    RAW=$(api GET /ingredientes/$ING_ID)
    parse_response "$RAW"
    assert_status "GET /ingredientes/$ING_ID" 200 "$STATUS" "$BODY"

    # Update stock
    RAW=$(api PATCH /ingredientes/$ING_ID/stock -d '{"cantidad":50,"tipo_movimiento":"entrada","motivo":"Compra"}')
    parse_response "$RAW"
    assert_status "PATCH /ingredientes/$ING_ID/stock (entrada)" 200 "$STATUS" "$BODY"

    # Update stock salida
    RAW=$(api PATCH /ingredientes/$ING_ID/stock -d '{"cantidad":10,"tipo_movimiento":"salida","motivo":"Consumo"}')
    parse_response "$RAW"
    assert_status "PATCH /ingredientes/$ING_ID/stock (salida)" 200 "$STATUS" "$BODY"

    # Stock insuficiente
    RAW=$(api PATCH /ingredientes/$ING_ID/stock -d '{"cantidad":99999,"tipo_movimiento":"salida","motivo":"Test"}')
    parse_response "$RAW"
    assert_status "PATCH stock insuficiente" 400 "$STATUS" "$BODY"

    # Delete
    RAW=$(api DELETE /ingredientes/$ING_ID)
    parse_response "$RAW"
    assert_status "DELETE /ingredientes/$ING_ID" 204 "$STATUS" "$BODY"
fi

echo ""

# --- 6. CATEGORIAS ---
bold "--- 6. Categorias ---"

RAW=$(api GET /categorias)
parse_response "$RAW"
assert_status "GET /categorias" 200 "$STATUS" "$BODY"

RAW=$(api POST /categorias -d '{"name":"TestCategory"}')
parse_response "$RAW"
assert_status "POST /categorias" 201 "$STATUS" "$BODY"

echo ""

# --- 7. MENU ---
bold "--- 7. Menu ---"

RAW=$(api GET /menu)
parse_response "$RAW"
assert_status "GET /menu" 200 "$STATUS" "$BODY"
assert_json_field "Menu success" "$BODY" "['success']" "True"

echo ""

# --- 8. ORDENES ---
bold "--- 8. Ordenes ---"

# Crear orden
RAW=$(api POST /ordenes -d '{"mesa_id":"1","mesero_id":"1"}')
parse_response "$RAW"
assert_status "POST /ordenes" 201 "$STATUS" "$BODY"
assert_json_field "Orden success" "$BODY" "['success']" "True"
ORDER_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

if [ "$ORDER_ID" != "0" ]; then
    # Get orden
    RAW=$(api GET /ordenes/$ORDER_ID)
    parse_response "$RAW"
    assert_status "GET /ordenes/$ORDER_ID" 200 "$STATUS" "$BODY"
    assert_json_field "Orden estado" "$BODY" "['data']['estado']" "abierta"

    # Enviar a cocina
    RAW=$(api POST /ordenes/$ORDER_ID/enviar-cocina)
    parse_response "$RAW"
    assert_status "POST /ordenes/$ORDER_ID/enviar-cocina" 200 "$STATUS" "$BODY"

    # Verificar estado cambiado
    RAW=$(api GET /ordenes/$ORDER_ID)
    parse_response "$RAW"
    assert_json_field "Estado despues de enviar" "$BODY" "['data']['estado']" "enviada"

    # Dividir cuenta
    RAW=$(api POST /ordenes/$ORDER_ID/dividir -d '{"tipo_division":"partes_iguales","numero_partes":2}')
    parse_response "$RAW"
    assert_status "POST /ordenes/$ORDER_ID/dividir" 200 "$STATUS" "$BODY"
    assert_json_field "Division success" "$BODY" "['success']" "True"

    # Dividir por monto
    RAW=$(api POST /ordenes/$ORDER_ID/dividir -d '{"tipo_division":"por_monto","divisiones":[{"monto":30000},{"monto":20000}]}')
    parse_response "$RAW"
    assert_status "POST dividir por_monto" 200 "$STATUS" "$BODY"

    # Dividir tipo invalido
    RAW=$(api POST /ordenes/$ORDER_ID/dividir -d '{"tipo_division":"invalido"}')
    parse_response "$RAW"
    assert_status "POST dividir tipo invalido" 400 "$STATUS" "$BODY"

    # Update status
    RAW=$(api PATCH /ordenes/$ORDER_ID/status -d '{"status_id":4}')
    parse_response "$RAW"
    assert_status "PATCH /ordenes/$ORDER_ID/status" 200 "$STATUS" "$BODY"
fi

# Orden inexistente
RAW=$(api GET /ordenes/999999)
parse_response "$RAW"
assert_status "GET /ordenes/999999 (no existe)" 404 "$STATUS" "$BODY"

echo ""

# --- 9. PAGOS ---
bold "--- 9. Pagos ---"

if [ "$ORDER_ID" != "0" ]; then
    # Procesar pago
    RAW=$(api POST /pagos -d "{\"orden_id\":\"$ORDER_ID\",\"metodo_pago\":\"tarjeta\",\"monto\":50000,\"propina\":5000,\"detalles_pago\":{\"referencia_tarjeta\":\"REF-TEST-001\"}}")
    parse_response "$RAW"
    assert_status "POST /pagos" 200 "$STATUS" "$BODY"
    assert_json_field "Pago estado" "$BODY" "['data']['estado']" "aprobado"
    PAGO_ID=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "0")

    if [ "$PAGO_ID" != "0" ]; then
        # Factura
        RAW=$(api GET /pagos/$PAGO_ID/factura)
        parse_response "$RAW"
        assert_status "GET /pagos/$PAGO_ID/factura" 200 "$STATUS" "$BODY"
        assert_json_field "Factura success" "$BODY" "['success']" "True"

        # Verificar que la factura tiene los campos esperados
        TOTAL_CHECK=$((TOTAL + 1))
        HAS_FACTURA_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print('FACT-' in d['data']['factura_id'])" 2>/dev/null || echo "False")
        if [ "$HAS_FACTURA_ID" = "True" ]; then
            PASS=$((PASS + 1))
            TOTAL=$TOTAL_CHECK
            green "  [PASS] Factura tiene factura_id con formato FACT-"
        else
            FAIL=$((FAIL + 1))
            TOTAL=$TOTAL_CHECK
            red   "  [FAIL] Factura no tiene factura_id con formato FACT-"
        fi
    fi

    # Metodo de pago invalido
    RAW=$(api POST /pagos -d "{\"orden_id\":\"$ORDER_ID\",\"metodo_pago\":\"bitcoin\",\"monto\":1000}")
    parse_response "$RAW"
    assert_status "POST /pagos metodo invalido" 400 "$STATUS" "$BODY"
fi

# Factura inexistente
RAW=$(api GET /pagos/999999/factura)
parse_response "$RAW"
assert_status "GET /pagos/999999/factura (no existe)" 404 "$STATUS" "$BODY"

echo ""

# --- 10. REPORTES ---
bold "--- 10. Reportes ---"

RAW=$(api GET "/reportes/ventas?fecha_inicio=2026-01-01&fecha_fin=2026-12-31")
parse_response "$RAW"
assert_status "GET /reportes/ventas" 200 "$STATUS" "$BODY"
assert_json_field "Ventas success" "$BODY" "['success']" "True"

# Ventas sin parametros
RAW=$(api GET /reportes/ventas)
parse_response "$RAW"
assert_status "GET /reportes/ventas (sin fechas)" 400 "$STATUS" "$BODY"

RAW=$(api GET /reportes/inventario)
parse_response "$RAW"
assert_status "GET /reportes/inventario" 200 "$STATUS" "$BODY"
assert_json_field "Inventario success" "$BODY" "['success']" "True"

RAW=$(api GET "/reportes/propinas?fecha_inicio=2026-01-01&fecha_fin=2026-12-31")
parse_response "$RAW"
assert_status "GET /reportes/propinas" 200 "$STATUS" "$BODY"
assert_json_field "Propinas success" "$BODY" "['success']" "True"

# Propinas sin parametros
RAW=$(api GET /reportes/propinas)
parse_response "$RAW"
assert_status "GET /reportes/propinas (sin fechas)" 400 "$STATUS" "$BODY"

echo ""

# --- 11. LOGOUT ---
bold "--- 11. Logout ---"

RAW=$(api POST /auth/logout)
parse_response "$RAW"
assert_status "POST /auth/logout" 200 "$STATUS" "$BODY"
assert_json_field "Logout message" "$BODY" "['success']" "True"

echo ""

# =============================================================================
bold "============================================"
bold "  RESULTADOS"
bold "============================================"
echo ""
echo "  Total:    $TOTAL"
green "  Pasaron:  $PASS"
if [ "$FAIL" -gt 0 ]; then
    red   "  Fallaron: $FAIL"
else
    green "  Fallaron: 0"
fi
echo ""

if [ "$FAIL" -gt 0 ]; then
    red "Hay $FAIL prueba(s) fallida(s)."
    exit 1
else
    green "Todas las pruebas pasaron."
    exit 0
fi
