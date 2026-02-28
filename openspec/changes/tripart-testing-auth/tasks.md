## 1. Setup Base y Nivel 1 (Happy Path)

- [x] 1.1 Preparar el esqueleto `TestUsecase_Login_Basic` en `internal/usecase/auth/usecase_test.go`.
- [x] 1.2 Implementar los `mocks` de `userRepo` y `sessionRepo` si no están aún en su estructura completa para Testing.
- [x] 1.3 Desarrollar las pruebas de Login exitoso, verificando que devuelve `AccessToken` y `RefreshToken`.

## 2. Nivel 2 (Edge / Chaos Path)

- [x] 2.1 Preparar la estructura Table-Driven `TestUsecase_Login_Edge` para evaluar casos extremos.
- [x] 2.2 Añadir una prueba de límite de contraseña (e.g. payload inmensamente largo superior al byte limit real del bcrypt).
- [x] 2.3 Añadir caso donde el Token de Refresco tiene expiración al milisegundo límite.

## 3. Nivel 3 (Adversariales / OWASP)

- [x] 3.1 Inicializar `TestUsecase_Token_Adversarial` específicamente enfocado en validación JWT.
- [x] 3.2 Implementar test intentando usar el algoritmo `"none"` (`{"alg":"none"}`) simulando un atacante.
- [x] 3.3 Implementar test de manipulación de Payload: Crear un token válido, parsear las claims por debajo, cambiar su Rol a admin/owner, recomponer pero con firma falsificada o usando un default secret para comprobar el rechazo nativo.
