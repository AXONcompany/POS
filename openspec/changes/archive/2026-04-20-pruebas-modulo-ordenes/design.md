## Context

El módulo de órdenes tiene dos capas con lógica relevante: el usecase (`internal/usecase/order/usecase.go`) y el handler HTTP (`internal/infrastructure/rest/order/handler.go`). El `usecase_test.go` existente usa `testify/mock` con mocks manuales y cubre los happy paths principales. Los handlers no tienen ningún test.

## Goals / Non-Goals

**Goals:**
- Cubrir error paths faltantes en el usecase: `CancelOrderItem` (repo errors), `DivideOrder` (tipos no triviales, repo error, tipo inválido), `GetDivisionsByOrder`
- Cubrir `AddProductToOrder` en estado SENT y errores de precio/receta
- Crear `handler_test.go` con httptest que verifique mapeos HTTP: 400, 409, 422, 500 para los handlers críticos

**Non-Goals:**
- Tests de integración con BD real (ya cubierto por `pos_happy_path_integration_test.go`)
- Tests de los handlers de auth, payment, table u otros módulos
- Refactorizar los mocks existentes

## Decisions

**Usar httptest.NewRecorder + gin.New() para tests de handler**
Los handlers dependen de `*gin.Context`, así que se instancian con `gin.New()`, se registra la ruta directamente, y se dispara con `httptest.NewRecorder`. Esto permite testear el handler sin levantar servidor real. Alternativa descartada: mockear `*gin.Context` directamente — demasiado frágil.

**Inyectar usecase como interfaz en handler tests**
Para los handler tests se define una interfaz mínima `orderUsecase` que el handler acepta, y se provee un stub simple (sin testify/mock) que devuelve valores controlados. Esto evita acoplar los handler tests a la implementación del usecase.

**Reutilizar los mocks existentes en usecase_test.go**
Los `MockRepository`, `MockProductInventoryRepository` y `MockAuditRepository` ya están declarados en el archivo de test. Los nuevos sub-tests los reutilizan directamente.

## Risks / Trade-offs

- Los handler tests verifican el mapeo de error→HTTP pero no la lógica de negocio (eso ya lo hace el usecase test) — duplicación mínima, separación clara.
- Si el handler cambia el formato de response, los tests fallan: aceptable, es el comportamiento deseado.

## Open Questions

_(ninguna)_
