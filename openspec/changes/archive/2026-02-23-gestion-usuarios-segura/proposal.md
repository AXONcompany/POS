## Why

El sistema requiere soportar diferentes tipos de usuarios (meseros, cajeros y propietarios) con distintos niveles de acceso y responsabilidades. Es fundamental implementar mecanismos de seguridad robustos ("blindados") para proteger las operaciones financieras del POS (meseros/cajeros) y la administración de licenciamiento y datos del restaurante (propietario).

## What Changes

- Implementación de un sistema de autenticación seguro y control de acceso basado en roles (RBAC).
- Creación de modelos de datos para usuarios, roles, sesiones y restaurantes (licenciamiento/multi-tenant).
- Flujos de operación interna para el POS (lógica de negocio para meseros y cajeros).
- Panel de administración externa para el propietario, enfocado en el licenciamiento y configuración del restaurante.

## Capabilities

### New Capabilities
- `user-auth`: Autenticación segura, gestión de sesiones y control de acceso basado en roles (RBAC).
- `pos-operations`: Operaciones del sistema de punto de venta específicas para los roles de mesero y cajero.
- `restaurant-management`: Gestión de licenciamiento, suscripción y datos del restaurante exclusivo para el rol de propietario.

### Modified Capabilities


## Impact

- **Backend**: Integración de middlewares de autenticación/autorización seguro.
- **Base de Datos**: Nuevos esquemas/tablas para la gestión de usuarios, roles, permisos y licenciamiento por restaurante.
