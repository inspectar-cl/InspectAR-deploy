# Etapa de construcción
FROM golang:1.22-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git

WORKDIR /app

# Copiar todo el código fuente
COPY . .

# Descargar dependencias y compilar
RUN go mod download
RUN go mod tidy
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o ia-service ./cmd/main.go

# Etapa de producción
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/ia-service .

# Copiar archivos de configuración
COPY --from=builder /app/config ./config

# Exponer puerto
EXPOSE 8095

# Comando de inicio
CMD ["./ia-service"]
