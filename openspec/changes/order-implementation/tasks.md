## 1. Base de Datos y Modelos

- [x] 1.1 Crear migración `000004_create_orders_table` (SQL de subida y bajada) para `estados_pedido`, `pedidos` y `pedido_productos`.
- [x] 1.2 Definir los structs de dominio en `internal/core/domain/order.go` (`Order`, `OrderStatus`, `OrderItem`).

## 2. Capa de Abstracción y Casos de Uso

- [x] 2.1 Definir la interfaz `OrderRepository` en `internal/core/ports/order_repository.go`.
- [x] 2.2 Implementar el repositorio en `internal/infrastructure/db/postgres/order_repository.go` cuidando el uso de transacciones.
- [x] 2.3 Implementar la capa de lógica en `internal/usecase/order/usecase.go` para coordinar alta, actualización y listado de pedidos.
- [x] 2.4 Agregar pruebas unitarias a `usecase.go` verificando los flujos de éxito y error.

## 3. Endpoints REST (API)

- [x] 3.1 Crear controlador `internal/infrastructure/rest/order/handler.go` con manejadores para crear pedidos, actualizar estado y listarlos por mesa.
- [x] 3.2 Modificar `cmd/server/routes.go` para exponer las rutas `/api/v1/orders` bajo el sistema de autenticación RBAC existente.
