## Context

El módulo de Autenticación (`Auth`) es vital para la seguridad general del Punto de Venta (POS). Maneja sesiones, validación de contraseñas (hash) y generación de tokens JWT con claims de configuración de entorno y rol de usuario.
Dado su peso sobre toda la operación (solo un cajero o mesero autorizado puede interactuar con el sistema dependiente del endpoint), requerimos una capa de pruebas rigurosa dividida en tres estratos (tripartición).

## Goals / Non-Goals

**Goals:**
- **Nivel 1 (Básicos)**: Validar el camino feliz (`Happy Path`) del inicio de sesión (Login) y refresco de tokens.
- **Nivel 2 (Borde)**: Validar escenarios inusuales como caducidades de tokens justo al ser usados, peticiones asíncronas masivas concurrentes intentando crear/actualizar sesión.
- **Nivel 3 (Adversariales)**: Identificar vulnerabilidades como el intento de usar algoritmos "none" en JWT (`alg: none`), manipulación de claims, o `Timing Attacks` en la validación de contraseñas.
- Establecer un estándar de pruebas basado en Table-Driven Tests.

**Non-Goals:**
- Cambiar la librería criptográfica existente (e.g. `jwt-go` a otra o `bcrypt` a `argon2`), el objetivo es *probar* lo que ya existe. Si la prueba falla, el fix pertinente será abordado en otro "change".

## Decisions

- **Framework**: Se usará el marco de pruebas estándar de Go `testing`.
- **Estructura**: `table-driven tests` para enviar un gran abanico de inputs maliciosos de forma concisa.
- **Mocks Strictos**: Mocks personalizados o base que permitan emular respuestas de base de datos tardías o manipuladas (ej. para probar timeouts).
- **Separación de Archivos**: Agrupación lógica dentro de `internal/usecase/auth/usecase_test.go` en tres bloques grandes (Basic, Edge, Adversarial) o en funciones correspondientes `TestUsecase_Login_Basic`, `TestUsecase_Login_Edge`, `TestUsecase_Login_Adversarial`.

## Risks / Trade-offs

- **Risk**: Complejidad y lentitud en los tests si se implementan esperas asíncronas reales.
  - **Mitigation**: Paralelización controlada, usar `time` inyectado (mock) siempre que sea posible o duraciones muy cortas para pruebas de timeout.
