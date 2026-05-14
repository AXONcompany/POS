## ADDED Requirements

### Requirement: Gestión de Comandas por Meseros
El sistema DEBE permitir exclusivamente a los usuarios con el rol `MESERO` la creación, modificación y el borrado de órdenes de consumo atadas a mesas o pedidos directos. Las órdenes se vinculan a una `venue_id` (extraída del JWT) y opcionalmente a un `pos_terminal_id`. Al agregar productos a una orden, el sistema DEBE validar la disponibilidad de stock de los ingredientes requeridos por la receta técnica de cada producto y descontarlos atómicamente.

#### Scenario: Creación de orden exitosa
- **WHEN** el usuario autenticado tiene el rol `MESERO` y envía una nueva orden al backend.
- **THEN** el sistema registra la orden con `venue_id` del JWT del mesero y opcionalmente la vincula a un `pos_terminal_id`.

#### Scenario: Adición de producto con stock disponible
- **WHEN** el usuario autenticado tiene el rol `MESERO` y agrega un producto a una orden en estado `PENDING` o `SENT`, y los ingredientes requeridos tienen stock suficiente
- **THEN** el sistema inserta el item con precio real, descuenta el stock de ingredientes y actualiza el total de la orden en una sola transacción

#### Scenario: Adición de producto con stock insuficiente
- **WHEN** el usuario autenticado tiene el rol `MESERO` y agrega un producto cuyos ingredientes no tienen stock suficiente
- **THEN** el sistema rechaza la operación con error `ErrInsufficientStock` (HTTP 409) sin modificar la orden ni el inventario

### Requirement: Cobro y Cierre de Órdenes por Cajeros
El sistema DEBE restringir la finalización de órdenes (facturación, recibo de pagos) a los usuarios con rol `CAJERO`. Los pagos se registran con `venue_id` y `pos_terminal_id` del cajero. Los pagos pueden referenciar opcionalmente una división de cuenta específica mediante `division_ref`.

#### Scenario: Facturación de mesa
- **WHEN** el usuario autenticado tiene el rol `CAJERO` e intenta cerrar y pagar una orden existente.
- **THEN** el sistema liquida la orden, registra los pagos con `venue_id` y `pos_terminal_id`, y emite el recibo de venta.

#### Scenario: Pago vinculado a una división
- **WHEN** el cajero registra un pago con un `division_ref` válido apuntando a una división de la orden
- **THEN** el sistema registra el pago, marca la división como `is_paid: true` y retorna confirmación

#### Scenario: Intento de facturación por mesero
- **WHEN** un usuario con rol `MESERO` intenta cerrar el pago de una orden.
- **THEN** el sistema rechaza la operación por falta de privilegios (Forbidden).
