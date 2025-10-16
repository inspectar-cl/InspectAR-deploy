# Etapa de construcción
FROM golang:1.24 AS builder

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
FROM golang:1.24 AS final

WORKDIR /app

# Instalar dependencias del sistema incluyendo Chrome
RUN apt-get update && \
    apt-get install -y \
    wget \
    curl \
    gnupg \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libdrm2 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libgbm1 \
    libxss1 \
    libnss3 \
    && rm -rf /var/lib/apt/lists/*

# Instalar Google Chrome para chromedp
RUN wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | gpg --dearmor -o /usr/share/keyrings/google-linux-signing-key.gpg && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/google-linux-signing-key.gpg] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list && \
    apt-get update && \
    apt-get install -y google-chrome-stable && \
    rm -rf /var/lib/apt/lists/*

# Copiar el binario compilado
COPY --from=builder /app/main .

# Copiar la carpeta de configuración
COPY config/config.docker.yaml /app/config/config.yaml

# Copiar templates
COPY templates/ /app/templates/

# Copiar directorio de storage para firmas digitales
COPY storage/ /app/storage/
RUN mkdir -p /app/storage/firmas && chmod -R 755 /app/storage

# Copiar script de inicialización de firmas
COPY init-firmas.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/init-firmas.sh

# Crear un script de inicio que ejecute ambos
RUN echo '#!/bin/sh' > /usr/local/bin/startup.sh && \
    echo '/usr/local/bin/init-firmas.sh' >> /usr/local/bin/startup.sh && \
    echo 'exec /app/main' >> /usr/local/bin/startup.sh && \
    chmod +x /usr/local/bin/startup.sh

# Exponer el puerto
EXPOSE 8092

# Usar entrypoint para inicializar firmas antes de ejecutar la aplicación
ENTRYPOINT ["/usr/local/bin/startup.sh"]

# Ejecutar la aplicación
CMD ["./main"]