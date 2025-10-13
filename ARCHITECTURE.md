# Sistema POS Backend (Go)

## Descripción general

El backend del sistema POS (Point of Sale) es una API RESTful desarrollada en Go (Golang) que proporciona servicios para la gestión de ventas, inventario, usuarios, reportes y sincronización con sistemas externos (por ejemplo, contabilidad o facturación electrónica).

El objetivo principal es ofrecer un servicio rápido, modular y escalable, que pueda ejecutarse tanto en local (tienda) como en la nube, manteniendo sincronización asíncrona en caso de pérdida de conectividad.

## Arquitectura general

El sistema sigue una arquitectura Clean Architecture / Hexagonal, separando la lógica de negocio del acceso a datos y del transporte HTTP.
Esto permite mantener el dominio desacoplado y fácilmente testeable.

```
+---------------------------------------------------------+
|                     Presentation Layer                  |
|   (HTTP handlers, GraphQL resolvers, gRPC endpoints)    |
+---------------------------+-----------------------------+
|        Application Layer  |                            |
|  (Use cases / Services)   |                            |
+---------------------------+-----------------------------+
|           Domain Layer                                   |
|   (Entities, Aggregates, Value Objects, Interfaces)      |
+---------------------------------------------------------+
|          Infrastructure Layer                            |
|   (DB adapters, Repositories, External APIs, Logger)     |
+---------------------------------------------------------+
```

## Estructura de carpetas

```
pos-backend/
├── cmd/
│   └── server/
│       └── main.go          # Punto de entrada de la aplicación
│
├── internal/
│   ├── domain/              # Entidades y lógica de negocio pura
│   │   ├── sale/
│   │   │   ├── sale.go
│   │   │   └── sale_test.go
│   │   ├── product/
│   │   └── user/
│   │
│   ├── usecase/             # Casos de uso (application layer)
│   │   ├── sale_service.go
│   │   ├── product_service.go
│   │   └── user_service.go
│   │
│   ├── infrastructure/      # Adaptadores hacia sistemas externos
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   └── sqlite/
│   │   ├── rest/
│   │   │   ├── handler_sale.go
│   │   │   ├── handler_user.go
│   │   │   └── router.go
│   │   ├── messaging/       # Pub/Sub, RabbitMQ, Kafka, etc.
│   │   └── logging/
│   │
│   └── config/
│       └── config.go        # Manejo centralizado de configuración
│
├── pkg/                     # Librerías reutilizables (shared utils)
│   ├── errors/
│   ├── jwt/
│   └── pagination/
│
├── migrations/              # Scripts SQL
├── docs/                    # Documentación (Swagger, OpenAPI)
└── go.mod
```

## Tecnologías principales

| Componente | Tecnología |
|------------|------------|
| Lenguaje | Go 1.23+ |
| Framework HTTP | Gin |
| ORM / SQL builder |sqlx |
| Base de datos | PostgreSQL |
| Mensajería | RabbitMQ / NATS |
| Auth | JWT con HMAC / RSA |
| Configuración | viper |
| Logs | zerolog |
| Tests | testing, testify, dockertest |
