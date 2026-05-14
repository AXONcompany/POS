openapi: 3.0.0
info:
  title: AXONPOS
  description: Endpoints basicos
  version: 1.0.0
  
servers:
  - url: https://api.AxonPos.com/v1
    description: Respiralo ahí

tags:
  - name: Autenticación
    description: Login, registro y sesiones
  - name: Propietario
    description: Gestión del perfil del propietario
  - name: Sedes
    description: Gestión de sedes (venues)
  - name: Terminales POS
    description: Gestión de terminales POS por sede
  - name: Mesas
    description: Gestion de mesas y salones
  - name: Ordenes
    description: Gestion de ordenes y pedidos
  - name: División de cuenta
    description: División y consulta de divisiones de una orden
  - name: Pagos
    description: Procesamiento de pagos
  - name: Menu
    description: Gestion del menu
  - name: Inventario
    description: Control de inventario
  - name: Reportes
    description: Reportes basicos

paths:
  /mesas:
    get:
      tags:
        - Mesas
      summary: Listar todas las mesas
      description: Obtiene la lista de mesas con su estado actual
      responses:
        '200':
          description: Lista de mesas (revisar si es mas conveniente paginacion)
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Mesa'
    
    post:
      tags:
        - Mesas
      summary: Crear una mesa
      description: Crea una nueva mesa
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - numero
                - capacidad
              properties:
                numero:
                  type: string
                  example: "1"
                capacidad:
                  type: integer
                  example: 4
      responses:
        '201':
          description: Mesa creada exitosamenet
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Mesa'

  /mesas/{id}:
    get:
      tags:
        - Mesas
      summary: Obtener mesa
      description: Obtiene los detalles de una mesa por ID
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Detalles de la mesa
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Mesa'
        '404':
          description: Mesa no encontrada
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

    delete:
      tags:
        - Mesas
      summary: Eliminar mesa
      description: Elimina una mesa (solo CAJERO o PROPIETARIO)
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Mesa eliminada
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Mesa eliminada exitosamente"
        '404':
          description: Mesa no encontrada
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /mesas/{id}/estado:
    patch:
      tags:
        - Mesas
      summary: Cambiar estado de mesa
      description: Actualiza el estado de una mesa (libre, ocupada, reservada, limpiando)
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - estado
              properties:
                estado:
                  type: string
                  enum: [libre, ocupada, reservada, limpiando]
                  example: "ocupada"
      responses:
        '200':
          description: Estado actualizado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Mesa'
  /auth/register:
    post:
      tags:
        - Autenticación
      summary: Registrar usuario
      description: Registra un nuevo usuario en el sistema (solo ADMIN puede crear usuarios)
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - nombre
                - email
                - password
                - rol
              properties:
                nombre:
                  type: string
                  example: "Juan Pérez"
                email:
                  type: string
                  format: email
                  example: "juan.perez@gmail.com"
                password:
                  type: string
                  format: password
                  example: "123456"
                  minLength: 8
                rol:
                  type: string
                  enum: [CAJERO, MESERO]
                  example: "MESERO"
                telefono:
                  type: string
                  example: "3001234567"
      responses:
        '201':
          description: Usuario registrado exitosamente
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Usuario registrado exitosamente"
                  data:
                    $ref: '#/components/schemas/Usuario'
        '400':
          description: Datos invalids
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: No autorizado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '403':
          description: Sin permisos
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '409':
          description: El email ya esta registrado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: false
                  error:
                    type: string
                    example: "El email ya esta registrado en el sistema"

  /auth/register-owner:
    post:
      tags:
        - Autenticación
      summary: Registrar propietario
      description: Crea una cuenta de propietario con su primera sede. Endpoint público.
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - nombre
                - email
                - password
                - nombre_sede
              properties:
                nombre:
                  type: string
                  example: "Carlos Dueño"
                email:
                  type: string
                  format: email
                  example: "carlos@restaurante.com"
                password:
                  type: string
                  format: password
                  example: "secreto123"
                nombre_sede:
                  type: string
                  example: "Sede Centro"
                direccion:
                  type: string
                  example: "Calle 1 # 2-3"
                telefono:
                  type: string
                  example: "3001234567"
      responses:
        '201':
          description: Propietario y sede creados
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      access_token:
                        type: string
                      refresh_token:
                        type: string
                      user:
                        $ref: '#/components/schemas/Usuario'
        '409':
          description: El email ya está registrado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/login:
    post:
      tags:
        - Autenticación
      summary: Iniciar sesion
      description: Autentica un usuario y devuelve un token JWT
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - email
                - password
              properties:
                email:
                  type: string
                  format: email
                  example: "juan.perez@gmail.com"
                password:
                  type: string
                  format: password
                  example: "123456"
      responses:
        '200':
          description: Login exitoso
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      token:
                        type: string
                        example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                      usuario:
                        $ref: '#/components/schemas/Usuario'
                      expires_in:
                        type: integer
                        example: 3600
                        description: "Tiempo de expiración del token"
        '400':
          description: Datos invalidos
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: Credenciales incorrectas
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: false
                  error:
                    type: string
                    example: "Email o contraseña incorrectos"

  /auth/me:
    get:
      tags:
        - Autenticación
      summary: Obtener usuario actual
      description: Devuelve la informacion del usuario autenticado
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Informacion del usuario
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Usuario'
        '401':
          description: Token invalido o expirado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /auth/logout:
    post:
      tags:
        - Autenticación
      summary: Cerrar sesion
      description: Invalida el token actual del usuario
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Sesion cerrada exitosamente
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Sesion cerrada exitosamente"
        '401':
          description: No autorizado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /usuarios:
    get:
      tags:
        - Autenticación
      summary: Listar usuarios
      description: Lista todos los usuarios del sistema (solo ADM)
      security:
        - BearerAuth: []
      parameters:
        - name: rol
          in: query
          schema:
            type: string
            enum: [PROPIETARIO, CAJERO, MESERO]
          description: Filtrar por rol
        - name: activo
          in: query
          schema:
            type: boolean
          description: Filtrar por estado activo/inactivo
      responses:
        '200':
          description: Lista de usuarios
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Usuario'
        '401':
          description: No autorizado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '403':
          description: Sin permisos (solo ADMIN)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /usuarios/{id}:
    get:
      tags:
        - Autenticación
      summary: Obtener usuario
      description: Obtiene informacion de un usuario especifico
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Informacion del usuario
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Usuario'
        '404':
          description: Usuario no encontrado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

    patch:
      tags:
        - Autenticación
      summary: actualizar usuario
      description: actualiza la información de un usuario (ADMIN)
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                nombre:
                  type: string
                email:
                  type: string
                  format: email
                rol:
                  type: string
                  enum: [ADMIN, MESERO, CAJA]
                activo:
                  type: boolean
                telefono:
                  type: string
      responses:
        '200':
          description: Usuario actualizado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Usuario'
        '403':
          description: Sin permisos
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: Usuario no encontrado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

    delete:
      tags:
        - Autenticación
      summary: Eliminar usuario
      description: Desactiva un usuario del sistema (solo ADMIN)
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Usuario eliminado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Usuario desactivado exitosamente"
        '403':
          description: Sin permisos
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: Usuario no encontrado
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /usuarios/mesero:
    post:
      tags:
        - Autenticación
      summary: Registrar mesero
      description: El PROPIETARIO o CAJERO crea un nuevo mesero. Retorna las credenciales temporales generadas.
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - nombre
                - email
              properties:
                nombre:
                  type: string
                  example: "Ana Mesero"
                email:
                  type: string
                  format: email
                  example: "ana@restaurante.com"
      responses:
        '201':
          description: Mesero creado con contraseña temporal
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      usuario:
                        $ref: '#/components/schemas/Usuario'
                      password_temporal:
                        type: string
                        example: "xK9mP2qR"

  /ordenes:
    get:
      tags:
        - Ordenes
      summary: Listar ordenes por mesa
      description: Lista todas las ordenes de una mesa
      security:
        - BearerAuth: []
      parameters:
        - name: table_id
          in: query
          required: true
          schema:
            type: integer
          description: ID de la mesa
      responses:
        '200':
          description: Lista de ordenes
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Orden'

    post:
      tags:
        - Ordenes
      summary: Crear una orden
      description: Crea una nueva orden para una mesa
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - mesa_id
                - mesero_id
              properties:
                mesa_id:
                  type: string
                  example: "mesa 1"
                mesero_id:
                  type: string
                  example: "777"
      responses:
        '201':
          description: Orden creada
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Orden'

  /ordenes/{id}:
    get:
      tags:
        - Ordenes
      summary: Obtener orden
      description: Obtiene los detalles de una orden
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Detalles de la orden
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Orden'

  /ordenes/{id}/items:
    post:
      tags:
        - Ordenes
      summary: Agregar items a orden
      description: Agrega items del menu a una orden 
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - items
              properties:
                items:
                  type: array
                  items:
                    type: object
                    required:
                      - menu_item_id
                      - cantidad
                    properties:
                      menu_item_id:
                        type: string
                        example: "item 1"
                      cantidad:
                        type: integer
                        example: 2
                      notas:
                        type: string
                        example: "Sin cebolla soy gay"
      responses:
        '200':
          description: Items agregados
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Orden'

  /ordenes/{id}/items/{item_id}:
    delete:
      tags:
        - Ordenes
      summary: Cancelar item de orden
      description: Cancela un item de una orden 
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
        - name: item_id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Item cancelado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Item cancelado exitosamente"

  /ordenes/{id}/status:
    patch:
      tags:
        - Ordenes
      summary: Actualizar estado de orden
      description: Cambia el estado de una orden. Transiciones válidas — PENDING(1)→SENT(2)→PREPARING(3)→READY(4)→PAID(5) o CANCELLED(6)
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - status_id
              properties:
                status_id:
                  type: integer
                  enum: [1, 2, 3, 4, 5, 6]
                  example: 3
      responses:
        '200':
          description: Estado actualizado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Estado de orden actualizado"
        '422':
          description: Transición de estado inválida
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /ordenes/{id}/checkout:
    post:
      tags:
        - Pagos
      summary: Checkout de orden
      description: Marca la orden como PAID (estado 5). Requiere que la orden esté en estado READY(4). Solo CAJERO o PROPIETARIO.
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Pago procesado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Pago procesado"
        '422':
          description: La orden no está lista para pago
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /ordenes/{id}/enviar-cocina:
    post:
      tags:
        - Ordenes
      summary: Enviar orden a cocina
      description: Envia la orden a la cocina xd
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Orden enviada a cocina
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  message:
                    type: string
                    example: "Orden enviada a cocina y bar"

  /menu:
    get:
      tags:
        - Menu
      summary: Listar menú
      description: Obtiene todos los items del menu organizado por cetegorias
      parameters:
        - name: categoria
          in: query
          schema:
            type: string
            enum: [entradas, platos_fuertes, bebidas, postres]
      responses:
        '200':
          description: Lista del menu
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/MenuItem'

    post:
      tags:
        - Menu
      summary: Crear item de menu
      description: Crea un nuevo item en el men
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - nombre
                - categoria
                - precio
              properties:
                nombre:
                  type: string
                  example: "Hamburguesa Placa blanca"
                categoria:
                  type: string
                  enum: [entradas, platos_fuertes, bebidas, postres]
                  example: "platos_fuertes"
                precio:
                  type: number
                  example: 25000
                descripcion:
                  type: string
                  example: "Hamburguesa de carne"
                disponible:
                  type: boolean
                  example: true
      responses:
        '201':
          description: Item creado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/MenuItem'

  /menu/{id}:
    patch:
      tags:
        - Menu
      summary: Actualizar item de menu
      description: Actualiza un item existente
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                nombre:
                  type: string
                precio:
                  type: number
                disponible:
                  type: boolean
      responses:
        '200':
          description: Item actualizado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/MenuItem'

  /ordenes/{id}/divisiones:
    get:
      tags:
        - División de cuenta
      summary: Consultar divisiones
      description: Retorna las divisiones activas de una orden con su estado de pago
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Lista de divisiones
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Division'

  /ordenes/{id}/dividir:
    post:
      tags:
        - División de cuenta
      summary: Dividir cuenta
      description: Divide la cuenta de una orden (revisar), cada division representa lo que va a pagar cada parte cuando seqa requerido (por item o por monto)
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - tipo_division
              properties:
                tipo_division:
                  type: string
                  enum: [partes_iguales, por_item, por_monto]
                  example: "partes_iguales"
                numero_partes:
                  type: integer
                  example: 2
                divisiones:
                  type: array
                  items:
                    type: object
                    properties:
                      items:
                        type: array
                        items:
                          type: string
                      monto:
                        type: number
      responses:
        '200':
          description: Cuenta dividida
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      type: object
                      properties:
                        division_id:
                          type: string
                          example: "div_1_1"
                        subtotal:
                          type: number
                          example: 35000
                        impuestos:
                          type: number
                          example: 6650
                        total:
                          type: number
                          example: 41650
                        is_paid:
                          type: boolean
                          example: false

  /pagos:
    post:
      tags:
        - Pagos
      summary: Procesar pago
      description: Procesa el pago de una orden
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - orden_id
                - metodo_pago
                - monto
              properties:
                orden_id:
                  type: string
                  example: "orden_001"
                division_id:
                  type: string
                  example: "div_001"
                metodo_pago:
                  type: string
                  enum: [efectivo, tarjeta, multiple]
                  example: "tarjeta"
                monto:
                  type: number
                  example: 41650
                propina:
                  type: number
                  example: 4165
                detalles_pago:
                  type: object
                  properties:
                    efectivo:
                      type: number
                    tarjeta:
                      type: number
                    referencia_tarjeta:
                      type: string
      responses:
        '200':
          description: Pago procesado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Pago'

  /pagos/{id}/factura:
    get:
      tags:
        - Pagos
      summary: Generar factura
      description: Genera la factura xd
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Factura generada
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      factura_id:
                        type: string
                        example: "FACT-2024-001234"
                      fecha:
                        type: string
                        format: date-time
                      codigo:
                        type: string
                        example: "abc123..."
                      url_pdf:
                        type: string
                        example: "factura.pdf"

  /ingredientes:
    get:
      tags:
        - Inventario
      summary: Listar ingredientes
      description: listado de ingredientes_bajo_stock
      responses:
        '200':
          description: Lista de ingredientes
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Ingrediente'

    post:
      tags:
        - Inventario
      summary: Crear ingrediente
      description: Registra un nuevo ingrediente
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - nombre
                - unidad_medida
                - stock_actual
                - stock_minimo
              properties:
                nombre:
                  type: string
                  example: "Carne molida"
                unidad_medida:
                  type: string
                  example: "kg"
                stock_actual:
                  type: number
                  example: 15.5
                stock_minimo:
                  type: number
                  example: 5
                costo_unitario:
                  type: number
                  example: 18000
      responses:
        '201':
          description: Ingrediente creado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Ingrediente'

  /ingredientes/{id}/stock:
    patch:
      tags:
        - Inventario
      summary: Actualizar stock
      description: Actualiza la cantidad en stock de un ingrediente
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - cantidad
                - tipo_movimiento
              properties:
                cantidad:
                  type: number
                  example: 10
                tipo_movimiento:
                  type: string
                  enum: [entrada, salida]
                  example: "entrada"
                motivo:
                  type: string
                  example: "Compra a proveedor"
      responses:
        '200':
          description: Stock actualizado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Ingrediente'


  /reportes/ventas:
    get:
      tags:
        - Reportes
      summary: Reporte de ventas
      description: Genera reporte de ventas por periodo
      parameters:
        - name: fecha_inicio
          in: query
          required: true
          schema:
            type: string
            format: date
        - name: fecha_fin
          in: query
          required: true
          schema:
            type: string
            format: date
        - name: tipo
          in: query
          schema:
            type: string
            enum: [por_item, por_dia, por_hora]
            default: por_dia
      responses:
        '200':
          description: Reporte generado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      periodo:
                        type: object
                        properties:
                          inicio:
                            type: string
                            format: date
                          fin:
                            type: string
                            format: date
                      total_ventas:
                        type: number
                        example: 2450000
                      total_ordenes:
                        type: integer
                        example: 87
                      ticket_promedio:
                        type: number
                        example: 28161
                      detalle:
                        type: array
                        items:
                          type: object

  /reportes/inventario:
    get:
      tags:
        - Reportes
      summary: Reporte de inventario
      description: Muestra el estado actual del inventario con alertas
      responses:
        '200':
          description: Reporte de inventario
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: object
                    properties:
                      ingredientes_bajo_stock:
                        type: array
                        items:
                          $ref: '#/components/schemas/Ingrediente'
                      valor_total_inventario:
                        type: number
                        example: 3500000

  /reportes/propinas:
    get:
      tags:
        - Reportes
      summary: Reporte de propinas
      description: Reporte de propinas por mesero y perioso
      parameters:
        - name: fecha_inicio
          in: query
          required: true
          schema:
            type: string
            format: date
        - name: fecha_fin
          in: query
          required: true
          schema:
            type: string
            format: date
      responses:
        '200':
          description: Reporte de propinas
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      type: object
                      properties:
                        mesero_id:
                          type: string
                        mesero_nombre:
                          type: string
                        total_propinas:
                          type: number
                          example: 125000
                        numero_ordenes:
                          type: integer
                          example: 15

  /propietario:
    get:
      tags:
        - Propietario
      summary: Obtener perfil del propietario
      description: Retorna los datos del propietario autenticado
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Perfil del propietario
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Propietario'

    patch:
      tags:
        - Propietario
      summary: Actualizar perfil del propietario
      security:
        - BearerAuth: []
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                nombre:
                  type: string
                  example: "Carlos Dueño"
      responses:
        '200':
          description: Perfil actualizado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Propietario'

  /sedes:
    get:
      tags:
        - Sedes
      summary: Listar sedes
      description: Lista todas las sedes del propietario autenticado
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Lista de sedes
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Sede'

    post:
      tags:
        - Sedes
      summary: Crear sede
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - nombre
              properties:
                nombre:
                  type: string
                  example: "Sede Norte"
                direccion:
                  type: string
                  example: "Carrera 15 # 80-10"
                telefono:
                  type: string
                  example: "3109876543"
      responses:
        '201':
          description: Sede creada
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Sede'

  /sedes/{id}:
    get:
      tags:
        - Sedes
      summary: Obtener sede
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Detalles de la sede
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Sede'

    patch:
      tags:
        - Sedes
      summary: Actualizar sede
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                nombre:
                  type: string
                direccion:
                  type: string
                telefono:
                  type: string
                activo:
                  type: boolean
      responses:
        '200':
          description: Sede actualizada
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Sede'

  /terminales:
    get:
      tags:
        - Terminales POS
      summary: Listar terminales
      description: Lista los terminales POS de la sede del propietario
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Lista de terminales
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Terminal'

    post:
      tags:
        - Terminales POS
      summary: Crear terminal
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - terminal_name
                - venue_id
              properties:
                terminal_name:
                  type: string
                  example: "Caja 1"
                venue_id:
                  type: integer
                  example: 1
      responses:
        '201':
          description: Terminal creado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Terminal'

  /terminales/{id}:
    get:
      tags:
        - Terminales POS
      summary: Obtener terminal
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Detalles del terminal
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Terminal'

    patch:
      tags:
        - Terminales POS
      summary: Actualizar terminal
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                terminal_name:
                  type: string
                activo:
                  type: boolean
      responses:
        '200':
          description: Terminal actualizado
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                    example: true
                  data:
                    $ref: '#/components/schemas/Terminal'

components:
  schemas:
    Mesa:
      type: object
      properties:
        id:
          type: integer
          example: 4
        number:
          type: string
          example: "Table 4"
        state:
          type: string
          enum: [free, occupied, reserved, cleaning]
          example: "occupied"
        capacity:
          type: integer
          example: 4
        waiter:
          type: object
          nullable: true
          properties:
            id:
              type: integer
              example: 15
            name:
              type: string
              example: "Luis"
        guests:
          type: integer
          nullable: true
          example: 2
          description: "Número de comensales actuales"
        occupiedMinutes:
          type: integer
          nullable: true
          example: 12
          description: "Minutos desde que la mesa fue ocupada"
        currentOrder:
          type: object
          nullable: true
          properties:
            id:
              type: string
              example: "ORD-12784"
            total:
              type: number
              format: float
              example: 32.00
    Usuario:
      type: object
      properties:
        id:
          type: string
          example: "usr_001"
        nombre:
          type: string
          example: "Juan Pérez"
        email:
          type: string
          format: email
          example: "juan.perez@gmail.com"
        rol:
          type: string
          enum: [PROPIETARIO, CAJERO, MESERO]
          example: "MESERO"
        telefono:
          type: string
          example: "3001234567"
        activo:
          type: boolean
          example: true
        fecha_creacion:
          type: string
          format: date-time
          example: "2024-12-01T10:30:00Z"
        ultimo_acceso:
          type: string
          format: date-time
          example: "2024-12-09T14:20:00Z"

    Error:
      type: object
      properties:
        success:
          type: boolean
          example: false
        error:
          type: string
          example: "Mensaje de error descriptivo"
        code:
          type: string
          example: "UNAUTHORIZED"

    MenuItem:
      type: object
      properties:
        id:
          type: string
          example: "item_001"
        nombre:
          type: string
          example: "Hamburguesa Clásica"
        categoria:
          type: string
          example: "platos_fuertes"
        precio:
          type: number
          example: 25000
        descripcion:
          type: string
        disponible:
          type: boolean
          example: true

    Orden:
      type: object
      properties:
        id:
          type: string
          example: "orden #1"
        mesa_id:
          type: string
          example: "mesa 1"
        mesero_id:
          type: string
          example: "777"
        estado:
          type: string
          enum: [abierta, enviada, pagada, cancelada]
          example: "abierta"
        items:
          type: array
          items:
            type: object
            properties:
              id:
                type: string
              menu_item_id:
                type: string
              nombre:
                type: string
              cantidad:
                type: integer
              precio_unitario:
                type: number
              notas:
                type: string
              estado:
                type: string
                enum: [pendiente, en_cocina, listo, entregado, cancelado]
        subtotal:
          type: number
          example: 50000
        impuestos:
          type: number
          example: 9500
        total:
          type: number
          example: 59500
        fecha_creacion:
          type: string
          format: date-time

    Pago:
      type: object
      properties:
        id:
          type: string
          example: "pago_001"
        orden_id:
          type: string
        metodo_pago:
          type: string
          example: "tarjeta"
        monto:
          type: number
          example: 59500
        propina:
          type: number
          example: 5950
        total:
          type: number
          example: 65450
        estado:
          type: string
          enum: [pendiente, aprobado, rechazado]
          example: "aprobado"
        referencia:
          type: string
        fecha:
          type: string
          format: date-time

    Ingrediente:
      type: object
      properties:
        id:
          type: string
          example: "ing_001"
        nombre:
          type: string
          example: "Carne molida"
        unidad_medida:
          type: string
          example: "kg"
        stock_actual:
          type: number
          example: 15.5
        stock_minimo:
          type: number
          example: 5
        costo_unitario:
          type: number
          example: 18000

    Division:
      type: object
      properties:
        division_id:
          type: string
          example: "div_1_1"
        subtotal:
          type: number
          example: 35000
        impuestos:
          type: number
          example: 6650
        total:
          type: number
          example: 41650
        is_paid:
          type: boolean
          example: false

    Propietario:
      type: object
      properties:
        id:
          type: integer
          example: 1
        nombre:
          type: string
          example: "Carlos Dueño"
        email:
          type: string
          format: email
          example: "carlos@restaurante.com"
        activo:
          type: boolean
          example: true

    Sede:
      type: object
      properties:
        id:
          type: integer
          example: 1
        nombre:
          type: string
          example: "Sede Centro"
        direccion:
          type: string
          example: "Calle 1 # 2-3"
        telefono:
          type: string
          example: "3001234567"
        activo:
          type: boolean
          example: true

    Terminal:
      type: object
      properties:
        id:
          type: integer
          example: 1
        terminal_name:
          type: string
          example: "Caja 1"
        venue_id:
          type: integer
          example: 1
        activo:
          type: boolean
          example: true

  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - BearerAuth: []