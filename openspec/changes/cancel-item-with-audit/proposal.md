## Why

Actualmente, el sistema no permite la cancelación de ítems individuales de una orden una vez que han sido enviados, ni realiza la restauración automática del inventario asociada a esa cancelación. Además, carecemos de un registro de auditoría que permita rastrear quién realizó cambios críticos y revertirlos en caso de error. Esta funcionalidad es vital para la integridad operativa y financiera del POS.

## What Changes

- **Cancelación de Ítems**: Capacidad para marcar un ítem de orden como cancelado (`cancelled_at`) en lugar de eliminarlo físicamente.
- **Restauración de Inventario**: Al cancelar un ítem, el sistema devolverá automáticamente las cantidades de ingredientes al stock basándose en la receta del producto.
- **Registro de Auditoría (Snapshots)**: Implementación de un sistema de auditoría que guarda el estado "antes" y "después" de cada entidad modificada en formato JSONB.
- **Ajuste de Totales**: El total de la orden se recalculará automáticamente al cancelar ítems.

## Capabilities

### New Capabilities
- `audit-logging`: Registro de cambios en entidades con soporte para snapshots (old_value/new_value) para trazabilidad y recuperación.

### Modified Capabilities
- `order-item-addition`: Se modifica para permitir la reversión de las adiciones mediante la cancelación y restauración de stock.

## Impact

- **Base de Datos**: Nuevas columnas en `order_items` y nueva tabla `audit_log`.
- **Repositorio**: Nueva lógica transaccional para cancelación y restauración de stock.
- **Usecase**: Inclusión de lógica de auditoría (snapshots) en el flujo de cancelación.
- **API**: Nuevo endpoint para cancelar ítems de una orden.
