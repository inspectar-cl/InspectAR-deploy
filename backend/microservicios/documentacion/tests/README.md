# 🧪 DIRECTORIO DE PRUEBAS - MICROSERVICIO DOCUMENTACIÓN

Este directorio contiene todos los archivos relacionados con testing y validación del microservicio de documentación.

## 📁 CONTENIDO DEL DIRECTORIO

### 🧪 **Scripts de Prueba Funcionales**
- **`test_upload_and_query.sh`** - Test rápido de carga y consulta de documentos (⚡ Recomendado)
- **`test_complete_workflow.sh`** - Test exhaustivo del flujo completo de trabajo
- **`test_docker.sh`** - Script para probar el microservicio en entorno Docker
- **`test_routes.sh`** - Script bash automatizado para probar todas las rutas del API

### 📋 **Documentación de Pruebas**
- **`README.md`** - Esta guía de uso (archivo actual)
- **`API_TESTS.md`** - Documentación técnica de tests del API
- **`PRUEBAS_COMPLETADAS.md`** - Reporte ejecutivo de pruebas completadas

## 🚀 CÓMO EJECUTAR LAS PRUEBAS

### 🐳 **REQUISITOS PREVIOS: DOCKER**

Antes de ejecutar cualquier test, asegúrate de que los servicios Docker estén funcionando:

```bash
# Desde el directorio backend/
cd /home/joytan/repo/InspectAR/backend

# Iniciar servicios Docker
docker-compose up documentacion-db documentacion-service -d

# Verificar que estén funcionando
docker-compose ps
```

### 🧪 **TESTS DISPONIBLES**

#### 1. **Test Rápido de Funcionalidad** ⚡
```bash
cd microservicios/documentacion/
./tests/test_upload_and_query.sh
```
**Qué hace:**
- Sube 3 documentos de prueba
- Realiza consultas básicas
- Prueba funcionalidad IA
- Tarda ~30 segundos

#### 2. **Test Completo del Flujo de Trabajo** 🔄
```bash
cd microservicios/documentacion/
./tests/test_complete_workflow.sh
```
**Qué hace:**
- Crea documentos con contenido detallado
- Prueba todos los endpoints del API
- Valida funcionalidad IA avanzada
- Ejecuta casos edge y validaciones
- Tarda ~2-3 minutos

#### 3. **Test de Docker** 🐳
```bash
cd microservicios/documentacion/
./tests/test_docker.sh
```
**Qué hace:**
- Verifica que el microservicio esté disponible
- Prueba endpoints básicos
- Valida conectividad con base de datos

#### 4. **Test de Rutas API** 🛣️
```bash
cd microservicios/documentacion/
./tests/test_routes.sh
```
**Qué hace:**
- Prueba sistemáticamente todos los endpoints
- Valida códigos de respuesta HTTP
- Ejecuta casos de éxito y error

### 📊 **INTERPRETACIÓN DE RESULTADOS**

#### ✅ **Indicadores de Éxito:**
- `✅ PASS` - Test exitoso
- `HTTP 200/201` - Respuestas correctas
- `🎉 PRUEBA COMPLETA EXITOSA` - Todo funcionando

#### ❌ **Indicadores de Error:**
- `❌ FAIL` - Test fallido  
- `HTTP 4xx/5xx` - Errores de servidor
- `Connection refused` - Servicio no disponible

#### ⚠️ **Advertencias:**
- `⚠️ Warning` - Funcionalidad parcial
- `Expected vs Got` - Diferencias menores

### � **SOLUCIÓN DE PROBLEMAS**

#### **Servicio no disponible:**
```bash
# Verificar estado de contenedores
docker-compose ps

# Reiniciar servicios
docker-compose restart documentacion-service

# Ver logs
docker logs documentacion-service
```

#### **Error de permisos en archivos:**
```bash
# Hacer ejecutables los scripts
chmod +x tests/*.sh
```

#### **Error de base de datos:**
```bash
# Verificar conexión a BD
docker exec documentacion-db pg_isready -U documentacion_user

# Ver logs de BD
docker logs documentacion-db
```

### 🎯 **TESTS RECOMENDADOS POR SITUACIÓN**

#### **Desarrollo Diario:** 
- `test_upload_and_query.sh` (rápido)

#### **Antes de Deploy:**
- `test_complete_workflow.sh` (exhaustivo)

#### **Verificación Docker:**
- `test_docker.sh` (conectividad)

#### **Debug de API:**
- `test_routes.sh` (endpoints específicos)

## � **COBERTURA ACTUAL**

- ✅ **Carga de Documentos**: 100% validada
- ✅ **Consultas y Búsquedas**: 100% funcionando  
- ✅ **Integración IA (Gemini)**: Completamente testada
- ✅ **Base de Datos PostgreSQL**: Validada
- ✅ **Docker Compose**: Funcionando
- ✅ **Casos Edge**: Cubiertos
- ✅ **Validaciones**: Implementadas

## 🎯 **ESTADO ACTUAL**

**🟢 TODAS LAS PRUEBAS PASANDO** ✅

El microservicio está completamente validado y listo para producción.

### **Último Test Ejecutado:**
- **Fecha:** 31 de Agosto, 2025
- **Documentos en Sistema:** 21+ documentos
- **Fichas Técnicas:** 6 fichas
- **IA Funcionando:** ✅ Gemini respondiendo
- **Base de Datos:** ✅ PostgreSQL operativa
