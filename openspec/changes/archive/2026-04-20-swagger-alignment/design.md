# Alineación de endpoints con swagger.md

## Resumen

Ajustar todos los endpoints existentes y crear los faltantes para que coincidan exactamente con los esquemas, rutas y formatos de respuesta definidos en `swagger.md`.

---

## Grupos de trabajo

### Grupo 1 — Fundación: rutas y wrapper de respuesta

> Base que todo lo demás depende. Hacer primero.

- [ ] **T1.1** Crear helper `rest.SuccessResponse(data)` y `rest.ErrorResponse(message, code)` en un archivo `internal/infrastructure/rest/response.go`
- [ ] **T1.2** Renombrar prefijos de rutas en `routes.go`:
  - `/api/v1/tables` → `/mesas`
  - `/api/v1/orders` → `/ordenes`
  - `/api/v1/ingredients` → `/ingredientes`
  - `/api/v1/menu` → `/menu`
  - `/api/v1/categories` → `/categorias`
  - Mover `/auth/login` a que tenga el mismo prefijo que el resto de auth

---

### Grupo 2 — Mesas

- [ ] **T2.1** Actualizar `table/mapper.go`: producir `Response` con los campos del schema `Mesa`:
  ```
  id, number, state (free/occupied/reserved/cleaning), capacity,
  waiter {id, name}, guests, occupiedMinutes, currentOrder {id, total}
  ```
- [ ] **T2.2** Aplicar el wrapper `{ success, data }` en los handlers `GetAll`, `GetByID` y `Create`
- [ ] **T2.3** Cambiar `POST /mesas` para que acepte `numero` y `capacidad` en el body (según swagger)
- [ ] **T2.4** Renombrar `PATCH /mesas/:id` → `PATCH /mesas/:id/estado` y que solo acepte `{ estado }` como body

---

### Grupo 3 — Auth

- [ ] **T3.1** Ajustar respuesta de `POST /auth/login` al schema exacto del swagger:
  ```json
  { "success": true, "data": { "token": "...", "usuario": {...}, "expires_in": 900 } }
  ```
- [ ] **T3.2** Crear handler `POST /auth/register` (registrar usuario con rol, solo ADMIN/PROPIETARIO)
- [ ] **T3.3** Crear handler `GET /auth/me` (retorna `Usuario` del token actual)
- [ ] **T3.4** Crear handler `POST /auth/logout` (revoca el refresh token de la cookie)

---

### Grupo 4 — Usuarios

- [ ] **T4.1** Crear `user/mapper.go` con `ToUsuarioResponse` que produzca el schema `Usuario`:
  ```
  id, nombre, email, rol (nombre del rol, no ID), telefono, activo, fecha_creacion, ultimo_acceso
  ```
  Requiere lookup de nombre de rol: puede agregarse al dominio `User` o resolverse en el mapper.
- [ ] **T4.2** Registrar rutas en `routes.go`: `GET /usuarios`, `GET /usuarios/:id`, `PATCH /usuarios/:id`, `DELETE /usuarios/:id`
- [ ] **T4.3** Crear handler `GET /usuarios` — lista usuarios del restaurante (solo PROPIETARIO)
- [ ] **T4.4** Crear handler `GET /usuarios/:id` — obtiene un usuario por ID
- [ ] **T4.5** Crear handler `PATCH /usuarios/:id` — actualiza `nombre`, `email`, `rol`, `activo`, `telefono`
- [ ] **T4.6** Crear handler `DELETE /usuarios/:id` — desactiva usuario (`is_active = false`)

---

### Grupo 5 — Órdenes

- [ ] **T5.1** Crear `order/mapper.go` con `ToOrdenResponse` que produzca el schema `Orden`:
  ```
  id, mesa_id, mesero_id, estado (abierta/enviada/pagada/cancelada),
  items [{id, menu_item_id, nombre, cantidad, precio_unitario, notas, estado}],
  subtotal, impuestos (19%), total, fecha_creacion
  ```
- [ ] **T5.2** Ajustar `POST /ordenes` para aceptar `{ mesa_id, mesero_id }` (sin `items` obligatorios al crear)
- [ ] **T5.3** Crear handler `GET /ordenes/:id` — detalle de orden con su schema completo
- [ ] **T5.4** Crear handler `POST /ordenes/:id/items` — agregar items a una orden existente
- [ ] **T5.5** Crear handler `DELETE /ordenes/:id/items/:item_id` — cancelar un item de la orden
- [ ] **T5.6** Crear handler `POST /ordenes/:id/enviar-cocina` — cambia estado de la orden a `enviada`
- [ ] **T5.7** Aplicar wrapper `{ success, data }` a `POST /ordenes` y `GET /ordenes`

---

### Grupo 6 — Menú

- [ ] **T6.1** Actualizar `product/mapper.go` o crear `menu/mapper.go` con `ToMenuItemResponse`:
  ```
  id, nombre, categoria, precio, descripcion, disponible
  ```
- [ ] **T6.2** Aplicar wrapper `{ success, data }` a todos los handlers de menú
- [ ] **T6.3** Verificar que `POST /menu` acepta `nombre`, `categoria`, `precio`, `descripcion`, `disponible`

---

### Grupo 7 — Ingredientes

- [ ] **T7.1** Aplicar wrapper `{ success, data }` en handlers `GetAll`, `Create`, `GetByID`, `Update`
- [ ] **T7.2** Crear handler `PATCH /ingredientes/:id/stock` — actualiza stock con `{ cantidad, tipo_movimiento (entrada|salida), motivo }`

---

### Grupo 8 — Pagos *(módulo nuevo)*

- [ ] **T8.1** Crear dominio `internal/domain/payment/payment.go` con struct `Pago`
- [ ] **T8.2** Agregar migration SQL para tabla `payments`
- [ ] **T8.3** Crear `POST /pagos` handler — procesa pago de una orden (vincula con checkout existente)
- [ ] **T8.4** Crear `GET /pagos/:id/factura` handler — genera datos de factura
- [ ] **T8.5** Registrar rutas en `routes.go`

---

### Grupo 9 — División de cuenta *(módulo nuevo)*

- [ ] **T9.1** Crear handler `POST /ordenes/:id/dividir` que acepte `{ tipo_division, numero_partes, divisiones }`
- [ ] **T9.2** Implementar lógica de división en el usecase de órdenes

---

### Grupo 10 — Reportes *(módulo nuevo)*

- [ ] **T10.1** Crear handler `GET /reportes/ventas` con parámetros `fecha_inicio`, `fecha_fin`, `tipo`
- [ ] **T10.2** Crear handler `GET /reportes/inventario` — ingredientes bajo stock + valor total
- [ ] **T10.3** Crear handler `GET /reportes/propinas` con parámetros `fecha_inicio`, `fecha_fin`
- [ ] **T10.4** Registrar rutas en `routes.go`

---

## Orden de implementación sugerido

```
T1 → T2 → T3 → T4 → T5 → T6 → T7 → T8 → T9 → T10
```

Los grupos 1–7 ajustan funcionalidad ya existente.
Los grupos 8–10 son módulos nuevos que se pueden abordar en iteraciones separadas.

---

## Mapa de rutas: actual → swagger

| Actual                          | Swagger                              | Estado        |
|---------------------------------|--------------------------------------|---------------|
| `POST /auth/login`              | `POST /auth/login`                   | Parcial       |
| —                               | `POST /auth/register`                | Falta         |
| —                               | `GET /auth/me`                       | Falta         |
| —                               | `POST /auth/logout`                  | Falta         |
| `GET /api/v1/tables`            | `GET /mesas`                         | Renombrar     |
| `POST /api/v1/tables`           | `POST /mesas`                        | Renombrar     |
| `PATCH /api/v1/tables/:id`      | `PATCH /mesas/:id/estado`            | Renombrar + ajustar |
| —                               | `GET /usuarios`                      | Falta         |
| —                               | `GET /usuarios/:id`                  | Falta         |
| —                               | `PATCH /usuarios/:id`                | Falta         |
| —                               | `DELETE /usuarios/:id`               | Falta         |
| `POST /api/v1/orders`           | `POST /ordenes`                      | Renombrar + ajustar |
| `GET /api/v1/orders`            | `GET /ordenes/:id`                   | Renombrar + ajustar |
| —                               | `POST /ordenes/:id/items`            | Falta         |
| —                               | `DELETE /ordenes/:id/items/:item_id` | Falta         |
| —                               | `POST /ordenes/:id/enviar-cocina`    | Falta         |
| `GET /api/v1/menu`              | `GET /menu`                          | Renombrar     |
| `POST /api/v1/menu`             | `POST /menu`                         | Renombrar     |
| `PATCH /api/v1/menu/:id`        | `PATCH /menu/:id`                    | Renombrar     |
| `GET /api/v1/ingredients`       | `GET /ingredientes`                  | Renombrar     |
| `POST /api/v1/ingredients`      | `POST /ingredientes`                 | Renombrar     |
| —                               | `PATCH /ingredientes/:id/stock`      | Falta         |
| —                               | `POST /ordenes/:id/dividir`          | Falta         |
| —                               | `POST /pagos`                        | Falta         |
| —                               | `GET /pagos/:id/factura`             | Falta         |
| —                               | `GET /reportes/ventas`               | Falta         |
| —                               | `GET /reportes/inventario`           | Falta         |
| —                               | `GET /reportes/propinas`             | Falta         |
