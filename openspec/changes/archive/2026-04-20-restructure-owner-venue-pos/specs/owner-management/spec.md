## ADDED Requirements

### Requirement: Registro y Autenticacion de Propietarios
El sistema DEBE permitir crear cuentas de propietario (Owner) con email unico y password hasheado con bcrypt. Los propietarios se almacenan en la tabla `owners`, separada de los usuarios operativos.

#### Scenario: Creacion exitosa de owner
- **WHEN** se inserta un nuevo owner con email unico y password hasheado.
- **THEN** el sistema crea el registro en `owners` con `is_active = true` y timestamps automaticos.

#### Scenario: Email duplicado de owner
- **WHEN** se intenta crear un owner con un email que ya existe en `owners`.
- **THEN** el sistema rechaza la operacion con error de constraint unique.

### Requirement: Gestion Multi-Sede por Propietario
El sistema DEBE permitir que un propietario posea multiples sedes (venues). Cada venue tiene un FK `owner_id` a la tabla `owners`.

#### Scenario: Owner crea una nueva sede
- **WHEN** un owner registra una nueva venue con nombre y direccion.
- **THEN** el sistema crea la venue con `owner_id` del propietario y `is_active = true`.

#### Scenario: Listado de sedes por owner
- **WHEN** se consultan las venues con un `owner_id` especifico.
- **THEN** el sistema retorna solo las venues que pertenecen a ese owner.
