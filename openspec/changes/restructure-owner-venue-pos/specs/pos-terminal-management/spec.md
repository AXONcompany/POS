## ADDED Requirements

### Requirement: Registro de Terminal POS por Sede
El sistema DEBE permitir crear terminales POS dentro de una venue. Cada terminal tiene un FK `venue_id` y un nombre identificador.

#### Scenario: Creacion de terminal POS
- **WHEN** se crea un pos_terminal con `venue_id` valido y nombre de terminal.
- **THEN** el sistema registra el terminal con `is_active = true` vinculado a la venue.

#### Scenario: Listado de terminales por venue
- **WHEN** se consultan los terminales con un `venue_id` especifico.
- **THEN** el sistema retorna solo los terminales que pertenecen a esa venue.

### Requirement: Vinculacion de Operaciones a Terminal
El sistema DEBE permitir que las ordenes y pagos referencien opcionalmente un `pos_terminal_id` para identificar que terminal proceso la operacion.

#### Scenario: Orden asociada a un terminal POS
- **WHEN** se crea una orden con `pos_terminal_id` valido.
- **THEN** la orden queda vinculada al terminal y a su venue correspondiente.
