# Multi-stage build para optimizar imagen final
FROM golang:1.22-alpine AS builder

# Información de mantenimiento
LABEL maintainer="InspectAR Team"
LABEL description="Microservicio de Documentación Técnica - InspectAR"

# Instalar dependencias necesarias para compilación
RUN apk add --no-cache git ca-certificates tzdata

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar aplicación con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o documentacion-service ./cmd/main.go

# Imagen final mínima
FROM alpine:latest

# Instalar certificados CA y crear usuario no-root
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Crear directorio para almacenamiento local
RUN mkdir -p /app/storage && \
    chown -R appuser:appuser /app

# Establecer directorio de trabajo
WORKDIR /app

# Copiar binario compilado
COPY --from=builder /app/documentacion-service .

# Crear directorio de configuración y copiar archivos si existen
RUN mkdir -p ./config

# Cambiar a usuario no-root
USER appuser

# Exponer puerto del microservicio
EXPOSE 8092

# Variables de entorno por defecto
ENV GIN_MODE=release
ENV PORT=8092

# Healthcheck para verificar estado del servicio
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8092/health || exit 1

# Comando por defecto
CMD ["./documentacion-service"]
