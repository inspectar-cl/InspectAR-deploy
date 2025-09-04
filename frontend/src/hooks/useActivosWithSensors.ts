import { useState, useEffect, useCallback } from 'react';
import Services from '@/modules/Services';

// Tipos del backend
interface ActivoBackend {
  id: number;
  activo_id: string;
  nombre: string;
  tipo: string;
  estado: string;
  ubicacion: string;
  edificio_id: number;
  creado_en: string;
}

interface SensorBackend {
  sensor_id: string;
  tipo: string;
  unidad: string;
  estado: 'connected' | 'disconnected' | 'never_connected';
  is_active: boolean;
  last_seen?: string;
  first_seen?: string;
  total_reports: number;
  created_at?: string;
  updated_at?: string;
}

interface ActivoWithSensorsBackend {
  id: string;
  activo_id: string;
  nombre: string;
  ubicacion: string;
  estado: string;
  id_edificio: string;
  total_sensores: number;
  sensores: SensorBackend[];
  resumen: {
    sensores_activos: number;
    sensores_desconectados: number;
    sensores_nunca_conectados: number;
  };
}

interface UltimosValores {
  [sensor_id: string]: {
    tiempo: string;
    valor: number;
  } | null;
}

// Tipos para el frontend
export interface SensorSample {
  ts: Date;
  value: number;
  unit: string;
}

export interface SensorRow {
  id: string;
  name: string;       // "Sensor Caudal 1"
  type: string;       // "caudal / presión / temp"
  unit: string;       // "L/s", "bar", "°C"
  lastValue: number;  // último valor recibido
  lastSeen: Date;     // timestamp último envío
  status: 'connected' | 'disconnected' | 'never_connected'; // estado del backend
  history24h?: SensorSample[];
}

export interface ActivoWithSensors {
  assetName: string;
  imageUrl: string;
  sensores: SensorRow[];
}

const gs = new Services();

// Mapeo de tipos a imágenes
const TIPO_IMAGENES: Record<string, string> = {
  "bomba de agua": "https://www.iprecom.com/wp-content/uploads/2019/07/bomba-de-agua.jpg",
  "transformador": "https://www.eabel.com/wp-content/uploads/2024/04/Electrical-Control-Panel-0401.webp",
  "caldera": "https://www.cerney.es/wp-content/uploads/20210429-10213168222-scaled-1-1.jpg",
  "ascensor": "https://elevabalear.com/wp-content/uploads/2024/11/caracteristicas-de-los-ascensores-electricos.jpg"
};

// Función para generar nombres descriptivos de sensores
function generarNombreSensor(tipo: string, sensor_id: string): string {
  const tipoNombres: Record<string, string> = {
    'temperatura': 'Sensor de Temperatura',
    'presion': 'Sensor de Presión',
    'caudal': 'Caudalímetro',
    'vibracion': 'Sensor de Vibración',
    'velocidad': 'Sensor de Velocidad',
    'humedad': 'Sensor de Humedad',
    'nivel': 'Sensor de Nivel'
  };
  
  const nombreBase = tipoNombres[tipo.toLowerCase()] || `Sensor ${tipo}`;
  // Extraer número del ID si existe
  const numeroMatch = sensor_id.match(/\d+/);
  const numero = numeroMatch ? ` #${numeroMatch[0]}` : '';
  
  return `${nombreBase}${numero}`;
}

export function useActivosWithSensors() {
  const [activos, setActivos] = useState<ActivoWithSensors[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchActivosWithSensors = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      // 1. Obtener todos los tipos de activos disponibles
      const tiposActivos = ['bomba de agua', 'caldera', 'ascensor', 'transformador'];
      const activosPromises = tiposActivos.map(async (tipo) => {
        try {
          const response = await gs.get(`/gestion/activos/tipo/${encodeURIComponent(tipo)}`);
          console.log(`📊 Respuesta activos tipo ${tipo}:`, response);
          return response.activos || [];
        } catch (error) {
          console.warn(`No se pudieron obtener activos de tipo ${tipo}:`, error);
          return [];
        }
      });

      const activosArrays = await Promise.all(activosPromises);
      const todosLosActivos = activosArrays.flat() as ActivoBackend[];
      
      console.log(`📋 Total activos encontrados:`, todosLosActivos.length, todosLosActivos);

      if (todosLosActivos.length === 0) {
        setActivos([]);
        setLoading(false);
        return;
      }

      // 2. Para cada activo, obtener sensores y últimos valores
      const activosConSensores = await Promise.all(
        todosLosActivos.map(async (activo) => {
          try {
            // Obtener estado de sensores desde el parser service
            const sensoresResponse = await gs.get(`/parser/activo/${activo.activo_id}/sensores/estado`);
            console.log(`🔍 Respuesta sensores para ${activo.activo_id}:`, sensoresResponse);
            
            // gs.get() ya parsea el JSON, no necesitas .json()
            const activoConSensores = sensoresResponse as ActivoWithSensorsBackend;

            // Obtener últimos valores
            let ultimosValores: UltimosValores = {};
            console.log("activo_id" , activo.activo_id);
            try {
              const valoresResponse = await gs.get(`/parser/lectura/${activo.activo_id}/datos/ultimo`);
              console.log(`📈 Respuesta últimos valores para ${activo.activo_id}:`, valoresResponse);
              // gs.get() ya parsea el JSON, usar directamente
              ultimosValores = valoresResponse || {};
            } catch (error) {
              console.warn(`No se pudieron obtener últimos valores para ${activo.activo_id}:`, error);
            }

            // Transformar sensores al formato del frontend
            const sensoresTransformados: SensorRow[] = activoConSensores.sensores.map((sensor) => {
              const ultimoValor = ultimosValores[sensor.sensor_id];
              
              return {
                id: sensor.sensor_id,
                name: generarNombreSensor(sensor.tipo, sensor.sensor_id),
                type: sensor.tipo,
                unit: sensor.unidad,
                lastValue: ultimoValor?.valor || 0,
                lastSeen: sensor.last_seen ? new Date(sensor.last_seen) : new Date(Date.now() - 10 * 60 * 1000), // 10 min ago si no hay fecha
                status: sensor.estado, // Usar el estado directamente del backend
                // history24h se puede agregar más tarde si es necesario
              };
            });

            return {
              assetName: `${activo.nombre} — ${activo.ubicacion}`,
              imageUrl: TIPO_IMAGENES[activo.tipo.toLowerCase()] || TIPO_IMAGENES["transformador"],
              sensores: sensoresTransformados,
            } as ActivoWithSensors;

          } catch (error) {
            console.error(`❌ Error procesando activo ${activo.activo_id}:`, error);
            return null;
          }
        })
      );

      // Filtrar activos que no pudieron ser procesados
      const activosValidos = activosConSensores.filter(
        (activo): activo is ActivoWithSensors => activo !== null
      );

      console.log(`✅ Activos válidos procesados:`, activosValidos.length, activosValidos);
      setActivos(activosValidos);

    } catch (error) {
      console.error('Error fetching activos with sensors:', error);
      setError('Error al cargar los datos de activos y sensores');
    } finally {
      setLoading(false);
    }
  }, []);

  // Refresh automático cada 30 segundos
  useEffect(() => {
    fetchActivosWithSensors();
    
    const interval = setInterval(fetchActivosWithSensors, 30000);
    return () => clearInterval(interval);
  }, [fetchActivosWithSensors]);

  return {
    activos,
    loading,
    error,
    refresh: fetchActivosWithSensors
  };
}
