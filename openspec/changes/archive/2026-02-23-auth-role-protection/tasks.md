# Tasks: auth-role-protection

- [x] **Task 1: Update Auth Middleware for Hierarchical Roles**
  - Modify `RequireRole(roleID int)` in `internal/infrastructure/rest/middleware/auth.go` to accept variadic arguments `RequireRoles(roles ...int)` OR build a role hierarchy check where `RolePropietario` (1) has permission for anything `RoleCajero` (2) or `RoleMesero` (3) can do.

- [x] **Task 2: Protect Public Endpoints behind AuthMiddleware**
  - In `internal/infrastructure/rest/routes.go`, create a generalized API group (e.g., `protected := r.Group("/api/v1")`) that uses `middleware.AuthMiddleware(jwtSecret)`.
  - Move the currently public routes for `/tables`, `/ingredients`, `/categories`, `/products`, and `/menu` inside this protected block.

- [x] **Task 3: Apply Role-Based Access Control to Menus and Inventory**
  - Apply `RequireRoles(RolePropietario)` to the destructive and modification methods of `/ingredients`, `/categories`, and `/products`.
  - Apply `RequireRoles(RolePropietario)` to the POST/PATCH methods of `/menu`.
  - Apply `RequireRoles(RoleMesero, RoleCajero, RolePropietario)` (or equivalent hierarchy) to the GET method of `/menu`.

- [x] **Task 4: Apply Role-Based Access Control to Tables**
  - Apply `RequireRoles(RoleCajero, RolePropietario)` to the POST, PATCH, DELETE, and AssignWaitress methods of `/tables`.
  - Apply `RequireRoles(RoleMesero, RoleCajero, RolePropietario)` to the GET method of `/tables`.

- [x] **Task 5: Refactor existing Orders Endpoint to use updated Middleware**
  - Update the checkout route inside `/orders` to use the new `RequireRoles(RoleCajero, RolePropietario)` structure instead of the exact match `RequireRole(RoleCajero)`.

- [x] **Task 6: Verification and Testing**
  - Run `go run cmd/server/main.go` and simulate requests using different JWTs (Waiter, Cashier, Owner) to ensure 403 Forbidden is appropriately thrown when limits are exceeded, and 200 OK when permitted.
