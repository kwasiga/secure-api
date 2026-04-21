# Secure API

A production-ready REST API in Go featuring JWT authentication, role-based access control, per-IP rate limiting, and PostgreSQL persistence.

---

## Stack

| Package | Purpose |
|---|---|
| Go | Application language |
| PostgreSQL | Primary database |
| sqlc | Type-safe SQL query generation |
| chi | HTTP router with middleware grouping |
| golang-jwt/jwt | JWT signing and validation (HS256) |
| bcrypt | Password hashing (cost 14) |
| golang.org/x/time/rate | Token-bucket rate limiting |
| godotenv | Optional `.env` file loading |

---

## Project Structure

```
.
├── cmd/api/            # Application entrypoint (main.go)
├── config/             # Environment config loading
├── db/
│   ├── queries/        # Raw SQL queries (input to sqlc)
│   ├── schema.sql      # Table definitions and enums
│   └── sqlc/           # sqlc-generated type-safe Go code
└── internal/
    ├── auth/           # JWT generation/validation, password hashing
    ├── handler/        # HTTP handlers (auth, user, admin)
    ├── middleware/     # Auth, RBAC, and rate-limiter middleware
    ├── model/          # App-layer data types
    ├── repository/     # Database access layer
    └── validator/      # Request decoding and struct validation
```

---

## Setup

### Option A — Docker (recommended)

Requires [Docker](https://docs.docker.com/get-docker/) and Docker Compose.

```bash
# Clone the repo
git clone https://github.com/kwasiga/secure-api
cd secure-api

# Start the API and a Postgres instance together
JWT_SECRET=your-secret-key docker compose up --build
```

The database schema is applied automatically on first boot — no migration step needed. The API will be available at `http://localhost:8080`.

To stop:

```bash
docker compose down
```

To wipe the database volume as well:

```bash
docker compose down -v
```

---

### Option B — Local (Go + Postgres)

#### Prerequisites

- Go 1.22+
- PostgreSQL running locally

#### Steps

1. **Clone the repo**

   ```bash
   git clone https://github.com/kwasiga/secure-api
   cd secure-api
   ```

2. **Create the database and apply the schema**

   ```bash
   createdb secure_api
   psql secure_api < db/schema.sql
   ```

3. **Configure environment**

   Copy the example env file and fill in your values:

   ```bash
   cp .env.example .env
   ```

   ```env
   PORT=8080
   DATABASE_URL=postgres://user:password@localhost:5432/secure_api?sslmode=disable
   JWT_SECRET=your-secret-key
   ```

4. **Install dependencies**

   ```bash
   go mod tidy
   ```

5. **Run the server**

   ```bash
   go run cmd/api/main.go
   ```

---

## API Reference

### Public routes

| Method | Path | Description |
|---|---|---|
| `POST` | `/register` | Create a new user account |
| `POST` | `/login` | Authenticate and receive tokens |

#### POST /register

```json
{
  "email": "user@example.com",
  "first_name": "Jane",
  "last_name": "Doe",
  "password": "password123"
}
```

Returns `201 Created` with the new user object.

#### POST /login

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Returns:

```json
{
  "access_token": "<jwt>",
  "refresh_token": "<jwt>"
}
```

The **access token** expires in 15 minutes. The **refresh token** expires in 7 days.

---

### Protected routes (JWT required)

Include the access token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

| Method | Path | Description |
|---|---|---|
| `GET` | `/me` | Get the authenticated user's profile |
| `PUT` | `/me` | Update first and last name |

#### PUT /me

```json
{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

---

### Admin routes (JWT + `admin` role required)

| Method | Path | Description |
|---|---|---|
| `GET` | `/admin/users` | List all registered users |

---

## Middleware

### AuthMiddleware

Reads the `Authorization: Bearer <token>` header, validates the JWT with HS256, and stores the parsed claims (user ID and role) in the request context. Returns `401` if the token is missing or invalid.

### RequireRole

Reads the claims already stored by `AuthMiddleware` and checks that the user's role is in the allowed list. Returns `401` if claims are absent, `403` if the role is not permitted.

```go
r.Use(middleware.RequireRole("admin"))
```

### RateLimiter

Per-IP token-bucket limiter backed by `golang.org/x/time/rate`. Configured globally at 10 requests/second with a burst of 20. Clients inactive for more than 3 minutes are evicted from memory. Returns `429 Too Many Requests` when the bucket is empty.

---

## Auth Details

- Passwords are hashed with **bcrypt** at cost 14 before storage
- Login issues two tokens signed with **HS256**:
  - **Access token** — short-lived (15 min), used for API requests
  - **Refresh token** — long-lived (7 days), intended for obtaining new access tokens
- Both tokens carry `user_id` and `role` claims
