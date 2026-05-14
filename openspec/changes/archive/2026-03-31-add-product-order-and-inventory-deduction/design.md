## Context

El sistema actualmente tiene `AddProductToOrder` como stub. El repositorio de orden ya usa el patrón `pgx.Tx` abierta dentro del método del repositorio (ver `Create`). El proyecto mezcla sqlc y SQL raw; para operaciones con transacciones multi-tabla se usa SQL raw. La deducción de stock debe ser atómica con la inserción de items para evitar inconsistencias: si falla cualquier paso, ningún cambio persiste.

Restricción de arquitectura: la TX no puede exponerse a la capa usecase (clean architecture). El repositorio encapsula la transacción completa.

## Goals / Non-Goals

**Goals:**
- Implementar `AddProductToOrder` con precio real desde BD y deducción de ingredientes.
- Garantizar atomicidad total: items + stock + total de orden en una sola TX.
- Retornar `ErrInsufficientStock` cuando un ingrediente no tiene stock suficiente.
- Mantener la interfaz `Repository` existente como base y agregar una nueva interfaz `ProductInventoryRepository` definida donde se usa (usecase), no donde se implementa.

**Non-Goals:**
- Integración con KDS o notificaciones en tiempo real.
- Validación de permisos por rol (ya está en middleware).
- Cancelación de items (`CancelOrderItem`) — es el paso 3 del plan.
- División de cuenta o checkout — pasos posteriores.

## Decisions

### 1. Nueva interfaz `ProductInventoryRepository` en el usecase, no en el dominio

**Decisión:** Definir la interfaz en `internal/usecase/order/usecase.go`.

**Rationale:** Interface Segregation — el usecase define solo lo que necesita. El `postgres.ProductRepository` la implementará implícitamente. Esto sigue el patrón ya establecido en el proyecto con `Repository`.

**Alternativa descartada:** Definirla en el dominio — contamina el dominio con dependencias de infraestructura.

---

### 2. `AddItemsWithInventory` en una sola `pgx.Tx` dentro del repositorio

**Decisión:** El método del repositorio abre TX, ejecuta INSERT en `order_items`, UPDATE de stock por cada ingrediente, y UPDATE del total de la orden, y hace commit. El rollback es automático vía `defer tx.Rollback`.

**Rationale:** Consistencia total. Si el stock de un ingrediente es insuficiente (0 rows afectadas en el UPDATE), se hace rollback inmediato y se retorna `ErrInsufficientStock`. Sigue el patrón de `Create`.

**Alternativa descartada:** Hacer los UPDATEs de stock en goroutines paralelas — complica el manejo de errores y hace imposible rollback ordenado.

---

### 3. Validación de estado de la orden antes de agregar items

**Decisión:** `AddProductToOrder` usa el `GetStatusByID` ya implementado y valida que el estado sea `PENDING(1)` o `SENT(2)`. Para estados posteriores retorna `ErrInvalidStatusTransition`.

**Rationale:** Reutiliza infraestructura ya existente (paso 1). No se puede agregar platos a una orden que ya está en cocina o pagada.

---

### 4. Acumulación de deducciones por ingrediente antes de la TX

**Decisión:** El usecase recorre los items, consulta la receta de cada producto y acumula `map[ingredientID]totalQty` antes de llamar al repositorio. El repositorio recibe un slice de `StockDeduction` ya consolidado.

**Rationale:** Evita múltiples UPDATEs al mismo ingrediente dentro de la TX si dos items del pedido usan el mismo ingrediente. Un solo UPDATE por ingrediente es más eficiente y evita deadlocks.

## Risks / Trade-offs

- **[Riesgo] Receta vacía:** Si un producto no tiene receta técnica, `GetRecipeLines` retorna slice vacío. El item se agrega sin descontar stock. → Mitigación: aceptable para el MVP; es responsabilidad del propietario configurar la receta al crear el producto.

- **[Riesgo] Race condition de stock:** Entre el `GetStatusByID` y el `UPDATE ingredients SET stock = stock - $qty`, otro request podría tomar el último stock. → Mitigación: el UPDATE con `AND stock >= $qty` actúa como guard atómico. Si retorna 0 rows, la TX hace rollback con `ErrInsufficientStock`. No se necesita `SELECT FOR UPDATE` en el ingrediente.

- **[Trade-off] Dos queries al repo antes de la TX (precio + receta):** El usecase hace N queries para obtener precios y recetas antes de abrir la TX. → Aceptable para el MVP; el volumen de items por orden es bajo. La TX solo ejecuta las escrituras.

## Migration Plan

Sin migraciones de esquema. La columna `cancelled_at` (migración 002) no es necesaria para este paso. Solo cambios en código.

Despliegue: reemplazar binario. Sin rollback especial — si falla, el stub anterior retornaba `nil`; la nueva versión retorna error explícito en casos de fallo, lo cual es un comportamiento más correcto.

## Open Questions

- ¿El `UnitPrice` que se guarda en `order_items` debe ser el precio de venta del producto en ese momento (snapshot), o debe consultarse cada vez? → **Decisión actual:** snapshot al momento del pedido (se obtiene de `GetProductPrice` y se guarda en el item). Correcto para auditoría.
