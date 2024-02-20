# Money Forward interview

## Requirements
The server api should serve for this endpoints:

* `GET /users/{id}`
* `GET /users/{id}/accounts`
* `GET /accounts/{id}`
  

All endpoints should return json payloads matching the following criteria.
`/users/{id}`: A detail of `users` including list `account_ids` that have been owned by the user that are equal to the passing `{id}` from the database.

Example:
```json
{
  "id": 1,
  "name": "Alice",
  "account_ids": [
    1,
    3,
    5,
    ...
  ]
}

``` 

`/users/{id}/accounts`: A list detail of `accounts` that have been owned by the user that are equal to the passing `{id}` from the database.

Example:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "A銀行",
    "balance": 20000
  },
  {
    "id": 3,
    "user_id": 1,
    "name": "C信用金庫",
    "balance": 120000
  },
  {
    "id": 5,
    "user_id": 1,
    "name": "E銀行",
    "balance": 5000
  },
  ...
]
``` 
`/accounts/{id}`: A detail of `accounts` that are equal to the passing `{id}` from the database.

Example:
```json
{
  "id": 2,
  "user_id": 2,
  "name": "Bカード",
  "balance": 200
}
``` 

## Features (the APIs that not occurs in requirements):

- [x] Finish three APIs in **Requirements**.
- [x] Adding mutation APIs for `accounts`, `users`.
- [x] Adding login flow.
- [x] Apply caching for optimize API.
- [x] Adding some middlewares hor http handle steps like `recovery`,`cors`,`authenticate`,`rbac`.


# Architecture: 
The Architecture using the [clean architecture](https://raw.githubusercontent.com/phungvandat/clean-architecture/dev/images/clean-arch.png) that control the system split to 3 layers.

![clean architecture](https://raw.githubusercontent.com/phungvandat/clean-architecture/dev/images/clean-arch.png)


# Built With

This section list any major frameworks/libraries used to bootstrap project.
* ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
* ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) 

## Libraries:

- Using [snowflake](github.com/bwmarrin/snowflake) engine for generate id.
- Using [paseto](github.com/o1egl/paseto) engine for generate token.
- Using [lru](github.com/hashicorp/golang-lru/v2) for in memory caching.
- Using [pq](github.com/lib/pq) for postgres driver.
- Using [cobra](github.com/spf13/cobra) for generate command line.

# Folder structure
```sh
.
├── LICENSE
├── Makefile   #for quickly command line
├── README.md
├── app-exe     # binary of go build
├── cmd         # contain command line for running
│   ├── root.go
│   ├── srv.go
│   └── start.go
├── configs       # contain config object and environment loader
│   ├── address.go
│   ├── configs.go
│   └── database.go
├── deployments   # using for deployments
├── developments  # using for developments include docker, env file
│   ├── dev.env     # environment file for dev environment
│   └── docker-compose.yml   # docker compose file for dev environment
├── go.mod
├── go.sum
├── internal  # contain layers that applying clean architecture without export to outside this module.
│   ├── deliveries # contain delivery/transport layer of clean architecture
│   │   └── http
│   │       ├── account.go
│   │       ├── auth.go
│   │       └── user.go
│   ├── entities # contain entities (data transfer object)
│   │   ├── account.go
│   │   └── user.go
│   ├── models # contain models that including request, response.
│   │   ├── account.go
│   │   ├── auth.go
│   │   ├── common.go
│   │   └── user.go
│   ├── repositories # contain repository/store layer of clean architecture
│   │   ├── account.go
│   │   └── user.go
│   └── services # contain service/domain layer of clean architecture
│       ├── account.go
│       ├── auth.go
│       └── user.go
├── main.go
├── migrations # contain migration files for database
│   ├── 00001_migrate.up.sql
│   └── 00002_migrate.up.sql
└── pkg    
    ├── cache # contain interface of cache pattern
    │   └── cache.go
    ├── crypto_utils # contain password util 
    │   └── util.go
    ├── database # contain database util
    │   ├── executor.go
    │   ├── type.go
    │   └── util.go
    ├── http_server # contain http server that follow native http lib by go
    │   ├── common.go
    │   ├── http.go
    │   ├── http_test.go
    │   ├── middleware.go
    │   ├── response.go
    │   ├── util.go
    │   ├── util_test.go
    │   └── xcontext  # contain context of http handler
    │       ├── context.go
    │       └── ctx.go
    ├── id_utils  # for id utility
    │   ├── id.go
    │   └── snowflake.go  # snowflake id generator
    ├── logger  # for logger
    │   └── logger.go
    ├── lru # for lru cache
    │   └── cache.go
    ├── postgres_client # postgres client
    │   ├── client.go
    │   └── tx.go
    ├── processor
    │   └── processor.go
    ├── reflect_utils # contain reflect utility
    │   ├── util.go
    │   └── util_test.go
    └── token_utils    # contain token utility
        ├── authenticator.go
        ├── jwt.go
        └── paseto.go
```


# Prerequisites

- Make sure you have Go installed ([download](https://golang.org/dl/)). Version `1.21` or higher is required.
- Docker (version `20.10.22+`)

# Getting started
First of all, we should set up environments:

```sh
  copy ./developments/.env.example ./developments/.env
  copy ./developments/dev.env.example ./developments/dev.env
```

1. Starting database:

```sh
  make start-db
```

2. Migrate:

```sh
  make migrate
```

3. Simple start server:

```sh
  make start
```

# Action flows:

We already migrate a default super admin user by **admin** and **donkihote**.

```sh
  curl --location 'localhost:8080/auth/login' \
    --header 'Content-Type: application/json' \
    --data '{
        "user_name":"admin",
        "password": "donkihote"
    }'
```


Create user:

```sh
  curl --location 'localhost:8080/users' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer ${given_token}' \
    --data '{
        "name": "Le Duy Dat",
        "user_name": "duyledat197",
        "password": "1234567",
        "role": "USER"
    }'
```

Get user detail (include account ids) by id:

```sh
curl --location 'localhost:8080/users/{id}' \
  --header 'Content-Type: application/json'
```

Create account by user id:

```sh
  curl --location 'localhost:8080/users/{id}/accounts' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer ${given_token}' \
    --data '{
      "name": "Account 1",
      "balance": 100
    }'
```

Get accounts by user id:

```sh
curl --location 'localhost:8080/users/{id}/accounts' \
  --header 'Content-Type: application/json'
```

Get account detail by id:

```sh
curl --location 'localhost:8080/accounts/{id}' \
  --header 'Content-Type: application/json'
```