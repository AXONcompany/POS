## ADDED Requirements

### Requirement: Cobertura HTTP del handler de órdenes
El sistema DEBE tener tests que verifiquen que cada handler del módulo de órdenes retorna el código HTTP correcto y un body JSON estructurado para los escenarios de éxito y error de dominio. Los tests usan `httptest` y un stub del usecase sin tocar la BD.

#### Scenario: CreateOrder retorna 201 con body de orden
- **WHEN** el handler `CreateOrder` recibe un request válido y el usecase retorna una orden
- **THEN** el handler responde con HTTP 201 y un JSON con el campo `data`

#### Scenario: CreateOrder retorna 500 si el usecase falla
- **WHEN** el usecase de CreateOrder retorna un error interno
- **THEN** el handler responde con HTTP 500

#### Scenario: AddItems retorna 409 por stock insuficiente
- **WHEN** el usecase de AddProductToOrder retorna `ErrInsufficientStock`
- **THEN** el handler responde con HTTP 409 y código `INSUFFICIENT_STOCK`

#### Scenario: AddItems retorna 422 por estado inválido
- **WHEN** el usecase de AddProductToOrder retorna `ErrInvalidStatusTransition`
- **THEN** el handler responde con HTTP 422 y código `INVALID_TRANSITION`

#### Scenario: CancelItem retorna 409 si el item ya está cancelado
- **WHEN** el usecase de CancelOrderItem retorna `ErrItemAlreadyCancelled`
- **THEN** el handler responde con HTTP 409 y código `ITEM_ALREADY_CANCELLED`

#### Scenario: CancelItem retorna 422 por estado inválido
- **WHEN** el usecase de CancelOrderItem retorna `ErrInvalidStatusTransition`
- **THEN** el handler responde con HTTP 422

#### Scenario: DivideOrder retorna 409 por división ya pagada
- **WHEN** el usecase de DivideOrder retorna `ErrDivisionAlreadyPaid`
- **THEN** el handler responde con HTTP 409 y código `DIVISION_ALREADY_PAID`

#### Scenario: DivideOrder retorna 200 con lista de divisiones
- **WHEN** el usecase de DivideOrder retorna divisiones calculadas
- **THEN** el handler responde con HTTP 200 y un JSON con el array de divisiones

#### Scenario: GetDivisions retorna 200 con divisiones existentes
- **WHEN** el usecase de GetDivisionsByOrder retorna divisiones para la orden
- **THEN** el handler responde con HTTP 200 y el array de divisiones

#### Scenario: CheckoutOrder retorna 422 por estado inválido
- **WHEN** el usecase de CheckoutOrder retorna `ErrInvalidStatusTransition`
- **THEN** el handler responde con HTTP 422
