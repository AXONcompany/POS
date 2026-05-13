# Plan de Reestructuración — AxonPOS Backend

**Fecha:** 2026-05-12  
**Estado:** Decisiones cerradas — listo para ejecución

### Decisiones confirmadas

| # | Decisión | Resolución |
|---|---|---|
| 1 | Mercado de pagos | **MercadoPago** — Colombia, arrancando en Santa Marta |
| 2 | Auth hosting | **Supabase Cloud** con path documentado a self-hosted |
| 3 | Formato instalador | **Windows (.msi)** como primera plataforma |
| 4 | Período de gracia | **5 días** tras vencimiento de licencia |
| 5 | Trial gratuito | **7 días por defecto**, configurable: habilitado/deshabilitado y duración editable |

---

## 1. Contexto y motivación

El backend actual (`POS/`) es un monolito Go que asume conectividad constante y un único cliente web. El nuevo modelo de producto exige:

- **App instalada por terminal** (Tauri + SQLite + servidor Go embebido)
- **Meseros en celulares** que se conectan a la terminal vía WiFi local
- **Operación offline-first**: ninguna operación se interrumpe si cae internet
- **Múltiples terminales** por sede sincronizando al mismo cloud
- **Licenciamiento** por plan con validación offline mediante JWT firmado

Esto requiere descomponer el monolito en tres piezas claramente separadas.

---

## 2. Arquitectura objetivo

```
┌─────────────────────────────────────────────────────┐
│                    CLOUD                            │
│                                                     │
│  ┌─────────────┐    ┌──────────────────────────┐   │
│  │ Auth Service│    │     POS Cloud Backend    │   │
│  │ (Supabase/  │    │   (Go actual, refactor)  │   │
│  │  Clerk)     │    │   PostgreSQL             │   │
│  └─────────────┘    └──────────────────────────┘   │
│         │                      │                    │
│  ┌─────────────────────────────────────────────┐   │
│  │           License Service (Go nuevo)         │   │
│  │           PostgreSQL propio                  │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
              │  Internet (sync + licencia)
┌─────────────────────────────────────────────────────┐
│               INSTALACIÓN LOCAL (Tauri)             │
│                                                     │
│  ┌─────────────────────────────────────────────┐   │
│  │         Local POS Server (Go binario)        │   │
│  │         SQLite  │  Sync Client               │   │
│  │         QR Code │  mDNS                      │   │
│  └─────────────────────────────────────────────┘   │
│                      │  WiFi local                  │
│          📱 Mesero 1        📱 Mesero 2             │
│          (browser)          (browser)               │
└─────────────────────────────────────────────────────┘
```

---

## 3. Descomposición de servicios

### 3.1 Auth Service — NO construir, usar Supabase Auth

**Responsabilidades:**
- Login de propietarios con Google OAuth
- Emisión de tokens de acceso (JWT cortos, 15 min)
- Refresh tokens (7 días)
- Gestión de sesiones de propietarios en el portal web

**Por qué no construirlo:** implementar OAuth correctamente toma semanas y no aporta ventaja competitiva. Supabase Auth es open source, self-hostable, tiene Go SDK y se integra con PostgreSQL nativamente.

**Cambio en el backend actual:** el módulo `internal/infrastructure/rest/auth` se simplifica — deja de manejar OAuth y se convierte en un validador de tokens emitidos por Supabase.

---

### 3.2 License Service — nuevo microservicio Go

**Responsabilidades:**
- Recibir webhooks de pago (Stripe / MercadoPago) y activar licencias
- Emitir **JWTs de licencia** firmados con RSA-256
- Renovar tokens antes del vencimiento
- Gestionar planes y cantidad de terminales permitidas por sede
- Registrar y autenticar terminales (emitir `terminal_token`)

**Base de datos:** PostgreSQL propio (no compartido con POS Backend)

**Tablas nuevas:**

```sql
-- Planes disponibles
CREATE TABLE plans (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,   -- starter, pro, enterprise
    max_terminals INT NOT NULL,
    price_month NUMERIC(10,2),
    features    JSONB
);

-- Licencias por sede
CREATE TABLE licenses (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id      UUID NOT NULL,         -- referencia al owner en Auth Service
    venue_id      UUID NOT NULL,
    plan_id       INT REFERENCES plans(id),
    status        VARCHAR(20) NOT NULL,  -- trial, active, expired, cancelled, blocked
    starts_at     TIMESTAMPTZ NOT NULL,
    expires_at    TIMESTAMPTZ NOT NULL,
    grace_days    INT DEFAULT 5,
    created_at    TIMESTAMPTZ DEFAULT now(),
    updated_at    TIMESTAMPTZ DEFAULT now()
);

-- Tokens de terminal (credencial de máquina, no de usuario)
CREATE TABLE terminal_tokens (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_id    UUID REFERENCES licenses(id),
    venue_id      UUID NOT NULL,
    terminal_name VARCHAR(100),
    token_hash    TEXT NOT NULL,         -- hash del terminal_token
    last_seen_at  TIMESTAMPTZ,
    version       VARCHAR(20),           -- versión del POS instalado
    is_active     BOOLEAN DEFAULT true,
    created_at    TIMESTAMPTZ DEFAULT now()
);

-- Configuración global del trial (singleton)
CREATE TABLE trial_config (
    id              INTEGER PRIMARY KEY CHECK (id = 1),
    is_enabled      BOOLEAN NOT NULL DEFAULT true,
    duration_days   INT NOT NULL DEFAULT 7,
    grace_days      INT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT now(),
    updated_by      TEXT
);

INSERT INTO trial_config (id, is_enabled, duration_days, grace_days)
VALUES (1, true, 7, 0);

-- Historial de pagos
CREATE TABLE payment_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_id      UUID REFERENCES licenses(id),
    provider        VARCHAR(20) NOT NULL,  -- stripe, mercadopago
    provider_ref    TEXT NOT NULL,
    amount          NUMERIC(10,2),
    status          VARCHAR(20) NOT NULL,  -- confirmed, failed, refunded
    raw_payload     JSONB,
    created_at      TIMESTAMPTZ DEFAULT now()
);
```

**JWT de licencia — payload:**
```json
{
  "iss": "axonpos-license",
  "sub": "terminal_id",
  "venue_id": "uuid",
  "owner_id": "uuid",
  "plan": "pro",
  "max_terminals": 5,
  "features": ["offline", "multi_terminal", "reporting"],
  "iat": 1234567890,
  "exp": 1234567890,
  "grace_days": 5
}
```
Firmado con **RSA-256**. La clave pública va embebida en el binario Tauri. La terminal verifica sin red.

**Endpoints:**
```
POST /license/activate          → activar licencia tras pago confirmado
POST /license/renew             → renovar JWT próximo a vencer
POST /terminals/register        → registrar nueva terminal, retorna terminal_token
POST /terminals/heartbeat       → actualizar last_seen_at y versión
DELETE /terminals/:id           → desactivar terminal

POST /webhooks/stripe
POST /webhooks/mercadopago
```

---

### 3.3 POS Cloud Backend — refactorizar el actual

**Responsabilidades:**
- Recibir y consolidar eventos de sincronización de todas las terminales
- Resolver conflictos entre terminales
- Servir snapshots de menú y layout de mesas
- Reportes y contabilidad (datos consolidados de todas las terminales)
- WebSocket para estado de mesas en tiempo real (cuando terminales están online)

**Cambios respecto al estado actual:**

#### Módulo de auth
- Eliminar `register-owner` y manejo de contraseñas de propietarios
- El endpoint `/auth/login` pasa a validar tokens de Supabase Auth (verificar JWT externo)
- Los meseros siguen con PIN, pero el PIN se valida contra SQLite local — el cloud solo sincroniza el catálogo de staff

#### Nuevas migraciones

**`000015_add_sync_infrastructure.up.sql`**
```sql
-- Eventos de sincronización
CREATE TABLE sync_events (
    id              UUID PRIMARY KEY,
    terminal_id     UUID NOT NULL,
    venue_id        UUID NOT NULL,
    entity          VARCHAR(50) NOT NULL,  -- order, order_item, table, payment, staff
    entity_id       TEXT NOT NULL,
    operation       VARCHAR(20) NOT NULL,  -- CREATE, UPDATE, DELETE
    payload         JSONB NOT NULL,
    client_created_at TIMESTAMPTZ NOT NULL, -- timestamp del cliente
    synced_at       TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_sync_events_terminal ON sync_events(terminal_id, synced_at);
CREATE INDEX idx_sync_events_venue    ON sync_events(venue_id, synced_at);

-- Tracking de sincronización por terminal
ALTER TABLE pos_terminals ADD COLUMN terminal_token_hash TEXT;
ALTER TABLE pos_terminals ADD COLUMN last_sync_at        TIMESTAMPTZ;
ALTER TABLE pos_terminals ADD COLUMN pos_version         VARCHAR(20);

-- Trazabilidad de origen en órdenes y pagos
ALTER TABLE orders   ADD COLUMN terminal_id UUID REFERENCES pos_terminals(id);
ALTER TABLE orders   ADD COLUMN synced_at   TIMESTAMPTZ;
ALTER TABLE payments ADD COLUMN terminal_id UUID REFERENCES pos_terminals(id);

-- Reclamo de mesa por terminal (para conflictos multi-terminal)
ALTER TABLE tables ADD COLUMN claimed_by_terminal UUID REFERENCES pos_terminals(id);
ALTER TABLE tables ADD COLUMN claimed_at          TIMESTAMPTZ;
```

**`000016_add_pagination_support.up.sql`**
```sql
-- Los list endpoints actualmente devuelven todo. Agregar soporte de paginación
-- Se resuelve a nivel de queries sqlc, no requiere cambios de schema.
-- Esta migración es un placeholder para documentar la tarea.
```

#### Nuevo endpoint de sincronización
```
POST /sync
```
**Request:**
```json
{
  "terminal_id": "uuid",
  "last_sync_at": "2024-01-01T10:00:00Z",
  "events": [
    {
      "id": "local-uuid",
      "entity": "order",
      "entity_id": "local-order-uuid",
      "operation": "CREATE",
      "payload": { ... },
      "client_created_at": "2024-01-01T10:01:00Z"
    }
  ]
}
```
**Response:**
```json
{
  "synced_until": "2024-01-01T10:05:00Z",
  "remote_events": [ ... ],
  "conflicts": [
    {
      "local_event_id": "...",
      "resolution": "rejected",
      "reason": "table_already_claimed",
      "details": "Mesa 5 fue reclamada por terminal B a las 10:02"
    }
  ]
}
```

**Reglas de resolución de conflictos:**

| Entidad | Conflicto | Resolución |
|---|---|---|
| Mesa (reclamo) | Dos terminales abren misma mesa offline | Primer sync gana, segundo notificado |
| Order items | Dos terminales agregan items a misma orden | Merge — todos los items se aplican |
| Estado de orden | Estado contradictorio (e.g., PAID desde dos lados) | Primero en sync gana, segundo rechazado |
| Pago | Dos terminales cobran misma orden | Primero gana, segundo rechazado con alerta |
| Menú | Terminal modifica menú offline | Rechazado — menú es pull-only desde cloud |
| Layout de mesas | Terminal modifica layout offline | Rechazado — layout es pull-only desde cloud |

#### Nuevos endpoints de snapshot (para sincronización inicial)
```
GET /snapshot/menu        → catálogo completo (productos, categorías)
GET /snapshot/tables      → layout de mesas de la sede
GET /snapshot/staff       → lista de meseros/cajeros de la sede
```

#### WebSocket para estado de mesas en tiempo real
```
WS /ws/venue/:venue_id
```
Cuando una terminal está online, recibe push inmediato de cambios de estado de mesas (check-in, orden creada, cobrado). Esto reduce la latencia percibida en ambientes multi-terminal.

---

### 3.4 Local POS Server — nuevo binario Go embebido en Tauri

Este es el componente más nuevo. Es un servidor Go que vive dentro del paquete Tauri como **sidecar** y arranca junto con la app.

**Responsabilidades:**
- Servir el React compilado en el puerto `3000` (acceso de meseros)
- Exponer API REST local en el puerto `3001`
- Leer y escribir SQLite local
- Mantener log de cambios (`change_log`)
- Background sync loop: subir eventos, bajar cambios remotos
- Validar JWT de licencia al arrancar y periódicamente
- Generar QR con IP local + puerto
- Broadcast mDNS como `axonpos-{venue}.local`

**Estructura sugerida para el nuevo módulo:**
```
cmd/
  server/main.go        → modo cloud (PostgreSQL, comportamiento actual)
  local/main.go         → modo local (SQLite, nuevo)

internal/
  domain/               → sin cambios (compartido entre modos)
  usecase/              → sin cambios (compartido)
  infrastructure/
    persistence/
      postgres/         → repositorios actuales
      sqlite/           → nuevos repositorios SQLite (mismas interfaces)
    rest/               → handlers actuales (compatibles con ambos modos)
    sync/               → cliente de sincronización (nuevo)
      client.go         → POST /sync al cloud
      changlog.go       → escritura y lectura del change_log local
      resolver.go       → aplicar remote_events al SQLite local
    license/
      validator.go      → verificar JWT con clave pública embebida
    discovery/
      qrcode.go         → generar QR con IP:puerto
      mdns.go           → broadcast mDNS
```

**El truco clave:** como la arquitectura limpia ya separa dominio/usecase/infraestructura, los **usecases no cambian** al cambiar de PostgreSQL a SQLite. Solo cambia el repositorio concreto que se inyecta.

**Schema SQLite local** — espejo reducido del cloud:
```sql
-- Tablas locales (subset del schema cloud)
CREATE TABLE products    (...); -- read-only, se carga desde snapshot
CREATE TABLE categories  (...); -- read-only
CREATE TABLE staff       (...); -- sincronizado, PIN-based auth local
CREATE TABLE tables      (...); -- sincronizado, con claimed_by_terminal
CREATE TABLE orders      (...); -- escritura local, sync al cloud
CREATE TABLE order_items (...); -- escritura local
CREATE TABLE payments    (...); -- escritura local

-- Infraestructura de sync
CREATE TABLE change_log (
    id              TEXT PRIMARY KEY,    -- UUID generado por el cliente
    entity          TEXT NOT NULL,
    entity_id       TEXT NOT NULL,
    operation       TEXT NOT NULL,       -- CREATE, UPDATE, DELETE
    payload         TEXT NOT NULL,       -- JSON
    created_at      TEXT NOT NULL,       -- ISO 8601
    synced_at       TEXT                 -- NULL hasta confirmar sync
);

-- Licencia almacenada localmente
CREATE TABLE license_cache (
    id          INTEGER PRIMARY KEY CHECK (id = 1),  -- singleton
    jwt         TEXT NOT NULL,
    expires_at  TEXT NOT NULL,
    cached_at   TEXT NOT NULL
);

-- Configuración de terminal
CREATE TABLE terminal_config (
    id              INTEGER PRIMARY KEY CHECK (id = 1),
    terminal_id     TEXT NOT NULL,
    terminal_token  TEXT NOT NULL,
    venue_id        TEXT NOT NULL,
    terminal_name   TEXT
);
```

---

## 4. Gestión de terminales — flujo completo

```
PRIMERA ACTIVACIÓN
──────────────────
1. Owner instala el POS en el dispositivo
2. Tauri arranca → Local POS Server levanta en puertos 3000/3001
3. App detecta que no hay terminal_config → muestra pantalla de activación
4. Owner hace login con Google (OAuth → Supabase Auth)
5. App llama License Service: POST /terminals/register { venue_id, terminal_name }
6. License Service verifica:
   a. El token de Google pertenece a un owner válido
   b. La licencia del venue está activa
   c. No se superó max_terminals
7. Si OK → genera terminal_id + terminal_token, retorna license JWT
8. App guarda terminal_config + license_cache en SQLite
9. Terminal operativa

OPERACIÓN NORMAL
────────────────
1. Local POS Server arranca → valida license JWT (RSA local, sin red)
2. Si JWT válido (o en grace period) → app inicia
3. Muestra QR con http://{ip_local}:3000
4. Meseros escanean QR → entran al React en sus celulares
5. Cada mutación (orden, pago, estado de mesa) → escribe SQLite + inserta en change_log
6. Sync loop cada 30s → POST /sync al cloud con eventos pendientes

RENOVACIÓN DE LICENCIA
──────────────────────
1. License validator detecta JWT expira en < 3 días
2. Si hay red → POST /license/renew al License Service
3. Recibe nuevo JWT → actualiza license_cache
4. Si no hay red → corre con grace_days hasta reconectar

MULTI-TERMINAL SYNC
───────────────────
1. Terminal A toma orden en Mesa 5 → change_log: { CLAIM_TABLE, table_id: 5, terminal_id: A }
2. Terminal B (offline) también toma Mesa 5 → change_log local: { CLAIM_TABLE, table_id: 5, terminal_id: B }
3. Terminal A sincroniza primero → cloud registra: Mesa 5 = Terminal A
4. Terminal B sincroniza → cloud detecta conflicto
5. Cloud responde: { conflict: "table_already_claimed", resolution: "rejected" }
6. Local POS Server de B muestra alerta al cajero: "Mesa 5 fue tomada por Terminal A"
7. Cajero decide: fusionar órdenes manualmente o transferir
```

---

## 5. Portal web — alcance

Un proyecto Next.js separado (fuera del directorio `POS/`).

**Páginas:**
- Landing / pricing
- Registro (Google OAuth vía Supabase Auth)
- Dashboard del owner:
  - Sedes y estado de licencias
  - Terminales activas (última vez online, versión)
  - Descargar instalador (enlace al release de GitHub o S3)
  - Historial de pagos y facturas
- Checkout de plan (Stripe / MercadoPago)

---

## 6. Plan de ejecución por fases

### Fase 1 — Fundamentos de sync (2 semanas)
- [ ] Migración `000015`: agregar columnas de sync al schema PostgreSQL
- [ ] Nuevo módulo `internal/infrastructure/sync/` en el backend Go
- [ ] Endpoint `POST /sync` con resolución de conflictos básica
- [ ] Endpoints de snapshot: `/snapshot/menu`, `/snapshot/tables`, `/snapshot/staff`
- [ ] Tests de integración para sync

### Fase 2 — Local POS Server (2 semanas)
- [ ] Nuevo `cmd/local/main.go` con modo SQLite
- [ ] Repositorios SQLite (`internal/infrastructure/persistence/sqlite/`)
- [ ] Schema SQLite con `change_log` y `terminal_config`
- [ ] Sync client (goroutine background)
- [ ] QR code generator
- [ ] mDNS broadcast

### Fase 3 — License Service (1.5 semanas)
- [ ] Nuevo microservicio Go en `../license-service/`
- [ ] Schema PostgreSQL propio (licenses, terminal_tokens, payment_events, trial_config)
- [ ] Endpoints de registro de terminal y emisión de JWT
- [ ] Integración con MercadoPago (webhook + checkout) — único proveedor inicial
- [ ] Sistema de trial: activación/desactivación y edición de duración
- [ ] Validator en Local POS Server con RSA public key embebida

### Fase 4 — Auth refactor (1 semana)
- [ ] Configurar Supabase Auth con Google OAuth
- [ ] Simplificar `internal/infrastructure/rest/auth`: validar tokens de Supabase en vez de emitirlos
- [ ] Auth de meseros: PIN local en SQLite (no requiere cloud)
- [ ] Eliminar tabla `sessions` del POS Backend (pasa a Supabase)

### Fase 5 — Tauri app (2 semanas)
- [ ] Setup del proyecto Tauri sobre el frontend React existente
- [ ] Bundling del Local POS Server como sidecar
- [ ] IPC Tauri ↔ React para estado de licencia y sync
- [ ] Pantalla de activación de terminal
- [ ] Display de QR en pantalla principal

### Fase 6 — Portal web (2 semanas)
- [ ] Proyecto Next.js
- [ ] Integración con Supabase Auth (Google OAuth)
- [ ] Dashboard del propietario
- [ ] Checkout de plan (MercadoPago)
- [ ] Enlace de descarga del instalador (.msi Windows)

### Fase 7 — Hardening y WebSocket (1 semana)
- [ ] WebSocket `/ws/venue/:venue_id` para estado de mesas en tiempo real
- [ ] Paginación en todos los endpoints de lista
- [ ] Rate limiting
- [ ] CORS restrictivo (solo dominios propios + red local)
- [ ] Tests E2E del flujo completo offline → online → sync

---

## 7. Resumen de tecnologías

| Componente | Tecnología | Estado |
|---|---|---|
| POS Cloud Backend | Go + Gin + PostgreSQL | Refactorizar |
| License Service | Go + Gin + PostgreSQL | Nuevo |
| Auth | Supabase Auth (Google OAuth) | Nuevo (externo) |
| Local POS Server | Go + SQLite (embebido en Tauri) | Nuevo |
| POS App | Tauri + React + TypeScript | Migrar frontend actual |
| Portal web | Next.js + Supabase Auth | Nuevo |
| Pagos | MercadoPago (Colombia) | Nuevo |
| Descubrimiento local | QR code + mDNS | Nuevo |
| Sync | Event log + REST batch | Nuevo |
| Licencias offline | JWT RSA-256 | Nuevo |

---

## 8. Trial gratuito — diseño detallado

### Comportamiento
- 7 días por defecto, completamente configurable desde un panel de administración interno
- Se puede habilitar o deshabilitar globalmente en cualquier momento
- La duración es editable sin redespliegue (vive en base de datos, no en código)
- Al registrarse, si el trial está habilitado, el venue recibe una licencia con `status = 'trial'` automáticamente
- Al vencer el trial: el JWT local entra en `grace_days = 0` (no hay gracia para trials, se bloquea inmediatamente o se puede configurar)
- El propietario recibe emails de recordatorio: a los 5 días, al día 7 (vencimiento), y al bloquearse

### Schema adicional en License Service

```sql
-- Configuración global del trial (singleton — solo una fila)
CREATE TABLE trial_config (
    id              INTEGER PRIMARY KEY CHECK (id = 1),
    is_enabled      BOOLEAN NOT NULL DEFAULT true,
    duration_days   INT NOT NULL DEFAULT 7,
    grace_days      INT NOT NULL DEFAULT 0,       -- gracia tras vencer trial
    updated_at      TIMESTAMPTZ DEFAULT now(),
    updated_by      TEXT                          -- email del admin que lo modificó
);

-- Insertar configuración inicial
INSERT INTO trial_config (id, is_enabled, duration_days, grace_days)
VALUES (1, true, 7, 0);
```

La tabla `licenses` ya existente maneja el trial con el campo `status`:

```
status: 'trial' | 'active' | 'expired' | 'cancelled' | 'blocked'
```

Cuando `status = 'trial'`:
- `expires_at = created_at + trial_config.duration_days`
- `grace_days = trial_config.grace_days`
- El JWT de licencia incluye `"plan": "trial"` — el frontend puede mostrar un banner de trial con días restantes

### Flujo de conversión trial → pago

```
1. Owner se registra → licencia trial creada automáticamente (si trial habilitado)
2. POS emite JWT con plan=trial y expires_at
3. Frontend muestra banner: "Te quedan X días de prueba"
4. Owner paga plan → License Service actualiza licencia: status=active, plan=starter|pro|enterprise
5. Se emite nuevo JWT sin restricción de trial
6. Si trial vence sin pago → status=blocked, terminal muestra pantalla de bloqueo con CTA de pago
```

### Endpoints de administración del trial

```
GET  /admin/trial-config          → ver configuración actual
PUT  /admin/trial-config          → editar duración, habilitar/deshabilitar
                                    { is_enabled: bool, duration_days: int, grace_days: int }
```

Estos endpoints requieren un token de **super-admin** (credencial interna de AxonPOS, no del propietario del restaurante).

---

## 9. Path a Supabase self-hosted

Supabase Cloud para el MVP. Cuando el negocio justifique el costo operativo o haya requerimientos de datos en Colombia (Ley 1581), migrar a self-hosted es directo:

1. El Auth Service se abstrae detrás de una interfaz en Go — el POS Backend no llama a Supabase directamente, llama a la interfaz
2. Supabase self-hosted corre en Docker con el mismo API compatible
3. Solo cambia la `SUPABASE_URL` en las variables de entorno
4. Cero cambios de código

---

## 10. Notas de implementación Windows (.msi)

Tauri genera el `.msi` nativamente con el comando `tauri build --target x86_64-pc-windows-msvc`. Requerimientos:
- El sidecar Go debe compilarse para `GOOS=windows GOARCH=amd64`
- Incluir el `manifest` de permisos de red en Tauri para que Windows no bloquee el servidor local
- El firewall de Windows pedirá permiso la primera vez que el servidor local abra el puerto — documentar esto para el instalador
- Certificado de firma de código (`code signing`) para evitar advertencia de SmartScreen — recomendado antes de distribución pública
