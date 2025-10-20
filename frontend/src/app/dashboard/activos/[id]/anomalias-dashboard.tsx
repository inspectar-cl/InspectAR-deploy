/* eslint-disable @typescript-eslint/explicit-function-return-type -- tipo inferido */
'use client'

import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Grid,
  CircularProgress,
  Divider,
  Chip,
} from '@mui/material';
import { useUserToken } from '@/hooks/use-usertoken';
import { Activity as ActivityIcon, AlertTriangle as AlertIcon } from 'lucide-react';
import { type Prediccion } from '@/types/prediccion';
import { getColorBySeverity } from '@/utils/SeverityUtils';
import MetricCard from '@/components/dashboard/prediccion/MetricCard';
import PrediccionCards from '@/components/dashboard/prediccion/PrediccionCards';
import AnomalyChart from '@/components/dashboard/prediccion/AnomalyChart';
import Services from '@/modules/Services';

const gs = new Services();

// Interfaces para anomalías
interface AnomaliaAPI {
  id: number;
  activo_id: number;
  sensor_id: string;
  timestamp: string;
  anomaly_score: number;
  anomaly_likelihood: number;
  severidad: string;
  descripcion: string;
  threshold: number;
  is_anomaly: number;
  created_at: string;
}

interface AnomaliasPaginadas {
  activo_id: number;
  anomalies: AnomaliaAPI[];
  total: number;
  page: number;
  limit: number;
  has_more: boolean;
}

interface AnomaliasDashboardProps {
  activoId: number;
  activoNombre: string;
}

// Función para calcular severidad desde likelihood (igual que el API)
const severityFromLikelihood = (likelihood: number): "Baja" | "Media" | "Alta" => {
  if (likelihood < 40) {
    return "Baja";
  } else if (likelihood < 70) {
    return "Media";
  } else {
    return "Alta";
  }
};

// Función para obtener descripción según severidad (igual que el API)
const getDescripcionPorSeveridad = (severidad: "Baja" | "Media" | "Alta"): string => {
  const descripciones = {
    "Baja": "Funcionamiento dentro del rango esperado.",
    "Media": "Comportamiento irregular detectado. Revisar condiciones operativas.",
    "Alta": "Anomalía crítica detectada. Atención prioritaria requerida."
  };
  return descripciones[severidad];
};

// Función para generar datos mock (lógica del API)
const generateMockData = (activoId: number): Prediccion[] => {
  const now = new Date();
  const predicciones: Prediccion[] = [];

  // Generar 50 registros de ejemplo
  for (let i = 0; i < 50; i++) {
    const timestamp = new Date(now.getTime() - i * 60000); // Cada minuto hacia atrás
    
    // 1. Generar anomaly_score (0-100)
    const anomalyScore = Math.floor(Math.random() * 101); // 0-100
    
    // 2. anomaly_likelihood >= anomaly_score (como en el modelo real)
    const anomalyLikelihood = anomalyScore + Math.floor(Math.random() * (101 - anomalyScore));
    
    // 3. Calcular severidad basándose en likelihood (lógica del API)
    const severidad = severityFromLikelihood(anomalyLikelihood);
    
    // 4. Threshold aleatorio entre 0-100
    const threshold = Math.floor(Math.random() * 101);
    
    predicciones.push({
      id: i + 1,
      activoId,
      timestamp: timestamp.toISOString(),
      anomalyScore,
      anomalyLikelihood,
      severidad,
      descripcion: getDescripcionPorSeveridad(severidad),
      threshold,
      is_anomaly: anomalyLikelihood > threshold, // Compara likelihood con threshold
    });
  }

  return predicciones;
};

export default function AnomaliasDashboard({ activoId, activoNombre }: AnomaliasDashboardProps) {
  const [predicciones, setPredicciones] = useState<Prediccion[]>([]);
  const [loading, setLoading] = useState(true);
  const [usingMockData, setUsingMockData] = useState(false);
  const { user, isLoading} = useUserToken();

  useEffect(() => {
    const fetchAnomalias = async (retryCount = 0, maxRetries = 3) => {
      try {
        setLoading(true);
        if (isLoading || !user) {return;}
                
        // Intentar obtener datos reales
        const response = await gs.authorizedGet(
          `/ML/anomalies/activo/${activoId}?page=1&limit=50`, user.token
        ) as AnomaliasPaginadas;

        if (response && Array.isArray(response.anomalies) && response.anomalies.length > 0) {
          // Si hay datos reales, usarlos
          
          const prediccionesMapeadas = response.anomalies.map((anomalia) => ({
            id: anomalia.id,
            activoId: anomalia.activo_id,
            timestamp: anomalia.timestamp,
            anomalyScore: Math.round(anomalia.anomaly_score),
            anomalyLikelihood: Math.round(anomalia.anomaly_likelihood),
            severidad: anomalia.severidad as "Baja" | "Media" | "Alta",
            descripcion: anomalia.descripcion,
            threshold: Math.round(anomalia.threshold),
            is_anomaly: anomalia.is_anomaly === 1,
            sensorId: anomalia.sensor_id,
            createdAt: anomalia.created_at,
          }));

          // Ordenar por timestamp descendente
          prediccionesMapeadas.sort((a, b) => 
            new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
          );

          setPredicciones(prediccionesMapeadas);
          setUsingMockData(false);
        } else {
          // Si no hay datos, lanzar error para reintentar
          throw new Error('Respuesta vacía del servidor');
        }
      } catch (err) {
        
        // Si aún quedan reintentos, esperar y volver a intentar
        if (retryCount < maxRetries) {
          const delayMs = 1000 * (retryCount + 1); // Espera incremental: 1s, 2s, 3s
          
          setTimeout(() => {
            void fetchAnomalias(retryCount + 1, maxRetries);
          }, delayMs);
          return; // No ejecutar el finally todavía
        }
        
        // Si se agotaron los reintentos, usar datos mock        
        const mockData = generateMockData(activoId);
        mockData.sort((a, b) => 
          new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
        );
        
        setPredicciones(mockData);
        setUsingMockData(true);
      } finally {
        // Solo quitar loading cuando termine completamente (o se agoten reintentos)
        if (retryCount >= maxRetries || predicciones.length > 0) {
          setLoading(false);
        }
      }
    };

    void fetchAnomalias();

    // Actualizar anomalías cada 30 segundos (sin reintentos en el polling)
    const interval = setInterval(() => {
      void fetchAnomalias(0, 0); // Sin reintentos en el polling automático
    }, 30000);

    return () => { clearInterval(interval); };
  }, [activoId]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
        <Typography ml={2}>Cargando predicciones...</Typography>
      </Box>
    );
  }

  if (predicciones.length === 0) {
    return (
      <Box 
        sx={{ 
          p: 3, 
          textAlign: 'center', 
          backgroundColor: '#f5f5f5', 
          borderRadius: 2 
        }}
      >
        <Typography variant="body1" color="text.secondary">
          No hay anomalías detectadas para este activo.
        </Typography>
      </Box>
    );
  }

  const ultimaPrediccion = predicciones[0];
  const totalAnomalias = predicciones.filter(p => p.is_anomaly).length;

  return (
    <Box>
      <Divider sx={{ my: 4 }} />
      
      <Box display="flex" alignItems="center" gap={2} mb={3}>
        <Typography variant="h4">
          Predicciones de Anomalías (ML)
        </Typography>
        {usingMockData && (
          <Chip 
            label="Datos de Simulación" 
            color="warning" 
            size="small"
            icon={<AlertIcon size={16} />}
          />
        )}
        {!usingMockData && (
          <Chip 
            label="Datos Reales" 
            color="success" 
            size="small"
            icon={<ActivityIcon size={16} />}
          />
        )}
      </Box>

      <Grid container spacing={3}>
        {/* Métricas principales de la última predicción */}
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <MetricCard
            title="Último Anomaly Score"
            value={ultimaPrediccion.anomalyScore}
            color={getColorBySeverity(ultimaPrediccion.severidad)}
            icon={<ActivityIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <MetricCard
            title="Último Anomaly Likelihood"
            value={ultimaPrediccion.anomalyLikelihood}
            color={getColorBySeverity(ultimaPrediccion.severidad)}
            icon={<AlertIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <MetricCard
            title="Total Anomalías"
            value={totalAnomalias}
            color="#2196f3"
            icon={<AlertIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <Box
            sx={{
              p: 2,
              borderRadius: 2,
              boxShadow: 1,
              display: 'flex',
              flexDirection: 'column',
              gap: 1
            }}
          >
            <Box display="flex" alignItems="center" gap={1}>
              <ActivityIcon size={20} style={{ color: '#4caf50' }} />
              <Typography variant="body2" color="text.secondary">
                Última Actualización
              </Typography>
            </Box>
            <Typography variant="h5" fontWeight="bold">
              {new Date(ultimaPrediccion.timestamp).toLocaleTimeString('es-ES', {
                hour: '2-digit',
                minute: '2-digit'
              })}
            </Typography>
          </Box>
        </Grid>

        {/* Gráfico de anomalías */}
        <Grid size={{ xs: 12 }}>
          <AnomalyChart 
            predicciones={predicciones} 
            activoNombre={activoNombre} 
          />
        </Grid>

        {/* Tarjetas de predicciones detalladas */}
        <Grid size={{ xs: 12 }}>
          <Typography variant="h6" mb={2}>
            Últimas 10 Anomalías Detectadas
          </Typography>
          <PrediccionCards predicciones={predicciones.slice(0, 10)} />
        </Grid>
      </Grid>
    </Box>
  );
}