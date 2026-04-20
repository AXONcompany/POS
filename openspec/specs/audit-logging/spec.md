## ADDED Requirements

### Requirement: Registro de auditoría con snapshots
El sistema DEBE registrar cada cambio crítico en las entidades (órdenes, ítems de orden, productos, etc.) capturando el estado anterior y posterior en formato JSONB para permitir la trazabilidad y recuperación de datos.

#### Scenario: Registro exitoso de cambio
- **WHEN** un usuario realiza una acción de modificación (ej: cancelar ítem)
- **THEN** el sistema guarda en la tabla `audit_log` el `old_value` (estado antes del cambio), el `new_value` (estado después del cambio), el `user_id`, la `action` y el `venue_id` en la misma transacción que el cambio original.

#### Scenario: Falla en el registro de auditoría
- **WHEN** ocurre un error al intentar persistir el log de auditoría
- **THEN** el sistema DEBE realizar un rollback de toda la operación (incluyendo el cambio en la entidad original) para garantizar la consistencia.
