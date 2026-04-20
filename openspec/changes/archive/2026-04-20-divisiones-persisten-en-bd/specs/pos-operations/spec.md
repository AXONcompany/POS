## MODIFIED Requirements

### Requirement: Cobro y Cierre de Órdenes por Cajeros
El sistema DEBE restringir la finalización de órdenes (facturación, recibo de pagos) a los usuarios con rol `CAJERO`. Los pagos pueden referenciar opcionalmente una división de cuenta específica mediante `division_ref`.

#### Scenario: Facturación de mesa
- **WHEN** el usuario autenticado tiene el rol `CAJERO` e intenta cerrar y pagar una orden existente.
- **THEN** el sistema liquida la orden, registra los pagos y emite el recibo electrónico de venta para ese restaurante.

#### Scenario: Pago vinculado a una división
- **WHEN** el cajero registra un pago con un `division_ref` válido apuntando a una división de la orden
- **THEN** el sistema registra el pago, marca la división como `is_paid: true` y retorna confirmación

#### Scenario: Intento de facturación por mesero
- **WHEN** un usuario con rol `MESERO` intenta cerrar el pago de una orden.
- **THEN** el sistema rechaza la operación por falta de privilegios (Forbidden).
