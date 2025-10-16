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
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      // Usar el nuevo endpoint con autenticación
      const response = await gs.authorizedGet('/obtener-todos-activos?sensores=true', user.token) as ApiResponse;
      
      if (!response.activos || response.activos.length === 0) {
        setActivos([]);
        setLoading(false);
        return;
      }

      // Transformar los datos al formato del frontend
      const activosTransformados = response.activos.map((activo) => {
        // Transformar sensores al formato del frontend
        const sensoresTransformados: SensorRow[] = activo.sensores.map((sensor: SensorBackend) => {
          // Obtener el último valor de los datos históricos si existen
          let lastValue = 0;
          let history24h: SensorSample[] = [];
          
          if (sensor.datos && Array.isArray(sensor.datos) && sensor.datos.length > 0) {
            // El último valor es el más reciente (último elemento del array)
            const ultimoDato = sensor.datos[sensor.datos.length - 1];
            lastValue = ultimoDato.valor;
            
            // Convertir todos los datos al formato SensorSample
            history24h = sensor.datos.map((dato) => ({
              ts: new Date(dato.tiempo),
              value: dato.valor,
              unit: sensor.unidad
            }));
            
            // Filtrar solo las últimas 24 horas
            const hace24h = new Date(Date.now() - 24 * 60 * 60 * 1000);
            history24h = history24h.filter(sample => sample.ts >= hace24h);
            
            // Ordenar por timestamp descendente (más reciente primero)
            history24h.sort((a, b) => b.ts.getTime() - a.ts.getTime());
            
          }
          
          return {
            id: sensor.sensor_id,
            name: generarNombreSensor(sensor.tipo, sensor.sensor_id),
            type: sensor.tipo,
            unit: sensor.unidad,
            lastValue,
            lastSeen: sensor.last_seen ? new Date(sensor.last_seen) : new Date(Date.now() - 10 * 60 * 1000),
            status: sensor.estado,
            history24h
          };
        });

        return {
          assetName: `${activo.nombre} — ${activo.ubicacion}`,
          imageUrl: TIPO_IMAGENES[activo.tipo.toLowerCase()] || TIPO_IMAGENES.transformador,
          sensores: sensoresTransformados,
        } as ActivoWithSensors;
      });

      setActivos(activosTransformados);

    } catch (err) {
      setError('Error al cargar los datos de activos y sensores');
    } finally {
      setLoading(false);
    }
  }, [authLoading, user]);

  // Refresh automático cada 30 segundos
  useEffect(() => {
    void fetchActivosWithSensors();
    
    const interval = setInterval(() => {
      void fetchActivosWithSensors();
    }, 30000);
    return () => { 
      clearInterval(interval); 
    };
  }, [fetchActivosWithSensors]);

  return {
    activos,
    loading,
    error,
    refresh: fetchActivosWithSensors
  };
}
