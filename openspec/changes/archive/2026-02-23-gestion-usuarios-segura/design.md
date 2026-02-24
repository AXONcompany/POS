## Context

El sistema POS actual necesita una actualización estructural para soportar múltiples roles de usuario con niveles de acceso estrictamente delimitados. Se requiere diferenciar entre la operativa interna del restaurante (meseros y cajeros) y la administración comercial/licenciamiento (propietarios). Esto requiere una arquitectura segura ("blindada") para prevenir accesos no autorizados a datos financieros y configuración de licencias.

## Goals / Non-Goals

**Goals:**
- Implementar Autenticación Segura mediante JWT (JSON Web Tokens) o mecanismo equivalente estándar de la industria.
- Establecer un Control de Acceso Basado en Roles (RBAC) robusto en el backend.
- Diseñar un modelo de datos multi-tenant donde la entidad `Restaurant` aisle la información de cada negocio.
- Asegurar que el rol `Propietario` gestione el restaurante (pagos de licencia, métricas globales) sin afectar la sesión de POS operativa activa.

**Non-Goals:**
- No se construirá un sistema de recursos humanos completo (nóminas, horarios complejos).
- No se implementará SSO (Single Sign-On) con proveedores externos (Google, Facebook) en esta primera fase.

## Decisions

- **Mecanismo de Autenticación**: JWT con "Short-Lived Access Tokens" (15m) y "HttpOnly Refresh Tokens" (7d) para mitigar el riesgo de robo de tokens (XSS) mientras se mantiene la usabilidad.
- **Hashing de Contraseñas**: Uso de `bcrypt` o `argon2` con salting automático provisto por la librería para un almacenamiento seguro de credenciales.
- **Autorización (RBAC)**: Middleware global en las rutas protegidas que verifique el rol del usuario (Propietario, Cajero, Mesero) antes de procesar el controlador.
- **Modelo de Licenciamiento**: Los propietarios estarán vinculados a una o más entidades `Restaurant`. Los usuarios internos (meseros, cajeros) existirán únicamente dentro de un `Restaurant` específico (aislamiento de datos usando `restaurant_id` en todas las consultas operativas).

## Risks / Trade-offs

- **[Riesgo de fuga de datos entre restaurantes]** → **Mitigación**: Implementar un tenant-resolver en el middleware que inyecte el `restaurant_id` en el contexto de la solicitud. Todas las consultas a la BD deben filtrar obligatoriamente por este ID.
- **[Aumento de latencia por validación de tokens]** → **Mitigación**: Firmar el estado del rol directamente en el JWT, para evitar consultas a la base de datos en cada endpoint, validando solo la firma y expiración criptográficamente.
- **[Complejidad de sesión para usuarios en múltiples dispositivos]** → **Mitigación**: Mantener un registro de sesiones activas (Refresh Tokens almacenados en BD con hash) para permitir la revocación remota.
