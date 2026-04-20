## Context

El handler `DivideOrder` calcula las divisiones al vuelo y las retorna al cliente, pero no las persiste. Esto significa que cada vez que el cajero consulta la división tiene que recalcularla, y no hay forma de vincular un pago a una parte específica de la cuenta. La tabla `order_divisions` y la columna `payments.division_ref` ya están diseñadas en el plan de implementación pero no tienen migración ni código todavía.

## Goals / Non-Goals

**Goals:**
- Persistir las divisiones de cuenta en `order_divisions` al momento de dividir.
- Permitir que un pago referencie opcionalmente una división (`division_ref`).
- Exponer un endpoint para obtener las divisiones activas de una orden.
- Mantener compatibilidad: si no hay divisiones, el flujo de pago sigue igual.

**Non-Goals:**
- Validar que todas las divisiones estén pagadas antes del checkout (eso es el ítem 5 del plan).
- Calcular propinas por división.
- Modificar o eliminar divisiones ya creadas.

## Decisions

### 1. Divisiones se reemplazan al re-dividir (no se acumulan)

- **Decisión**: `DivideOrder` borra las divisiones existentes de la orden e inserta las nuevas en una sola TX.
- **Razón**: Un cajero puede querer cambiar la forma de dividir antes de cobrar. Acumular divisiones crearía inconsistencias (ej: 3 partes iguales + 2 por monto).
- **Alternativa**: Bloquear la re-división si alguna parte ya fue pagada. Descartado por complejidad excesiva para MVP.

### 2. `division_ref` es nullable en payments

- **Decisión**: La columna `payments.division_ref` es `VARCHAR(50) NULL` con FK a `order_divisions.id`.
- **Razón**: La gran mayoría de órdenes no se dividen. Hacer el campo obligatorio rompería el flujo de pago normal.
- **Alternativa**: Tabla separada `payment_divisions`. Descartado — complejidad innecesaria.

### 3. IDs de división como UUID generados en Go

- **Decisión**: El ID de cada división se genera en Go con `fmt.Sprintf("div_%d_%d", orderID, index)` (determinista y legible).
- **Razón**: Permite idempotencia y facilita el debugging. PostgreSQL usa `VARCHAR(50)` como PK.
- **Alternativa**: `uuid.New()` (random). Descartado — menos legible, no aporta beneficio real en este contexto.

### 4. Transacción en repositorio (patrón existente)

- **Decisión**: `CreateDivisions` abre una `pgx.Tx` que borra las divisiones previas e inserta las nuevas.
- **Razón**: Coherente con el patrón ya establecido en `order_repository.go::Create` y `AddItemsWithInventory`.
- **Alternativa**: Manejar TX en el usecase. Descartado para mantener persistencia aislada.

## Risks / Trade-offs

- **[Riesgo] Re-división con parte ya pagada**: Si un cajero re-divide una cuenta donde una división ya tiene un pago asociado, esa división será eliminada y el pago quedará con `division_ref` apuntando a un registro inexistente. **Mitigación MVP**: En la implementación, verificar que no existan pagos vinculados antes de borrar divisiones — retornar error si los hay.
- **[Trade-off] IDs deterministas**: Si dos requests simultáneos dividen la misma orden, ambos generan los mismos IDs y la segunda TX fallará por PK duplicado. **Mitigación**: La operación no es concurrente en un escenario MVP de un cajero por terminal.
