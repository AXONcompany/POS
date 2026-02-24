# Access Control Spec Delta

## MODIFIED Requirements

### Requirement: Require Roles by endpoint hierarchy
The system MUST assign access to the modules in `routes.go` depending on the current user's role. Waiters (Role 3), Cashiers (Role 2), and Owners/Admins (Role 1) will have nested permissions. Waiters only have read logic, Cashiers have table allocation rights, and Owners have full CRUD logic on everything.

#### Scenario: Waiter trying to access endpoints
- Given a logged-in user with Role 3 (Waitress/Mesero)
- When they make a `GET` request to `/tables` or `/menu`
- Then the system allows the request
- When they make a `POST` request to `/orders`
- Then the system allows the request
- But when they make a `POST` request to `/tables` or `/ingredients`
- Then the system returns a 403 Forbidden.

#### Scenario: Cashier trying to access endpoints
- Given a logged-in user with Role 2 (Cashier/Cajero)
- When they make a `POST` or `PATCH` request to `/tables`
- Then the system allows the request
- When they make a `DELETE` request to `/tables/:id`
- Then the system allows the request
- But when they make a `POST` request to `/menu` or `/ingredients`
- Then the system returns a 403 Forbidden.

#### Scenario: Owner accessing endpoints
- Given a logged-in user with Role 1 (Owner/Propietario)
- When they make a `DELETE`, `PUT`, or `POST` request to `/ingredients`, `/products`, `/categories`, or `/menu`
- Then the system allows the request.
