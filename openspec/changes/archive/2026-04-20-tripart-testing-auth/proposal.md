## Why

El módulo de Autenticación (`Auth`) es el componente más crítico del sistema en términos de seguridad. Un fallo aquí expone todas las demás entidades (órdenes, productos, usuarios). Actualmente, las pruebas existentes se enfocan principalmente en el *Happy Path* (creación exitosa, login exitoso). Necesitamos asegurar la robustez del sistema sometiéndolo a una estrategia de pruebas tripartita: validación funcional básica, resistencia a casos borde, y protección contra manipulación malintencionada (ataques de seguridad, inyección, fuerza bruta en lógica).

## What Changes

- Implementación de un nuevo set de pruebas para el `usecase` de `auth`.
- Desarrollo de **Nivel 1 (Básicos)**: Login, registro, creación de sesión y verificación de token.
- Desarrollo de **Nivel 2 (Borde)**: Manejo de sesiones expiradas concurrentemente, revocación en carrera, múltiples logins simultáneos, manipulación de horas locales vs UTC.
- Desarrollo de **Nivel 3 (Malintencionados)**: Intentos de validación con tokens JWT manipulados (algoritmo modificado, sin firma), inyección de SQL en campos de login, e intentos de acceso a sesiones de otros usuarios (suplantación).

## Capabilities

### New Capabilities
- `auth-resilience-testing`: Suite de pruebas exhaustiva para validar la durabilidad y seguridad del sistema de autenticación bajo estrés y actores maliciosos.

### Modified Capabilities

## Impact

- **Código Afectado**: Se creará/modificará extensamente `internal/usecase/auth/usecase_test.go` o un paquete dedicado a pruebas de integración.
- **Sistemas**: Mejorará la confianza general en la capa de seguridad, reduciendo el riesgo de vulnerabilidades OWASP en el POS.
