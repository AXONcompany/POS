## 1. Bugfixes Criticos (Pre-reestructuracion)

- [x] 1.1 Corregir bug en `RefreshToken` de `internal/usecase/auth/usecase.go`: generar tokens directamente sin llamar a Login con el hash
- [x] 1.2 Alinear nombres de roles en `internal/infrastructure/rest/auth/handler.go`: ADMIN->PROPIETARIO, CAJA->CAJERO en el map `roleNames` y en el handler `Register`

## 2. Migracion de Base de Datos

- [x] 2.1 Crear `db/migrations/000011_restructure_owner_venue_pos.up.sql` con todas las sentencias DDL en orden de dependencia
- [x] 2.2 Crear tabla `owners` con campos id, name, email, password_hash, is_active, timestamps
- [x] 2.3 Crear tabla `venues` con FK `owner_id` a `owners`, campos name, address, phone, is_active, timestamps
- [x] 2.4 Crear tabla `pos_terminals` con FK `venue_id` a `venues`, campos terminal_name, is_active, timestamps
- [x] 2.5 Migrar datos existentes: crear owner y venue a partir de datos en `restaurants` y `users`
- [x] 2.6 Agregar columna `venue_id` a `users` y eliminar `restaurant_id`
- [x] 2.7 Agregar columna `venue_id` a `ingredients`, `products`, `categories`, `tables`
- [x] 2.8 Cambiar constraint UNIQUE de `tables.table_number` a `(venue_id, table_number)`
- [x] 2.9 Actualizar `orders`: reemplazar `restaurant_id` con `venue_id`, agregar `pos_terminal_id`
- [x] 2.10 Actualizar `payments`: reemplazar `restaurant_id` con `venue_id`, agregar `pos_terminal_id`
- [x] 2.11 Eliminar tablas `table_waitress`, `waitress`, `restaurants`
- [x] 2.12 Crear `db/migrations/000011_restructure_owner_venue_pos.down.sql` para rollback

## 3. Dominio (Entidades Go)

- [ ] 3.1 Crear `internal/domain/owner/owner.go` con struct Owner
- [ ] 3.2 Crear `internal/domain/venue/venue.go` con struct Venue
- [ ] 3.3 Crear `internal/domain/pos/terminal.go` con struct Terminal
- [ ] 3.4 Modificar `internal/domain/user/user.go`: RestaurantID -> VenueID
- [ ] 3.5 Modificar `internal/domain/table/table.go`: agregar VenueID, eliminar TableWaitress
- [ ] 3.6 Modificar `internal/domain/order/order.go`: RestaurantID -> VenueID, agregar POSTerminalID
- [ ] 3.7 Modificar `internal/domain/payment/payment.go`: RestaurantID -> VenueID, agregar POSTerminalID
- [ ] 3.8 Agregar VenueID a `internal/domain/ingredient/ingredient.go`
- [ ] 3.9 Agregar VenueID a `internal/domain/product/product.go` y `category.go`
- [ ] 3.10 Eliminar `internal/domain/restaurant/restaurant.go`

## 4. Queries SQL y sqlc

- [ ] 4.1 Crear `db/queries/owners.sql` con CRUD de owners
- [ ] 4.2 Crear `db/queries/venues.sql` con CRUD de venues filtrado por owner_id
- [ ] 4.3 Crear `db/queries/pos_terminals.sql` con CRUD filtrado por venue_id
- [ ] 4.4 Actualizar `db/queries/ingredients.sql`: agregar venue_id a todas las queries
- [ ] 4.5 Actualizar `db/queries/products.sql`: agregar venue_id a todas las queries
- [ ] 4.6 Actualizar `db/queries/categories.sql`: agregar venue_id a todas las queries
- [ ] 4.7 Actualizar `db/queries/tables.sql`: agregar venue_id, eliminar queries de waitress
- [ ] 4.8 Actualizar `db/queries/users.sql`: restaurant_id -> venue_id
- [ ] 4.9 Eliminar `db/queries/waitress.sql`
- [ ] 4.10 Ejecutar `sqlc generate` y verificar que no hay errores

## 5. Repositorios (Persistence Layer)

- [ ] 5.1 Crear `postgres/owner_repository.go` con CRUD
- [ ] 5.2 Crear `postgres/venue_repository.go` con CRUD filtrado por owner_id
- [ ] 5.3 Crear `postgres/pos_terminal_repository.go` con CRUD filtrado por venue_id
- [ ] 5.4 Actualizar `postgres/ingredient_repository.go`: agregar venueID a todos los metodos
- [ ] 5.5 Actualizar `postgres/product_repository.go`: agregar venueID
- [ ] 5.6 Actualizar `postgres/category_repository.go`: agregar venueID
- [ ] 5.7 Actualizar `postgres/table_repository.go`: agregar venueID, eliminar metodos waitress
- [ ] 5.8 Actualizar `postgres/order_repository.go`: restaurantID -> venueID
- [ ] 5.9 Actualizar `postgres/payment_repository.go`: restaurantID -> venueID
- [ ] 5.10 Actualizar `postgres/report_repository.go`: restaurantID -> venueID
- [ ] 5.11 Actualizar `postgres/user_repository.go`: restaurantID -> venueID
- [ ] 5.12 Eliminar `postgres/restaurant_repository.go`

## 6. Interfaces de Repositorio y Usecases

- [ ] 6.1 Crear `usecase/owner/usecase.go` con CRUD y autenticacion de owners
- [ ] 6.2 Crear `usecase/venue/usecase.go` con CRUD de venues por owner
- [ ] 6.3 Crear `usecase/pos/usecase.go` con CRUD de terminales por venue
- [ ] 6.4 Actualizar `usecase/ingredient/usecase.go`: agregar venueID a interface y metodos
- [ ] 6.5 Actualizar `usecase/product/usecase.go`: agregar venueID a interfaces y metodos
- [ ] 6.6 Actualizar `usecase/table/usecase.go`: agregar venueID, eliminar metodos waitress
- [ ] 6.7 Actualizar `usecase/order/usecase.go`: restaurantID -> venueID
- [ ] 6.8 Actualizar `usecase/payment/usecase.go`: restaurantID -> venueID
- [ ] 6.9 Actualizar `usecase/report/usecase.go`: restaurantID -> venueID
- [ ] 6.10 Actualizar `usecase/user/usecase.go`: ListByRestaurant -> ListByVenue
- [ ] 6.11 Actualizar `usecase/auth/usecase.go`: restaurant_id -> venue_id en RegisterUser y generateToken

## 7. Handlers REST y Middleware

- [ ] 7.1 Actualizar `middleware/auth.go`: RestaurantIDKey -> VenueIDKey, agregar OwnerIDKey
- [ ] 7.2 Actualizar `auth/handler.go`: nombres de roles, venue_id en Register
- [ ] 7.3 Actualizar handlers existentes (ingredient, product, table, order, payment, report, user): extraer venue_id del JWT
- [ ] 7.4 Crear handler `rest/owner/handler.go`
- [ ] 7.5 Crear handler `rest/venue/handler.go`
- [ ] 7.6 Crear handler `rest/pos/handler.go`

## 8. Rutas y Wiring

- [ ] 8.1 Actualizar `routes.go`: agregar rutas para owner, venue, pos_terminal
- [ ] 8.2 Actualizar `cmd/server/main.go`: instanciar nuevos repos, usecases, handlers

## 9. Seed Data

- [ ] 9.1 Crear seed data post-reestructuracion con owner, venue, terminal, cajero, mesero y mesa de ejemplo

## 10. Verificacion

- [ ] 10.1 Compilar sin errores: `go build ./cmd/server/`
- [ ] 10.2 Ejecutar sqlc generate sin errores
- [ ] 10.3 Ejecutar migraciones en BD local
- [ ] 10.4 Actualizar tests existentes (auth, order, product, user, middleware)
- [ ] 10.5 Crear tests nuevos para owner, venue, pos_terminal, aislamiento
- [ ] 10.6 Ejecutar todos los tests: `go test ./internal/... -v`
