## ADDED Requirements

### Requirement: Aislamiento de Ingredientes por Sede
El sistema DEBE vincular cada ingrediente a una `venue_id`. Las queries de ingredientes DEBEN filtrar por `venue_id` para garantizar que una sede no acceda al inventario de otra.

#### Scenario: Ingrediente creado con venue_id
- **WHEN** se crea un ingrediente con `venue_id = V1`.
- **THEN** el ingrediente solo es visible en consultas que filtren por `venue_id = V1`.

#### Scenario: Intento de acceso a ingrediente de otra sede
- **WHEN** se consulta un ingrediente con `venue_id = V2` desde el contexto de `venue_id = V1`.
- **THEN** el sistema no retorna el ingrediente (query no lo encuentra).

### Requirement: Aislamiento de Productos y Categorias por Sede
El sistema DEBE vincular cada producto y categoria a una `venue_id`. Cada sede tiene su catalogo independiente.

#### Scenario: Producto visible solo en su sede
- **WHEN** se listan productos con filtro `venue_id = V1`.
- **THEN** el sistema retorna solo los productos de venue V1.

### Requirement: Aislamiento de Mesas por Sede
El sistema DEBE vincular cada mesa a una `venue_id`. El constraint UNIQUE de `table_number` pasa a ser `UNIQUE(venue_id, table_number)` para permitir numeros repetidos en sedes distintas.

#### Scenario: Mesas con mismo numero en sedes diferentes
- **WHEN** la venue V1 tiene mesa numero 1 y la venue V2 crea mesa numero 1.
- **THEN** el sistema permite ambas porque el constraint es por `(venue_id, table_number)`.

#### Scenario: Mesa duplicada en la misma sede
- **WHEN** la venue V1 ya tiene mesa numero 5 e intenta crear otra mesa numero 5.
- **THEN** el sistema rechaza la operacion por violacion de constraint unique `(venue_id, table_number)`.

### Requirement: Eliminacion de Tabla Waitress
El sistema DEBE eliminar las tablas `waitress` y `table_waitress`. Los meseros se identifican como `users` con `role_id = 3` (MESERO) vinculados a una `venue_id`.

#### Scenario: Mesero como user con rol
- **WHEN** se consultan los meseros de una venue.
- **THEN** el sistema filtra `users WHERE venue_id = X AND role_id = 3`.
