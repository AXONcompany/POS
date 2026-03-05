# AXON POS API Documentation

Esta es la documentación oficial de los endpoints de la API REST del sistema AXON POS.

## Base URL
Todas las peticiones en producción van dirigidas al host raíz (ej. `http://localhost:8080`). 
No hay prefijo `/api/v1` en la versión actual del router.

## Autenticación y Autorización
La API utiliza JSON Web Tokens (JWT) pasados mediante el header `Authorization: Bearer <token>`.
Los endpoints protegidos validan el esquema RBAC (Control de Acceso Basado en Roles):
- `Role 1`: PROPIETARIO
- `Role 2`: CAJERO
- `Role 3`: MESERO

---

## 1. Salud y Monitoreo (Health)
| Método | Endpoint  | Acceso  | Descripción                               |
|--------|-----------|---------|-------------------------------------------|
| `GET`  | `/health` | Público | Devuelve estado `200 OK` si el server vive|
| `GET`  | `/ping`   | Público | Endpoint de prueba simple (pong)          |

---

## 2. Autenticación y Usuarios (`/auth`, `/usuarios`)
| Método | Endpoint               | Rol Permitido         | Descripción                                                                 |
|--------|------------------------|-----------------------|-----------------------------------------------------------------------------|
| `POST` | `/auth/login`          | Público               | Login para cualquier usuario. Retorna JWT access+refresh token.             |
| `POST` | `/auth/register-owner` | Público               | Registra un PROPIETARIO y su primera Sede (`Venue`). Crea sesión autologin. |
| `POST` | `/auth/register`       | PROPIETARIO, CAJERO*  | Crear usuario genérico (*Cajeros solo pueden crear meseros).              |
| `GET`  | `/auth/me`             | Sesión Válida         | Obtiene la información del usuario en base al token JWT.                    |
| `POST` | `/auth/logout`         | Sesión Válida         | Revoca la sesión (invalidación del _refresh_token_ en BD).                  |
| `GET`  | `/usuarios`            | PROPIETARIO           | Lista todos los usuarios de la Sede.                                        |
| `POST` | `/usuarios/mesero`     | PROPIETARIO, CAJERO   | Crea un mesero con contraseña autogenerada (se devuelve 1 sola vez en texto)|
| `GET`  | `/usuarios/:id`        | PROPIETARIO           | Obtiene detalles de un usuario específico.                                  |
| `PATCH`| `/usuarios/:id`        | PROPIETARIO           | Actualiza datos de un usuario de la sede.                                   |

---

## 3. Configuración Administrativa (Sedes y Terminales)
| Método | Endpoint          | Rol Permitido         | Descripción                                          |
|--------|-------------------|-----------------------|------------------------------------------------------|
| `GET`  | `/propietario`    | PROPIETARIO           | Detalles de la cuenta Padre del Franquiciado.        |
| `PATCH`| `/propietario`    | PROPIETARIO           | Edita perfil del propietario.                        |
| `GET`  | `/sedes`          | PROPIETARIO           | Lista sedes asociadas al propietario autenticado.    |
| `POST` | `/sedes`          | PROPIETARIO           | Crea una nueva Sede.                                 |
| `PATCH`| `/sedes/:id`      | PROPIETARIO           | Edita nombre, teléfono o dirección de una Sede.      |
| `GET`  | `/terminales`     | PROPIETARIO           | Lista Cajas registradoras (POS) asignadas a la sede.|
| `POST` | `/terminales`     | PROPIETARIO           | Crea una nueva caja.                                 |
| `PATCH`| `/terminales/:id` | PROPIETARIO           | Edita nombre o suspende caja.                        |

---

## 4. Mesas (`/mesas`)
| Método  | Endpoint           | Rol Permitido         | Descripción                                         |
|---------|--------------------|-----------------------|-----------------------------------------------------|
| `GET`   | `/mesas`           | MESERO, CAJERO, PROP  | Obtiene inventario de mesas de la Sede y su estado. |
| `POST`  | `/mesas`           | CAJERO, PROPIETARIO   | Crea una mesa. Rq: `{numero, capacidad}`            |
| `PATCH` | `/mesas/:id/estado`| CAJERO, PROPIETARIO   | Cambia la mesa a libre, ocupada, sucia, etc.        |

---

## 5. Control de Inventario y Menú (`/ingredientes`, `/categorias`, `/products`, `/menu`)
| Método | Endpoint                | Rol Permitido         | Descripción                                          |
|--------|-------------------------|-----------------------|------------------------------------------------------|
| `POST` | `/ingredientes`         | PROPIETARIO           | Registra insumo de bodega.                           |
| `PATCH`| `/ingredientes/:id/stock`| PROPIETARIO          | Suman/Restan al inventario bodega.                   |
| `POST` | `/categorias`           | PROPIETARIO           | Crea categoría de menú (Bebidas, Fuertes, etc.).     |
| `POST` | `/menu`                 | PROPIETARIO           | Crea ítem de venta, asociándolo a la receta técnica.|
| `GET`  | `/menu`                 | MESERO, CAJERO, PROP  | Catálogo con precios de venta y descripción.         |

---

## 6. Comandas y Ventas (`/ordenes`, `/pagos`)
| Método  | Endpoint                       | Rol Permitido         | Descripción                                                |
|---------|--------------------------------|-----------------------|------------------------------------------------------------|
| `POST`  | `/ordenes`                     | MESERO, CAJERO, PROP  | Inicia factura/pedido asociando mesa.                      |
| `POST`  | `/ordenes/:id/items`           | MESERO, CAJERO, PROP  | Agrega platos/bebidas al pedido.                           |
| `POST`  | `/ordenes/:id/enviar-cocina`   | MESERO, CAJERO, PROP  | Cierra adición y transfiere al Kitchen Display.            |
| `POST`  | `/ordenes/:id/dividir`         | CAJERO, PROPIETARIO   | Divide cuenta (por montos o partes iguales).               |
| `POST`  | `/pagos`                       | CAJERO, PROPIETARIO   | Realiza el cobro, marcando cierre de venta.                |
| `GET`   | `/pagos/:id/factura`           | CAJERO, PROPIETARIO   | Emite recibo (Ticket Fiscal local o ticket digital).       |

---

## 7. Reportes Gerenciales (`/reportes`)
Requieren Rol **PROPIETARIO**. El CAJERO no puede acceder a reportes consolidados, solo la consola gerencial.
| Método | Endpoint                 | Parámetros `Querystring` | Descripción |
|--------|--------------------------|--------------------------|-------------|
| `GET`  | `/reportes/ventas`       | `fecha_inicio`, `fecha_fin` | Resumen de ventas netas, consolidadas por categorías. |
| `GET`  | `/reportes/inventario`   | ninguno                  | Cantidades disponibles de todos los ingredientes crudos. |
| `GET`  | `/reportes/propinas`     | `fecha_inicio`, `fecha_fin` | Total de propinas procesadas y asiganción de tickets. |
