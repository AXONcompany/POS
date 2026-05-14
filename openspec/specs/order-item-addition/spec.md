## ADDED Requirements

### Requirement: Adición de productos a orden existente con deducción de inventario

El sistema DEBE permitir agregar uno o más productos a una orden existente, obteniendo el precio real de cada producto desde la base de datos, calculando las deducciones de ingredientes según la receta técnica del producto, y persistiendo los items junto con el descuento de stock en una única transacción atómica.

#### Scenario: Adición exitosa con stock suficiente
- **WHEN** un usuario autorizado agrega uno o más productos a una orden en estado `PENDING` o `SENT`, y todos los ingredientes requeridos tienen stock suficiente
- **THEN** el sistema inserta los items en `order_items` con el precio real del producto, descuenta los ingredientes en `ingredients` y actualiza `orders.total_amount`, todo en una sola transacción; retorna HTTP 200 con la orden actualizada

#### Scenario: Stock insuficiente de un ingrediente
- **WHEN** un usuario intenta agregar un producto cuya receta requiere más cantidad de un ingrediente de la que hay disponible en stock
- **THEN** el sistema rechaza toda la operación con HTTP 409 (`ErrInsufficientStock`), no persiste ningún cambio (rollback completo) y retorna un mensaje de error descriptivo

#### Scenario: Orden en estado no editable
- **WHEN** un usuario intenta agregar productos a una orden en estado `PREPARING`, `READY`, `PAID` o `CANCELLED`
- **THEN** el sistema rechaza la operación con HTTP 422 (`ErrInvalidStatusTransition`) sin modificar la orden ni el inventario

#### Scenario: Producto sin receta técnica configurada
- **WHEN** se agrega un producto que no tiene ingredientes asociados en su receta
- **THEN** el sistema agrega el item a la orden con precio real sin descontar ningún ingrediente, y retorna HTTP 200 con la orden actualizada

#### Scenario: Múltiples items que comparten ingrediente
- **WHEN** se agregan dos o más productos en un mismo request que comparten el mismo ingrediente en su receta
- **THEN** el sistema acumula la deducción total de ese ingrediente y ejecuta un único UPDATE en `ingredients` para ese ingrediente, garantizando consistencia y evitando conflictos de concurrencia

### Requirement: Cancelación de ítem con restauración de inventario
El sistema DEBE permitir cancelar un ítem de una orden existente marcándolo como cancelado (`cancelled_at`), recalcular el total de la orden y restaurar el stock de ingredientes correspondientes a la receta de dicho ítem en una única transacción atómica.

#### Scenario: Cancelación exitosa
- **WHEN** un usuario autorizado cancela un ítem activo de una orden en estado `PENDING` o `SENT`
- **THEN** el sistema actualiza `cancelled_at` en `order_items`, restaura el stock en `ingredients` basándose en la receta del ítem y reduce el `total_amount` de la orden, todo en una sola transacción; retorna HTTP 200 con la orden actualizada.

#### Scenario: Ítem ya cancelado
- **WHEN** un usuario intenta cancelar un ítem que ya tiene `cancelled_at IS NOT NULL`
- **THEN** el sistema rechaza la operación con HTTP 409 (`ErrItemAlreadyCancelled`) sin modificar el inventario ni la orden.

#### Scenario: Orden en estado no editable para cancelación
- **WHEN** un usuario intenta cancelar un ítem de una orden en estado `PREPARING`, `READY`, `PAID` o `CANCELLED`
- **THEN** el sistema rechaza la operación con HTTP 422 (`ErrInvalidStatusTransition`) sin modificar el ítem ni el inventario.
