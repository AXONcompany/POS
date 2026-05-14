## MODIFIED Requirements

### Requirement: Cancelación de ítem con restauración de inventario
El sistema DEBE permitir cancelar un ítem de una orden existente marcándolo como cancelado (`cancelled_at`), recalcular el total de la orden y restaurar el stock de ingredientes correspondientes a la receta de dicho ítem en una única transacción atómica.

#### Scenario: Cancelación exitosa
- **WHEN** un usuario autorizado cancela un ítem activo de una orden en estado `PENDING` o `SENT`
- **THEN** el sistema actualiza `cancelled_at` en `order_items`, restaura el stock en `ingredients` basándose en la receta del ítem y reduce el `total_amount` de la orden, todo en una sola transacción; retorna HTTP 200 con la orden actualizada.

#### Scenario: Ítem ya cancelado
- **WHEN** un usuario intenta cancelar un ítem que ya tiene `cancelled_at IS NOT NULL`
- **THEN** el sistema rechaza la operación con HTTP 409 (`ErrItemAlreadyCancelled`) sin modificar el inventario ni la orden.

#### Scenario: Orden en estado no editable para cancelación
- **WHEN** un usuario intenta cancelar un ítem de una orden en estado `PREPARING`, `READY`, `PAID` o `CANCELLED`
- **THEN** el sistema rechaza la operación con HTTP 422 (`ErrInvalidStatusTransition`) sin modificar el ítem ni el inventario.
