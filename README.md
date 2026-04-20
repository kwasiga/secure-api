# Secure API

A production-ready REST API in Go with JWT authentication, role-based access control, per-IP rate limiting, and PostgreSQL persistence.

## Stack

| Package | Purpose |
|---|---|
| Go | Application language |
| PostgreSQL | Primary database |
| sqlc | Type-safe SQL query generation |
| golang-jwt/jwt | JWT signing and validation (HS256) |
| bcrypt | Password hashing |
| golang.org/x/time/rate | Token-bucket rate limiting |
| godotenv | Environment variable loading |

## Project Structure

```
.
├── cmd/api/            # Application entrypoint
├── config/             # Environment config loading
├── db/sqlc/            # sqlc-generated database code
└── internal/
    ├── auth/           # JWT generation/validation, password hashing
    ├── middleware/     # Auth, RBAC, and rate-limiter middleware
    ├── model/          # App-layer data types
    └── repository/     # Database access layer
```

## Setup

1. Copy `.env.example` to `.env` and fill in values:

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/secure_api
JWT_SECRET=your-secret-key
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run the server:

```bash
go run cmd/api/main.go
```

## Middleware

### Auth (`AuthMiddleware`)

Validates the JWT on every protected route. Reads the `Authorization: Bearer <token>` header, calls `auth.ValidateToken`, and attaches the parsed claims to the request context. Returns `401` if the token is missing or invalid.

### RBAC (`RequireRole`)

Wraps a handler and restricts access to users whose `role` claim matches one of the allowed roles. Returns `401` if claims are missing, `403` if the role is not permitted.

```go
mux.Handle("/admin", middleware.RequireRole("admin")(adminHandler))
```

### Rate Limiter (`RateLimiter`)

Per-IP token-bucket rate limiter backed by `golang.org/x/time/rate`. Inactive clients are evicted after 3 minutes. Returns `429 Too Many Requests` when the bucket is empty.

```go
rl := middleware.NewRateLimiter(10, 20) // 10 req/s, burst of 20
mux.Handle("/", rl.Limit(handler))
```

## Auth Details

- Passwords hashed with **bcrypt** at cost 14
- Login issues a short-lived **access token** (15 min) and a long-lived **refresh token** (7 days), both signed HS256
- Tokens carry `user_id` and `role` claims
