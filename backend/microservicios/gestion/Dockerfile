# Etapa de construcción
FROM golang:1.22 AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN go build -o /app/main ./cmd/main.go

# Etapa final
FROM golang:1.22 AS final

WORKDIR /app

# Instalar curl para healthcheck
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Copiar el binario compilado
COPY --from=builder /app/main .

# Copiar la carpeta de configuración
COPY config/config.docker.yaml /app/config/config.yaml

# Exponer el puerto
EXPOSE 8092

# Ejecutar la aplicación
CMD ["./main"]
