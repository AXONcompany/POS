## REMOVED Requirements

### Requirement: Administracion de Entidad de Restaurante
**Reason**: La entidad `Restaurant` se reemplaza completamente por `Venue` (Sede) bajo la nueva jerarquia Owner > Venue.
**Migration**: Usar los nuevos endpoints de `venue-management` en lugar de los endpoints de restaurant. Los datos se migran automaticamente en la migracion 000011.

### Requirement: Control de Empleados del Restaurante
**Reason**: Se reemplaza por la gestion de empleados a nivel de Venue dentro de `venue-management`. Los empleados ahora se vinculan a `venue_id` en lugar de `restaurant_id`.
**Migration**: El campo `restaurant_id` en `users` se migra a `venue_id`. La logica de registro de empleados usa `venue_id` del contexto JWT del owner.
