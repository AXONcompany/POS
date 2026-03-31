## Why

`AddProductToOrder` es actualmente un stub que retorna `nil` sin hacer nada, lo que significa que ningún mesero puede agregar platos a una orden en producción. Tampoco existe deducción de inventario, por lo que el stock de ingredientes nunca se actualiza al vender. Esto bloquea el flujo operativo core del restaurante.

## What Changes

- Implementación real de `AddProductToOrder` en el usecase de orden, incluyendo validación de estado, cálculo de precio y deducción de stock.
- Nueva interfaz `ProductInventoryRepository` en el usecase de orden para obtener precio de producto y receta técnica.
- Nueva transacción atómica `AddItemsWithInventory` en el repositorio de orden que inserta items, descuenta ingredientes del stock y actualiza el total de la orden en una sola TX.
- Error centinela `ErrInsufficientStock` propagado hasta HTTP 409 cuando un ingrediente no tiene suficiente stock.
- Inyección de `productRepo` como segundo argumento en `uorder.NewUsecase` en `cmd/server/main.go`.

## Capabilities

### New Capabilities

- `order-item-addition`: Agregar uno o más productos a una orden existente, con validación de estado, precio real desde BD y deducción atómica de ingredientes del stock.

### Modified Capabilities

- `pos-operations`: Se añaden requisitos de validación de stock e integridad transaccional al flujo de adición de items a una orden.

## Impact

- **Usecase**: `internal/usecase/order/usecase.go` — nueva interfaz y lógica real en `AddProductToOrder`.
- **Dominio**: `internal/domain/order/order.go` — nuevos tipos `RecipeLine` y `StockDeduction`.
- **Repositorio**: `internal/infrastructure/persistence/postgres/order_repository.go` — método `AddItemsWithInventory` con TX multi-tabla.
- **Servidor**: `cmd/server/main.go` — pasar `productRepo` al usecase.
- **Handler**: `internal/infrastructure/rest/order/handler.go` — mapear `ErrInsufficientStock` → HTTP 409.
- **Dependencias**: Sin cambios en `go.mod`.
