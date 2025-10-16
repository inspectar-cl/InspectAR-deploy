# 📧 Microservicio de Notificaciones - InspectAR

## Descripción
Microservicio desarrollado en Go que gestiona notificaciones automáticas para el sistema InspectAR. Envía alertas por email a los usuarios de un edificio cuando se detectan problemas en los activos.

## Características
- 🗄️ **Base de datos PostgreSQL** con relaciones entre edificios, usuarios, activos y notificaciones
- 📧 **Envío automático de emails** con plantillas HTML responsivas
- 🔗 **API REST completa** con 7 endpoints
- 🐳 **Containerizado con Docker** y orquestado con Docker Compose
- ⚙️ **Configuración por variables de entorno** con fallback a archivos YAML
- 🎯 **Notificaciones inteligentes** que obtienen automáticamente los usuarios del edificio

## Arquitectura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Client    │───▶│  Notification   │───▶│   PostgreSQL    │
│                 │    │   Service       │    │   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Email SMTP    │
                       │   (Mailgun)     │
                       └─────────────────┘
```

## Base de Datos

### Esquema de tablas:
- **edificios**: Información de edificios
- **usuarios**: Usuarios del sistema con email
- **activos**: Activos/sensores de cada edificio
- **notificaciones**: Registro de todas las notificaciones

### Relaciones:
- Un edificio tiene muchos usuarios y activos
- Una notificación pertenece a un activo
- Cada notificación se envía a todos los usuarios del edificio

## API Endpoints

### 📧 Notificaciones

#### 1. Crear nueva notificación
**Ruta:** `POST /notification`

**JSON que recibe:**
```json
{
  "activo_id": 1
}
```

**JSON de respuesta (éxito):**
```json
{
  "message": "Notificación creada exitosamente",
  "notificacion": {
    "id": 15,
    "activo_id": 1,
    "activo": {
      "id": 1,
      "nombre": "Sensor Temperatura Planta Baja",
      "edificio_id": 1,
      "edificio": {
        "id": 1,
        "direccion": "Av. Libertador 1234, CABA",
        "numero_activos": 3,
        "created_at": "2025-08-08T03:50:10Z",
        "updated_at": "2025-08-08T03:50:10Z"
      },
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    },
    "usuario_id": null,
    "mensaje": "Alerta detectada en Sensor Temperatura Planta Baja",
    "enviado": true,
    "created_at": "2025-08-08T03:50:10Z",
    "sent_at": "2025-08-08T03:50:10Z"
  }
}
```

**Ejemplo curl:**
```bash
curl -X POST http://localhost:8091/notification \
  -H "Content-Type: application/json" \
  -d '{"activo_id": 1}'
```

---

#### 2. Obtener notificaciones por activo
**Ruta:** `GET /notification/:activo_id`

**JSON que recibe:** Ninguno (parámetro en URL)

**JSON de respuesta:**
```json
{
  "notificaciones": [
    {
      "id": 15,
      "activo_id": 1,
      "activo": {
        "id": 1,
        "nombre": "Sensor Temperatura Planta Baja",
        "edificio_id": 1,
        "created_at": "2025-08-08T03:50:10Z",
        "updated_at": "2025-08-08T03:50:10Z"
      },
      "usuario_id": null,
      "mensaje": "Alerta detectada en Sensor Temperatura Planta Baja",
      "enviado": true,
      "created_at": "2025-08-08T03:50:10Z",
      "sent_at": "2025-08-08T03:50:10Z"
    }
  ]
}
```

**Ejemplo curl:**
```bash
curl http://localhost:8091/notification/1
```

---

#### 3. Reenviar notificación por email
**Ruta:** `PUT /notification/:notification_id/send`

**JSON que recibe:** Ninguno (parámetro en URL)

**JSON de respuesta:**
```json
{
  "message": "Notificación enviada exitosamente"
}
```

**Ejemplo curl:**
```bash
curl -X PUT http://localhost:8091/notification/15/send
```

---

### 🏢 Edificios

#### 4. Listar todos los edificios
**Ruta:** `GET /edificios`

**JSON que recibe:** Ninguno

**JSON de respuesta:**
```json
{
  "edificios": [
    {
      "id": 1,
      "direccion": "Av. Libertador 1234, CABA",
      "numero_activos": 3,
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    },
    {
      "id": 2,
      "direccion": "Corrientes 5678, CABA",
      "numero_activos": 2,
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    }
  ]
}
```

**Ejemplo curl:**
```bash
curl http://localhost:8091/edificios
```

---

#### 5. Obtener edificio específico
**Ruta:** `GET /edificio/:id`

**JSON que recibe:** Ninguno (parámetro en URL)

**JSON de respuesta:**
```json
{
  "edificio": {
    "id": 1,
    "direccion": "Av. Libertador 1234, CABA",
    "numero_activos": 3,
    "created_at": "2025-08-08T03:50:10Z",
    "updated_at": "2025-08-08T03:50:10Z"
  }
}
```

**Ejemplo curl:**
```bash
curl http://localhost:8091/edificio/1
```

---

#### 6. Obtener usuarios de un edificio
**Ruta:** `GET /edificio/:id/usuarios`

**JSON que recibe:** Ninguno (parámetro en URL)

**JSON de respuesta:**
```json
{
  "usuarios": [
    {
      "id": 1,
      "scope": "admin",
      "usuario": "admin_user",
      "correo": "joytan33334@gmail.com",
      "numero": "+54911234567",
      "edificio_id": 1,
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    },
    {
      "id": 2,
      "scope": "operator",
      "usuario": "operator1",
      "correo": "joytan33334@gmail.com",
      "numero": "+54911234568",
      "edificio_id": 1,
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    }
  ]
}
```

**Ejemplo curl:**
```bash
curl http://localhost:8091/edificio/1/usuarios
```

---

#### 7. Obtener activos de un edificio
**Ruta:** `GET /edificio/:id/activos`

**JSON que recibe:** Ninguno (parámetro en URL)

**JSON de respuesta:**
```json
{
  "activos": [
    {
      "id": 1,
      "nombre": "Sensor Temperatura Planta Baja",
      "edificio_id": 1,
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    },
    {
      "id": 2,
      "nombre": "Sensor Humedad Primer Piso",
      "edificio_id": 1,
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    },
    {
      "id": 3,
      "nombre": "Sistema HVAC Principal",
      "edificio_id": 1,
      "created_at": "2025-08-08T03:50:10Z",
      "updated_at": "2025-08-08T03:50:10Z"
    }
  ]
}
```

**Ejemplo curl:**
```bash
curl http://localhost:8091/edificio/1/activos
```

---

### ❌ Respuestas de Error

Todos los endpoints pueden devolver errores en el siguiente formato:

**Error 400 (Bad Request):**
```json
{
  "error": "Datos inválidos"
}
```

**Error 500 (Internal Server Error):**
```json
{
  "error": "Error al crear notificación: record not found"
}
```

## Configuración

### Variables de entorno requeridas:
```bash
# Configuración de correo
EMAIL_SMTP_HOST=smtp.mailgun.org
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-username
EMAIL_SMTP_PASSWORD=your-password
EMAIL_FROM_EMAIL=notifications@yourdomain.com

# Configuración de base de datos (automática con Docker Compose)
DB_HOST=notification-db
DB_PORT=5432
DB_USER=notification_user
DB_PASSWORD=notification_pass
DB_NAME=notification_db
```

### Alternativa Gmail (para testing):
```bash
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=tu-email@gmail.com
EMAIL_SMTP_PASSWORD=tu-app-password  # App Password de Gmail
EMAIL_FROM_EMAIL=tu-email@gmail.com
```

## Instalación y Ejecución

### Con Docker Compose (recomendado):
```bash
# Desde la raíz del proyecto backend
docker-compose up --build notification-service notification-db
```

### Manual:
```bash
# Instalar dependencias
go mod tidy

# Ejecutar
go run cmd/main.go
```

## Uso

### Ejemplo: Crear notificación crítica
```bash
curl -X POST http://localhost:8091/notification 
  -H "Content-Type: application/json" 
  -d '{
    "activo_id": 1,
    "tipo": "critica",
    "mensaje": "Temperatura crítica detectada",
    "detalles": "El sensor ha registrado 85°C, se requiere intervención inmediata."
  }'
```

### Ejemplo: Obtener usuarios de un edificio
```bash
curl http://localhost:8091/edificio/1/usuarios
```

### Ejemplo: Reenviar notificación
```bash
curl -X PUT http://localhost:8091/notification/1/send
```

## Plantilla de Email

El sistema genera emails HTML con:
- **Header profesional** con logo y branding
- **Información del activo** y edificio afectado
- **Mensaje y detalles** de la alerta
- **Timestamp** de la notificación
- **Footer** con información de contacto
- **Diseño responsivo** para mobile y desktop

## Logs y Monitoreo

El servicio genera logs detallados:
```
✅ Usando configuración de correo desde variables de entorno
✅ Conectado a PostgreSQL
📧 Configuración de correo cargada correctamente
📨 Enviando notificación por email a: usuario@example.com
✅ Email enviado exitosamente
```

## Estructura del Proyecto

```
notification/
├── cmd/
│   └── main.go              # Punto de entrada
├── internal/
│   ├── handler/
│   │   └── notificationHandler.go  # Controladores API
│   ├── models/
│   │   └── models.go        # Modelos GORM
│   └── services/
│       └── emailService.go  # Servicio de email
├── config/
│   └── config.yaml          # Configuración por defecto
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Dependencias Principales

- **Gin**: Framework web HTTP
- **GORM**: ORM para PostgreSQL  
- **Viper**: Gestión de configuración
- **SMTP**: Cliente de email nativo de Go
- **PostgreSQL**: Base de datos principal

## Estado del Proyecto

✅ **Completado**:
- Base de datos PostgreSQL operativa
- API REST con 7 endpoints funcionales
- Sistema de emails con plantillas HTML
- Configuración por variables de entorno
- Contenedorización Docker
- Documentación completa

🔧 **Para producción**:
- Configurar credenciales SMTP válidas (Mailgun/SendGrid/Gmail)
- Ajustar configuración de logs
- Implementar autenticación/autorización si es requerida
- Configurar SSL/TLS para conexiones seguras

## Créditos

Desarrollado como parte del sistema InspectAR para monitoreo y alertas de activos en edificios.

Este microservicio maneja las notificaciones del sistema InspectAR con envío automático de correos electrónicos a usuarios del mismo edificio cuando se detecta una alerta en un activo.

## 🚀 Características

- ✅ **Notificaciones automáticas** por correo electrónico
- ✅ **Gestión de edificios, usuarios y activos**
- ✅ **Base de datos PostgreSQL** con relaciones
- ✅ **Envío masivo** a usuarios del mismo edificio
- ✅ **Templates HTML** profesionales para correos
- ✅ **Soporte para múltiples proveedores SMTP** (Mailgun, Gmail, Outlook)

## 🚀 Configuración de Correo Electrónico

### 1. Variables de Entorno

Copia el archivo a `.env` y configura tus credenciales de correo:

```bash
cp .env
```

### 2. Configuración de Gmail (Recomendado)

1. **Habilita la verificación en 2 pasos** en tu cuenta de Gmail
2. **Genera una contraseña de aplicación:**
   - Ve a [Configuración de Google](https://myaccount.google.com/)
   - Seguridad → Verificación en 2 pasos → Contraseñas de aplicaciones
   - Selecciona "Correo" y "Otro" → Escribe "InspectAR"
   - Copia la contraseña generada (16 caracteres)

3. **Configura el archivo `.env`:**
```bash
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=tu_correo@gmail.com
EMAIL_SMTP_PASSWORD=abcd_efgh_ijkl_mnop  # Contraseña de aplicación
EMAIL_FROM_EMAIL=tu_correo@gmail.com
```

### 3. Configuración de Outlook

```bash
EMAIL_SMTP_HOST=smtp-mail.outlook.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=tu_correo@outlook.com
EMAIL_SMTP_PASSWORD=tu_contraseña_normal
EMAIL_FROM_EMAIL=tu_correo@outlook.com
```

## 📧 Funcionalidad de Envío Automático

Cuando se crea una notificación (`POST /notification`), el sistema automáticamente:

1. **Identifica el edificio** del activo que generó la alerta
2. **Busca todos los usuarios** asignados a ese edificio
3. **Envía un correo HTML** profesional a todos esos usuarios
4. **Registra el resultado** en los logs

### Ejemplo de Flujo

```json
POST /notification
{
  "activo_id": 1
}
```

**Resultado:**
- Activo ID=1 pertenece al Edificio ID=1
- Usuarios del Edificio ID=1: admin_user, operator1
- ✅ Correo enviado a: admin_user@gmail.com, operator1@gmail.com

## Estructura de la Base de Datos

### Tablas:
- **edificios**: Direccion, numero de activos
- **usuarios**: Scope, usuario, correo, numero (pertenece a un edificio)
- **activos**: Nombre (pertenece a un edificio)
- **notificaciones**: Mensaje, estado de envío (relacionada a activo y usuario)

### Relaciones:
- Un edificio tiene muchos usuarios
- Un usuario pertenece a un edificio
- Un edificio tiene muchos activos  
- Un activo pertenece a un edificio
- Las notificaciones están relacionadas con activos

## Estructura del Proyecto

```
notification/
├── api/
│   └── router/
│       └── router.go          # Configuración de rutas HTTP
├── cmd/
│   └── main.go               # Punto de entrada de la aplicación
├── config/
│   └── config.yaml           # Configuración del servicio
├── internal/
│   ├── database/
│   │   └── postgres.go       # Conexión a PostgreSQL con GORM
│   ├── handler/
│   │   └── notificationHandler.go  # Handlers HTTP
│   ├── models/
│   │   └── models.go         # Modelos GORM
│   ├── repository/
│   │   └── notificationRepo.go     # Repositorio con GORM
│   └── services/
│       └── notificationService.go  # Lógica de negocio
├── Dockerfile                # Configuración Docker
└── go.mod                   # Dependencias Go
```

## Endpoints

### Notificaciones
- **POST /notification** - Crear notificación con `{"activo_id": number}`
- **GET /notification/:activo_id** - Obtener notificaciones por activo
- **PUT /notification/:notification_id/send** - Marcar notificación como enviada

### Edificios
- **GET /edificios** - Obtener todos los edificios
- **GET /edificio/:id** - Obtener edificio específico
- **GET /edificio/:id/usuarios** - Obtener usuarios de un edificio
- **GET /edificio/:id/activos** - Obtener activos de un edificio

## Configuración

El servicio se configura através del archivo `config/config.yaml`:

```yaml
postgres:
  host: "notification-db"
  port: 5432
  user: "notification_user"
  password: "notification_pass"
  dbname: "notification_db"
  sslmode: "disable"

server:
  port: "8091"
```

## Tecnologías

- **Backend**: Go con Gin framework
- **ORM**: GORM v2
- **Base de datos**: PostgreSQL 15
- **Contenedores**: Docker
- **Puerto**: 8091

## Construcción y Ejecución

```bash
# Construir y ejecutar con Docker Compose
docker-compose up --build notification-service notification-db

# O ejecutar todo el stack
docker-compose up --build
```

## Datos de Ejemplo

El servicio incluye datos de ejemplo:
- 3 edificios
- 5 usuarios distribuidos en los edificios
- 6 activos distribuidos en los edificios

## Endpoints

### POST /notification
Crea una nueva notificación para un activo específico.

**Request Body:**
```json
{
  "activo_id": "string"
}
```

### GET /notification/:activo_id
Obtiene todas las notificaciones para un activo específico.

### PUT /notification/:notification_id/send
Marca una notificación como enviada.

## Configuración

El servicio se configura através del archivo `config/config.yaml`:

```yaml
mongu:
  uri: "mongodb://mongu:27017"
  database: "notification_db"

server:
  port: "8091"
```

## Construcción y Ejecución

```bash
# Construir la imagen Docker
docker build -t notification-service .

# Ejecutar el servicio
docker run -p 8091:8091 notification-service
```
