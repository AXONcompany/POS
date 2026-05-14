## 1. Usecase: error paths faltantes

- [x] 1.1 Agregar sub-test `AddProductToOrder` en estado SENT (status_id=2): debe ser permitido
- [x] 1.2 Agregar sub-test `AddProductToOrder` con error en `GetProductPrice`
- [x] 1.3 Agregar sub-test `AddProductToOrder` con error en `GetRecipeLines`
- [x] 1.4 Agregar sub-test `AddProductToOrder` con producto sin receta (recipe vacía, sin deducciones)
- [x] 1.5 Agregar sub-test `CancelOrderItem` con orden en estado CANCELLED (status_id=6): debe retornar `ErrInvalidStatusTransition`
- [x] 1.6 Agregar sub-test `CancelOrderItem` con error en `GetOrderItem` (repo error)
- [x] 1.7 Agregar sub-test `CancelOrderItem` con error en `CancelItemWithInventoryRestore`
- [x] 1.8 Agregar sub-test `CancelOrderItem` con error en `auditRepo.SaveAudit`

## 2. Usecase: DivideOrder escenarios faltantes

- [x] 2.1 Agregar sub-test `DivideOrder` tipo `por_monto` calcula subtotal/impuesto por monto dado
- [x] 2.2 Agregar sub-test `DivideOrder` tipo `por_item` calcula partes iguales por número de elementos
- [x] 2.3 Agregar sub-test `DivideOrder` con tipo inválido retorna error
- [x] 2.4 Agregar sub-test `DivideOrder` con error en `repo.GetByID`

## 3. Usecase: GetDivisionsByOrder

- [x] 3.1 Agregar `TestGetDivisionsByOrder` con sub-test de éxito: retorna divisiones del repo
- [x] 3.2 Agregar sub-test de error de repo: propaga el error

## 4. Handler: setup del archivo de tests

- [x] 4.1 Crear `internal/infrastructure/rest/order/handler_test.go` con la interfaz `mockOrderUsecase` y su stub
- [x] 4.2 Implementar helper `setupRouter(stub)` que registra las rutas del handler con el stub inyectado y los middleware keys necesarios (venueID, userID) en el context de gin

## 5. Handler: tests HTTP

- [x] 5.1 Test `CreateOrder` retorna 201 con body JSON en happy path
- [x] 5.2 Test `CreateOrder` retorna 500 si el usecase retorna error
- [x] 5.3 Test `AddItems` retorna 409 si usecase retorna `ErrInsufficientStock`
- [x] 5.4 Test `AddItems` retorna 422 si usecase retorna `ErrInvalidStatusTransition`
- [x] 5.5 Test `CancelItem` retorna 409 si usecase retorna `ErrItemAlreadyCancelled`
- [x] 5.6 Test `CancelItem` retorna 422 si usecase retorna `ErrInvalidStatusTransition`
- [x] 5.7 Test `DivideOrder` retorna 200 con array de divisiones en happy path
- [x] 5.8 Test `DivideOrder` retorna 409 si usecase retorna `ErrDivisionAlreadyPaid`
- [x] 5.9 Test `GetDivisions` retorna 200 con array de divisiones
- [x] 5.10 Test `CheckoutOrder` retorna 422 si usecase retorna `ErrInvalidStatusTransition`
