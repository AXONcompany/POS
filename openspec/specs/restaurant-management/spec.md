## ADDED Requirements

### Requirement: Administración de Entidad de Restaurante
El sistema DEBE permitir a los usuarios con el rol `PROPIETARIO` gestionar el licenciamiento, la información del perfil del negocio y la configuración estructural (impuestos, monedas).

#### Scenario: Actualización de datos del restaurante
- **WHEN** el usuario autenticado tiene el rol `PROPIETARIO` y envía cambios al perfil del restaurante (nombre, dirección).
- **THEN** el sistema actualiza la información que se reflejará en todos los recibos y operaciones del tenant.

### Requirement: Control de Empleados del Restaurante
El sistema DEBE proveer un panel para que el `PROPIETARIO` dé de alta, suspenda o modifique las cuentas de sus `MESEROS` y `CAJEROS`. 

#### Scenario: Propietario registra un nuevo mesero
- **WHEN** el `PROPIETARIO` ingresa los datos y asigna la contraseña temporal (hasheada) a un nuevo empleado designándole el rol `MESERO`.
- **THEN** el sistema guarda el nuevo usuario forzando la asociación mediante el `restaurant_id` del propietario.

#### Scenario: Propietario suspende a un empleado
- **WHEN** el `PROPIETARIO` marca a un empleado existente como "Inactivo".
- **THEN** el sistema revoca inmediatamente los tokens del empleado e impide futuros inicios de sesión del mismo.
