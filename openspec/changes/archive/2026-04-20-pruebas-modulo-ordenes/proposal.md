## Why

El módulo de órdenes es el núcleo operativo del POS y ya tiene lógica de negocio compleja (deducción de inventario, cancelación con restauración de stock, división de cuenta, auditoría). Sin embargo, la cobertura de tests tiene huecos significativos: varios error paths del usecase no están cubiertos, `GetDivisionsByOrder` no tiene ningún test, y no existe ningún test a nivel HTTP handler — lo que significa que los mapeos de errores a códigos HTTP (409, 422, 500) nunca han sido verificados.

## What Changes

- Agregar sub-tests faltantes en `usecase_test.go` para cubrir error paths de `CancelOrderItem`, `DivideOrder` (tipos `por_monto`, `por_item`, tipo inválido, error de repo) y `GetDivisionsByOrder`
- Agregar sub-tests para `AddProductToOrder` en estado SENT (permitido) y errores de precio/receta
- Crear `internal/infrastructure/rest/order/handler_test.go` con tabla de tests HTTP para los handlers principales: `CreateOrder`, `AddItems`, `CancelItem`, `DivideOrder`, `GetDivisions`, `CheckoutOrder`

## Capabilities

### New Capabilities

- `order-handler-test-coverage`: Tests de capa HTTP del módulo de órdenes que verifican que los handlers mapean correctamente los errores de dominio a códigos HTTP

### Modified Capabilities

_(ninguna — no se modifican requisitos funcionales existentes)_

## Impact

- Solo archivos de test: `internal/usecase/order/usecase_test.go` (ampliar) y `internal/infrastructure/rest/order/handler_test.go` (nuevo)
- Sin cambios en código de producción
- Agrega dependencia de test: `net/http/httptest` (stdlib, ya disponible)
