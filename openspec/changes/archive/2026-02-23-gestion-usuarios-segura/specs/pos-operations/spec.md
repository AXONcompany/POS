## ADDED Requirements

### Requirement: Gestión de Comandas por Meseros
El sistema DEBE permitir exclusivamente a los usuarios con el rol `MESERO` la creación, modificación y el borrado (si está configurado) de órdenes de consumo atadas a mesas o pedidos directos.

#### Scenario: Creación de orden exitosa
- **WHEN** el usuario autenticado tiene el rol `MESERO` y envía una nueva orden al backend de su restaurante correspondiente.
- **THEN** el sistema registra la orden y actualiza el estado de la mesa a "Ocupada".

### Requirement: Cobro y Cierre de Órdenes por Cajeros
El sistema DEBE restringir la finalización de órdenes (facturación, recibo de pagos) a los usuarios con rol `CAJERO`.

#### Scenario: Facturación de mesa
- **WHEN** el usuario autenticado tiene el rol `CAJERO` e intenta cerrar y pagar una orden existente.
- **THEN** el sistema liquida la orden, registra los pagos y emite el recibo electrónico de venta para ese restaurante.

#### Scenario: Intento de facturación por mesero
- **WHEN** un usuario con rol `MESERO` intenta cerrar el pago de una orden.
- **THEN** el sistema rechaza la operación por falta de privilegios (Forbidden).
