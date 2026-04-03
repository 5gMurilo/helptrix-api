# helptrix-api

REST API for the Helptrix service marketplace platform, where Business users request services from Helper users.

---

## About

Helptrix is a service marketplace that connects two types of users: **Business** users, who need services performed, and **Helper** users, who offer and fulfill those services. This API handles authentication, user management, service listings, and service proposals between those two roles.

### Key Features

- User registration with role, categories, and address
- Paseto-based authentication with 8-hour token expiration
- Password hashing with SHA-256 and salt before storage
- Clean Architecture with strict separation of concerns
- GORM AutoMigrate for schema management
- OpenAPI (Swagger) documentation for all endpoints

### Tech Stack

| Technology         | Purpose                                   |
|--------------------|-------------------------------------------|
| Go 1.25            | Primary language                          |
| Gin                | HTTP router and request handling          |
| GORM               | ORM and database migrations               |
| PostgreSQL 16      | Primary database                          |
| Paseto v2          | Symmetric token authentication            |
| swaggo/swag        | OpenAPI / Swagger documentation generator |
| air                | Hot reload during development             |
| Docker / Compose   | Containerized local environment           |

---

## Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- [air](https://github.com/air-verse/air) for hot reload (`go install github.com/air-verse/air@latest`)
- [swag](https://github.com/swaggo/swag) for Swagger generation (`go install github.com/swaggo/swag/cmd/swag@latest`)

---

## Installation

```bash
# Clone the repository
git clone https://github.com/5gMurilo/helptrix-api.git
cd helptrix-api

# Install dependencies
go mod download

# Copy and configure environment variables
cp .env.example .env
```

Edit `.env` with your local configuration before running the application.

---

## Usage

### Development with hot reload

```bash
air
```

air watches `.go` and `.toml` files, rebuilds the binary to `./tmp/main`, and restarts the server on every change.

### Docker Compose (full stack)

```bash
docker compose up
```

This starts a PostgreSQL 16 container and the API container. The API waits for the database health check to pass before starting. The application is available at `http://localhost:8080`.

### Build

```bash
go build -o helptrix-api ./app/...
```

The Dockerfile uses a two-stage build: a `golang:1.23-alpine` builder stage compiles the binary, and an `alpine:3.19` runtime stage runs only the compiled binary.

---

## Available Commands

| Command                                    | Description                               |
|--------------------------------------------|-------------------------------------------|
| `air`                                      | Start development server with hot reload  |
| `docker compose up`                        | Start full stack with Docker Compose      |
| `go test ./...`                            | Run all tests                             |
| `go test ./modules/auth/... -run <name>`   | Run a specific test by name               |
| `swag init -g app/main.go`                 | Regenerate Swagger documentation          |
| `go build -o helptrix-api ./app/...`       | Build the production binary               |

---

## Environment Variables

| Variable                | Description                                      | Default value    |
|-------------------------|--------------------------------------------------|------------------|
| `SERVER_PORT`           | Port the HTTP server listens on                  | `8080`           |
| `DB_HOST`               | PostgreSQL host                                  | `localhost`      |
| `DB_PORT`               | PostgreSQL port                                  | `5432`           |
| `DB_USER`               | PostgreSQL user                                  | `helptrix`       |
| `DB_PASSWORD`           | PostgreSQL password                              | `helptrix`       |
| `DB_NAME`               | PostgreSQL database name                         | `helptrix_db`    |
| `DB_SSLMODE`            | PostgreSQL SSL mode                              | `disable`        |
| `PASETO_SYMMETRIC_KEY`  | 32-byte hex key used for Paseto token signing    | —                |
| `GIN_MODE`              | Gin mode (`debug` or `release`)                  | `debug`          |

> `PASETO_SYMMETRIC_KEY` must be at least 32 characters long and provided as a hex-encoded string.

---

## Project Structure

```
adapter/
  db/
    repository/     # GORM repositories — all database operations
  auth/             # Paseto token configuration, creation, and verification
  http/
    middleware/     # Authentication middleware
    router.go       # All route definitions
core/
  domain/           # Entities and DTOs per business entity
  interfaces/       # Contracts (interfaces) grouped by entity
  utils/            # Utility functions, constants, and regex
modules/            # Feature modules — services and controllers per entity
app/
  main.go           # Entry point: dependency injection and application bootstrap
docs/               # Swagger output, architecture specs, and app-specs
```

### Module Responsibilities

| Layer          | Responsibility                                                                                     |
|----------------|----------------------------------------------------------------------------------------------------|
| `adapter/http` | Registers routes, maps HTTP errors to status codes, returns JSON responses                         |
| `modules`      | Contains controllers and services per entity; orchestrates business logic                          |
| `core/domain`  | Defines entities, request DTOs, and response DTOs                                                  |
| `core/interfaces` | Declares repository, service, and controller contracts used for dependency inversion            |
| `adapter/db`   | Establishes the database connection; repositories execute all DB operations inside a single transaction |
| `adapter/auth` | Manages Paseto symmetric key setup and token lifecycle                                             |

---

## Architecture

The project follows **Clean Architecture** organized as **Screaming Architecture** — folders are named after business domains, not technical layers.

**Dependency flow:**

```
adapter/http  -->  modules  -->  core/interfaces  <--  adapter/db/repository
```

All concrete dependencies are injected in `app/main.go`. Modules depend only on abstractions defined in `core/interfaces`, never on concrete implementations.

**Key patterns:**

- **Transactions**: Every multi-step database operation uses a single GORM transaction. Rollback is executed on any failure; commit only on full success.
- **Error propagation**: Repositories return descriptive error messages. Services forward them without HTTP context. Controllers map them to appropriate HTTP status codes.
- **Soft deletes**: All delete operations use GORM soft delete via the `deleted_at` field.
- **AutoMigrate**: On startup, GORM automatically migrates the schema for `User`, `Address`, `Category`, and `UserCategory`.

**SOLID principles** are enforced throughout: single responsibility per file, open/closed extensions via interfaces, and dependency inversion through injected abstractions.

---

## API Endpoints

| Method | Path             | Auth required | Description                     |
|--------|------------------|---------------|---------------------------------|
| GET    | `/health`        | No            | Health check                    |
| POST   | `/auth/register` | No            | Register a new user              |
| POST   | `/auth/login`    | No            | Authenticate and receive a token |

Protected route groups (`/user`, `/service`, `/proposal`) are defined and require a valid Paseto token in the `Authorization` header. Their specific endpoints are under active development.

### POST /auth/register

Registers a new user. Accepts a JSON body with the following fields:

| Field          | Type             | Description                                 |
|----------------|------------------|---------------------------------------------|
| `name`         | string           | Full name                                   |
| `email`        | string           | Email address (unique per role)             |
| `password`     | string           | Plain-text password (hashed before storage) |
| `role_id`      | UUID             | User role identifier                        |
| `categories`   | array of numbers | Category IDs to associate with the user     |
| `phone`        | string           | Contact phone number                        |
| `document`     | string           | CPF or CNPJ                                 |
| `address`      | object           | Street, number, complement, neighborhood, city, state |

Returns `201 Created` on success. Returns `409 Conflict` if a user with the same email and role already exists. Returns `500` on internal errors.

### POST /auth/login

Authenticates a user by email and password. Accepts:

| Field      | Type   | Description           |
|------------|--------|-----------------------|
| `email`    | string | Registered email      |
| `password` | string | Plain-text password   |

Returns a Paseto token valid for 8 hours on success.

---

## Authentication

The API uses **Paseto v2 symmetric tokens**.

- Tokens are created upon successful login and contain `user_id`, `name`, and `email` as claims.
- Token expiration is **8 hours** from issuance.
- Passwords are hashed with **SHA-256** before being stored in the database.
- The symmetric key must be provided via the `PASETO_SYMMETRIC_KEY` environment variable and must be at least 32 characters long.
- Protected routes require the token to be sent in the `Authorization` header. The auth middleware validates the token and rejects expired or malformed tokens.

---

*Created: 2026-04-03 | Last updated: 2026-04-03*
