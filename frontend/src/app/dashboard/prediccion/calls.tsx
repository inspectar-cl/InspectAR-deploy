/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client'

import { useState, useEffect, useMemo } from "react";
import {
  Box,
  Typography,
  TextField,
  MenuItem,
  Select,
  InputLabel,
  FormControl,
  Grid,
  CircularProgress
} from "@mui/material";
import { Activity as ActivityIcon, AlertTriangle as AlertIcon } from "lucide-react";
import { type Prediccion } from "@/types/prediccion";
import { getColorBySeverity } from "@/utils/SeverityUtils";
import MetricCard from "@/components/dashboard/prediccion/MetricCard";
import PrediccionCards from "@/components/dashboard/prediccion/PrediccionCards";
import AlertasResumen from "@/components/dashboard/prediccion/AlertasResumen";
import AnomalyChart from "@/components/dashboard/prediccion/AnomalyChart";

import { useDriverTour } from "@/components/tutorial/use-driver-tour";
import type { TourKey } from "@/components/tutorial/tour-config";
import { TourButton } from "@/components/tutorial/tour-button";
import { useUserToken } from '@/hooks/use-usertoken';
import Services from '@/modules/Services';

// Tipos
interface Activo {
  id: number;
  nombre: string;
  tipo: string;
  ubicacion: string;
}

const gs = new Services();

export default function PrediccionClient() {
  const { user, isLoading: isUserLoading } = useUserToken();
  const [activos, setActivos] = useState<Activo[]>([]);
  const [predicciones, setPredicciones] = useState<Prediccion[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Configurar tutorial/tour
  const tourKey: TourKey = 'predicciones-filtros';
  const { startTour } = useDriverTour(tourKey);

  // Filtros
  const [filtroActivo, setFiltroActivo] = useState<string>("Todos");
  const [filtroSeveridad, setFiltroSeveridad] = useState<string>("Todos");
  const [filtroKeyword, setFiltroKeyword] = useState<string>("");

  // Fetch de activos y predicciones
  useEffect(() => {
    if (isUserLoading || !user) return;

    const fetchActivos = async () => {
      try {
        interface ActivoResponse {
          activo_id?: number;
          nombre?: string;
          tipo?: string;
          ubicacion?: string;
        }

        const response = await gs.authorizedGet('/obtener-todos-activos', user.token) as { activos?: ActivoResponse[] };
        
        const activosTransformados: Activo[] = (response.activos ?? []).map((a) => ({
          id: a.activo_id ?? 0,
          nombre: a.nombre ?? 'Sin nombre',
          tipo: a.tipo ?? 'Sin tipo',
          ubicacion: a.ubicacion ?? 'Sin ubicación',
        }));

        setActivos(activosTransformados);
      } catch (err) {
        console.error('Error al obtener activos:', err);
        setError('Error al cargar activos');
      }
    };

    const fetchPredicciones = async () => {
      try {
        setIsLoading(true);
        
        interface PrediccionResponse {
          id?: number;
          activo_id?: number;
          timestamp?: string;
          anomaly_score?: number;
          anomaly_likelihood?: number;
          severidad?: "Baja" | "Media" | "Alta";
          descripcion?: string;
          threshold?: number;
          is_anomaly?: boolean;
          most_influential_variable?: string;
          contribution_magnitude?: number;
        }

        const response = await gs.authorizedGet('/anomalia-activo/2', user.token) as { predicciones?: PrediccionResponse[] };
        
        const prediccionesTransformadas: Prediccion[] = (response.predicciones ?? []).map((p) => ({
          id: p.id ?? 0,
          activoId: p.activo_id ?? 0,
          timestamp: p.timestamp ?? new Date().toISOString(),
          anomalyScore: p.anomaly_score ?? 0,
          anomalyLikelihood: p.anomaly_likelihood ?? 0,
          severidad: p.severidad ?? "Baja",
          descripcion: p.descripcion ?? "Sin descripción",
          threshold: p.threshold ?? 50,
          is_anomaly: p.is_anomaly ?? false,
        }));

        // Ordenar por timestamp descendente
        prediccionesTransformadas.sort((a, b) => 
          new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
        );

        setPredicciones(prediccionesTransformadas);
        setError(null);
      } catch (err) {
        console.error('Error al obtener predicciones:', err);
        setError('Error al cargar predicciones');
      } finally {
        setIsLoading(false);
      }
    };

    const fetchAll = async () => {
      await fetchActivos();
      await fetchPredicciones();
    };

    // Llamado inicial inmediato
    void fetchAll();

    // Intervalo de actualización cada 30 segundos
    const interval = setInterval(() => {
      void fetchAll();
    }, 30000);

    // Limpieza del intervalo al desmontar componente
    return () => { clearInterval(interval); };
  }, [isUserLoading, user]);

  const prediccionesFiltradas = useMemo(() => {
    return predicciones.filter((p) => {
      const activoMatch =
        filtroActivo === "Todos" || activos.find((a) => a.id === p.activoId)?.nombre === filtroActivo;
      const severidadMatch = filtroSeveridad === "Todos" || p.severidad === filtroSeveridad;
      const keywordMatch = p.descripcion.toLowerCase().includes(filtroKeyword.toLowerCase());
      return activoMatch && severidadMatch && keywordMatch;
    });
  }, [predicciones, filtroActivo, filtroSeveridad, filtroKeyword, activos]);

  if (isUserLoading || isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
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

  return (
    <Box p={2}>
      <Box 
          display="flex" 
          justifyContent="space-between" 
          alignItems="center" 
          mb={3}
          id="tour-header" 
      >
          <Typography variant="h4">
              Predicciones de Anomalías
          </Typography>

          <TourButton 
            onClick={startTour}
            tooltipTitle="Iniciar Tutorial de Filtros"
            style="pulse 3s infinite"
          />
      </Box>

      {/* Filtros */}
      <Grid container spacing={2} mb={2}>
        <Grid size={{ xs: 12, sm: 4 }}>
          <FormControl fullWidth id="tour-filtro-activo">
            <InputLabel>Activo</InputLabel>
            <Select
              value={filtroActivo}
              label="Activo"
              onChange={(e) => { setFiltroActivo(e.target.value); }}
            >
              <MenuItem value="Todos">Todos ({activos.length} activos)</MenuItem>
              {activos.map((a, index) => (
                <MenuItem key={`activo-${a.id}-${index}`} value={a.nombre}>
                  {a.nombre}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid size={{ xs: 12, sm: 4 }}>
          <FormControl fullWidth id="tour-filtro-severidad">
            <InputLabel>Severidad</InputLabel>
            <Select
              value={filtroSeveridad}
              label="Severidad"
              onChange={(e) => { setFiltroSeveridad(e.target.value); }}
            >
              <MenuItem value="Todos">Todos</MenuItem>
              <MenuItem value="Baja">Baja</MenuItem>
              <MenuItem value="Media">Media</MenuItem>
              <MenuItem value="Alta">Alta</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        <Grid size={{ xs: 12, sm: 4 }}>
          <TextField
            fullWidth
            label="Buscar por palabra clave"
            value={filtroKeyword}
            onChange={(e) => { setFiltroKeyword(e.target.value); }}
          />
        </Grid>
      </Grid>

      {/* Resumen de alertas */}
      <AlertasResumen predicciones={prediccionesFiltradas} />

      {/* Dashboard por activo */}
      {activos.map((activo) => {
        const predActivo = prediccionesFiltradas.filter((p) => p.activoId === activo.id);
        if (predActivo.length === 0) return null;

        const ultimaPred = predActivo[0];
        const color = getColorBySeverity(ultimaPred.severidad);
        const totalAnomalias = predActivo.filter(p => p.is_anomaly).length;

        return (
          <Box key={`activo-${activo.id}`} mb={4} p={2} sx={{ borderRadius: 2 }}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
              <Typography variant="h6">
                {activo.nombre}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {predActivo.length} registros | {totalAnomalias} anomalías detectadas
              </Typography>
            </Box>

            {/* Métricas principales */}
            <Grid container spacing={3} mb={3} id="tour-metricas-principales">
              <Grid key={`score-${activo.id}`} size={{ xs: 12, sm: 6, md: 4 }} id="tour-anomaly-score">
                <MetricCard
                  title="Último Anomaly Score"
                  value={ultimaPred.anomalyScore}
                  color={color}
                  icon={<ActivityIcon />}
                  suffix=""
                />
              </Grid>
              <Grid key={`likelihood-${activo.id}`} size={{ xs: 12, sm: 6, md: 4 }} id="tour-anomaly-likelihood">
                <MetricCard
                  title="Último Anomaly Likelihood"
                  value={ultimaPred.anomalyLikelihood}
                  color={color}
                  icon={<AlertIcon />}
                  suffix=""
                />
              </Grid>
              <Grid key={`total-${activo.id}`} size={{ xs: 12, sm: 6, md: 4 }} id="tour-anomaly">
                <MetricCard
                  title="Total Anomalías"
                  value={totalAnomalias}
                  color="#2196f3"
                  icon={<AlertIcon />}
                  suffix=""
                />
              </Grid>
            </Grid>

            {/* Gráfico */}
            <Grid container spacing={2} id="tour-grafico-predicciones">
              <Grid key={`chart-${activo.id}`} size={{ xs: 12 }}>
                <AnomalyChart predicciones={predActivo} activoNombre={activo.nombre} />
              </Grid>

              {/* Tarjetas de predicciones */}
              <Grid key={`cards-${activo.id}`} size={{ xs: 12 }} id="tour-ultimas-anomalias">
                <Typography variant="h6" mb={2}>
                  Últimas 10 Anomalías
                </Typography>
                <PrediccionCards predicciones={predActivo.slice(0, 10)} />
              </Grid>
            </Grid>
          </Box>
        );
      })}
    </Box>
  );
}