## ADDED Requirements

### Requirement: Autenticación Segura (Secure Login)
El sistema DEBE permitir a los usuarios (Cajeros, Meseros, Propietarios) iniciar sesión de forma segura utilizando correo/usuario y contraseña encriptada, emitiendo tokens de acceso (JWT).

#### Scenario: Inicio de sesión exitoso
- **WHEN** un usuario provee credenciales válidas registradas en el sistema.
- **THEN** el sistema autentica al usuario y devuelve un token de acceso temporal y un token de actualización seguro (HttpOnly).

#### Scenario: Fallo de login por credenciales inválidas
- **WHEN** un usuario provee una contraseña incorrecta o correo no registrado.
- **THEN** el sistema deniega el acceso con un mensaje genérico (e.g. "Credenciales inválidas") para evitar enumeración de usuarios.

### Requirement: Control de Acceso Basado en Roles (RBAC)
El sistema DEBE verificar el rol de cada usuario antes de conceder acceso a cualquier endpoint protegido. Existen tres roles predeterminados: `PROPIETARIO`, `CAJERO` y `MESERO`.

#### Scenario: Acceso denegado a recurso no autorizado
- **WHEN** un usuario con rol `MESERO` intenta acceder al panel de administración de licencias.
- **THEN** el sistema bloquea la petición, devolviendo un error 403 Forbidden.

#### Scenario: Acceso concedido a recurso autorizado
- **WHEN** un usuario con rol `CAJERO` intenta acceder a los comandos de cobro.
- **THEN** el sistema valida su rol a través del middleware y permite la petición.
