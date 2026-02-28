## Context

El restaurante necesita un módulo para registrar y gestionar el ciclo de vida de los pedidos (orders). Actualmente existen módulos para catálogo de productos, mesas y ventas, pero falta la entidad central que vincule qué productos solicitó una mesa específica y quién (mesero) los está atendiendo. Este módulo es requerido para operar eficientemente y pasar de tomar la orden a preparación y finalmente a venta.

## Goals / Non-Goals

**Goals:**
- Implementar esquema de Base de Datos para `pedidos`, `pedido_productos` y `estados_pedido` según el diagrama ER provisto.
- Crear la capa de abstracción (UseCase, Repository) en Go para la gestión de órdenes.
- Exponer API REST para crear y administrar órdenes.
- Integrar transacciones para asegurar la integridad al momento de registrar productos en un pedido.

**Non-Goals:**
- Implementar la interfaz de usuario (Frontend) para las órdenes en este cambio.
- Lógica de facturación electrónica compleja (eso se delega al módulo de Ventas).

## Decisions

- **Patrón Arquitectónico**: Se seguirá utilizando Clean Architecture (`internal/infrastructure/rest`, `internal/usecase`, `internal/repository`).
- **Base de Datos**: Se crearán migraciones SQL (`000004_create_orders_table.up.sql`). La tabla `pedidos` tendrá las llaves foráneas correspondientes.
- **Relación con Mesero**: El diagrama asocia el mesero con la mesa a través de `MesaMesero`, por lo tanto el pedido, al estar asociado a la mesa, puede derivar el mesero o bien incluir un `mesero_id` opcional para auditoría si varios meseros atienden una zona. Añadiremos un campo directo a `pedidos` si es requerido por contexto operativo.
- **Manejo de Transacciones**: Como la inserción en `pedidos` y múltiples inserciones en `pedido_productos` (y potencial deducción de stock) deben ocurrir juntas, se usará `tx` en SQL.

## Risks / Trade-offs

- **Risk**: Condiciones de carrera (race conditions) si dos meseros modifican el mismo pedido simultáneamente. 
  - **Mitigation**: Implementar un sistema de versionado optimista o bloquear la fila en base de datos al actualizar la orden si es estrictamente necesario, o mantener operaciones atómicas para los estados.
