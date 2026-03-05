# POS Backend - API Reference

Esta documentación contiene todos los endpoints expuestos por el sistema POS para su consumo en el Frontend.
Todos los endpoints requieren un **Access Token** enviado a través del header `Authorization` a excepción de los marcados como "Público".

**Headers Globales Requeridos:**
```http
Content-Type: application/json
Authorization: Bearer <vuestro_token_jwt>
```

**Esquema de Roles soportados:**
- **Propietario** (Role ID: 1)
- **Cajero** (Role ID: 2)
- **Mesero** (Role ID: 3)

---

## 1. Salud y Monitoreo (Público)
Endpoints utilizados para monitorizar que el backend está corriendo.

- **GET `/health`**: Retorna el estado en JSON `{"status": "ok"}`
- **GET `/ping`**: Retorna string `"server say: pong"`

---

## 2. Autenticación (Público)

### POST `/auth/login`
Inicia sesión devolviendo el JWT necesario.

**Request:**
```json
{
  "email": "admin@test.com",
  "password": "admin"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "def456..."
}
```

---

## 3. Mesas (Tables)
Base URL: `/api/v1/tables`

### GET `/`
Obtiene la lista de todas las mesas.
*Roles:* `Mesero`, `Cajero`, `Propietario`

**Response:**
```json
[
  {
    "id": 1,
    "table_number": 1,
    "capacity": 4,
    "status": "LIBRE",
    "arrival_time": null,
    "created_at": "2026-03-02T21:54:46Z"
  }
]
```

### POST `/`
Crea una nueva mesa.
*Roles:* `Cajero`, `Propietario`

**Request:**
```json
{
  "table_number": 2,
  "capacity": 2,
  "status": "LIBRE"
}
```

### PATCH `/:id`
Actualiza el estado o detalles de la mesa (Ej: De `LIBRE` a `OCUPADO`).
*Roles:* `Cajero`, `Propietario`

**Request:**
```json
{
  "status": "OCUPADO"
}
```

### POST `/:id/assign`
Asigna un mesero a una mesa concreta.
*Roles:* `Cajero`, `Propietario`

**Request:**
```json
{
  "waitress_id": 2
}
```

---

## 4. Categorías de Productos
Base URL: `/api/v1/categories`

### GET `/`
Listar categorías.
*Roles:* `Propietario`

### POST `/`
Crear categoría de menú (Ej. Bebidas, Entradas).
*Roles:* `Propietario`

**Request:**
```json
{
  "name": "Bebidas",
  "description": "Bebidas frías y calientes"
}
```

---

## 5. Ingredientes / Inventario
Base URL: `/api/v1/ingredients`

### GET `/` | GET `/:id` | DELETE `/:id`
Manejo CRUD de los insumos.
*Roles:* `Propietario`

### POST `/`
Agregar materia prima.
*Roles:* `Propietario`

**Request:**
```json
{
  "name": "Azúcar",
  "unit_of_measure": "kg",
  "type": "dry",
  "stock": 100
}
```

---

## 6. Menú (Productos Finales)
Base URL: `/api/v1/menu`

### GET `/`
Listado del menú oficial para ser mostrado en la App del Mesero.
*Roles:* `Mesero`, `Cajero`, `Propietario`

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Coca Cola",
      "sales_price": 2.5,
      "is_active": true,
      "created_at": "..."
    }
  ],
  "page": 1,
  "page_size": 20
}
```

### POST `/`
Registrar elemento de menú compuesto de sus ingredientes necesarios (Receta).
*Roles:* `Propietario`

**Request:**
```json
{
  "name": "Coca Cola",
  "sales_price": 2.5,
  "description": "Lata de soda fría",
  "category_id": 1,
  "ingredients": [
    {
      "ingredient_id": 1,
      "quantity": 0.1
    }
  ]
}
```

### PATCH `/:id`
Actualizar precio/estado de un platillo.
*Roles:* `Propietario`

---

## 7. Órdenes (Operaciones de Venta)
Base URL: `/api/v1/orders`

### GET `/?table_id=1`
Lista las órdenes actuales, opcionalmente filtrando por mesa.
*Roles:* `Mesero`, `Cajero`, `Propietario`

### POST `/`
Nueva orden de cliente (Command originado en POS).
*Roles:* `Mesero`, `Cajero`, `Propietario`

**Request:**
```json
{
  "table_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "unit_price": 2.5
    }
  ]
}
```

### PATCH `/:id/status`
Mover la orden por sus estados (PENDIENDO -> COCINANDO -> LISTA)
*Roles:* `Mesero`, `Cajero`, `Propietario`

**Request:**
```json
{
  "status_id": 2
}
```

### POST `/:id/checkout`
Acción final de Caja cobrando el dinero de la orden (Pasando la orden a `PAID`).
*Roles:* `Cajero`, `Propietario`

**Response:**
```json
{
  "status": "PAID"
}
```

---
