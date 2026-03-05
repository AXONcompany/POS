## MODIFIED Requirements

### Requirement: Gestion de Comandas por Meseros
El sistema DEBE permitir exclusivamente a los usuarios con el rol `MESERO` la creacion, modificacion y el borrado de ordenes de consumo atadas a mesas o pedidos directos. Las ordenes se vinculan a una `venue_id` (en lugar de `restaurant_id`) y opcionalmente a un `pos_terminal_id`.

#### Scenario: Creacion de orden exitosa
- **WHEN** el usuario autenticado tiene el rol `MESERO` y envia una nueva orden al backend.
- **THEN** el sistema registra la orden con `venue_id` del JWT del mesero, actualiza el estado de la mesa a "Ocupada", y opcionalmente vincula la orden a un `pos_terminal_id`.

### Requirement: Cobro y Cierre de Ordenes por Cajeros
El sistema DEBE restringir la finalizacion de ordenes (facturacion, recibo de pagos) a los usuarios con rol `CAJERO`. Los pagos se registran con `venue_id` y `pos_terminal_id` del cajero.

#### Scenario: Facturacion de mesa
- **WHEN** el usuario autenticado tiene el rol `CAJERO` e intenta cerrar y pagar una orden existente.
- **THEN** el sistema liquida la orden, registra los pagos con `venue_id` y `pos_terminal_id`, y emite el recibo de venta.

#### Scenario: Intento de facturacion por mesero
- **WHEN** un usuario con rol `MESERO` intenta cerrar el pago de una orden.
- **THEN** el sistema rechaza la operacion por falta de privilegios (Forbidden).
