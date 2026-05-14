## MODIFIED Requirements

### Requirement: Autenticacion Segura (Secure Login)
El sistema DEBE permitir a los usuarios (Cajeros, Meseros) iniciar sesion de forma segura utilizando correo y contrasena encriptada, emitiendo tokens de acceso (JWT). El JWT DEBE incluir `venue_id` en lugar de `restaurant_id`. Los Propietarios (Owners) se autentican por separado a traves de la tabla `owners`.

#### Scenario: Inicio de sesion exitoso
- **WHEN** un usuario provee credenciales validas registradas en el sistema.
- **THEN** el sistema autentica al usuario y devuelve un token de acceso con claims `sub`, `email`, `role_id`, `venue_id` y un refresh token HttpOnly.

#### Scenario: Fallo de login por credenciales invalidas
- **WHEN** un usuario provee una contrasena incorrecta o correo no registrado.
- **THEN** el sistema deniega el acceso con un mensaje generico ("Credenciales invalidas") para evitar enumeracion de usuarios.

### Requirement: Control de Acceso Basado en Roles (RBAC)
El sistema DEBE verificar el rol de cada usuario antes de conceder acceso a cualquier endpoint protegido. Los roles operativos son `CAJERO` y `MESERO`. El rol `PROPIETARIO` se maneja desde la entidad `owners`. El middleware extrae `venue_id` del JWT para filtrar datos.

#### Scenario: Acceso denegado a recurso no autorizado
- **WHEN** un usuario con rol `MESERO` intenta acceder al panel de administracion.
- **THEN** el sistema bloquea la peticion, devolviendo un error 403 Forbidden.

#### Scenario: Acceso concedido a recurso autorizado
- **WHEN** un usuario con rol `CAJERO` intenta acceder a los comandos de cobro.
- **THEN** el sistema valida su rol y `venue_id` a traves del middleware y permite la peticion.

## ADDED Requirements

### Requirement: Correccion del Flujo RefreshToken
El sistema DEBE renovar tokens de acceso sin re-validar el password. El flujo de RefreshToken genera nuevos tokens directamente a partir del usuario almacenado, sin llamar a Login con el hash.

#### Scenario: Renovacion exitosa de token
- **WHEN** un usuario envia un refresh_token valido y no expirado.
- **THEN** el sistema revoca la sesion anterior, genera un nuevo access_token y refresh_token, y crea una nueva sesion.

#### Scenario: Refresh con token expirado
- **WHEN** un usuario envia un refresh_token expirado o revocado.
- **THEN** el sistema rechaza la peticion con error 401.

### Requirement: Consistencia de Nombres de Roles
El sistema DEBE usar los mismos nombres de rol en el handler REST y en la BD: `PROPIETARIO`, `CAJERO`, `MESERO`. Se elimina la inconsistencia donde el handler usaba `ADMIN`/`CAJA`.

#### Scenario: Registro con nombres de rol correctos
- **WHEN** se registra un usuario con rol `CAJERO`.
- **THEN** el sistema lo acepta y asigna `role_id = 2`.
