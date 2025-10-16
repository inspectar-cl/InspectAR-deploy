# 🚨 SISTEMA DE ALERTAS DE SENSOR - DOCUMENTACIÓN COMPLETA

## 📋 RESUMEN DE IMPLEMENTACIÓN

Se ha implementado un sistema completo de alertas de sensor que permite:

1. ✅ **Recibir alertas de sensores con diferentes tipos**
2. ✅ **Envío automático de notificaciones por correo electrónico**  
3. ✅ **Soporte para múltiples tipos de notificación**
4. ✅ **Integración entre ParserService y NotificationService**
5. ✅ **Base de datos actualizada con nuevos tipos y campos**

---

## 🛠️ CAMBIOS REALIZADOS

### 📊 **Base de Datos (notification)**

#### Nuevas Tablas:
```sql
-- Tipos de notificación
CREATE TABLE tipos_notificacion (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(50) UNIQUE NOT NULL,
    descripcion TEXT,
    prioridad INTEGER DEFAULT 1, -- 1=Low, 2=Medium, 3=High, 4=Critical
    activo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Tabla Notificaciones Actualizada:
```sql
CREATE TABLE notificaciones (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER REFERENCES activos(id) ON DELETE CASCADE,
    sensor_id VARCHAR(255), -- NUEVO: Para alertas específicas de sensor
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
    tipo_id INTEGER NOT NULL REFERENCES tipos_notificacion(id), -- NUEVO
    titulo VARCHAR(255) NOT NULL, -- NUEVO
    mensaje TEXT NOT NULL,
    datos_sensor JSONB, -- NUEVO: Para almacenar datos adicionales del sensor
    enviado BOOLEAN DEFAULT FALSE,
    enviado_email BOOLEAN DEFAULT FALSE, -- NUEVO
    enviado_sms BOOLEAN DEFAULT FALSE, -- NUEVO
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP NULL,
    email_sent_at TIMESTAMP NULL, -- NUEVO
    sms_sent_at TIMESTAMP NULL -- NUEVO
);
```

#### Tipos de Notificación Predefinidos:
- `alerta_sensor_desconexion` (Prioridad: 3 - High)
- `alerta_sensor_problema` (Prioridad: 2 - Medium)  
- `alerta_sensor_valor_anormal` (Prioridad: 2 - Medium)
- `alerta_activo_offline` (Prioridad: 4 - Critical)
- `alerta_activo_mantenimiento` (Prioridad: 1 - Low)
- `info_sistema_reinicio` (Prioridad: 1 - Low)
- `warning_sensor_bateria_baja` (Prioridad: 2 - Medium)
- `critical_sensor_error` (Prioridad: 4 - Critical)

---

## 🔌 **NUEVA API - ENDPOINT DE ALERTAS DE SENSOR**

### **POST /sensor/alert**

Crea una alerta de sensor con notificación automática por email.

#### **Request Body:**
```json
{
    "sensor_id": "TEMP_001_DISCONNECT_TEST",
    "activo_id": 1, // Opcional - si no se especifica, se envía a administradores
    "tipo_alerta": "desconexion", // "desconexion", "problema", "valor_anormal", etc.
    "titulo": "Sensor de Temperatura Desconectado",
    "mensaje": "El sensor TEMP_001 se ha desconectado inesperadamente. Última lectura hace 5 minutos.",
    "datos_sensor": { // Opcional - datos adicionales del sensor
        "ultima_lectura": "2025-09-03T04:30:00Z",
        "tipo": "temperatura",
        "unidad": "°C",
        "valor_anterior": 23.5,
        "ubicacion": "Sala Principal"
    },
    "prioridad": 3, // 1=Low, 2=Medium, 3=High, 4=Critical
    "enviar_email": true,
    "enviar_sms": false
}
```

#### **Response:**
```json
{
    "message": "Alerta de sensor creada y enviada exitosamente",
    "notificacion": {
        "id": 15,
        "sensor_id": "TEMP_001_DISCONNECT_TEST",
        "activo_id": 1,
        "tipo_id": 1,
        "titulo": "Sensor de Temperatura Desconectado",
        "mensaje": "El sensor TEMP_001 se ha desconectado inesperadamente...",
        "enviado": true,
        "enviado_email": true,
        "enviado_sms": false,
        "created_at": "2025-09-03T04:40:00Z",
        "email_sent_at": "2025-09-03T04:40:02Z"
    }
}
```

---

## 📡 **INTEGRACIÓN CON PARSERSERVICE**

El **ParserService** ahora envía automáticamente alertas al **NotificationService** cuando detecta sensores desconectados.

### **Flujo Automático:**
1. 🔍 ParserService detecta sensor desconectado (>5 min sin datos)
2. 📧 Genera console log local
3. 🚀 Envía HTTP POST a `notification-service:8091/sensor/alert`
4. ✉️ NotificationService envía emails automáticamente
5. 💾 Guarda registro en base de datos

### **Payload Enviado:**
```json
{
    "sensor_id": "TEMP_DISC_20250903_002522",
    "tipo_alerta": "desconexion",
    "titulo": "Sensor TEMP_DISC_20250903_002522 Desconectado",
    "mensaje": "El sensor TEMP_DISC_20250903_002522 se ha desconectado. Última actividad: 2025-09-03 04:25:30",
    "datos_sensor": {
        "ultima_actividad": "2025-09-03T04:25:30Z",
        "desconectado_en": "2025-09-03T04:30:40Z",
        "total_reportes": 4,
        "tiempo_desconexion": "5m10s"
    },
    "prioridad": 3,
    "enviar_email": true,
    "enviar_sms": false
}
```

---

## 📮 **SISTEMA DE CORREO ELECTRÓNICO**

### **Características:**
- ✅ Envío automático de emails HTML formateados
- ✅ Soporte para múltiples destinatarios
- ✅ Configuración via variables de entorno (.env)
- ✅ Logs detallados de envío
- ✅ Manejo de errores

### **Configuración (.env):**
```env
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=tu_email@gmail.com
EMAIL_SMTP_PASSWORD=tu_app_password
EMAIL_FROM_EMAIL=tu_email@gmail.com
```

### **Ejemplo de Email Enviado:**
```
🚨 ALERTA: Sensor TEMP_001_DISCONNECT_TEST Desconectado

📊 Datos del Sensor:
{
  "ultima_lectura": "2025-09-03T04:30:00Z",
  "tipo": "temperatura", 
  "unidad": "°C",
  "valor_anterior": 23.5,
  "ubicacion": "Sala Principal"
}

Fecha y Hora: 2025-09-03 04:40:00
```

---

## 🔗 **ENDPOINTS ADICIONALES**

### **GET /tipos-notificacion**
Obtiene todos los tipos de notificación disponibles.

```json
{
    "tipos_notificacion": [
        {
            "id": 1,
            "nombre": "alerta_sensor_desconexion",
            "descripcion": "Notificación cuando un sensor se desconecta",
            "prioridad": 3,
            "activo": true
        }
    ]
}
```

### **GET /notification/:activo_id**
Obtiene todas las notificaciones de un activo específico.

### **GET /edificios**
Lista todos los edificios con sus usuarios y activos.

### **GET /edificio/:id/usuarios**
Obtiene usuarios de un edificio específico (destinatarios de emails).

---

## 🧪 **TESTING**

### **Script de Pruebas:**
```bash
# Ejecutar tests de alertas de sensor
./microservicios/notification/tests/test_sensor_alerts.sh
```

### **Casos de Prueba Incluidos:**
1. ✅ Alerta de desconexión de sensor
2. ✅ Alerta de problema en sensor  
3. ✅ Alerta crítica sin activo específico
4. ✅ Alerta de valor anormal
5. ✅ Validación de datos inválidos
6. ✅ Verificación de notificaciones creadas
7. ✅ Consulta de usuarios que reciben emails

---

## 🚀 **DESPLIEGUE**

### **1. Reconstruir Servicios:**
```bash
cd /home/joytan/repo/InspectAR/backend
docker-compose down notification-service
docker-compose build notification-service
docker-compose up -d notification-service
```

### **2. Verificar Logs:**
```bash
# Logs del servicio de notificaciones
docker logs notification-service | grep -i "email\|alerta"

# Logs del parser service 
docker logs iot-service | grep -i "notificación\|alerta"
```

### **3. Probar Funcionalidad:**
```bash
# Test completo de alertas
./microservicios/notification/tests/test_sensor_alerts.sh

# Test de desconexión de sensor (ParserService)
./microservicios/ParserService/tests/test_full_disconnection_simulation.sh
```

---

## 💡 **EJEMPLOS DE USO**

### **1. Alerta Manual de Sensor:**
```bash
curl -X POST "http://localhost:8091/sensor/alert" \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP_001",
    "tipo_alerta": "problema", 
    "titulo": "Sensor con Problemas",
    "mensaje": "El sensor está reportando valores inconsistentes",
    "enviar_email": true
  }'
```

### **2. Consultar Notificaciones:**
```bash
curl "http://localhost:8091/notification/1" | jq .
```

### **3. Ver Tipos de Alerta:**
```bash
curl "http://localhost:8091/tipos-notificacion" | jq .
```

---

## 🎯 **PRÓXIMOS PASOS RECOMENDADOS**

1. **📱 SMS Integration:** Implementar envío real de SMS usando Twilio/AWS SNS
2. **🔔 Push Notifications:** Agregar notificaciones push para aplicaciones móviles  
3. **📊 Dashboard:** Panel de control para gestionar notificaciones
4. **⚙️ Configuración:** Permitir a usuarios configurar qué alertas recibir
5. **📈 Métricas:** Dashboard de estadísticas de alertas y tiempo de respuesta
6. **🔄 Webhooks:** Soporte para webhooks personalizados
7. **📅 Programación:** Notificaciones programadas para mantenimiento

---

## ✅ **ESTADO ACTUAL**

**SISTEMA COMPLETAMENTE FUNCIONAL** 🎉

- ✅ Base de datos actualizada
- ✅ API de alertas de sensor implementada
- ✅ Integración automática ParserService → NotificationService
- ✅ Envío de emails automático
- ✅ Tipos de notificación configurables
- ✅ Tests completos funcionando
- ✅ Documentación completa

El sistema está listo para producción y puede manejar alertas de sensores con notificaciones automáticas por correo electrónico.
