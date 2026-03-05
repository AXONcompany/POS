## Context

El POS actual tiene un modelo de datos plano donde `restaurants` es una tabla sin dueño, `users` se vinculan directamente a un `restaurant_id`, y las tablas de datos operativos (`ingredients`, `products`, `categories`, `tables`) son globales sin ningun filtro de tenant. Esto crea un sistema single-tenant de facto que no soporta el modelo de negocio real: un propietario que posee una o mas sedes, cada una con su POS, personal e inventario independiente.

La arquitectura Clean Architecture (cmd/internal/domain/usecase/infrastructure) esta bien establecida. El auth usa JWT HMAC-SHA256 con access+refresh tokens y RBAC con 3 roles. La BD es PostgreSQL con queries generadas por sqlc.

## Goals / Non-Goals

**Goals:**
- Establecer la jerarquia `Owner > Venue > POS Terminal` en el esquema de BD
- Aislar todos los datos operativos por `venue_id` (ingredientes, productos, categorias, mesas)
- Reemplazar `restaurant_id` por `venue_id` en `users`, `orders`, `payments`
- Crear una entidad `Owner` separada para propietarios con gestion multi-sede
- Crear la entidad `POS Terminal` para registrar que terminal procesa cada operacion
- Corregir el bug del RefreshToken y la inconsistencia de nombres de roles
- Mantener compatibilidad con la arquitectura Clean Architecture existente

**Non-Goals:**
- No se implementa UI/frontend en este cambio
- No se implementa catalogo de productos compartido entre sedes (cada sede es independiente)
- No se implementa sincronizacion offline entre sede y nube
- No se implementa facturacion electronica ni integraciones externas
- No se migran datos de produccion (solo se provee la migracion SQL; el operador la ejecuta)

## Decisions

### D1: Owner como entidad separada vs. User con rol

**Decision**: Owner es una tabla y entidad separada de `users`.

**Alternativas consideradas**:
- **A) User con role_id=1**: Simple, sin cambios de schema. Pero un "user" PROPIETARIO no puede poseer multiples venues de forma limpia, y mezcla la autenticacion del panel de gestion con la del POS.
- **B) Owner separado (elegida)**: Tabla `owners` independiente. `venues` tiene FK a `owners.id`. Los `users` (cajeros, meseros) pertenecen a una `venue`. Separacion limpia de autenticacion y permisos.

**Razon**: El Owner opera a nivel de gestion (crear sedes, ver reportes consolidados), mientras que los Users operan a nivel de sede (tomar pedidos, cobrar). Son contextos distintos.

### D2: Renombrar restaurants a venues

**Decision**: La tabla `restaurants` se elimina y se crea `venues` con FK a `owners`.

**Razon**: "Restaurant" implica un negocio completo. "Venue" (sede) implica una ubicacion fisica dentro de un negocio, que es el concepto correcto para soportar multi-sede.

### D3: Migracion atomica vs. incremental

**Decision**: Una sola migracion SQL (000011) que realiza todos los cambios de esquema.

**Alternativas consideradas**:
- **A) Multiples migraciones (011, 012, 013...)**: Mas granular, mas facil de debuggear. Pero introduce estados intermedios inconsistentes.
- **B) Una migracion atomica (elegida)**: Todos los cambios en una transaccion. Si falla, se hace rollback completo. No hay estados intermedios.

**Razon**: Los cambios son interdependientes (no puedes eliminar `restaurants` sin crear `venues` y actualizar las FKs). Una migracion atomica evita estados rotos.

### D4: Eliminar tabla waitress

**Decision**: Eliminar `waitress` y `table_waitress`. Los meseros son `users` con `role_id=3`.

**Razon**: La tabla `waitress` solo tiene un campo (`id_user` FK a users). Es completamente redundante. La asignacion de meseros a mesas se puede hacer directamente con `users` filtrando por rol.

### D5: venue_id en todas las tablas de datos

**Decision**: Agregar `venue_id` NOT NULL con FK a `venues` en: `ingredients`, `products`, `categories`, `tables`.

**Alternativas consideradas**:
- **A) Usar RLS (Row Level Security) de PostgreSQL**: Aislamiento transparente a nivel de BD. Pero agrega complejidad operativa y sqlc no genera soporte nativo para RLS.
- **B) FK explicita venue_id (elegida)**: Cada query filtra explicitamente por `venue_id`. Es explicito, entendible y compatible con sqlc.

## Risks / Trade-offs

| Riesgo | Mitigacion |
|--------|-----------|
| Migracion destructiva borra datos existentes | La migracion incluye `INSERT INTO` para migrar datos del esquema viejo al nuevo. Se provee migracion down para rollback. |
| Todos los endpoints cambian (breaking API) | Este es un sistema en desarrollo pre-produccion. No hay clientes externos que dependan de la API actual. |
| Owner sin autenticacion completa en esta fase | Se crea la tabla y el dominio. La autenticacion del Owner se implementa en el modulo 5 (usecases). En esta fase solo se establece el esquema. |
| Complejidad de queries con venue_id en todas partes | Es un trade-off aceptable: cada query es explicita sobre que datos consulta. Se puede abstraer en el middleware a futuro. |

## Migration Plan

1. Respaldar la BD actual (si hay datos relevantes)
2. Ejecutar `migrate -path db/migrations -database $DATABASE_URL up`
3. La migracion 000011 crea las nuevas tablas, migra datos, elimina las obsoletas
4. Regenerar codigo sqlc: `sqlc generate`
5. **Rollback**: Ejecutar la migracion down: `migrate -path db/migrations -database $DATABASE_URL down 1`

## Open Questions

- Ninguna para este modulo. Las decisiones de diseno cubren todos los aspectos del esquema de BD.
