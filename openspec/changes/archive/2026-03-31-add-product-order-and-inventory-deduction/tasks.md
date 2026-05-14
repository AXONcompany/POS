## 1. Tipos de dominio

- [x] 1.1 Agregar `RecipeLine{ IngredientID int64, QuantityRequired float64 }` a `internal/domain/order/order.go`
- [x] 1.2 Agregar `StockDeduction{ IngredientID int64, VenueID int, Quantity float64 }` a `internal/domain/order/order.go`
- [x] 1.3 Agregar error centinela `ErrInsufficientStock` a `internal/domain/order/order.go`

## 2. Interfaz y usecase

- [x] 2.1 Definir interfaz `ProductInventoryRepository` en `internal/usecase/order/usecase.go` con métodos `GetProductPrice` y `GetRecipeLines`
- [x] 2.2 Agregar campo `invRepo ProductInventoryRepository` a `Usecase` y actualizar `NewUsecase` para recibirlo como segundo argumento
- [x] 2.3 Implementar `AddProductToOrder`: obtener orden, validar estado PENDING/SENT, iterar items para obtener precio y receta, acumular deducciones por ingrediente, llamar a `repo.AddItemsWithInventory`

## 3. Repositorio de orden

- [x] 3.1 Agregar método `AddItemsWithInventory(ctx, orderID int64, venueID int, items []domainOrder.OrderItem, deductions []domainOrder.StockDeduction) error` a `internal/infrastructure/persistence/postgres/order_repository.go`
- [x] 3.2 Implementar la TX dentro de `AddItemsWithInventory`: INSERT batch en `order_items`, UPDATE de stock con guard `AND stock >= $qty` (0 rows → rollback + `ErrInsufficientStock`), UPDATE de `orders.total_amount`
- [x] 3.3 Agregar `AddItemsWithInventory` a la interfaz `Repository` en `internal/usecase/order/usecase.go`

## 4. Repositorio de productos (implementación de ProductInventoryRepository)

- [x] 4.1 Agregar método `GetProductPrice(ctx, productID int64, venueID int) (float64, error)` al repositorio de productos existente
- [x] 4.2 Agregar método `GetRecipeLines(ctx, productID int64) ([]domainOrder.RecipeLine, error)` al repositorio de productos existente

## 5. Servidor y handler

- [x] 5.1 Pasar `productRepo` como segundo argumento a `uorder.NewUsecase(orderRepo, productRepo)` en `cmd/server/main.go`
- [x] 5.2 Mapear `ErrInsufficientStock` → HTTP 409 en el handler `AddItems` de `internal/infrastructure/rest/order/handler.go`

## 6. Tests

- [x] 6.1 Agregar `GetStatusByID`, `AddItemsWithInventory` al `MockRepository` en `usecase_test.go`
- [x] 6.2 Crear `MockProductInventoryRepository` con `GetProductPrice` y `GetRecipeLines`
- [x] 6.3 Test: adición exitosa con deducción de stock
- [x] 6.4 Test: rechazo cuando la orden no está en estado editable (`ErrInvalidStatusTransition`)
- [x] 6.5 Test: propagación de `ErrInsufficientStock` desde el repositorio
- [x] 6.6 Test: acumulación correcta de deducciones cuando dos items comparten ingrediente
