## 1. Configuración Base de Datos y Modelos

- [x] 1.1 Crear script de migración para las tablas de `restaurants`, `users`, `roles` y `sessions`.
- [x] 1.2 Implementar los modelos en Go correspondientes a las nuevas tablas.
- [x] 1.3 Crear repositorios para la gestión de usuarios y restaurantes en el backend.

## 2. Sistema de Autenticación (user-auth)

- [x] 2.1 Implementar lógica de hashing de contraseñas (`bcrypt`/`argon2`) en el registro y creación de usuarios.
- [x] 2.2 Crear endpoint `/api/auth/login` que genere y retorne JWT Access Token y Refresh Token.
- [x] 2.3 Implementar middleware de validación JWT para proteger endpoints privados.
- [x] 2.4 Extender el middleware de autenticación para realizar Autorización (RBAC) basándose en los roles del JWT.

## 3. Operations de Propietario (restaurant-management)

- [x] 3.1 Crear endpoints para que el `PROPIETARIO` pueda gestionar su perfil y los datos del restaurante (`/api/restaurants`).
- [x] 3.2 Crear CRUD de empleados (`/api/users`) restringido al rol `PROPIETARIO`, que asocie los usuarios creados al mismo `restaurant_id`.

## 4. Opciones del POS (pos-operations)

- [x] 4.1 Modificar lógica de creación de órdenes para inyectar automáticamente el `restaurant_id` del empleado que hace la petición.
- [x] 4.2 Restringir endpoints de cobros y cierre financiero (`/api/orders/checkout`) requerir explícitamente el rol `CAJERO`.
- [x] 4.3 Restringir creación/modificación de órdenes (`/api/orders`) para `MESERO` y `CAJERO`.