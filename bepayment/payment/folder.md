internal/
├── domain/
│   ├── user/
│   │   ├── entity.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── errors.go
│   └── order/
│       ├── entity.go
│       ├── model.go
│       └── repository.go
│
├── application/
│   └── user/
│       ├── create_user.go
│       ├── get_user.go
│       └── dto.go
│
├── infrastructure/
│   ├── persistence/
│   │   └── postgres/
│   │       └── user_repository.go
│   └── sqlc/
│       ├── db.go
│       └── queries.sql.go
│
├── interfaces/
│   └── http/
│       ├── handler/
│       │   └── user_handler.go
│       └── router.go
│
├── conv/
│   └── pg.go
│
└── main.go


internal/
├── user/
│   ├── domain/
│   │   └── user.go
│   ├── repository/
│   │   └── postgres/
│   │       ├── query.sql
│   │       ├── models.go      (generated)
│   │       ├── query.sql.go   (generated)
│   │       └── sqlc.yaml
│
├── order/
│   ├── domain/
│   │   └── order.go
│   ├── repository/
│   │   └── postgres/
│   │       ├── query.sql
│   │       ├── models.go      (generated)
│   │       ├── query.sql.go   (generated)
│   │       └── sqlc.yaml
│
└── shared/
    └── schema/
        └── schema.sql

