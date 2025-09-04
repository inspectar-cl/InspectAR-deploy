#!/bin/bash

# ===================================================================
# SCRIPT DE CARGA COMPLETA DE DATOS A MONGODB (IoT Service)
# ===================================================================
# Este script envía todos los datos necesarios a MongoDB para
# que el servicio IoT tenga información completa y sincronizada

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 INICIANDO CARGA COMPLETA DE DATOS A MONGODB${NC}"
echo "================================================"

# Verificar si MongoDB está corriendo
if ! docker ps | grep -q "mongu"; then
    echo -e "${RED}❌ Error: MongoDB no está ejecutándose${NC}"
    echo "Ejecuta: docker-compose up -d mongu"
    exit 1
fi

echo -e "${YELLOW}📊 Enviando datos a MongoDB...${NC}"

# ===================================================================
# 1. LIMPIAR COLECCIONES EXISTENTES
# ===================================================================
echo -e "\n${BLUE}🧹 Limpiando colecciones existentes...${NC}"

docker exec mongu mongosh iot_db --eval "
print('🧹 Limpiando colecciones existentes...');
db.activos.deleteMany({});
db.sensores.deleteMany({});
db.sensor_status.deleteMany({});
db.device_registry.deleteMany({});
print('✅ Colecciones limpiadas');
" --quiet

# ===================================================================
# 2. INSERTAR ACTIVOS PRINCIPALES
# ===================================================================
echo -e "\n${BLUE}📦 Insertando activos principales...${NC}"

docker exec mongu mongosh iot_db --eval "
print('📦 Insertando activos principales...');

db.activos.insertMany([
  {
    activo_id: 'AC-1001',
    nombre: 'Caldera Principal',
    estado: 'operativo',
    ubicacion: 'Av. Providencia 123, Santiago',
    edificio_id: 1,
    tipo: 'caldera',
    marca: 'Bosch',
    modelo: 'Universal U042',
    año_instalacion: 2020,
    sensores: [
      { sensor_id: 'SENSOR_AC1001_01', tipo: 'temperatura', unidad: '°C' },
      { sensor_id: 'SENSOR_AC1001_02', tipo: 'presion', unidad: 'bar' }
    ]
  },
  {
    activo_id: 'AC-1002',
    nombre: 'Bomba Centrífuga A',
    estado: 'operativo',
    ubicacion: 'Av. Providencia 123, Santiago',
    edificio_id: 1,
    tipo: 'bomba',
    marca: 'Grundfos',
    modelo: 'CM 3-4',
    año_instalacion: 2021,
    sensores: [
      { sensor_id: 'SENSOR_AC1002_01', tipo: 'vibracion', unidad: 'mm/s' },
      { sensor_id: 'SENSOR_AC1002_02', tipo: 'caudal', unidad: 'L/min' }
    ]
  },
  {
    activo_id: 'AC-1004',
    nombre: 'Bomba Hidráulica 1',
    estado: 'operativo',
    ubicacion: 'Av. Las Condes 456, Las Condes',
    edificio_id: 2,
    tipo: 'bomba',
    marca: 'Wilo',
    modelo: 'TOP-S 25/7',
    año_instalacion: 2019,
    sensores: [
      { sensor_id: 'SENSOR_AC1004_01', tipo: 'caudal', unidad: 'L/min' },
      { sensor_id: 'SENSOR_AC1004_02', tipo: 'presion', unidad: 'bar' }
    ]
  },
  {
    activo_id: 'AC-1006',
    nombre: 'Transformador Principal',
    estado: 'operativo',
    ubicacion: 'Av. Vicuña Mackenna 789, La Florida',
    edificio_id: 3,
    tipo: 'transformador',
    marca: 'ABB',
    modelo: 'ONAN 500 kVA',
    año_instalacion: 2018,
    sensores: [
      { sensor_id: 'SENSOR_AC1006_01', tipo: 'voltaje', unidad: 'V' },
      { sensor_id: 'SENSOR_AC1006_02', tipo: 'temperatura', unidad: '°C' }
    ]
  },
  {
    activo_id: 'AC-1007',
    nombre: 'Ascensor Central',
    estado: 'operativo',
    ubicacion: 'Av. Providencia 123, Santiago',
    edificio_id: 1,
    tipo: 'ascensor',
    marca: 'Otis',
    modelo: 'Gen2 Premier',
    año_instalacion: 2022,
    sensores: [
      { sensor_id: 'SENSOR_AC1007_01', tipo: 'vibracion', unidad: 'mm/s' },
      { sensor_id: 'SENSOR_AC1007_02', tipo: 'voltaje', unidad: 'V' }
    ]
  },
  {
    activo_id: 'AC-1008',
    nombre: 'Bomba de Emergencia',
    estado: 'mantenimiento',
    ubicacion: 'Ruta 68 Km 15, Melipilla',
    edificio_id: 4,
    tipo: 'bomba',
    marca: 'Pedrollo',
    modelo: 'JSWm 2AX',
    año_instalacion: 2020,
    sensores: [
      { sensor_id: 'SENSOR_AC1008_01', tipo: 'presion', unidad: 'bar' },
      { sensor_id: 'SENSOR_AC1008_02', tipo: 'caudal', unidad: 'L/min' }
    ]
  }
]);

print('✅ Insertados ' + db.activos.countDocuments() + ' activos');
" --quiet

# ===================================================================
# 3. INSERTAR CONFIGURACIÓN DE SENSORES
# ===================================================================
echo -e "\n${BLUE}🔧 Insertando configuración de sensores...${NC}"

docker exec mongu mongosh iot_db --eval "
print('🔧 Insertando configuración de sensores...');

db.sensores.insertMany([
  // Sensores Caldera Principal (AC-1001)
  {
    sensor_id: 'SENSOR_AC1001_01',
    activo_id: 'AC-1001',
    tipo: 'temperatura',
    unidad: '°C',
    rango_min: 15,
    rango_max: 85,
    umbral_alerta: 80,
    umbral_critico: 90,
    frecuencia_lectura: 60, // segundos
    ubicacion_sensor: 'Cámara de combustión'
  },
  {
    sensor_id: 'SENSOR_AC1001_02',
    activo_id: 'AC-1001',
    tipo: 'presion',
    unidad: 'bar',
    rango_min: 0,
    rango_max: 10,
    umbral_alerta: 8,
    umbral_critico: 9.5,
    frecuencia_lectura: 30,
    ubicacion_sensor: 'Línea principal'
  },
  
  // Sensores Bomba Centrífuga A (AC-1002)
  {
    sensor_id: 'SENSOR_AC1002_01',
    activo_id: 'AC-1002',
    tipo: 'vibracion',
    unidad: 'mm/s',
    rango_min: 0,
    rango_max: 25,
    umbral_alerta: 20,
    umbral_critico: 23,
    frecuencia_lectura: 120,
    ubicacion_sensor: 'Cojinete principal'
  },
  {
    sensor_id: 'SENSOR_AC1002_02',
    activo_id: 'AC-1002',
    tipo: 'caudal',
    unidad: 'L/min',
    rango_min: 100,
    rango_max: 800,
    umbral_alerta: 150,
    umbral_critico: 120,
    frecuencia_lectura: 60,
    ubicacion_sensor: 'Salida de bomba'
  },
  
  // Sensores Bomba Hidráulica 1 (AC-1004)
  {
    sensor_id: 'SENSOR_AC1004_01',
    activo_id: 'AC-1004',
    tipo: 'caudal',
    unidad: 'L/min',
    rango_min: 50,
    rango_max: 500,
    umbral_alerta: 80,
    umbral_critico: 60,
    frecuencia_lectura: 45,
    ubicacion_sensor: 'Línea de descarga'
  },
  {
    sensor_id: 'SENSOR_AC1004_02',
    activo_id: 'AC-1004',
    tipo: 'presion',
    unidad: 'bar',
    rango_min: 1,
    rango_max: 15,
    umbral_alerta: 12,
    umbral_critico: 14,
    frecuencia_lectura: 30,
    ubicacion_sensor: 'Cabezal de presión'
  },
  
  // Sensores Transformador Principal (AC-1006)
  {
    sensor_id: 'SENSOR_AC1006_01',
    activo_id: 'AC-1006',
    tipo: 'voltaje',
    unidad: 'V',
    rango_min: 200,
    rango_max: 240,
    umbral_alerta: 235,
    umbral_critico: 250,
    frecuencia_lectura: 15,
    ubicacion_sensor: 'Salida secundaria'
  },
  {
    sensor_id: 'SENSOR_AC1006_02',
    activo_id: 'AC-1006',
    tipo: 'temperatura',
    unidad: '°C',
    rango_min: 20,
    rango_max: 60,
    umbral_alerta: 55,
    umbral_critico: 65,
    frecuencia_lectura: 90,
    ubicacion_sensor: 'Devanado primario'
  },
  
  // Sensores Ascensor Central (AC-1007)
  {
    sensor_id: 'SENSOR_AC1007_01',
    activo_id: 'AC-1007',
    tipo: 'vibracion',
    unidad: 'mm/s',
    rango_min: 0,
    rango_max: 15,
    umbral_alerta: 12,
    umbral_critico: 14,
    frecuencia_lectura: 180,
    ubicacion_sensor: 'Motor de tracción'
  },
  {
    sensor_id: 'SENSOR_AC1007_02',
    activo_id: 'AC-1007',
    tipo: 'voltaje',
    unidad: 'V',
    rango_min: 380,
    rango_max: 420,
    umbral_alerta: 410,
    umbral_critico: 430,
    frecuencia_lectura: 60,
    ubicacion_sensor: 'Tablero principal'
  },
  
  // Sensores Bomba de Emergencia (AC-1008)
  {
    sensor_id: 'SENSOR_AC1008_01',
    activo_id: 'AC-1008',
    tipo: 'presion',
    unidad: 'bar',
    rango_min: 0,
    rango_max: 8,
    umbral_alerta: 6,
    umbral_critico: 7.5,
    frecuencia_lectura: 30,
    ubicacion_sensor: 'Salida principal'
  },
  {
    sensor_id: 'SENSOR_AC1008_02',
    activo_id: 'AC-1008',
    tipo: 'caudal',
    unidad: 'L/min',
    rango_min: 0,
    rango_max: 300,
    umbral_alerta: 50,
    umbral_critico: 30,
    frecuencia_lectura: 45,
    ubicacion_sensor: 'Medidor de flujo'
  }
]);

print('✅ Insertados ' + db.sensores.countDocuments() + ' sensores');
" --quiet

# ===================================================================
# 4. INSERTAR ESTADOS ACTUALES DE SENSORES
# ===================================================================
echo -e "\n${BLUE}📊 Insertando estados actuales de sensores...${NC}"

docker exec mongu mongosh iot_db --eval "
print('📊 Insertando estados actuales de sensores...');

var now = new Date();
var hace5min = new Date(now.getTime() - 5*60000);
var hace2min = new Date(now.getTime() - 2*60000);

db.sensor_status.insertMany([
  // Estados Caldera Principal (AC-1001)
  {
    sensor_id: 'SENSOR_AC1001_01',
    activo_id: 'AC-1001',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 65.5,
    alarma_activa: false,
    calidad_señal: 95,
    nivel_bateria: null, // sensor cableado
    total_lecturas: 1440 // 24 horas * 60 lecturas/hora
  },
  {
    sensor_id: 'SENSOR_AC1001_02',
    activo_id: 'AC-1001',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 6.2,
    alarma_activa: false,
    calidad_señal: 92,
    nivel_bateria: null,
    total_lecturas: 2880
  },
  
  // Estados Bomba Centrífuga A (AC-1002)
  {
    sensor_id: 'SENSOR_AC1002_01',
    activo_id: 'AC-1002',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 12.8,
    alarma_activa: false,
    calidad_señal: 88,
    nivel_bateria: 85,
    total_lecturas: 720
  },
  {
    sensor_id: 'SENSOR_AC1002_02',
    activo_id: 'AC-1002',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 425.0,
    alarma_activa: false,
    calidad_señal: 90,
    nivel_bateria: null,
    total_lecturas: 1440
  },
  
  // Estados Bomba Hidráulica 1 (AC-1004)
  {
    sensor_id: 'SENSOR_AC1004_01',
    activo_id: 'AC-1004',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 320.5,
    alarma_activa: false,
    calidad_señal: 94,
    nivel_bateria: 78,
    total_lecturas: 1920
  },
  {
    sensor_id: 'SENSOR_AC1004_02',
    activo_id: 'AC-1004',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 8.5,
    alarma_activa: false,
    calidad_señal: 91,
    nivel_bateria: null,
    total_lecturas: 2880
  },
  
  // Estados Transformador Principal (AC-1006)
  {
    sensor_id: 'SENSOR_AC1006_01',
    activo_id: 'AC-1006',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 220.8,
    alarma_activa: false,
    calidad_señal: 97,
    nivel_bateria: null,
    total_lecturas: 5760 // 4 lecturas/min * 24h * 60min
  },
  {
    sensor_id: 'SENSOR_AC1006_02',
    activo_id: 'AC-1006',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 42.3,
    alarma_activa: false,
    calidad_señal: 89,
    nivel_bateria: null,
    total_lecturas: 960
  },
  
  // Estados Ascensor Central (AC-1007)
  {
    sensor_id: 'SENSOR_AC1007_01',
    activo_id: 'AC-1007',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 8.2,
    alarma_activa: false,
    calidad_señal: 86,
    nivel_bateria: 92,
    total_lecturas: 480
  },
  {
    sensor_id: 'SENSOR_AC1007_02',
    activo_id: 'AC-1007',
    estado: 'conectado',
    ultima_lectura: hace2min,
    valor_actual: 395.0,
    alarma_activa: false,
    calidad_señal: 93,
    nivel_bateria: null,
    total_lecturas: 1440
  },
  
  // Estados Bomba de Emergencia (AC-1008) - EN MANTENIMIENTO
  {
    sensor_id: 'SENSOR_AC1008_01',
    activo_id: 'AC-1008',
    estado: 'desconectado',
    ultima_lectura: hace5min,
    valor_actual: 0,
    alarma_activa: true,
    calidad_señal: 0,
    nivel_bateria: 45,
    total_lecturas: 2400,
    motivo_desconexion: 'Mantenimiento programado'
  },
  {
    sensor_id: 'SENSOR_AC1008_02',
    activo_id: 'AC-1008',
    estado: 'desconectado',
    ultima_lectura: hace5min,
    valor_actual: 0,
    alarma_activa: true,
    calidad_señal: 0,
    nivel_bateria: 45,
    total_lecturas: 1920,
    motivo_desconexion: 'Mantenimiento programado'
  }
]);

print('✅ Insertados ' + db.sensor_status.countDocuments() + ' estados de sensores');
" --quiet

# ===================================================================
# 5. INSERTAR REGISTRO DE DISPOSITIVOS
# ===================================================================
echo -e "\n${BLUE}🔌 Insertando registro de dispositivos...${NC}"

docker exec mongu mongosh iot_db --eval "
print('🔌 Insertando registro de dispositivos...');

var now = new Date();
var ayer = new Date(now.getTime() - 24*60*60*1000);
var hace5min = new Date(now.getTime() - 5*60000);

db.device_registry.insertMany([
  {
    device_id: 'GATEWAY_PROV_001',
    tipo: 'gateway',
    ubicacion: 'Av. Providencia 123',
    edificio_id: 1,
    ip_address: '192.168.1.100',
    mac_address: '00:1B:44:11:3A:B7',
    firmware_version: '2.1.4',
    estado: 'activo',
    ultima_conexion: now,
    sensores_conectados: ['SENSOR_AC1001_01', 'SENSOR_AC1001_02', 'SENSOR_AC1002_01', 'SENSOR_AC1002_02', 'SENSOR_AC1007_01', 'SENSOR_AC1007_02'],
    fecha_instalacion: ayer
  },
  {
    device_id: 'GATEWAY_CONDE_001',
    tipo: 'gateway',
    ubicacion: 'Av. Las Condes 456',
    edificio_id: 2,
    ip_address: '192.168.2.100',
    mac_address: '00:1B:44:11:3A:C8',
    firmware_version: '2.1.4',
    estado: 'activo',
    ultima_conexion: now,
    sensores_conectados: ['SENSOR_AC1004_01', 'SENSOR_AC1004_02'],
    fecha_instalacion: ayer
  },
  {
    device_id: 'GATEWAY_VICU_001',
    tipo: 'gateway',
    ubicacion: 'Av. Vicuña Mackenna 789',
    edificio_id: 3,
    ip_address: '192.168.3.100',
    mac_address: '00:1B:44:11:3A:D9',
    firmware_version: '2.1.3',
    estado: 'activo',
    ultima_conexion: now,
    sensores_conectados: ['SENSOR_AC1006_01', 'SENSOR_AC1006_02'],
    fecha_instalacion: ayer
  },
  {
    device_id: 'GATEWAY_MELI_001',
    tipo: 'gateway',
    ubicacion: 'Ruta 68 Km 15',
    edificio_id: 4,
    ip_address: '192.168.4.100',
    mac_address: '00:1B:44:11:3A:EA',
    firmware_version: '2.0.8',
    estado: 'mantenimiento',
    ultima_conexion: hace5min,
    sensores_conectados: ['SENSOR_AC1008_01', 'SENSOR_AC1008_02'],
    fecha_instalacion: ayer,
    observaciones: 'Gateway en mantenimiento junto con bomba de emergencia'
  }
]);

print('✅ Insertados ' + db.device_registry.countDocuments() + ' dispositivos');
" --quiet

# ===================================================================
# 6. VERIFICACIÓN FINAL
# ===================================================================
echo -e "\n${BLUE}🔍 Verificación final de datos...${NC}"

docker exec mongu mongosh iot_db --eval "
print('🔍 Verificación final de datos...');
print('');
print('📊 RESUMEN DE DATOS INSERTADOS:');
print('   • Activos: ' + db.activos.countDocuments());
print('   • Sensores: ' + db.sensores.countDocuments());
print('   • Estados de sensores: ' + db.sensor_status.countDocuments());
print('   • Dispositivos registrados: ' + db.device_registry.countDocuments());
print('');

print('🏭 ACTIVOS POR EDIFICIO:');
db.activos.aggregate([
  { \$group: { _id: '\$edificio_id', count: { \$sum: 1 }, activos: { \$push: '\$activo_id' } } },
  { \$sort: { _id: 1 } }
]).forEach(function(result) {
  print('   Edificio ' + result._id + ': ' + result.count + ' activos (' + result.activos.join(', ') + ')');
});

print('');
print('🔧 SENSORES POR TIPO:');
db.sensores.aggregate([
  { \$group: { _id: '\$tipo', count: { \$sum: 1 } } },
  { \$sort: { _id: 1 } }
]).forEach(function(result) {
  print('   ' + result._id + ': ' + result.count + ' sensores');
});

print('');
print('📡 ESTADO DE CONECTIVIDAD:');
db.sensor_status.aggregate([
  { \$group: { _id: '\$estado', count: { \$sum: 1 } } }
]).forEach(function(result) {
  print('   ' + result._id + ': ' + result.count + ' sensores');
});

print('');
print('✅ Carga de datos completada exitosamente');
" --quiet

echo -e "\n${GREEN}✅ CARGA COMPLETA DE DATOS FINALIZADA${NC}"
echo "================================================"
echo -e "${YELLOW}📋 Resumen:${NC}"
echo "   • 6 activos insertados con información completa"
echo "   • 12 sensores configurados con parámetros técnicos"
echo "   • 12 estados de sensores con datos en tiempo real"
echo "   • 4 gateways registrados por edificio"
echo ""
echo -e "${GREEN}🎯 MongoDB está listo para el servicio IoT${NC}"
echo -e "${BLUE}💡 Puedes verificar los datos con: docker exec mongu mongosh iot_db${NC}"
