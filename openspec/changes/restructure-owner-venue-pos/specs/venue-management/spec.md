## ADDED Requirements

### Requirement: Entidad Sede (Venue) con Aislamiento de Datos
El sistema DEBE representar cada ubicacion fisica como una `venue` con su propio inventario, mesas, productos y personal. Cada venue pertenece a un `owner` mediante FK `owner_id`.

#### Scenario: Creacion de venue con datos completos
- **WHEN** se crea una venue con nombre, direccion, telefono y owner_id valido.
- **THEN** el sistema registra la venue con `is_active = true` y la vincula al owner.

#### Scenario: Usuarios vinculados a una venue
- **WHEN** se crea un usuario (cajero o mesero) con un `venue_id`.
- **THEN** el usuario queda asociado exclusivamente a esa venue y solo puede operar datos de ella.

### Requirement: Reemplazo de Restaurants por Venues
El sistema DEBE eliminar la tabla `restaurants` y reemplazarla por `venues` con el campo adicional `owner_id`. Todas las referencias a `restaurant_id` se actualizan a `venue_id`.

#### Scenario: Migracion de datos de restaurants a venues
- **WHEN** se ejecuta la migracion 000011.
- **THEN** los datos existentes en `restaurants` se migran a `venues` con un `owner_id` asignado, y la tabla `restaurants` se elimina.
