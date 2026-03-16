# POS API Documentation

## Base URL
`http://localhost:8080`

## Response Format

All successful responses follow the wrapper:
```json
{ "success": true, "data": { ... } }
```

Message-only responses:
```json
{ "success": true, "message": "..." }
```

Error responses:
```json
{ "success": false, "error": "...", "code": "BAD_REQUEST | UNAUTHORIZED | NOT_FOUND | CONFLICT | INTERNAL_ERROR" }
```

## Authentication
All endpoints (except `/health`, `/ping`, `/auth/login`, `/auth/register-owner`) require `Authorization: Bearer <token>`.

---

### Health & Monitoring

| Method | Route | Description |
|---|---|---|
| GET | `/health` | Returns `{"status": "ok"}` |
| GET | `/ping` | Returns `"server say: pong"` |

---

### Auth

| Method | Route | Auth | Description |
|---|---|---|---|
| POST | `/auth/login` | Public | Login, returns JWT |
| POST | `/auth/register-owner` | Public | Register owner + venue |
| POST | `/auth/refresh` | Public | Refresh access token |
| POST | `/auth/register` | Propietario, Cajero | Register user (role-restricted) |
| GET | `/auth/me` | Any | Current user info |
| POST | `/auth/logout` | Any | Revoke refresh token |

#### POST `/auth/login`
```json
// Request
{
  "email": "admin@test.com",       // required, valid email
  "password": "admin"              // required
}
// Response 200
{
  "success": true,
  "data": {
    "token": "eyJhbGciOi...",
    "usuario": {
      "id": 1,
      "nombre": "Admin",
      "email": "admin@test.com",
      "rol": "PROPIETARIO",        // PROPIETARIO | CAJERO | MESERO
      "activo": true,
      "fecha_creacion": "2024-01-01T00:00:00Z",
      "telefono": "+57300...",     // optional
      "ultimo_acceso": "..."       // optional
    },
    "expires_in": 900
  }
}
```

#### POST `/auth/register-owner`
```json
// Request
{
  "nombre": "Juan",                // required
  "email": "juan@test.com",        // required, valid email
  "password": "12345678",          // required, min 8 chars
  "nombre_sede": "Mi Restaurante", // required
  "direccion": "Calle 1 #2-3",    // optional
  "telefono": "+57300..."          // optional
}
// Response 201
{
  "success": true,
  "data": {
    "token": "eyJhbGciOi...",
    "usuario": {
      "id": 1,
      "nombre": "Juan",
      "email": "juan@test.com",
      "rol": "PROPIETARIO"
    },
    "message": "Propietario y sede creados exitosamente",
    "expires_in": 900
  }
}
```

#### POST `/auth/refresh`
```json
// Request (via body, si no se envia como cookie)
{
  "refresh_token": "eyJhbGciOi..."  // optional si se envia como cookie
}
// Response 200
{
  "success": true,
  "data": {
    "token": "eyJhbGciOi...",      // nuevo access token (15min)
    "expires_in": 900,
    "usuario": {
      "id": 1,
      "nombre": "Juan",
      "email": "juan@test.com",
      "rol": "PROPIETARIO"
    }
  }
}
```
> El refresh token se puede enviar como cookie HTTP-only (`refresh_token`) o en el body JSON.
> Al refrescar, el refresh token anterior se revoca y se genera uno nuevo (rotacion de tokens).
> El access token dura 15 minutos, el refresh token dura 24 horas.

#### POST `/auth/register`
```json
// Request
{
  "nombre": "Carlos",              // required
  "email": "carlos@test.com",      // required, valid email
  "password": "12345678",          // required, min 8 chars
  "rol": "CAJERO",                 // required: PROPIETARIO | CAJERO | MESERO
  "telefono": "+57300..."          // optional
}
// Response 201
{
  "success": true,
  "message": "Usuario registrado exitosamente",
  "data": {
    "id": 2,
    "nombre": "Carlos",
    "email": "carlos@test.com",
    "rol": "CAJERO",
    "activo": true,
    "fecha_creacion": "2024-01-01T00:00:00Z"
  }
}
```

#### GET `/auth/me`
```json
// Response 200
{
  "success": true,
  "data": {
    "id": 1,
    "nombre": "Admin",
    "email": "admin@test.com",
    "rol": "PROPIETARIO",
    "activo": true,
    "fecha_creacion": "2024-01-01T00:00:00Z",
    "telefono": "+57300...",       // optional
    "ultimo_acceso": "..."         // optional
  }
}
```

#### POST `/auth/logout`
```json
// Response 200
{ "success": true, "message": "Sesion cerrada exitosamente" }
```

---

### Propietario (Owner)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/propietario` | Propietario | Get authenticated owner info |
| PATCH | `/propietario` | Propietario | Update owner data |

#### GET `/propietario`
```json
// Response 200
{
  "success": true,
  "data": {
    "id": 1,
    "nombre": "Juan",
    "email": "juan@test.com",
    "activo": true,
    "creado_en": "2024-01-01T00:00:00Z"
  }
}
```

#### PATCH `/propietario`
```json
// Request
{
  "nombre": "Juan actualizado",    // optional
  "email": "nuevo@test.com"        // optional
}
// Response 200
{
  "success": true,
  "data": {
    "id": 1,
    "nombre": "Juan actualizado",
    "email": "nuevo@test.com",
    "activo": true,
    "creado_en": "2024-01-01T00:00:00Z"
  }
}
```

---

### Sedes (Venues)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/sedes` | Propietario | Create venue |
| GET | `/sedes` | Propietario | List owner's venues |
| GET | `/sedes/:id` | Propietario | Get venue by ID |
| PATCH | `/sedes/:id` | Propietario | Update venue |

#### POST `/sedes`
```json
// Request
{
  "nombre": "Sede Norte",          // required
  "direccion": "Calle 100 #15",   // optional
  "telefono": "+57300..."          // optional
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "nombre": "Sede Norte",
    "direccion": "Calle 100 #15",
    "telefono": "+57300...",
    "activo": true,
    "creado_en": "2024-01-01T00:00:00Z"
  }
}
```

#### GET `/sedes`
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nombre": "Sede Norte",
      "direccion": "Calle 100 #15",
      "telefono": "+57300...",
      "activo": true,
      "creado_en": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### GET `/sedes/:id`
```json
// Response 200 (same object as above)
```

#### PATCH `/sedes/:id`
```json
// Request
{
  "nombre": "Sede Sur",            // optional
  "direccion": "Carrera 7 #40",   // optional
  "telefono": "+57301..."          // optional
}
// Response 200
{
  "success": true,
  "data": {
    "id": 1,
    "nombre": "Sede Sur",
    "direccion": "Carrera 7 #40",
    "telefono": "+57301...",
    "activo": true
  }
}
```

---

### Terminales POS (POS Terminals)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/terminales` | Propietario | Create terminal |
| GET | `/terminales` | Propietario | List venue terminals |
| GET | `/terminales/:id` | Propietario | Get terminal by ID |
| PATCH | `/terminales/:id` | Propietario | Update terminal |

#### POST `/terminales`
```json
// Request
{
  "nombre": "Caja 1"              // required
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "nombre": "Caja 1",
    "venue_id": 1,
    "activo": true
  }
}
```

#### GET `/terminales`
```json
// Response 200
{
  "success": true,
  "data": [
    { "id": 1, "nombre": "Caja 1", "venue_id": 1, "activo": true }
  ]
}
```

#### GET `/terminales/:id`
```json
// Response 200 (same terminal object)
```

#### PATCH `/terminales/:id`
```json
// Request
{
  "nombre": "Caja Principal",     // optional
  "activo": false                  // optional
}
// Response 200 (updated terminal object)
```

---

### Mesas (Tables)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/mesas` | All | List tables |
| GET | `/mesas/:id` | All | Get table |
| POST | `/mesas` | Cajero, Propietario | Create table |
| PATCH | `/mesas/:id/estado` | Mesero, Cajero, Propietario | Update status |
| DELETE | `/mesas/:id` | Cajero, Propietario | Delete table |
| POST | `/mesas/:id/asignar` | Cajero, Propietario | Assign waiter |
| GET | `/mesas/:id/asignaciones` | Cajero, Propietario | Assignment history |

**Estados validos:** `LIBRE` (default), `OCUPADA`, `RESERVADA`

#### POST `/mesas`
```json
// Request
{
  "numero": 5,                     // required, > 0
  "capacidad": 4,                  // required, > 0
  "status": "LIBRE"                // optional, default: LIBRE
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "number": 5,
    "capacity": 4,
    "state": "LIBRE",
    "arrival_time": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET `/mesas` / GET `/mesas/:id`
```json
// Response 200
{
  "success": true,
  "data": {
    "id": 1,
    "number": 5,
    "capacity": 4,
    "state": "LIBRE",
    "arrival_time": null,          // optional
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### PATCH `/mesas/:id/estado`
```json
// Request
{
  "estado": "OCUPADA"              // required: LIBRE | OCUPADA | RESERVADA
}
// Response 200 (updated table object)
```

#### POST `/mesas/:id/asignar`
```json
// Request
{
  "user_id": 5                     // required, > 0 (ID del mesero)
}
// Response 200
{
  "success": true,
  "data": {
    "id": 1,
    "table_id": 3,
    "user_id": 5,
    "asignado_en": "2024-01-01T12:00:00Z"
  }
}
```
> Si la mesa ya tenia un mesero asignado, se cierra la asignacion anterior automaticamente.

#### GET `/mesas/:id/asignaciones`
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "id": 1,
      "table_id": 3,
      "user_id": 5,
      "nombre_mesero": "Pedro",
      "asignado_en": "2024-01-01T12:00:00Z",
      "desasignado_en": "2024-01-01T16:00:00Z"
    },
    {
      "id": 2,
      "table_id": 3,
      "user_id": 7,
      "nombre_mesero": "Carlos",
      "asignado_en": "2024-01-01T16:00:00Z",
      "desasignado_en": null
    }
  ]
}
```

---

### Usuarios (Users)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/usuarios` | Propietario | List users |
| GET | `/usuarios/:id` | Propietario | Get user |
| PATCH | `/usuarios/:id` | Propietario | Update user |
| DELETE | `/usuarios/:id` | Propietario | Deactivate user |
| POST | `/usuarios/mesero` | Propietario, Cajero | Quick-register waiter |

#### GET `/usuarios`
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "id": 2,
      "nombre": "Carlos",
      "email": "carlos@test.com",
      "rol": "CAJERO",
      "activo": true,
      "fecha_creacion": "2024-01-01T00:00:00Z",
      "telefono": "+57300...",     // optional
      "ultimo_acceso": "..."       // optional
    }
  ]
}
```

#### GET `/usuarios/:id`
```json
// Response 200 (same user object)
```

#### PATCH `/usuarios/:id`
```json
// Request (all fields optional)
{
  "nombre": "Carlos Gomez",
  "email": "nuevo@test.com",
  "rol_id": 2,                    // 1=Propietario, 2=Cajero, 3=Mesero
  "activo": true,
  "telefono": "+57301..."
}
// Response 200 (updated user object)
```

#### DELETE `/usuarios/:id`
```json
// Response 200
{ "success": true, "message": "Usuario desactivado exitosamente" }
```

#### POST `/usuarios/mesero`
```json
// Request
{
  "nombre": "Pedro",               // required
  "email": "pedro@test.com"        // required, valid email
}
// Response 201
{
  "success": true,
  "data": {
    "message": "Mesero creado exitosamente. Guarde las credenciales, no se mostraran de nuevo.",
    "credenciales": {
      "email": "pedro@test.com",
      "password": "auto-generated-password"
    },
    "usuario": {
      "id": 5,
      "nombre": "Pedro",
      "email": "pedro@test.com",
      "rol": "MESERO",
      "activo": true
    }
  }
}
```

---

### Ingredientes (Ingredients)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/ingredientes` | Propietario | List ingredients |
| POST | `/ingredientes` | Propietario | Create ingredient |
| GET | `/ingredientes/:id` | Propietario | Get ingredient |
| PUT | `/ingredientes/:id` | Propietario | Update ingredient |
| PATCH | `/ingredientes/:id/stock` | Propietario | Stock movement |
| DELETE | `/ingredientes/:id` | Propietario | Delete ingredient |
| GET | `/ingredientes/report` | Propietario | Inventory report |

#### POST `/ingredientes`
```json
// Request
{
  "name": "Sugar",                 // required
  "unit_of_measure": "kg",        // required
  "type": "dry",                   // required
  "stock": 50                      // optional, default 0
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Sugar",
    "unit_of_measure": "kg",
    "type": "dry",
    "stock": 50,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET `/ingredientes` / GET `/ingredientes/:id`
```json
// Response 200 (same ingredient object or array)
```

#### PUT `/ingredientes/:id`
```json
// Request (all fields optional)
{
  "name": "Brown Sugar",
  "unit_of_measure": "kg",
  "type": "dry",
  "stock": 100
}
// Response 200 (updated ingredient object)
```

#### PATCH `/ingredientes/:id/stock`
```json
// Request
{
  "cantidad": 50,                  // required
  "tipo_movimiento": "entrada",   // required: "entrada" | "salida"
  "motivo": "Weekly purchase"      // optional
}
// Response 200
{ "success": true, "message": "..." }
```

---

### Categorias (Categories)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/categorias` | Propietario | Create category |
| GET | `/categorias` | Propietario | List categories |

#### POST `/categorias`
```json
// Request
{
  "name": "Bebidas"                // required
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Bebidas",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": null             // optional
  }
}
```

#### GET `/categorias`
```json
// Response 200 (array of category objects)
```

---

### Products (Internal)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/products` | Propietario | Create product |
| GET | `/products` | Propietario | List products |
| POST | `/products/:id/ingredients` | Propietario | Add ingredient to recipe |
| GET | `/products/:id/ingredients` | Propietario | Get product recipe |

#### POST `/products`
```json
// Request
{
  "name": "Hamburguesa",           // required
  "sales_price": 12.00,           // required, >= 0
  "is_active": true                // optional
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Hamburguesa",
    "sales_price": 12.00,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": null
  }
}
```

#### POST `/products/:id/ingredients`
```json
// Request
{
  "ingredient_id": 1,             // required, > 0
  "quantity": 0.2                  // required, > 0
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "product_id": 1,
    "ingredient_id": 1,
    "ingredient_name": "Sugar",    // optional
    "unit_of_measure": "kg",       // optional
    "quantity": 0.2
  }
}
```

#### GET `/products/:id/ingredients`
```json
// Response 200 (array of recipe item objects)
```

---

### Menu

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/menu` | All | List menu items |
| POST | `/menu` | Propietario | Create menu item with recipe |
| PATCH | `/menu/:id` | Propietario | Update menu item |

#### POST `/menu`
```json
// Request
{
  "name": "Hamburguesa",           // required
  "sales_price": 12.00,           // required, >= 0
  "ingredients": [                 // required
    { "ingredient_id": 1, "quantity": 0.2 }
  ]
}
// Response 201
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Hamburguesa",
    "sales_price": 12.00,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### PATCH `/menu/:id`
```json
// Request (all fields optional)
{
  "name": "Hamburguesa Doble",
  "sales_price": 18.00
}
// Response 200 (updated product object)
```

---

### Ordenes (Orders)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/ordenes` | All | Create order |
| GET | `/ordenes?table_id=1` | All | List by table |
| GET | `/ordenes/:id` | All | Get order |
| POST | `/ordenes/:id/items` | All | Add items |
| DELETE | `/ordenes/:id/items/:item_id` | All | Cancel item |
| POST | `/ordenes/:id/enviar-cocina` | All | Send to kitchen |
| PATCH | `/ordenes/:id/status` | All | Update status |
| POST | `/ordenes/:id/dividir` | Cajero, Propietario | Split bill |
| POST | `/ordenes/:id/checkout` | Cajero, Propietario | Checkout |

#### Order Response Object
```json
{
  "id": 1,
  "mesa_id": 1,
  "mesero_id": 3,
  "estado": "abierta",            // abierta | enviada | en_preparacion | lista | pagada | cancelada
  "items": [
    {
      "id": 1,
      "menu_item_id": 1,
      "cantidad": 2,
      "precio_unitario": 12.00,
      "notas": "Sin cebolla"
    }
  ],
  "subtotal": 24.00,
  "impuestos": 4.56,
  "total": 28.56,
  "fecha_creacion": "2024-01-01T00:00:00Z"
}
```

#### POST `/ordenes`
```json
// Request
{
  "mesa_id": "1",                  // optional
  "mesero_id": "3"                 // optional
}
// Response 201 (order object wrapped in success)
```

#### POST `/ordenes/:id/items`
```json
// Request
{
  "items": [                       // required
    {
      "menu_item_id": "1",        // required
      "cantidad": 2,              // required, > 0
      "notas": "Sin cebolla"      // optional
    }
  ]
}
// Response 200 (updated order object)
```

#### DELETE `/ordenes/:id/items/:item_id`
```json
// Response 200
{ "success": true, "message": "Item cancelado exitosamente" }
```

#### POST `/ordenes/:id/enviar-cocina`
```json
// Response 200
{ "success": true, "message": "Orden enviada a cocina y bar" }
```

#### PATCH `/ordenes/:id/status`
```json
// Request
{
  "status_id": 3                   // required: 1=abierta, 2=enviada, 3=en_preparacion, 4=lista, 5=pagada, 6=cancelada
}
// Response 200
{ "success": true, "message": "Estado de orden actualizado" }
```

#### POST `/ordenes/:id/dividir`
```json
// Request (option A: equal split)
{
  "tipo_division": "partes_iguales",  // required
  "numero_partes": 3
}
// Request (option B: by amount)
{
  "tipo_division": "por_monto",
  "divisiones": [
    { "items": [], "monto": 15000 },
    { "items": [], "monto": 10000 }
  ]
}
// Request (option C: by item)
{
  "tipo_division": "por_item",
  "divisiones": [
    { "items": ["1", "2"], "monto": 0 }
  ]
}
// Response 200
{
  "success": true,
  "data": [
    {
      "division_id": "div_1",
      "subtotal": 8000,
      "impuestos": 1520,
      "total": 9520
    }
  ]
}
```

#### POST `/ordenes/:id/checkout`
```json
// Response 200
{ "success": true, "message": "Pago procesado" }
```

---

### Pagos (Payments)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/pagos` | Cajero, Propietario | Process payment |
| GET | `/pagos/:id/factura` | Cajero, Propietario | Generate invoice |

#### POST `/pagos`
```json
// Request
{
  "orden_id": "1",                 // required
  "division_id": "div_1",         // optional
  "metodo_pago": "efectivo",      // required: efectivo | tarjeta | multiple
  "monto": 41650,                 // required
  "propina": 4165,                // optional
  "detalles_pago": {              // optional (for "multiple")
    "efectivo": 20000,
    "tarjeta": 21650,
    "referencia_tarjeta": "TXN123"
  }
}
// Response 200
{
  "success": true,
  "data": {
    "id": 1,
    "orden_id": 1,
    "metodo_pago": "efectivo",
    "monto": 41650,
    "propina": 4165,
    "total": 45815,
    "estado": "completado",
    "referencia": "",
    "fecha": "2024-01-01T00:00:00Z"
  }
}
```

#### GET `/pagos/:id/factura`
```json
// Response 200
{
  "success": true,
  "data": { ... }                  // Invoice object (varies)
}
```

---

### Reportes (Reports)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/reportes/ventas` | Propietario | Sales report |
| GET | `/reportes/inventario` | Propietario | Inventory report |
| GET | `/reportes/propinas` | Propietario | Tips report |

**Query params:** `fecha_inicio` (date), `fecha_fin` (date), `tipo` (por_dia/por_item/por_hora)

```json
// Response 200
{
  "success": true,
  "data": { ... }                  // Report data (varies by type)
}
```
