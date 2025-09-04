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
    activo_id: 'AC-1001',
    nombre: 'Caldera Principal',
    estado: 'operativo',
    ubicacion: 'Av. Providencia 123, Santiago',
    edificio_id: 1,
    tipo: 'caldera'
  },
  {
    activo_id: 'AC-1002',
    nombre: 'Bomba Centrífuga A',
    estado: 'operativo',
    ubicacion: 'Av. Providencia 123, Santiago',
    edificio_id: 1,
    tipo: 'bomba'
  },
  {
    activo_id: 'AC-1003',
    nombre: 'Bomba Hidráulica 1',
    estado: 'operativo',
    ubicacion: 'Av. Las Condes 456, Las Condes',
    edificio_id: 2,
    tipo: 'bomba'
  },
  {
    activo_id: 'AC-1004',
    nombre: 'Ascensor Norte',
    estado: 'operativo',
    ubicacion: 'Av. Las Condes 456, Las Condes',
    edificio_id: 2,
    tipo: 'ascensor'
  },
  {
    activo_id: 'AC-1005',
    nombre: 'Transformador Secundario',
    estado: 'mantenimiento',
    ubicacion: 'Av. Vicuña Mackenna 789, La Florida',
    edificio_id: 3,
    tipo: 'transformador'
  },
  {
    activo_id: 'AC-1006',
    nombre: 'Transformador Principal',
    estado: 'operativo',
    ubicacion: 'Av. Vicuña Mackenna 789, La Florida',
    edificio_id: 3,
    tipo: 'transformador'
  },
  {
    activo_id: 'AC-1007',
    nombre: 'Ascensor Central',
    estado: 'operativo',
    ubicacion: 'Av. Providencia 123, Santiago',
    edificio_id: 1,
    tipo: 'ascensor'
  },
  {
    activo_id: 'AC-1008',
    nombre: 'Bomba de Emergencia',
    estado: 'mantenimiento',
    ubicacion: 'Ruta 68 Km 15, Melipilla',
    edificio_id: 4,
    tipo: 'bomba'
  }
]);

print('🔧 Insertando sensores...');
db.sensores.insertMany([
  {sensor_id: 'SENSOR_AC1001_01', activo_id: 'AC-1001', tipo: 'temperatura', unidad: '°C'},
  {sensor_id: 'SENSOR_AC1001_02', activo_id: 'AC-1001', tipo: 'presion', unidad: 'bar'},
  {sensor_id: 'SENSOR_AC1002_01', activo_id: 'AC-1002', tipo: 'vibracion', unidad: 'mm/s'},
  {sensor_id: 'SENSOR_AC1002_02', activo_id: 'AC-1002', tipo: 'caudal', unidad: 'L/min'},
  {sensor_id: 'SENSOR_AC1003_01', activo_id: 'AC-1003', tipo: 'caudal', unidad: 'L/min'},
  {sensor_id: 'SENSOR_AC1003_02', activo_id: 'AC-1003', tipo: 'presion', unidad: 'bar'},
  {sensor_id: 'SENSOR_AC1004_01', activo_id: 'AC-1004', tipo: 'vibracion', unidad: 'mm/s'},
  {sensor_id: 'SENSOR_AC1004_02', activo_id: 'AC-1004', tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1005_01', activo_id: 'AC-1005', tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1005_02', activo_id: 'AC-1005', tipo: 'temperatura', unidad: '°C'},
  {sensor_id: 'SENSOR_AC1006_01', activo_id: 'AC-1006', tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1006_02', activo_id: 'AC-1006', tipo: 'temperatura', unidad: '°C'},
  {sensor_id: 'SENSOR_AC1007_01', activo_id: 'AC-1007', tipo: 'vibracion', unidad: 'mm/s'},
  {sensor_id: 'SENSOR_AC1007_02', activo_id: 'AC-1007', tipo: 'voltaje', unidad: 'V'},
  {sensor_id: 'SENSOR_AC1008_01', activo_id: 'AC-1008', tipo: 'presion', unidad: 'bar'},
  {sensor_id: 'SENSOR_AC1008_02', activo_id: 'AC-1008', tipo: 'caudal', unidad: 'L/min'}
]);

print('📊 Insertando estados de sensores...');
var now = new Date();
db.sensor_status.insertMany([
  {sensor_id: 'SENSOR_AC1001_01', activo_id: 'AC-1001', estado: 'conectado', ultima_lectura: now, valor_actual: 65.5},
  {sensor_id: 'SENSOR_AC1001_02', activo_id: 'AC-1001', estado: 'conectado', ultima_lectura: now, valor_actual: 6.2},
  {sensor_id: 'SENSOR_AC1002_01', activo_id: 'AC-1002', estado: 'conectado', ultima_lectura: now, valor_actual: 12.8},
  {sensor_id: 'SENSOR_AC1002_02', activo_id: 'AC-1002', estado: 'conectado', ultima_lectura: now, valor_actual: 425.0},
  {sensor_id: 'SENSOR_AC1003_01', activo_id: 'AC-1003', estado: 'conectado', ultima_lectura: now, valor_actual: 280.5},
  {sensor_id: 'SENSOR_AC1003_02', activo_id: 'AC-1003', estado: 'conectado', ultima_lectura: now, valor_actual: 7.8},
  {sensor_id: 'SENSOR_AC1004_01', activo_id: 'AC-1004', estado: 'conectado', ultima_lectura: now, valor_actual: 9.1},
  {sensor_id: 'SENSOR_AC1004_02', activo_id: 'AC-1004', estado: 'conectado', ultima_lectura: now, valor_actual: 400.0},
  {sensor_id: 'SENSOR_AC1005_01', activo_id: 'AC-1005', estado: 'desconectado', ultima_lectura: now, valor_actual: 0},
  {sensor_id: 'SENSOR_AC1005_02', activo_id: 'AC-1005', estado: 'desconectado', ultima_lectura: now, valor_actual: 0},
  {sensor_id: 'SENSOR_AC1006_01', activo_id: 'AC-1006', estado: 'conectado', ultima_lectura: now, valor_actual: 220.8},
  {sensor_id: 'SENSOR_AC1006_02', activo_id: 'AC-1006', estado: 'conectado', ultima_lectura: now, valor_actual: 42.3},
  {sensor_id: 'SENSOR_AC1007_01', activo_id: 'AC-1007', estado: 'conectado', ultima_lectura: now, valor_actual: 8.2},
  {sensor_id: 'SENSOR_AC1007_02', activo_id: 'AC-1007', estado: 'conectado', ultima_lectura: now, valor_actual: 395.0},
  {sensor_id: 'SENSOR_AC1008_01', activo_id: 'AC-1008', estado: 'desconectado', ultima_lectura: now, valor_actual: 0},
  {sensor_id: 'SENSOR_AC1008_02', activo_id: 'AC-1008', estado: 'desconectado', ultima_lectura: now, valor_actual: 0}
]);

print('✅ Datos cargados exitosamente');
print('📊 Resumen:');
print('   • Activos: ' + db.activos.countDocuments());
print('   • Sensores: ' + db.sensores.countDocuments());
print('   • Estados: ' + db.sensor_status.countDocuments());
print('🔗 Mapeo PostgreSQL ↔ MongoDB:');
print('   ID 1-8 ↔ AC-1001 a AC-1008');
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
echo -e "   • MongoDB IoT: AC-1001 a AC-1008"
