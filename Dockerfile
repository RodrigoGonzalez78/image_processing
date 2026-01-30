
# Etapa 1: build
FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o app ./cmd/server

# Etapa 2: imagen final
FROM debian:bookworm-slim

# Instalar dependencias necesarias para SQLite
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Crear usuario no-root por seguridad
RUN useradd -r -s /bin/false appuser && \
    chown -R appuser:appuser /app

COPY --from=builder /app/app .

# Cambiar a usuario no-root
USER appuser

EXPOSE 8080

CMD ["./app"]