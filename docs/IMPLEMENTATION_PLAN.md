# Plan de Implementación AXON POS

> Generado: 2026-03-23
> Alcance: Flujo de orden, integridad de pagos (solo efectivo), reportes, seguridad base.
> Fuera de alcance: Integración con pasarela de pago, KDS.

---

## Migraciones necesarias

| # | Archivo | Contenido |
|---|---|---|
| 002 | `002_order_items_cancel.sql` | Columna `cancelled_at TIMESTAMPTZ NULL` en `order_items` |
| 003 | `003_divisions.sql` | Tabla `order_divisions` + FK en `payments.division_ref` |
| 004 | `004_shift_closings.sql` | Tabla `shift_closings` para Z-report |
| 005 | `005_audit_log.sql` | Tabla `audit_log` |

### Detalle de migraciones

**002 — `cancelled_at` en order_items**
```sql
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMPTZ NULL;
```

**003 — order_divisions**
```sql
CREATE TABLE IF NOT EXISTS order_divisions (
    id             VARCHAR(50) PRIMARY KEY,
    order_id       BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    venue_id       INTEGER NOT NULL REFERENCES venues(id),
    division_type  VARCHAR(20) NOT NULL,   -- partes_iguales | por_monto | por_item
    amount         DECIMAL(10,2) NOT NULL,
    tax            DECIMAL(10,2) NOT NULL,
    total          DECIMAL(10,2) NOT NULL,
    is_paid        BOOLEAN NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_order_divisions_order_id ON order_divisions(order_id);

ALTER TABLE payments ADD COLUMN IF NOT EXISTS division_ref VARCHAR(50)
    REFERENCES order_divisions(id);
```

**004 — shift_closings**
```sql
CREATE TABLE IF NOT EXISTS shift_closings (
    id               BIGSERIAL PRIMARY KEY,
    venue_id         INTEGER NOT NULL REFERENCES venues(id),
    user_id          INTEGER NOT NULL REFERENCES users(id),
    closed_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    report_snapshot  JSONB NOT NULL
);
```

**005 — audit_log**
```sql
CREATE TABLE IF NOT EXISTS audit_log (
    id           BIGSERIAL PRIMARY KEY,
    venue_id     INTEGER REFERENCES venues(id),
    user_id      INTEGER REFERENCES users(id),
    action       VARCHAR(50) NOT NULL,    -- CREATE, UPDATE, DELETE, LOGIN
    entity_type  VARCHAR(50) NOT NULL,   -- order, payment, ingredient, etc.
    entity_id    VARCHAR(50),
    old_value    JSONB,
    new_value    JSONB,
    ip_address   VARCHAR(45),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_venue_id ON audit_log(venue_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_user_id  ON audit_log(user_id);
```

---

## Fase 1 — Flujo de orden (semana 1)

### 1. Máquina de estados de orden `S`

**Transiciones válidas:**
```
PENDING(1)    → SENT(2) | CANCELLED(6)
SENT(2)       → PREPARING(3) | CANCELLED(6)
PREPARING(3)  → READY(4)
READY(4)      → PAID(5)
PAID(5)       → (ninguna)
CANCELLED(6)  → (ninguna)
```

**Archivos a modificar:**

- `internal/domain/order/order.go`
  - Agregar mapa `validTransitions map[int][]int`
  - Agregar función `CanTransitionTo(current, next int) bool`
  - Agregar `ErrInvalidStatusTransition`

- `internal/usecase/order/usecase.go`
  - `UpdateOrderStatus`: obtener estado actual, validar transición, luego actualizar
  - `CheckoutOrder`: pasar por la misma validación (debe estar en READY → PAID)

- `internal/infrastructure/persistence/postgres/order_repository.go`
  - Agregar `GetStatusByID(ctx, orderID int64, venueID int) (int, error)`

- `internal/infrastructure/rest/order/handler.go`
  - Mapear `ErrInvalidStatusTransition` → HTTP 422

---

### 2. `AddProductToOrder` + deducción de inventario `L`

> Depende de: migración 002 (indirectamente); el diseño de interfaces.
> El stub actual retorna `nil` sin hacer nada. Un mesero no puede agregar platos en producción.

**Archivos a modificar:**

- `internal/domain/order/order.go`
  - Agregar `StockDeduction{ IngredientID int64, VenueID int, Quantity float64 }`

- `internal/usecase/order/usecase.go`
  - Agregar interfaz `ProductInventoryRepository`:
    ```go
    type ProductInventoryRepository interface {
        GetProductPrice(ctx context.Context, productID int64, venueID int) (float64, error)
        GetRecipeLines(ctx context.Context, productID int64) ([]domainOrder.RecipeLine, error)
    }
    ```
  - Agregar `RecipeLine{ IngredientID int64, QuantityRequired float64 }` al domain
  - Inyectar `invRepo ProductInventoryRepository` en `Usecase`
  - Implementar `AddProductToOrder`:
    1. Obtener orden, validar estado PENDING o SENT
    2. Para cada item: buscar precio y receta
    3. Acumular deducción de stock por ingrediente
    4. Llamar `repo.AddItemsWithInventory(ctx, orderID, venueID, items, deductions)`

- `internal/infrastructure/persistence/postgres/order_repository.go`
  - Implementar `AddItemsWithInventory` en una sola `pgx.Tx`:
    - `UPDATE ingredients SET stock = stock - $qty WHERE id = $id AND venue_id = $venueID AND stock >= $qty` (si 0 rows afectadas → rollback con `ErrInsufficientStock`)
    - INSERT en `order_items`
    - `UPDATE orders SET total_amount = total_amount + $delta, updated_at = NOW()`

- `cmd/server/main.go`
  - Pasar `productRepo` como segundo argumento a `uorder.NewUsecase(...)`

---

### 3. `CancelOrderItem` + restauración de inventario `M`

> Depende de: migración 002.

**Archivos a modificar:**

- `internal/usecase/order/usecase.go`
  - Agregar a `Repository`:
    - `GetOrderItem(ctx, itemID, orderID int64) (*domainOrder.OrderItem, error)`
    - `CancelItemWithInventoryRestore(ctx, item *domainOrder.OrderItem, venueID int, deductions []domainOrder.StockDeduction) error`
  - Implementar `CancelOrderItem`:
    1. Obtener el item de la orden
    2. Validar que la orden está en PENDING o SENT
    3. Verificar que `cancelled_at IS NULL`
    4. Calcular deducción inversa (restaurar stock)
    5. Llamar al repo

- `internal/infrastructure/persistence/postgres/order_repository.go`
  - Implementar `CancelItemWithInventoryRestore` en 1 TX:
    - `UPDATE order_items SET cancelled_at = NOW() WHERE id = $itemID AND cancelled_at IS NULL`
    - `UPDATE ingredients SET stock = stock + $qty WHERE id = $ingredientID AND venue_id = $venueID`
    - `UPDATE orders SET total_amount = total_amount - $delta, updated_at = NOW()`

---

## Fase 2 — Integridad de pagos (semana 2)

### 4. División de cuenta persiste en BD `M`

> Depende de: migración 003.

**Archivos a modificar:**

- `internal/domain/order/order.go`
  - Agregar struct `OrderDivision{ ID, OrderID, VenueID, DivisionType, Amount, Tax, Total, IsPaid, CreatedAt }`

- `internal/usecase/order/usecase.go`
  - Mover lógica de cálculo del handler al usecase
  - Agregar a `Repository`:
    - `CreateDivisions(ctx, divisions []domainOrder.OrderDivision) error`
    - `GetDivisionsByOrderID(ctx, orderID int64) ([]domainOrder.OrderDivision, error)`
    - `MarkDivisionPaid(ctx, divisionID string, orderID int64) error`

- `internal/infrastructure/persistence/postgres/order_repository.go`
  - Implementar los tres métodos anteriores

- `internal/infrastructure/rest/order/handler.go`
  - `DivideOrder` pasa a ser un wrapper simple sobre el usecase (eliminar lógica de cálculo del handler)

---

### 5. Checkout atómico `M`

> Depende de: fase 4 (divisiones).

**Archivos a modificar:**

- `internal/infrastructure/persistence/postgres/order_repository.go`
  - Implementar `CheckoutAtomic(ctx, orderID int64, venueID int) error` en 1 TX:
    1. `SELECT ... FOR UPDATE` en la orden (lock anti-concurrencia)
    2. Verificar que está en estado READY(4); si ya es 5 → `ErrAlreadyPaid`
    3. Si tiene divisiones → verificar que todas tienen `is_paid = true` → si no → `ErrDivisionsNotFullyPaid`
    4. Si no tiene divisiones → verificar que existe al menos 1 pago aprobado
    5. `UPDATE orders SET status_id = 5, updated_at = NOW()`

- `internal/usecase/order/usecase.go`
  - `CheckoutOrder` delega a `repo.CheckoutAtomic`

- `internal/infrastructure/rest/order/handler.go`
  - Mapear `ErrAlreadyPaid` y `ErrDivisionsNotFullyPaid` → HTTP 409

---

### 6. Validación de monto + idempotencia de pagos `M`

**Archivos a modificar:**

- `internal/usecase/payment/usecase.go`
  - Antes de crear pago:
    1. Si `divisionID != nil`: comparar `amount` con total de la división → `ErrAmountMismatch`
    2. Si no: comparar `amount` con `total_amount` de la orden → `ErrAmountMismatch`
    3. Verificar que no existe ya un pago aprobado para `(orderID, divisionID)` → `ErrAlreadyPaid`

- `internal/infrastructure/persistence/postgres/payment_repository.go`
  - Agregar `GetApprovedPaymentForOrder(ctx, orderID int64, divisionID *string) (*domainPayment.Payment, error)`
  - Agregar `GetOrderTotal(ctx, orderID int64, venueID int) (float64, error)`
  - Agregar `GetDivisionTotal(ctx, divisionID string, orderID int64) (float64, error)`
  - Envolver `Create` en TX que también llama `MarkDivisionPaid` si `divisionID != nil`

---

## Fase 3 — Reportes y seguridad (semana 3)

### 7. X-Report (corte de caja sin cerrar) `M`

Nueva ruta: `GET /reportes/x-report?desde=<ISO8601>` — solo PROPIETARIO y CAJERO.

Retorna: total órdenes, total ventas, total propinas, desglose por método de pago, propinas por empleado.

**Archivos a modificar:**

- `internal/usecase/report/usecase.go` — agregar `GetXReport(ctx, venueID int, since time.Time) (*XReportData, error)`
- `internal/infrastructure/persistence/postgres/report_repository.go` — query de agregados sobre `payments JOIN users WHERE created_at >= $since AND status = 'aprobado'`
- `internal/infrastructure/rest/report/handler.go` — agregar handler `XReport`
- `internal/infrastructure/rest/routes.go` — registrar ruta

---

### 8. Z-Report (cierre de día) `M`

> Depende de: migración 004.

Nueva ruta: `POST /reportes/z-report` — solo PROPIETARIO y CAJERO. No idempotente (cierra el turno).

Retorna: mismos datos que X-Report + efectivo esperado en caja. Persiste snapshot en `shift_closings.report_snapshot (JSONB)`.

**Archivos a modificar:**

- `internal/usecase/report/usecase.go` — agregar `CloseShiftZReport(ctx, venueID, userID int) (*ZReportData, error)`
- `internal/infrastructure/persistence/postgres/report_repository.go` — implementar query + INSERT en `shift_closings`
- `internal/infrastructure/rest/report/handler.go` — agregar handler `ZReport`
- `internal/infrastructure/rest/routes.go` — registrar ruta

---

### 9. Rate limiting en login `S`

**Archivos a crear:**

- `internal/infrastructure/rest/middleware/ratelimit.go`
  - Token bucket por IP usando `golang.org/x/time/rate`
  - 5 req/seg, burst 10
  - Map `IP → *rate.Limiter` protegido con `sync.Mutex`

**Archivos a modificar:**

- `internal/infrastructure/rest/routes.go` — aplicar middleware solo a `POST /auth/login`
- `go.mod` / `go.sum` — agregar `golang.org/x/time`

---

### 10. Audit log `M`

> Depende de: migración 005.

**Archivos a crear:**

- `internal/infrastructure/rest/middleware/audit.go` — middleware gin que corre `c.Next()` y luego loguea métodos mutantes (POST/PUT/PATCH/DELETE) a `audit_log`
- `internal/infrastructure/persistence/postgres/audit_repository.go`

**Archivos a modificar:**

- `internal/infrastructure/rest/routes.go` — aplicar al grupo `api` protegido

---

## Fase 4 — UX y limpieza (semana 3, paralelo)

### 11. Auto-actualización de estado de mesa `S`

**Archivos a modificar:**

- `internal/usecase/order/usecase.go`
  - Agregar interfaz `TableStatusUpdater{ UpdateTableStatus(ctx, tableID int64, venueID int, status string) error }`
  - Inyectar en `Usecase`
  - En `CreateOrderWithoutItems`: si `tableID != nil` → mesa pasa a `"occupied"`
  - Después de `CheckoutAtomic` exitoso: si la orden tiene `tableID` → mesa pasa a `"dirty"`

- `internal/infrastructure/persistence/postgres/table_repository.go`
  - Agregar `UpdateTableStatus(ctx, tableID int64, venueID int, status string) error`

- `cmd/server/main.go` — pasar `tableRepo` al usecase de orden

---

### 12. Cambio de contraseña del propietario `S`

Nueva ruta: `POST /propietario/password`

```json
{ "password_actual": "...", "password_nuevo": "...", "confirmar_password": "..." }
```

**Archivos a modificar:**

- `internal/usecase/owner/usecase.go` — agregar `ChangePassword(ctx, ownerID int, currentPass, newPass string) error`
- `internal/infrastructure/persistence/postgres/owner_repository.go` — agregar `UpdatePassword`
- `internal/infrastructure/rest/owner/handler.go` — agregar handler
- `internal/infrastructure/rest/routes.go` — registrar ruta

---

### 13. Soft delete consistente `S`

Revisar y parchear queries que no filtran `deleted_at IS NULL`:

- `db/queries/*.sql` — verificar `ingredients`, `products`, `categories`, `tables`
- SQL raw en `order_repository.go` — verificar en `GetByID`, `ListByTable`

---

### 14. CORS configurable `S`

**Archivos a modificar:**

- `internal/config/config.go`
  - Agregar `AllowedOrigins string` leído de `ALLOWED_ORIGINS` env var (default `"*"` para desarrollo)

- `internal/infrastructure/rest/server.go`
  - `origins := strings.Split(cfg.AllowedOrigins, ",")`
  - Usar `origins` en lugar del literal `"*"`

---

## Orden de ejecución recomendado

```
Día 1-2:   CORS (14) · Máquina de estados (1) · Migraciones 002 y 003
Día 3-5:   AddProductToOrder + inventario (2)
Día 6-7:   CancelOrderItem + restauración (3)
Día 8-9:   Divisiones persisten (4) · Checkout atómico (5)
Día 10:    Validación de monto + idempotencia (6)
Día 11-12: X/Z reports (7+8) · Migraciones 004 y 005
Día 13:    Rate limiting (9) · Audit log (10)
Día 14:    Table auto-status (11) · Password change (12) · Soft delete (13)
```

---

## Tabla resumen

| # | Feature | Complejidad | Migración | Archivos clave |
|---|---|---|---|---|
| 1 | Máquina de estados de orden | S | — | domain/order, usecase/order, order_repository |
| 2 | `AddProductToOrder` + inventario | L | 002 | usecase/order, order_repository, main.go |
| 3 | `CancelOrderItem` + restauración | M | 002 | usecase/order, order_repository |
| 4 | Divisiones persisten en BD | M | 003 | usecase/order, order_repository, order/handler |
| 5 | Checkout atómico | M | — | usecase/order, order_repository, order/handler |
| 6 | Validación monto + idempotencia | M | — | usecase/payment, payment_repository |
| 7 | X-Report | M | — | usecase/report, report_repository, report/handler |
| 8 | Z-Report | M | 004 | usecase/report, report_repository, report/handler |
| 9 | Rate limiting login | S | — | middleware/ratelimit.go, routes.go |
| 10 | Audit log | M | 005 | middleware/audit.go, audit_repository.go |
| 11 | Table auto-status | S | — | usecase/order, table_repository, main.go |
| 12 | Password change propietario | S | — | usecase/owner, owner_repository, owner/handler |
| 13 | Soft delete consistente | S | — | db/queries/*.sql, order_repository |
| 14 | CORS configurable | S | — | config/config.go, rest/server.go |

---

## Notas de diseño

**Transacciones cross-tabla:** Los items 2, 3 y 5 requieren operaciones atómicas sobre múltiples tablas. El patrón del proyecto es abrir `pgx.Tx` dentro del método del repositorio (ver `order_repository.go::Create`). Seguir ese patrón — no exponer `pgx.Tx` a la capa usecase para mantener clean architecture.

**sqlc vs SQL raw:** El proyecto mezcla ambos. Para las operaciones nuevas con transacciones complejas o múltiples statements, usar SQL raw (patrón de `order_repository.go` y `payment_repository.go`). Para CRUD simple nuevo, cualquiera funciona.

**IVA hardcodeado al 19%:** Aparece en múltiples lugares del handler de divisiones. Candidato futuro para configuración por venue.
