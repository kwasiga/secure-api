# ── Build stage ──────────────────────────────────────────────────────────────
# Uses the full Go toolchain to compile a statically linked binary.
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Download dependencies first so Docker can cache this layer independently.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 produces a fully static binary that can run in a scratch/alpine image.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# ── Runtime stage ─────────────────────────────────────────────────────────────
# Minimal alpine image — only the binary and essential runtime libs are included.
FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
