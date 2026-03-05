# POS API Documentation

## Base URL
`http://localhost:8080`

## Response Format
All responses follow: `{ "success": true/false, "data": {...} }`

## Authentication
All endpoints (except `/health`, `/ping`, `/auth/login`) require `Authorization: Bearer <token>`.

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
| POST | `/auth/register` | Propietario | Register new user |
| GET | `/auth/me` | Any | Current user info |
| POST | `/auth/logout` | Any | Revoke refresh token |

#### POST `/auth/login`
```json
// Request
{ "email": "admin@test.com", "password": "admin" }
// Response
{ "success": true, "data": { "token": "...", "usuario": {...}, "expires_in": 900 } }
```

---

### Mesas (Tables)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/mesas` | All | List tables |
| GET | `/mesas/:id` | All | Get table |
| POST | `/mesas` | Cajero, Propietario | Create table |
| PATCH | `/mesas/:id/estado` | Cajero, Propietario | Update status |
| DELETE | `/mesas/:id` | Cajero, Propietario | Delete table |
| POST | `/mesas/:id/assign` | Cajero, Propietario | Assign waiter |

```json
// POST /mesas
{ "numero": 5, "capacidad": 4 }
// PATCH /mesas/:id/estado
{ "estado": "occupied" }
```

---

### Usuarios (Users)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/usuarios` | Propietario | List users |
| GET | `/usuarios/:id` | Propietario | Get user |
| PATCH | `/usuarios/:id` | Propietario | Update user |
| DELETE | `/usuarios/:id` | Propietario | Deactivate user |

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

```json
// POST /ingredientes
{ "name": "Sugar", "unit_of_measure": "kg", "type": "dry", "stock": 50 }
// PATCH /ingredientes/:id/stock
{ "cantidad": 50, "tipo_movimiento": "entrada", "motivo": "Weekly purchase" }
```

---

### Categorias (Categories)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/categorias` | Propietario | Create category |
| GET | `/categorias` | Propietario | List categories |

---

### Menu

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/menu` | All | List menu items |
| POST | `/menu` | Propietario | Create menu item |
| PATCH | `/menu/:id` | Propietario | Update menu item |

```json
// POST /menu
{ "name": "Hamburguesa", "sales_price": 12.00, "ingredients": [{ "ingredient_id": 1, "quantity": 0.2 }] }
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

```json
// POST /ordenes
{ "mesa_id": "1", "mesero_id": "3" }
// POST /ordenes/:id/items
{ "items": [{ "menu_item_id": "1", "cantidad": 2, "notas": "Sin cebolla" }] }
// POST /ordenes/:id/dividir
{ "tipo_division": "partes_iguales", "numero_partes": 3 }
```

---

### Pagos (Payments)

| Method | Route | Roles | Description |
|---|---|---|---|
| POST | `/pagos` | Cajero, Propietario | Process payment |
| GET | `/pagos/:id/factura` | Cajero, Propietario | Generate invoice |

```json
// POST /pagos
{ "orden_id": "1", "metodo_pago": "tarjeta", "monto": 41650, "propina": 4165 }
```

---

### Reportes (Reports)

| Method | Route | Roles | Description |
|---|---|---|---|
| GET | `/reportes/ventas` | Propietario | Sales report |
| GET | `/reportes/inventario` | Propietario | Inventory report |
| GET | `/reportes/propinas` | Propietario | Tips report |

**Query params:** `fecha_inicio` (date), `fecha_fin` (date), `tipo` (por_dia/por_item/por_hora)
