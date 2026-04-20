## ADDED Requirements

### Requirement: División de cuenta persistida
El sistema DEBE permitir dividir el total de una orden en partes (iguales, por monto o por ítem) y persistir cada división en la base de datos, vinculándola a la orden y marcando si ha sido pagada.

#### Scenario: División en partes iguales exitosa
- **WHEN** un usuario autorizado solicita dividir una orden en N partes iguales
- **THEN** el sistema calcula el subtotal, impuestos y total de cada parte, persiste N registros en `order_divisions` y retorna la lista de divisiones con sus IDs y estado `is_paid: false`

#### Scenario: División por monto exitosa
- **WHEN** un usuario autorizado envía una lista de montos para dividir la cuenta
- **THEN** el sistema persiste una división por cada monto especificado con su desglose de impuestos y retorna la lista de divisiones

#### Scenario: Re-división de orden sin pagos vinculados
- **WHEN** un usuario re-divide una orden que ya tiene divisiones pero ninguna ha sido pagada
- **THEN** el sistema elimina las divisiones previas e inserta las nuevas en una única transacción atómica, retornando las divisiones actualizadas

#### Scenario: Intento de re-división con pago vinculado
- **WHEN** un usuario intenta re-dividir una orden que tiene al menos una división con un pago asociado
- **THEN** el sistema rechaza la operación con HTTP 409 (`ErrDivisionAlreadyPaid`) sin modificar las divisiones existentes

#### Scenario: Consulta de divisiones activas
- **WHEN** un usuario autorizado consulta las divisiones de una orden
- **THEN** el sistema retorna la lista de divisiones de esa orden con su estado de pago actual (`is_paid`)
