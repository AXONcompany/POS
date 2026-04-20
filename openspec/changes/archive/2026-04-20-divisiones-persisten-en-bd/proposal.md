## Why

Actualmente la división de cuentas es un cálculo efímero que el handler devuelve al cliente pero nunca persiste. Esto impide vincular pagos parciales a divisiones específicas, verificar que todas las partes de una cuenta dividida estén saldadas antes de cerrar la orden, y tener trazabilidad de cómo se dividió una cuenta.

## What Changes

- **Nueva tabla `order_divisions`**: Persiste cada división de cuenta con su monto, impuesto, total y estado de pago (`is_paid`).
- **Nueva columna `payments.division_ref`**: Permite vincular un pago a una división concreta.
- **Lógica de cálculo migrada al usecase**: El handler `DivideOrder` pasa a ser un wrapper delegando a `usecase.DivideOrder`, que calcula, persiste y retorna las divisiones.
- **Nuevos métodos en repositorio**: `CreateDivisions`, `GetDivisionsByOrderID`, `MarkDivisionPaid`.
- **Migración 003**: Crea la tabla `order_divisions` y agrega la FK en `payments`.

## Capabilities

### New Capabilities
- `order-division`: Creación, consulta y marcado de pago de divisiones de cuenta vinculadas a una orden.

### Modified Capabilities
- `pos-operations`: El flujo de cobro se modifica para que un pago pueda referenciar una división; la validación de checkout tendrá en cuenta si la orden está dividida.

## Impact

- **Base de datos**: Nueva tabla `order_divisions`; nueva columna nullable `payments.division_ref`.
- **Repositorio**: Nuevos métodos en `order_repository.go`; `payment_repository.go` actualiza `Create` para aceptar `division_ref`.
- **Usecase**: `DivideOrder` en `usecase/order/usecase.go` implementado de verdad; `ProcessPayment` en `usecase/payment/usecase.go` acepta `division_ref` opcional.
- **Handler**: `DivideOrder` y `ProcessPayment` actualizados para pasar el nuevo campo.
- **Rutas**: Sin cambios en la estructura de rutas.
