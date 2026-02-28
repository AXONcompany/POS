## Why

El restaurante requiere un sistema para gestionar las órdenes (pedidos) de los clientes, vinculando los productos solicitados con las mesas y el personal (meseros) que los atiende. Esto es fundamental para el funcionamiento del punto de venta y para mantener un control preciso de las ventas y la preparación de los platillos.

## What Changes

- Creación de las tablas de base de datos para `pedidos`, `pedido_productos`, `estados_pedido` y la relación con `Mesa`, `Venta` y `Mesero` (a través de `MesaMesero`).
- Implementación de la lógica de negocio (Use Cases) para crear, actualizar, listar y cancelar pedidos.
- Desarrollo de los endpoints REST para la gestión de las órdenes.
- Pruebas unitarias y de integración para el nuevo módulo de órdenes.

## Capabilities

### New Capabilities
- `order-management`: Creación, actualización, listado y gestión del ciclo de vida de los pedidos, incluyendo su relación con productos, mesas y meseros.

### Modified Capabilities

## Impact

- **Base de Datos**: Nuevas migraciones para el esquema de órdenes.
- **Backend API**: Nuevos endpoints bajo `/api/v1/orders`.
- **Sistemas**: Se integra con el módulo existente de productos, mesas y ventas.
