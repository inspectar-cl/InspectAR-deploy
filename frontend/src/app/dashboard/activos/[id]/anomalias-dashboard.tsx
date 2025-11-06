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
import { useUserToken } from '@/hooks/use-usertoken';
import { Activity as ActivityIcon, AlertTriangle as AlertIcon } from 'lucide-react';
import { type Prediccion } from '@/types/prediccion';
import { getColorBySeverity } from '@/utils/SeverityUtils';
import MetricCard from '@/components/dashboard/prediccion/MetricCard';
import PrediccionCards from '@/components/dashboard/prediccion/PrediccionCards';
import AnomalyChart from '@/components/dashboard/prediccion/AnomalyChart';
import Services from '@/modules/Services';
import { useDriverTour } from "@/components/tutorial/use-driver-tour";
import type { TourKey } from "@/components/tutorial/tour-config";
import { TourButton } from "@/components/tutorial/tour-button";

const gs = new Services();

interface AnomaliasDashboardProps {
  activoId: number;
  activoNombre: string;
}

export default function AnomaliasDashboard({ activoId, activoNombre }: AnomaliasDashboardProps) {
  const { user, isLoading: isUserLoading } = useUserToken();
  const [predicciones, setPredicciones] = useState<Prediccion[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const tourKey: TourKey = 'prediccion-activo';
  const { startTour } = useDriverTour(tourKey);

  useEffect(() => {
<<<<<<< HEAD
    
=======
>>>>>>> origin/feature-front
    if (isUserLoading || !user) {
      return;
    }

    const fetchAnomalias = async () => {
      try {
        setIsLoading(true);
        
        interface AnomaliasPaginadas {
          activo_id?: number;
          data?: {
            id?: number;
            activo_id?: number;
            timestamp?: string;
            anomaly_score?: number;
            anomaly_likelihood?: number;
            severidad?: "Baja" | "Media" | "Alta";
            descripcion?: string;
            threshold?: number;
            is_anomaly?: number | boolean;
            most_influential_variable?: string;
            contribution_magnitude?: number;
          }[];
          count?: number;
          message?: string;
          limit?: number;
          offset?: number;
        }

        // Ajusta la ruta según tu API
        const url = `/anomalia-activo/${activoId}`;
<<<<<<< HEAD
        
=======
>>>>>>> origin/feature-front
        const response = await gs.authorizedGet(
          url, 
          user.token
        ) as AnomaliasPaginadas;

<<<<<<< HEAD

=======
>>>>>>> origin/feature-front
        const prediccionesTransformadas: Prediccion[] = (response.data ?? []).map((p) => ({
          id: p.id ?? 0,
          activoId: p.activo_id ?? activoId,
          timestamp: p.timestamp ?? new Date().toISOString(),
          anomalyScore: Math.round(p.anomaly_score ?? 0),
          anomalyLikelihood: Math.round(p.anomaly_likelihood ?? 0),
          severidad: p.severidad ?? "Baja",
          descripcion: p.descripcion ?? "Sin descripción",
          threshold: Math.round(p.threshold ?? 50),
          is_anomaly: p.is_anomaly === 1 || p.is_anomaly === true,
          most_influential_variable: p.most_influential_variable ?? "N/A",
          contribution_magnitude: p.contribution_magnitude ?? 0,
        }));


        // Ordenar por timestamp descendente
        prediccionesTransformadas.sort((a, b) => 
          new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
        );

        setPredicciones(prediccionesTransformadas);
        setError(null);
      } catch (err) {
        setError('Error al cargar anomalías del activo');
      } finally {
        setIsLoading(false);
      }
    };

    // Llamado inicial inmediato
    void fetchAnomalias();

    // Actualización periódica cada 5 minutos
    const interval = setInterval(() => {
      void fetchAnomalias();
    }, 300000);

    // Limpieza del intervalo al desmontar componente
    return () => { 
      clearInterval(interval); 
    };
  }, [activoId, isUserLoading, user]);

  if (isUserLoading || isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
        <Typography ml={2}>Cargando predicciones...</Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Box p={2}>
        <Typography color="error">{error}</Typography>
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
  const anomaliasDetectadas = predicciones.filter(p => p.is_anomaly);

  return (
    <Box>
      <Divider sx={{ my: 4 }} />
      
      <Box display="flex" alignItems="center" gap={2} mb={3} id="tour-resumen-predicciones-activo">
        <Typography variant="h4">
          Predicciones de Anomalías (ML)
        </Typography>
        <TourButton 
            onClick={startTour}
            tooltipTitle="Iniciar Tutorial de Predicciones"
            style="pulse 3s infinite"
          />
      </Box>

      <Grid container spacing={3}>
        {/* Métricas principales de la última predicción */}
        <Grid size={{ xs: 12, sm: 6, md: 3 }} id="tour-metricas-anomaly-score">
          <MetricCard
            title="Último Anomaly Score"
            value={ultimaPrediccion.anomalyScore}
            color={getColorBySeverity(ultimaPrediccion.severidad)}
            icon={<ActivityIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }} id='tour-metricas-anomaly-likelihood'>
          <MetricCard
            title="Último Anomaly Likelihood"
            value={ultimaPrediccion.anomalyLikelihood}
            color={getColorBySeverity(ultimaPrediccion.severidad)}
            icon={<AlertIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }} id='tour-metricas-total-anomalias'>
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
        <Grid size={{ xs: 12 }} id='tour-grafico-anomalias'>
          <AnomalyChart 
            predicciones={predicciones} 
            activoNombre={activoNombre} 
          />
        </Grid>

        {/* Tarjetas de predicciones detalladas */}
        <Grid size={{ xs: 12 }}>
          <Typography variant="h6" mb={2} id="tour-ultimas-anomalias">
            Últimas 10 Anomalías Detectadas
          </Typography>
          {anomaliasDetectadas.length > 0 ? (
            <PrediccionCards predicciones={anomaliasDetectadas.slice(0, 10)} />
          ) : (
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
          )}
        </Grid>
      </Grid>
    </Box>
  );
}