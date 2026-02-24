# Design: Auth Role Protection

## Overview
The goal is to protect the endpoints in `routes.go` by grouping them under a generalized `/api/v1` block secured by `middleware.AuthMiddleware`. Within this protected group, we will apply `middleware.RequireRole()` on a per-subgroup basis to ensure Waiters, Cashiers, and Owners can only perform actions suited to their roles.

## Architecture & Implementation
The change focuses exclusively on `internal/infrastructure/rest/routes.go`. The existing endpoints and models remain the same; only their accessibility changes.

### 1. General Protected Group
We will wrap almost all routes (except `/health`, `/ping`, and the login route) in an authenticated base group.

### 2. Role Assignments & Subgroups
Roles are defined locally as constants in the routes file:
- `RolePropietario = 1`
- `RoleCajero = 2`
- `RoleMesero = 3`

The routing logic will enforce the following matrix:

| Endpoint | Method | Required Role | Rationale |
|----------|--------|---------------|-----------|
| `/menu` | GET | `RoleMesero` (and above) | Waiters need to see the menu to take orders. |
| `/menu` | POST, PATCH | `RolePropietario` | Only admins can change the menu structure. |
| `/tables` | GET | `RoleMesero` (and above) | Waiters need to see table statuses and capacities. |
| `/tables` | POST, PATCH, DELETE | `RoleCajero` (and above) | Cashiers or managers assign tables and update statuses. |
| `/orders` | POST | `RoleMesero` (and above) | Waiters create orders. |
| `/orders/:id/checkout` | POST | `RoleCajero` (and above) | Only Cashiers can receive payments. |
| `/ingredients` | GET, POST, PUT, DELETE | `RolePropietario` | Inventory manipulation is restricted to admins. |
| `/products` | GET, POST, PUT, DELETE | `RolePropietario` | Raw product management is restricted to admins. |
| `/categories` | GET, POST, PUT, DELETE | `RolePropietario` | Category management is restricted to admins. |

### 3. Middleware Strategy
A hierarchical approach will be utilized in Gin:
- Waiters naturally get the lowest requirement (`RequireRole(RoleMesero)`), meaning `RoleCajero` and `RolePropietario` implicitly inherit the access. *Note: If `RequireRole` is currently an exact match, it must be refactored to allow roles with higher privileges, or the authorization logic in `routes.go` must attach multiple accepted roles (e.g., `RequireRoles(RoleMesero, RoleCajero, RolePropietario)`).*

### 4. Technical Risks & Mitigation
- **Risk:** The current `RequireRole` middleware might only check for an exact role ID match rather than doing a hierarchical "greater than or equal to" permissions check. 
- **Mitigation:** We will inspect `internal/infrastructure/rest/middleware/auth.go`. If needed, we will modify `RequireRole` to accept variadic arguments `RequireRoles(roles ...int)` or create a hierarchy check to allow owners to do everything waiters can do.
