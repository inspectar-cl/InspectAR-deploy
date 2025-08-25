# Dockerfile para microservicio de documentación
FROM golang:1.22-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Imagen final
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Crear directorio de trabajo
WORKDIR /root/

# Copiar el binario compilado
COPY --from=builder /app/main .

# Copiar archivos de configuración
COPY --from=builder /app/config ./config

# Crear directorio para storage local
RUN mkdir -p /root/storage/uploads

# Exponer puerto
EXPOSE 8093

# Variables de entorno por defecto
ENV DOC_SERVER_PORT=8093
ENV DOC_SERVER_HOST=0.0.0.0

# Comando para ejecutar la aplicación
CMD ["./main"]
