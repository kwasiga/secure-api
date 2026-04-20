# secure-api

A secure REST API written in Go with JWT-based authentication and PostgreSQL persistence.

## Stack

- **Go** — application language
- **PostgreSQL** — database
- **sqlc** — type-safe SQL query generation
- **golang-jwt/jwt** — JWT access and refresh token signing
- **bcrypt** — password hashing
- **godotenv** — environment variable loading

## Project Structure

```
.
├── cmd/api/          # Application entrypoint
├── config/           # Environment config loading
├── db/sqlc/          # sqlc-generated database code
└── internal/
    ├── auth/         # JWT and password utilities
    ├── model/        # App-layer data types
    └── repository/   # Database access layer
```

## Setup

1. Copy `.env.example` to `.env` and fill in the values:

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

## Auth

- Passwords are hashed with **bcrypt** (cost 14)
- Login returns a short-lived **access token** (15 min) and a long-lived **refresh token** (7 days), both signed with HS256
- Tokens embed `user_id` and `role` claims
