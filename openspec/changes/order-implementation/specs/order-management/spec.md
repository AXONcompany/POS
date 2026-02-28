## ADDED Requirements

### Requirement: Crear Pedido
El sistema MUST permitir a los usuarios autorizados crear un nuevo pedido asociado a una mesa, especificando los productos, sus cantidades, y asociando la venta transaccional.

#### Scenario: Creación exitosa de pedido
- **WHEN** un usuario envía una solicitud válida con lista de productos, mesa ID
- **THEN** el sistema valida que los productos existan y tengan disponibilidad
- **THEN** el sistema registra el pedido en estado inicial y devuelve el ID del pedido creado

### Requirement: Listar Pedidos Activos
El sistema MUST proveer un endpoint para listar pedidos filtrados para seguimiento operativo en restaurante.

#### Scenario: Listar pedidos pendientes de una mesa
- **WHEN** una solicitud busca listar los pedidos de una mesa en curso
- **THEN** el sistema responde con una lista estructurada que incluye el total y detalle de los productos ordenados

### Requirement: Actualizar Estado del Pedido
El sistema MUST proveer un mecanismo explícito para avanzar un pedido en su ciclo de vida (e.g. Recibido -> Preparando -> Listo -> Entregado).

#### Scenario: Cambio de estado de pedido
- **WHEN** el usuario actualiza el estado de un pedido enviando el nuevo identificador de estado
- **THEN** el sistema actualiza el registro atómicamente y se refleja en subsecuentes listas de pedidos
