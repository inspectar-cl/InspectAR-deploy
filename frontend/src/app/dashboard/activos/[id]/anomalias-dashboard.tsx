/* eslint-disable @typescript-eslint/explicit-function-return-type -- tipo inferido */
'use client'

import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Grid,
  CircularProgress,
  Divider,
} from '@mui/material';
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

export default function AnomaliasDashboard({ activoId, activoNombre }: AnomaliasDashboardProps) {
  const [predicciones, setPredicciones] = useState<Prediccion[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAnomalias = async () => {
      try {
        setLoading(true);
        const response = await gs.get(
          `/ia/anomalies/activo/${activoId}?page=1&limit=50`
        ) as AnomaliasPaginadas;

        if (response && Array.isArray(response.anomalies)) {
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

          // Ordenar por timestamp descendente (más recientes primero)
          prediccionesMapeadas.sort((a, b) => 
            new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
          );

          setPredicciones(prediccionesMapeadas);
        }
      } catch (err) {
        setPredicciones([]);
      } finally {
        setLoading(false);
      }
    };

    void fetchAnomalias();

    // Actualizar anomalías cada 30 segundos
    const interval = setInterval(() => {
      void fetchAnomalias();
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
      
      <Typography variant="h4" mb={3}>
        Predicciones de Anomalías (ML)
      </Typography>

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