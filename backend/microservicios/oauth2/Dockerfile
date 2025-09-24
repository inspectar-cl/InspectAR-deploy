# Etapa de construcción
FROM golang:1.23 AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN go build -o /app/main ./cmd/main.go

# Etapa final para una imagen más liviana
FROM golang:1.23 AS final

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/main .

# Copiar la carpeta de configuración
COPY config/* /app/config/

# Exponer el puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["./main"]
