## MODIFIED Requirements

### Requirement: Gestión de Comandas por Meseros

El sistema DEBE permitir exclusivamente a los usuarios con el rol `MESERO` la creación, modificación y el borrado (si está configurado) de órdenes de consumo atadas a mesas o pedidos directos. Al agregar productos a una orden, el sistema DEBE validar la disponibilidad de stock de los ingredientes requeridos por la receta técnica de cada producto y descontarlos atómicamente.

#### Scenario: Creación de orden exitosa
- **WHEN** el usuario autenticado tiene el rol `MESERO` y envía una nueva orden al backend de su restaurante correspondiente.
- **THEN** el sistema registra la orden y actualiza el estado de la mesa a "Ocupada".

#### Scenario: Adición de producto con stock disponible
- **WHEN** el usuario autenticado tiene el rol `MESERO` y agrega un producto a una orden en estado `PENDING` o `SENT`, y los ingredientes requeridos tienen stock suficiente
- **THEN** el sistema inserta el item con precio real, descuenta el stock de ingredientes y actualiza el total de la orden en una sola transacción

#### Scenario: Adición de producto con stock insuficiente
- **WHEN** el usuario autenticado tiene el rol `MESERO` y agrega un producto cuyos ingredientes no tienen stock suficiente
- **THEN** el sistema rechaza la operación con error `ErrInsufficientStock` (HTTP 409) sin modificar la orden ni el inventario
