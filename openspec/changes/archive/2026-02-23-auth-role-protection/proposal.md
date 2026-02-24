# Proposal: auth-role-protection

## Description
This change introduces Role-Based Access Control (RBAC) across all API endpoints in the POS backend architecture. Currently, only the `/orders` endpoints are protected by authentication, leaving sensitive operations like managing the menu, ingredients, tables, and categories completely exposed. This proposal aims to secure the system using JWT-based authentication combined with specific role requirements (Owner, Cashier, Waiter) tailored to the daily operations of a restaurant.

## Goals
- Secure all API endpoints under a general protected group.
- Implement strict RBAC according to the restaurant's hierarchy:
  - **Waiters (Meseros - Role 3):** Can view tables, view the menu, and create orders.
  - **Cashiers (Cajeros - Role 2):** Can do everything a Waiter can do + manage table statuses/assignments, and checkout/close orders.
  - **Owners/Admins (Propietarios - Role 1):** Can do everything a Cashier can do + manage the menu, products, ingredients, categories, and view inventory reports.
- Prevent unauthorized users from performing destructive actions (DELETE, POST, PUT, PATCH) on core domain entities.

## Non-Goals
- Changing the existing JWT token generation format.
- Adding new roles beyond the currently defined `RolePropietario`, `RoleCajero`, and `RoleMesero`.
- Implementing fine-grained row-level security in PostgreSQL (access control is handled at the HTTP routing layer).

## Rationale
In a restaurant setting, it's critical that staff can only access features relevant to their job. Exposing `/ingredients` to a waiter could result in accidental deletion of stock, ruining inventory reports. By applying the existing `AuthMiddleware` and `RequireRole` middleware consistently across `routes.go`, we can leverage the robust Clean Architecture we recently implemented to safely deliver a production-ready MVP.
