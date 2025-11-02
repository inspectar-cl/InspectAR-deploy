/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
"use client";

import { useState, useMemo } from "react";
import {
  Box,
  Typography,
  TextField,
  MenuItem,
  Select,
  InputLabel,
  FormControl,
  Grid
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

// Tipos
interface Activo {
  id: number;
  nombre: string;
  tipo: string;
  ubicacion: string;
}

// Función para calcular severidad desde likelihood (igual que el API)
const severityFromLikelihood = (likelihood: number): "Baja" | "Media" | "Alta" => {
  if (likelihood < 40) {
    return "Baja";
  } else if (likelihood < 70) {
    return "Media";
  }
  return "Alta";
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

// Función para generar datos mock realistas
const generateMockPredicciones = (): Prediccion[] => {
  const activos = [1, 2]; // IDs de activos
  const predicciones: Prediccion[] = [];
  
  let idCounter = 1;
  const now = new Date();

  // Generar datos para cada activo
  activos.forEach((activoId) => {
    // Generar entre 15-25 registros por activo
    const numRegistros = 15 + Math.floor(Math.random() * 10);
    
    for (let i = 0; i < numRegistros; i++) {
      // Timestamps en minutos hacia atrás
      const timestamp = new Date(now.getTime() - i * 60000);
      
      // Generar anomaly_score primero (0-100)
      const anomalyScore = Math.floor(Math.random() * 100);
      
      // anomaly_likelihood debe ser >= anomaly_score (como en el modelo real)
      const anomalyLikelihood = anomalyScore + Math.floor(Math.random() * (100 - anomalyScore + 1));
      
      // Calcular severidad basándose en likelihood (lógica del API)
      const severidad = severityFromLikelihood(anomalyLikelihood);
      
      // Threshold aleatorio entre 50-70
      const threshold = 50 + Math.floor(Math.random() * 20);
      
      predicciones.push({
        id: idCounter++, // ID único e incremental
        activoId,
        timestamp: timestamp.toISOString(),
        anomalyScore,
        anomalyLikelihood,
        severidad,
        descripcion: getDescripcionPorSeveridad(severidad), // Descripción según severidad
        threshold,
        is_anomaly: anomalyLikelihood > threshold, // Compara likelihood con threshold
      });
    }
  });

  // Ordenar por timestamp descendente (más recientes primero)
  predicciones.sort((a, b) => 
    new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  );

  return predicciones;
};

// Página principal
export default function Page() {
  // Activos mock
  const activos: Activo[] = useMemo(() => [
    { id: 1, nombre: "Bomba de Agua Principal", tipo: "Bomba", ubicacion: "Sótano B1" },
    { id: 2, nombre: "Caldera Industrial", tipo: "Caldera", ubicacion: "Cuarto de Máquinas" },
  ], []);

  // Predicciones mock generadas dinámicamente
  const predicciones: Prediccion[] = useMemo(() => generateMockPredicciones(), []);

  // Configurar tutorial/tour
  const tourKey: TourKey = 'predicciones-filtros';
  const { startTour } = useDriverTour(tourKey);

  // Filtros
  const [filtroActivo, setFiltroActivo] = useState<string>("Todos");
  const [filtroSeveridad, setFiltroSeveridad] = useState<string>("Todos");
  const [filtroKeyword, setFiltroKeyword] = useState<string>("");

  const prediccionesFiltradas = useMemo(() => {
    return predicciones.filter((p) => {
      const activoMatch =
        filtroActivo === "Todos" || activos.find((a) => a.id === p.activoId)?.nombre === filtroActivo;
      const severidadMatch = filtroSeveridad === "Todos" || p.severidad === filtroSeveridad;
      const keywordMatch = p.descripcion.toLowerCase().includes(filtroKeyword.toLowerCase());
      return activoMatch && severidadMatch && keywordMatch;
    });
  }, [predicciones, filtroActivo, filtroSeveridad, filtroKeyword, activos]);

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
              {activos.map((a) => (
                <MenuItem key={a.id} value={a.nombre}>
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

        const ultimaPred = predActivo[0]; // Ya ordenado por timestamp desc
        const color = getColorBySeverity(ultimaPred.severidad);
        const totalAnomalias = predActivo.filter(p => p.is_anomaly).length;

        return (
          <Box key={activo.id} mb={4} p={2} sx={{ borderRadius: 2 }}>
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
              <Grid size={{ xs: 12, sm: 6, md: 4 }} id="tour-anomaly-score">
                <MetricCard
                  title="Último Anomaly Score"
                  value={ultimaPred.anomalyScore}
                  color={color}
                  icon={<ActivityIcon />}
                  suffix=""
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6, md: 4 }} id="tour-anomaly-likelihood">
                <MetricCard
                  title="Último Anomaly Likelihood"
                  value={ultimaPred.anomalyLikelihood}
                  color={color}
                  icon={<AlertIcon />}
                  suffix=""
                />
              </Grid>
              <Grid size={{ xs: 12, sm: 6, md: 4 }} id="tour-anomaly">
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
              <Grid size={{ xs: 12 }}>
                <AnomalyChart predicciones={predActivo} activoNombre={activo.nombre} />
              </Grid>

              {/* Tarjetas de predicciones */}
              <Grid size={{ xs: 12 }} id="tour-ultimas-anomalias">
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