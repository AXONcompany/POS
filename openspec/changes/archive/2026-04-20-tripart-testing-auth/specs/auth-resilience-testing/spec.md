## ADDED Requirements

### Requirement: Pruebas Nivel 1 (Básicos)
El sistema MUST validar exitosamente el comportamiento principal bajo el uso normal ("Happy path").

#### Scenario: Login correcto de usuario válido
- **WHEN** un usuario provee correo y contraseña correctos
- **THEN** el sistema valida la firma hash y expide los correspondientes pares de Access y Refresh tokens

#### Scenario: Registro correcto de usuario
- **WHEN** un gestor con los permisos adecuados, o el sistema de setup registran un nuevo usuario
- **THEN** la contraseña es guardada procesada por un algoritmo de Hash (no texto plano)

### Requirement: Pruebas Nivel 2 (Borde)
El sistema MUST responder predeciblemente y sin fallos no manejados ante datos inesperados y condiciones de estado raras.

#### Scenario: Envío de Refresh Token a milisegundos de expirar
- **WHEN** se solicita un nuevo Access Token enviando un Refresh Token que asombrosamente ha expirado justo en el transcurso de la petición
- **THEN** el sistema evalúa la caducidad local de manera segura y deniega el acceso sin provocar pánico

#### Scenario: Contraseñas que exceden límites normales
- **WHEN** se intenta hacer Login con una contraseña de 1 Megabyte de longitud
- **THEN** el sistema previene el ataque de carga computacional denegando el string por exceso de longitud en vez de proceder al hashing

### Requirement: Pruebas Nivel 3 (Adversariales)
El sistema MUST repeler vectores de ataque comunes sobre su diseño.

#### Scenario: Manipulación de payload cambiando Role o Restaurant
- **WHEN** un usuario con sesión válida captura un JWT, modifica sus claims internos para asignarse rol 'PROPIETARIO' y lo reenvía
- **THEN** la firma criptográfica en JWT debe ser inválida forzando un rechazo del Token

#### Scenario: Algorithm None Attack
- **WHEN** un atacante envía un token manipulado cambiando la cabecera del JWT a `{"alg": "none"}` para evadir la verificación
- **THEN** el sistema detecta explícitamente algoritmos de baja de seguridad y fuerza el cierre de la solicitud con Unauthorized
