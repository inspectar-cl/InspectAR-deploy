/* eslint-disable camelcase -- Por consistencia */
/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import { useState, useEffect, useCallback } from 'react';
import Services from '@/modules/Services';
import { useUserToken } from '@/hooks/use-usertoken';

// Tipos del backend (actualizados para el nuevo endpoint)
interface ActivoBackend {
  id: number;
  nombre: string;
  tipo: string;
  estado: string;
  ubicacion: string;
  edificio_id: number;
  creado_en: string;
  sensores: SensorBackend[];
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
  datos?: {
    tiempo: string;
    valor: number;
  }[] | null;
}

interface ApiResponse {
  activos: ActivoBackend[];
  total: number;
  timestamp: string;
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
  const numeroMatch = /\d+/.exec(sensor_id);
  const numero = numeroMatch ? ` #${numeroMatch[0]}` : '';
  
  return `${nombreBase}${numero}`;
}

export function useActivosWithSensors() {
  const { user, isLoading: authLoading } = useUserToken();
  const [activos, setActivos] = useState<ActivoWithSensors[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchActivosWithSensors = useCallback(async () => {
    if (authLoading || !user) {
      console.log('user/token no disponibles para activos con sensores', { authLoading, user });
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      // Usar el nuevo endpoint con autenticación
      const response = await gs.authorizedGet('/obtener-todos-activos?sensores=true', user.token) as ApiResponse;
      console.log('📊 Respuesta completa de activos y sensores:', response);
      
      if (!response.activos || response.activos.length === 0) {
        setActivos([]);
        setLoading(false);
        return;
      }

      // Transformar los datos al formato del frontend
      const activosTransformados = response.activos.map((activo) => {
        // console.log(`🔧 Procesando activo: ${activo.nombre}`, activo);
        
        // Transformar sensores al formato del frontend
        const sensoresTransformados: SensorRow[] = activo.sensores.map((sensor: SensorBackend) => {
          // console.log(`📡 Procesando sensor: ${sensor.sensor_id}`, sensor);
          
          // Obtener el último valor de los datos históricos si existen
          let lastValue = 0;
          let history24h: SensorSample[] = [];
          
          if (sensor.datos && Array.isArray(sensor.datos) && sensor.datos.length > 0) {
            // console.log(`📊 Sensor ${sensor.sensor_id} tiene ${sensor.datos.length} datos`);
            
            // El último valor es el más reciente (primer elemento del array ya que viene ordenado descendente)
            const ultimoDato = sensor.datos[0];
            lastValue = ultimoDato.valor;
            // console.log(`🎯 Último valor para ${sensor.sensor_id}:`, lastValue);
            
            // Tomar los últimos 100 datos (o todos si hay menos de 100)
            // El backend ya los envía ordenados por tiempo descendente
            const datosRecientes = sensor.datos.slice(0, 100);
            
            // Convertir al formato SensorSample
            history24h = datosRecientes.map((dato) => ({
              ts: new Date(dato.tiempo),
              value: dato.valor,
              unit: sensor.unidad
            }));
            
            // console.log(`📈 History para ${sensor.sensor_id}:`, history24h.length, 'muestras (últimos 100 datos)');
          }
          
          // Determinar la última transmisión desde los datos reales
          let lastSeenDate: Date;
          if (sensor.datos && Array.isArray(sensor.datos) && sensor.datos.length > 0) {
            // Usar la fecha del primer dato (más reciente)
            lastSeenDate = new Date(sensor.datos[0].tiempo);
          } else if (sensor.last_seen) {
            // Fallback al campo last_seen si existe
            lastSeenDate = new Date(sensor.last_seen);
          } else {
            // Si no hay datos, usar una fecha antigua como fallback
            lastSeenDate = new Date(Date.now() - 10 * 60 * 1000);
          }
          
          return {
            id: sensor.sensor_id,
            name: generarNombreSensor(sensor.tipo, sensor.sensor_id),
            type: sensor.tipo,
            unit: sensor.unidad,
            lastValue,
            lastSeen: lastSeenDate,
            status: sensor.estado,
            history24h
          };
        });

        // console.log(`✅ Sensores transformados para ${activo.nombre}:`, sensoresTransformados);

        return {
          assetName: `${activo.nombre} — ${activo.ubicacion}`,
          imageUrl: TIPO_IMAGENES[activo.tipo.toLowerCase()] || TIPO_IMAGENES.transformador,
          sensores: sensoresTransformados,
        } as ActivoWithSensors;
      });

      // console.log(`✅ Activos procesados:`, activosTransformados.length, activosTransformados);
      setActivos(activosTransformados);

    } catch (err) {
      // console.error('Error fetching activos with sensors:', err);
      setError('Error al cargar los datos de activos y sensores');
    } finally {
      setLoading(false);
    }
  }, [authLoading, user]);

  // Refresh automático cada 30 segundos
  useEffect(() => {
    void fetchActivosWithSensors();
    
    // const interval = setInterval(() => {
    //   void fetchActivosWithSensors();
    // }, 30000);
    // return () => { 
    //   clearInterval(interval); 
    // };
  }, [fetchActivosWithSensors]);

  return {
    activos,
    loading,
    error,
    refresh: fetchActivosWithSensors
  };
}
