#!/bin/bash

# ===================================================================
# SCRIPT SIMPLE DE INICIALIZACIÓN DE DATOS
# ===================================================================
# Script simplificado para ejecutar desde el contenedor Docker

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 INICIALIZACIÓN AUTOMÁTICA DE DATOS${NC}"
echo "====================================="

# Función para esperar que un servicio esté listo
wait_for_service() {
    local host=$1
    local port=$2
    local service_name=$3
    local max_attempts=30
    local attempt=1

    echo -e "${YELLOW}⏳ Esperando $service_name ($host:$port)...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if nc -z $host $port 2>/dev/null; then
            echo -e "${GREEN}✅ $service_name está listo${NC}"
            return 0
        fi
        echo -e "${YELLOW}   Intento $attempt/$max_attempts...${NC}"
        sleep 2
        ((attempt++))
    done
    
    echo -e "${RED}❌ $service_name no respondió después de $((max_attempts * 2)) segundos${NC}"
    return 1
}

# Esperar servicios críticos
echo -e "\n${YELLOW}🔌 Verificando servicios...${NC}"
wait_for_service "mongu" 27017 "MongoDB"
wait_for_service "gestion-db" 5432 "Gestión DB"
wait_for_service "notification-db" 5432 "Notification DB"
wait_for_service "documentacion-db" 5432 "Documentación DB"

echo -e "\n${YELLOW}📊 Cargando datos en MongoDB...${NC}"

# Cargar datos en MongoDB
mongosh --host mongu:27017 iot_db --eval "
print('🧹 Limpiando colecciones existentes...');
db.activos.deleteMany({});
db.sensores.deleteMany({});
db.sensor_status.deleteMany({});

print('📦 Insertando activos sincronizados...');
db.activos.insertMany([
  {
    activo_id: 1,
    nombre: 'BombaDeAgua #1 ML',
    ubicacion: 'Planta B',
    estado: 'OK',
    id_edificio: 'EB',
    tipo: 'bomba de agua',
    // Campos adicionales para compatibilidad
    edificio_id: 1,
    _legacy_id: 1
  },
  {
  activo_id: 2,
    nombre: 'Bomba Centrífuga A',
    estado: 'operativo',
    edificio_id: 1,
    tipo: 'bomba'
  },
  {
  activo_id: 3,
    nombre: 'Bomba Hidráulica 1',
    estado: 'operativo',
    edificio_id: 2,
    tipo: 'bomba'
  },
  {
  activo_id: 4,
    nombre: 'Ascensor Norte',
    estado: 'operativo',
    edificio_id: 2,
    tipo: 'ascensor'
  },
  {
  activo_id: 5,
    nombre: 'Transformador Secundario',
    estado: 'mantenimiento',
    edificio_id: 3,
    tipo: 'transformador'
  },
  {
  activo_id: 6,
    nombre: 'Transformador Principal',
    estado: 'operativo',
    edificio_id: 3,
    tipo: 'transformador'
  },
  {
  activo_id: 7,
    nombre: 'Ascensor Central',
    estado: 'operativo',
    edificio_id: 1,
    tipo: 'ascensor'
  },
  {
  activo_id: 8,
    nombre: 'Bomba de Emergencia',
    estado: 'mantenimiento',
    edificio_id: 4,
    tipo: 'bomba'
  }
]);

print('🔧 Insertando sensores...');
db.sensores.insertMany([
  // Sensores para BombaDeAgua1 según especificación
  {sensor_id: 'temp1', activo_id: 1, tipo: 'Temperatura', unidad: '°C', _legacy_activo_id: 1},
  {sensor_id: 'pres1', activo_id: 1, tipo: 'Presión', unidad: 'Pa', _legacy_activo_id: 1},
  {sensor_id: 'caud1', activo_id: 1, tipo: 'Caudal', unidad: 'm^3/s', _legacy_activo_id: 1},
  // Sensores existentes para otros activos
  {sensor_id: 'SENSOR_AC1002_01', activo_id: 2, tipo: 'vibracion', unidad: 'mm/s'},
  {sensor_id: 'SENSOR_AC1002_02', activo_id: 2, tipo: 'caudal', unidad: 'L/min'},
  {sensor_id: 'SENSOR_AC1003_01', activo_id: 3, tipo: 'caudal', unidad: 'L/min'},
  {sensor_id: 'SENSOR_AC1003_02', activo_id: 3, tipo: 'presion', unidad: 'bar'},
  {sensor_id: 'SENSOR_AC1004_01', activo_id: 4, tipo: 'vibracion', unidad: 'mm/s'},
  {sensor_id: 'SENSOR_AC1004_02', activo_id: 4, tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1005_01', activo_id: 5, tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1005_02', activo_id: 5, tipo: 'temperatura', unidad: '°C'},
  {sensor_id: 'SENSOR_AC1006_01', activo_id: 6, tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1006_02', activo_id: 6, tipo: 'temperatura', unidad: '°C'},
  {sensor_id: 'SENSOR_AC1007_01', activo_id: 7, tipo: 'vibracion', unidad: 'mm/s'},
  {sensor_id: 'SENSOR_AC1007_02', activo_id: 7, tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1008_01', activo_id: 8, tipo: 'presion', unidad: 'bar'},
  {sensor_id: 'SENSOR_AC1008_02', activo_id: 8, tipo: 'caudal', unidad: 'L/min'}
]);

print('📊 Insertando estados de sensores...');
var now = new Date();
db.sensor_status.insertMany([
  // Estados para sensores de BombaDeAgua1
  {sensor_id: 'temp1', activo_id: 1, estado: 'conectado', ultima_lectura: now, valor_actual: 65.5, _legacy_activo_id: 1},
  {sensor_id: 'pres1', activo_id: 1, estado: 'conectado', ultima_lectura: now, valor_actual: 6200.0, _legacy_activo_id: 1}, // Convertido a Pa
  {sensor_id: 'caud1', activo_id: 1, estado: 'conectado', ultima_lectura: now, valor_actual: 0.425, _legacy_activo_id: 1}, // Convertido a m³/s
  // Estados existentes para otros activos
  {sensor_id: 'SENSOR_AC1002_01', activo_id: 2, estado: 'conectado', ultima_lectura: now, valor_actual: 12.8},
  {sensor_id: 'SENSOR_AC1002_02', activo_id: 2, estado: 'conectado', ultima_lectura: now, valor_actual: 425.0},
  {sensor_id: 'SENSOR_AC1003_01', activo_id: 3, estado: 'conectado', ultima_lectura: now, valor_actual: 280.5},
  {sensor_id: 'SENSOR_AC1003_02', activo_id: 3, estado: 'conectado', ultima_lectura: now, valor_actual: 7.8},
  {sensor_id: 'SENSOR_AC1004_01', activo_id: 4, estado: 'conectado', ultima_lectura: now, valor_actual: 9.1},
  {sensor_id: 'SENSOR_AC1004_02', activo_id: 4, estado: 'conectado', ultima_lectura: now, valor_actual: 400.0},
  {sensor_id: 'SENSOR_AC1005_01', activo_id: 5, estado: 'desconectado', ultima_lectura: now, valor_actual: 0},
  {sensor_id: 'SENSOR_AC1005_02', activo_id: 5, estado: 'desconectado', ultima_lectura: now, valor_actual: 0},
  {sensor_id: 'SENSOR_AC1006_01', activo_id: 6, estado: 'conectado', ultima_lectura: now, valor_actual: 220.8},
  {sensor_id: 'SENSOR_AC1006_02', activo_id: 6, estado: 'conectado', ultima_lectura: now, valor_actual: 42.3},
  {sensor_id: 'SENSOR_AC1007_01', activo_id: 7, estado: 'conectado', ultima_lectura: now, valor_actual: 8.2},
  {sensor_id: 'SENSOR_AC1007_02', activo_id: 7, estado: 'conectado', ultima_lectura: now, valor_actual: 395.0},
  {sensor_id: 'SENSOR_AC1008_01', activo_id: 8, estado: 'desconectado', ultima_lectura: now, valor_actual: 0},
  {sensor_id: 'SENSOR_AC1008_02', activo_id: 8, estado: 'desconectado', ultima_lectura: now, valor_actual: 0}
]);

print('✅ Datos cargados exitosamente');
print('📊 Resumen:');
print('   • Activos: ' + db.activos.countDocuments());
print('   • Sensores: ' + db.sensores.countDocuments());
print('   • Estados: ' + db.sensor_status.countDocuments());
print('🔗 Mapeo PostgreSQL ↔ MongoDB:');
print('   ID 1-8 ↔ 1001 a 1008');
" --quiet

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ MongoDB inicializado correctamente${NC}"
else
    echo -e "${RED}❌ Error inicializando MongoDB${NC}"
    exit 1
fi

echo -e "\n${GREEN}🎯 Inicialización automática completada${NC}"
echo -e "${BLUE}📋 Las bases de datos están sincronizadas:${NC}"
echo -e "   • Gestión DB: IDs 1-8"
echo -e "   • Notification DB: IDs 1-8"
echo -e "   • Documentación DB: IDs 1-8"
echo -e "   • MongoDB IoT: 1 a 8"
