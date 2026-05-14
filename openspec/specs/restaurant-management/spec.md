## REMOVED Requirements

> Este spec fue reemplazado por `venue-management` y `owner-management` en la migración de `restaurants` → `venues/owners`.

### ~~Requirement: Administración de Entidad de Restaurante~~
Reemplazado por `venue-management`: gestión de sedes (venues) vinculadas a un `owner`.

### ~~Requirement: Control de Empleados del Restaurante~~
Ahora los empleados se asocian mediante `venue_id` (no `restaurant_id`). Ver `venue-management` y `user-auth`.
