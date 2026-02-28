# POS Backend - API Reference

## Base URLs
- Public: `/` o `/auth`
- Protected API: `/api/v1`

Todas las rutas bajo `/api/v1` requieren un Access Token válido (Autenticación Bearer por JWT).
Roles de Autorización soportados por RBAC (Role-Based Access Control):
- **Propietario** (Role ID: 1)
- **Cajero** (Role ID: 2)
- **Mesero** (Role ID: 3)

---

## 1. Core & Health Checks
Rutas públicas y utilitarias.

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `GET` | `/health` | Chequeo de estado HTTP JSON. | Pública |
| `GET` | `/ping` | Chequeo de estado HTTP string. | Pública |

---

## 2. Authentication
Rutas relativas al manejo de sesiones y credenciales de usuario.

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `POST` | `/auth/login` | Iniciar sesión y retornar AccessToken & RefreshToken. | Pública |

*(Para refresco de sesión o registro de usuarios verificar especificaciones extra de Auth endpoints dependiendo de la implementación)*

---

## 3. Tables (Mesas)
Rutas protegidas bajo Auth y dependientes de roles del restaurante. Base URL: `/api/v1`.

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `GET` | `/tables` | Obtener todas las mesas. | Mesero, Cajero, Propietario |
| `GET` | `/tables/:id` | Recuperar datos de mesa específica. | Mesero, Cajero, Propietario |
| `POST` | `/tables` | Crear una nueva mesa. | Cajero, Propietario |
| `PATCH` | `/tables/:id` | Actualizar estado o características de la mesa. | Cajero, Propietario |
| `DELETE`| `/tables/:id` | Eliminar mesa del sistema. | Cajero, Propietario |
| `POST` | `/tables/:id/assign` | Asignar un mesero a una mesa activa. | Cajero, Propietario |

---

## 4. Ingredients (Ingredientes)
Base URL: `/api/v1/ingredients`. (Toda la suite está restringida principalmente a administración e inventarios)

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `GET` | `/ingredients` | Listar ingredientes. | Propietario |
| `GET` | `/ingredients/report`| Obtener reporte del stock general de inventario. | Propietario |
| `POST` | `/ingredients` | Registrar un nuevo ingrediente al catálogo. | Propietario |
| `GET` | `/ingredients/:id` | Detalles de un ingrediente. | Propietario |
| `PUT` | `/ingredients/:id` | Actualizar existencias o datos. | Propietario |
| `DELETE`| `/ingredients/:id` | Remover ingrediente. | Propietario |

---

## 5. Categories (Categorías de Productos)
Base URL: `/api/v1/categories`. (Utilizado para el menú y base de productos)

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `GET` | `/categories` | Listar todas las categorías (ej., Entradas, Bebidas). | Propietario |
| `POST` | `/categories` | Crear nueva categoría. | Propietario |

---

## 6. Products (Productos)
Base URL: `/api/v1/products`. Productos genéricos o materia prima con dependencias.

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `GET` | `/products` | Listar todos los productos base. | Propietario |
| `POST` | `/products` | Registrar producto base en inventario. | Propietario |
| `GET` | `/products/:id/ingredients`| Ver la receta/ingredientes usados en un producto. | Propietario |
| `POST`| `/products/:id/ingredients`| Vincular ingredientes a su receta (deducciones futuras). | Propietario |

---

## 7. Menu (Menú Público/Interno del Restaurante)
Base URL: `/api/v1/menu`. Listado curado diseñado para que los meseros levanten pedidos.

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `GET` | `/menu`| Obtiene todos los Ítems de Menú activos (Productos finales). | Mesero, Cajero, Propietario |
| `POST` | `/menu` | Empaqueta un producto como Ítem Oficial del Menú. | Propietario |
| `PATCH`| `/menu/:id`| Habilitar, deshabilitar, o cambiar precio del menú. | Propietario |

---

## 8. Orders (Órdenes de Cliente)
Base URL: `/api/v1/orders`. Motor central de operaciones del Punto de Venta.

| Método | Endpoint | Acción | Autorización |
| --- | --- | --- | --- |
| `GET` | `/orders` | Listar todas las ordenes. | Mesero, Cajero, Propietario |
| `POST` | `/orders` | Crear nueva orden adjuntando Items y opcional Mesa. | Mesero, Cajero, Propietario |
| `PATCH`| `/orders/:id/status`| Mover estados (PENDING -> PREPARING -> READY). | Mesero, Cajero, Propietario |
| `POST` | `/orders/:id/checkout`| Cerrar orden y pasarlo a "PAID", calculando total. | Cajero, Propietario |
