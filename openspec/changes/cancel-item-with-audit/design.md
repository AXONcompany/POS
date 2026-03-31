## Context

El sistema actual maneja órdenes e inventario, pero no tiene un mecanismo de "borrado lógico" para los ítems de las órdenes ni una forma de revertir el impacto de estos en el inventario de manera atómica. Además, la falta de un sistema de auditoría impide la trazabilidad de cambios críticos realizados por los usuarios.

## Goals / Non-Goals

**Goals:**
- Implementar la cancelación de ítems (`Soft Delete`) en `order_items`.
- Garantizar la restauración del stock de ingredientes al cancelar un ítem.
- Proveer un sistema de auditoría basado en instantáneas (Snapshots) JSONB.
- Asegurar la atomicidad de la cancelación y el registro de auditoría en una única transacción de base de datos.

**Non-Goals:**
- Implementar un motor de diferencias (diff) complejo en Go.
- Implementar triggers de base de datos para la auditoría.
- Manejar la restauración de stock para órdenes ya pagadas (PAID).

## Decisions

### 1. Snapshot-based Audit (Full JSON)
- **Decisión**: Guardar el objeto completo en las columnas `old_value` y `new_value` de tipo `JSONB`.
- **Razón**: Es más simple de implementar en el MVP y permite una recuperación total del estado anterior sin lógica adicional. PostgreSQL permite consultar campos específicos dentro del JSONB si es necesario en el futuro.
- **Alternativa**: Guardar solo los campos que cambiaron (diff). Descartado por complejidad excesiva para esta fase.

### 2. Transacción en la Capa de Repositorio
- **Decisión**: Envolver la actualización del ítem, la restauración del stock, el ajuste del total de la orden y la inserción del log de auditoría en un solo `pgx.Tx` dentro del repositorio.
- **Razón**: Evita estados inconsistentes (ej: ítem cancelado pero stock no restaurado).
- **Alternativa**: Manejar la transacción en la capa de Usecase. Descartado para mantener la persistencia aislada y evitar fugar tipos de `pgx` al dominio.

### 3. Registro de Auditoría en la Capa de Aplicación (Go)
- **Decisión**: El Usecase captura el estado actual, ejecuta el cambio y guarda el log de auditoría.
- **Razón**: Permite incluir contexto que la DB no conoce, como el ID del usuario real que realiza la acción a través de la API.

## Risks / Trade-offs

- **[Riesgo] Carga de Almacenamiento**: Guardar snapshots completos puede aumentar el tamaño de la base de datos si hay muchos cambios.
  - **Mitigación**: Los ítems de orden son pequeños. Se puede implementar una política de purga de logs antiguos en el futuro.
- **[Riesgo] Disciplina del Desarrollador**: Si se añade una nueva funcionalidad que modifica datos sin llamar a la auditoría, el cambio no se registrará.
  - **Mitigación**: Crear una interfaz de repositorio que obligue a pasar un contexto de auditoría o encapsular la lógica de guardado de auditoría.
