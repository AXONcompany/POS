## 1. Base de Datos y Modelos

- [x] 1.1 Crear migración `000012_add_order_items_cancelled_at` para añadir la columna `cancelled_at` a `order_items`.
- [x] 1.2 Crear migración `000013_add_audit_log` para crear la tabla `audit_log` con soporte para `JSONB`.
- [x] 1.3 Actualizar struct `OrderItem` en `internal/domain/order/order.go` para incluir el campo `CancelledAt`.
- [x] 1.4 Crear struct `AuditEntry` en `internal/domain/audit/audit.go`.

## 2. Repositorio (Persistencia)

- [x] 2.1 Implementar método `GetOrderItem(ctx, itemID, orderID)` en `order_repository.go`.
- [x] 2.2 Implementar método `CancelItemWithInventoryRestore(ctx, itemID, venueID, auditEntry)` en `order_repository.go` usando una transacción atómica.
- [x] 2.3 Implementar método `SaveAudit(ctx, entry)` en un nuevo `audit_repository.go` o dentro de `order_repository.go` de forma genérica.

## 3. Capa de Usecase

- [x] 3.1 Implementar lógica de snapshot ("before") en `CancelOrderItem`.
- [x] 3.2 Implementar cálculo de restauración de stock y llamada al repositorio transaccional.
- [x] 3.3 Implementar lógica de snapshot ("after") y registro de auditoría final.

## 4. API y Rutas

- [x] 4.1 Añadir endpoint `DELETE /ordenes/:id/items/:item_id` en el handler de órdenes.
- [x] 4.2 Registrar la nueva ruta en `internal/infrastructure/rest/routes.go`.
- [x] 4.3 Mapear nuevos errores (`ErrItemAlreadyCancelled`, `ErrInvalidStatusTransition`) a códigos HTTP adecuados.

## 5. Pruebas y Validación

- [x] 5.1 Crear prueba de integración que verifique la restauración de stock tras la cancelación.
- [x] 5.2 Verificar que el registro de auditoría contenga los JSON correctos.
