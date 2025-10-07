# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copia los archivos go.mod y go.sum y descarga dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Ejecuta tidy después de copiar todo el código
RUN go mod tidy

# Compila el binario
RUN go build -o apigateway ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copia el binario compilado desde el builder
COPY --from=builder /app/apigateway .

# Copiar la carpeta de configuración
COPY config/* /app/config/

# Copia el archivo .env si existe
COPY .env .env

EXPOSE 3500

CMD ["./apigateway"]