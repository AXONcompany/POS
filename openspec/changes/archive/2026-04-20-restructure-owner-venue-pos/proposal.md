## Why

El modelo de datos actual del POS es plano: una tabla `restaurants` sin dueño, usuarios vinculados directamente a un `restaurant_id`, y tablas criticas (`ingredients`, `products`, `categories`, `tables`) completamente globales sin aislamiento por tenant. Esto impide que un propietario gestione multiples sedes, que cada sede tenga su propio inventario/mesas/personal, y expone datos entre restaurantes distintos. Ademas, el flujo de `RefreshToken` tiene un bug critico que impide la renovacion de sesiones, y los nombres de roles en el handler no coinciden con la BD.

## What Changes

- **BREAKING**: Se elimina la tabla `restaurants` y se reemplaza por `venues` (sedes), propiedad de un `owner`.
- **BREAKING**: Se agrega `venue_id` a `ingredients`, `products`, `categories`, `tables` para aislamiento de datos por sede.
- **BREAKING**: Se reemplaza `restaurant_id` por `venue_id` en `users`, `orders`, `payments`.
- **BREAKING**: Se elimina la tabla `waitress` y `table_waitress` (redundante con users role MESERO).
- Se crea la tabla `owners` como entidad separada para propietarios con autenticacion propia.
- Se crea la tabla `venues` para representar sedes fisicas de un propietario.
- Se crea la tabla `pos_terminals` para representar terminales POS dentro de una sede.
- Se corrige el bug del flujo `RefreshToken` en el usecase de auth.
- Se corrige la inconsistencia de nombres de roles (handler vs BD).

## Capabilities

### New Capabilities
- `owner-management`: Gestion de propietarios como entidad independiente. CRUD de owners con autenticacion separada, capacidad de poseer multiples sedes.
- `venue-management`: Gestion de sedes (venues) por parte de un propietario. Cada sede es una ubicacion fisica con inventario, mesas y personal aislados.
- `pos-terminal-management`: Gestion de terminales POS dentro de una sede. Cada terminal registra las operaciones (ordenes, pagos) que procesa.
- `data-isolation`: Aislamiento completo de datos por sede. Todas las entidades operativas (ingredientes, productos, categorias, mesas) quedan vinculadas a una venue_id.

### Modified Capabilities
- `user-auth`: El JWT ahora incluye `venue_id` y `owner_id` en lugar de `restaurant_id`. Se corrige el flujo de RefreshToken. Se alinean nombres de roles.
- `restaurant-management`: Se reemplaza completamente por `venue-management`. El propietario ahora gestiona sedes en lugar de un restaurante plano.
- `pos-operations`: Las ordenes y pagos ahora referencian `venue_id` y opcionalmente `pos_terminal_id` en lugar de `restaurant_id`.

## Impact

- **Base de datos**: Migracion mayor que crea 3 tablas nuevas (`owners`, `venues`, `pos_terminals`), agrega `venue_id` a 4 tablas existentes, elimina 3 tablas (`restaurants`, `waitress`, `table_waitress`), y renombra columnas en `users`, `orders`, `payments`.
- **API REST**: Todos los endpoints que reciben o retornan `restaurant_id` cambian a `venue_id`. Nuevos endpoints para owners, venues y terminales POS.
- **JWT / Middleware**: Los claims del token cambian de `restaurant_id` a `venue_id` + `owner_id`. El middleware de autenticacion se actualiza.
- **Queries SQL / sqlc**: Todas las queries se actualizan para filtrar por `venue_id`. Se regenera el codigo generado por sqlc.
- **Tests**: 6 archivos de test existentes requieren actualizacion. Se crean tests nuevos para las nuevas entidades y para verificar aislamiento.
