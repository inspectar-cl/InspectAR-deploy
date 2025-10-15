# 🖊️ Guía Rápida: Firmas Digitales para Testing

## 📝 Resumen

Este sistema permite cargar 4 firmas JPG de prueba que se insertarán automáticamente en la base de datos para cada usuario.

## 🎯 Pasos Rápidos

### Opción 1: Generar Firmas Automáticamente (Recomendado)

```bash
# 1. Ir al directorio de firmas
cd /home/joytan/repo/InspectAR/backend/database/gestion/firmas_iniciales

# 2. Ejecutar el script generador
./generar_firmas.sh

# 3. Verificar que se crearon las firmas
ls -lh firma_usuario_*.jpg
```

### Opción 2: Usar tus Propias Firmas JPG

```bash
# 1. Copiar tus firmas al directorio
cp /ruta/a/tu/firma1.jpg /home/joytan/repo/InspectAR/backend/database/gestion/firmas_iniciales/firma_usuario_1.jpg
cp /ruta/a/tu/firma2.jpg /home/joytan/repo/InspectAR/backend/database/gestion/firmas_iniciales/firma_usuario_2.jpg
cp /ruta/a/tu/firma3.jpg /home/joytan/repo/InspectAR/backend/database/gestion/firmas_iniciales/firma_usuario_3.jpg
cp /ruta/a/tu/firma4.jpg /home/joytan/repo/InspectAR/backend/database/gestion/firmas_iniciales/firma_usuario_4.jpg
```

## 🚀 Reconstruir y Probar

```bash
# 1. Detener servicios
cd /home/joytan/repo/InspectAR/backend
docker-compose down

# 2. Reconstruir el servicio de gestión
docker-compose build gestion-service --no-cache

# 3. Iniciar servicios
docker-compose up -d gestion-db gestion-service

# 4. Esperar a que esté listo
sleep 10

# 5. Verificar que las firmas se insertaron
docker-compose exec gestion-db psql -U gestion_user -d gestion_db -c "SELECT id, usuario_id, nombre_archivo, formato FROM firmas_digitales;"
```

## ✅ Verificar Firmas

### Ver firmas en la base de datos
```bash
curl http://localhost:8092/firmas/usuario/1 | jq '.'
curl http://localhost:8092/firmas/usuario/2 | jq '.'
curl http://localhost:8092/firmas/usuario/3 | jq '.'
curl http://localhost:8092/firmas/usuario/4 | jq '.'
```

### Descargar una firma
```bash
curl http://localhost:8092/firmas/1/imagen --output firma_test.jpg
```

### Generar un reporte con firma
```bash
curl -X POST http://localhost:8092/reportes/activo/1 \
  -H "Content-Type: application/json" \
  -d '{
    "campos": ["ubicacion", "historial_mantenimientos"],
    "usar_firma_predeterminada": true,
    "usuario_id": 1
  }' \
  --output reporte_con_firma.pdf
```

## 📋 Estructura de Usuarios y Firmas

| Usuario ID | Username  | Email                  | Firma                     |
|-----------|-----------|------------------------|---------------------------|
| 1         | analista  | analista@example.com   | firma_usuario_1.jpg       |
| 2         | tecnico   | tecnico@example.com    | firma_usuario_2.jpg       |
| 3         | residente | residente@example.com  | firma_usuario_3.jpg       |
| 4         | admin     | admin@example.com      | firma_usuario_4.jpg       |

## 🔍 Solución de Problemas

### Las firmas no aparecen
```bash
# Verificar logs del contenedor
docker-compose logs gestion-service | grep -i firma

# Verificar que los archivos existen en el contenedor
docker-compose exec gestion-service ls -lh /app/storage/firmas/

# Verificar la base de datos
docker-compose exec gestion-db psql -U gestion_user -d gestion_db -c "SELECT COUNT(*) FROM firmas_digitales;"
```

### Regenerar firmas desde cero
```bash
# Detener y limpiar
docker-compose down -v

# Regenerar firmas
cd /home/joytan/repo/InspectAR/backend/database/gestion/firmas_iniciales
./generar_firmas.sh

# Reconstruir todo
cd /home/joytan/repo/InspectAR/backend
docker-compose build gestion-service --no-cache
docker-compose up -d
```

## 📸 Especificaciones de Firmas

- **Formato:** JPG (JPEG)
- **Dimensiones recomendadas:** 400x150 px
- **Tamaño máximo:** 2 MB
- **Fondo:** Preferiblemente blanco o transparente
- **Contenido:** Nombre o rúbrica del usuario

## 🎨 Personalizar Firmas Generadas

Edita el script `generar_firmas.sh` para cambiar:
- Colores
- Tamaños de fuente
- Texto de las firmas
- Dimensiones de imagen

## 📖 Más Información

Ver archivo completo de documentación: `FIRMAS_DIGITALES_API.md`
