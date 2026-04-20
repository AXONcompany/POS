## 1. Base de datos y dominio

- [x] 1.1 Crear migración `000014_add_order_divisions` con la tabla `order_divisions` (id VARCHAR(50) PK, order_id, venue_id, division_type, amount, tax, total, is_paid BOOLEAN, created_at) y FK de `payments.division_id` → `order_divisions.id`.
- [x] 1.2 Agregar struct `OrderDivision` a `internal/domain/order/order.go`.
- [x] 1.3 Agregar error centinela `ErrDivisionAlreadyPaid` a `internal/domain/order/order.go`.

## 2. Repositorio de ordenes

- [x] 2.1 Agregar `CreateDivisions(ctx, divisions []domainOrder.OrderDivision) error` a `order_repository.go`: borra divisiones previas sin pagos y luego inserta las nuevas en una TX; retorna `ErrDivisionAlreadyPaid` si hay pagos vinculados.
- [x] 2.2 Agregar `GetDivisionsByOrderID(ctx, orderID int64, venueID int) ([]domainOrder.OrderDivision, error)` a `order_repository.go`.
- [x] 2.3 Agregar `CreateDivisions` y `GetDivisionsByOrderID` a la interfaz `Repository` en `internal/usecase/order/usecase.go`.


## 3. Repositorio de pagos

- [x] 3.1 Convertir `PaymentRepository.Create` en `order_repository.go` a TX que, si `division_id != nil`, también ejecuta `UPDATE order_divisions SET is_paid = true WHERE id = $1` dentro de la misma transacción.

## 4. Usecase de ordenes

- [x] 4.1 Implementar `DivideOrder(ctx, venueID int, orderID int64, divisionType string, numParts int, customAmounts []float64) ([]domainOrder.OrderDivision, error)`: calcula subtotal/impuesto/total por parte según el tipo, genera IDs deterministas (`div_<orderID>_<i>`), llama a `repo.CreateDivisions` y retorna las divisiones.
- [x] 4.2 Implementar `GetDivisionsByOrder(ctx, venueID int, orderID int64) ([]domainOrder.OrderDivision, error)` en el usecase.

## 5. Handler y rutas

- [x] 5.1 Reemplazar la lógica de cálculo inline de `DivideOrder` en `internal/infrastructure/rest/order/handler.go` con una llamada al usecase; adaptar el request/response al nuevo tipo `OrderDivision`.
- [x] 5.2 Agregar handler `GetDivisions` en `handler.go` que llame a `uc.GetDivisionsByOrder`.
- [x] 5.3 Mapear `ErrDivisionAlreadyPaid` → HTTP 409 en el handler `DivideOrder`.
- [x] 5.4 Registrar ruta `GET /ordenes/:id/divisiones` en `internal/infrastructure/rest/routes.go`.

## 6. Pruebas

- [x] 6.1 Test: `DivideOrder` en partes iguales calcula y persiste correctamente.
- [x] 6.2 Test: re-división sin pagos vinculados reemplaza divisiones previas.
- [x] 6.3 Test: re-división con pago vinculado retorna `ErrDivisionAlreadyPaid`.
