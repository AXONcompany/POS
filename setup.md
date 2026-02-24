# POS — Setup (Docker + Postgres + Migraciones)
# a

Este documento deja el proyecto listo para correr en local con Docker Compose, aplicar migraciones y verificar que el backend esté funcionando.

## Requisitos

- Linux/macOS/WSL2
- Docker Engine + Docker Compose v2 (`docker compose version`)
- `curl` (para probar endpoints)
- Opcional: `psql` (si quieres inspeccionar la DB desde tu host)

## 1) Levantar servicios

Desde la raíz del repo:

```bash
docker compose up -d --build
```

Verifica estado:

```bash
docker compose ps
```

El backend expone HTTP en `http://localhost:8080`.

### Nota sobre Postgres y puertos

- El contenedor de Postgres escucha en **5432 dentro del contenedor**.
- En el host, este repo publica Postgres en **5433** (ver `docker-compose.yml`) para evitar conflicto si ya tienes Postgres local usando 5432.

## 2) Variables de entorno (DB)

El repo incluye un script que exporta variables requeridas para DB:

- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`

Ejecuta (en tu shell actual):

```bash
source scripts/load-db-env.sh
```

- Si existe un archivo `.env`, el script lo carga automáticamente.
- Si no existe `.env`, te pedirá los valores por consola.

Ejemplo de `.env` (opcional):

```env
DB_NAME=pos
DB_USER=postgres
DB_PASSWORD=postgres
```

> Importante: usa `source` para que las variables queden exportadas en tu terminal.

## 3) Migraciones

Las migraciones viven en `db/migrations` y se ejecutan con `migrate/migrate` vía Docker Compose.

### Aplicar migraciones

```bash
source scripts/load-db-env.sh
docker compose run --rm migrate up
```

Alternativa equivalente (scripts):

```bash
source scripts/load-db-env.sh
./scripts/migrate-up.sh
```

### Ver versión aplicada

```bash
source scripts/load-db-env.sh
docker compose run --rm migrate version
```

### Rollback (baja 1 por defecto)

```bash
source scripts/load-db-env.sh
docker compose run --rm migrate down 1
# o bajar N
docker compose run --rm migrate down 2
```

### Verificar tablas

Dentro del contenedor:

```bash
docker compose exec -T db psql -U "$DB_USER" -d "$DB_NAME" -c '\\dt'
```

Desde tu host (si tienes `psql`):

```bash
psql -h localhost -p 5433 -U "$DB_USER" -d "$DB_NAME" -c '\\dt'
```

## 4) Verificar que el backend está OK

Endpoints:

```bash
curl -sS http://localhost:8080/health && echo
curl -sS http://localhost:8080/ping && echo
```

Logs:

```bash
docker compose logs -f app
```

## 5) (Opcional) SQLC

El proyecto está configurado para generar código con sqlc:

- Schema para sqlc: `db/schema`
- Queries: `db/queries`
- Salida Go: `internal/infrastructure/persistence/postgres/sqlc`

Generar:

```bash
sqlc generate
```

### Ejemplo (1 query)

Query en:

- `db/queries/categories.sql`

Contenido (ya incluido en el repo):

```sql
-- name: CreateCategory :one
insert into categories (category_name)
values ($1)
returning id, created_at, updated_at, deleted_at, category_name;
```

Luego de correr `sqlc generate`, puedes llamarlo desde Go así (ejemplo conceptual):

```go
import (
	"context"

	apppg "github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres"
	dbsqlc "github.com/AXONcompany/POS/internal/infrastructure/persistence/postgres/sqlc"
)

func example(ctx context.Context, db *apppg.DB) error {
	q := dbsqlc.New(db.Pool)
	_, err := q.CreateCategory(ctx, "Bebidas")
	return err
}
```

> Si aún no tienes `sqlc` instalado, instálalo según tu entorno (Go toolchain / package manager).

## Troubleshooting

### A) `FATAL: role "-d" does not exist` en logs de Postgres

Suele pasar cuando el healthcheck intenta ejecutar `pg_isready` con variables inexistentes. En este repo el healthcheck usa `POSTGRES_USER/POSTGRES_DB` (variables internas del contenedor), así que si ves esto revisa que tu `docker-compose.yml` no haya sido modificado.

### B) `failed to bind host port 0.0.0.0:5432: address already in use`

Tienes un Postgres local usando el puerto 5432. Soluciones:

- Mantener este repo en `5433:5432` (recomendado), o
- Detener el Postgres local.

Ver quién ocupa el puerto:

```bash
sudo lsof -iTCP:5432 -sTCP:LISTEN -n -P
```

### C) `db error: ping postgres: context deadline exceeded` en logs del app

Indica que el backend no pudo conectar a Postgres. Revisa:

```bash
docker compose ps

docker compose logs --tail=200 db

docker compose logs --tail=200 app
```

Y confirma que el `db` está healthy antes de que arranque `app`.
