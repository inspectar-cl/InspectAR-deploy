"use client";

import React, { useState, useMemo } from "react";
import {
  Box,
  Chip,
  Typography,
  Grid,
  Card,
  CardContent,
  TextField,
  MenuItem,
  Select,
  InputLabel,
  FormControl,
  Avatar,
  LinearProgress,
} from "@mui/material";
import dynamic from "next/dynamic";
import { Activity as ActivityIcon, AlertTriangle as AlertIcon } from "lucide-react";
import { Prediccion } from "@/types/prediccion";
import { getColorBySeverity } from "@/utils/severityUtils";
import MetricCard from "@/components/dashboard/prediccion/MetricCard";
import PrediccionCards from "@/components/dashboard/prediccion/PrediccionCards";
import AlertasResumen from "@/components/dashboard/prediccion/AlertasResumen";
import AnomalyChart from "@/components/dashboard/prediccion/AnomalyChart";

// Dinámico porque ApexCharts no funciona en SSR
const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

// Tipos
type Activo = {
  id: number;
  nombre: string;
  tipo: string;
  ubicacion: string;
};


// Página principal
export default function Page() {
  // Datos de ejemplo
  const activos: Activo[] = [
    { id: 1, nombre: "Bomba de agua", tipo: "Bomba", ubicacion: "Sótano" },
    { id: 2, nombre: "Caldera", tipo: "Caldera", ubicacion: "Cuarto de máquinas" },
  ];

  const predicciones: Prediccion[] = [
    {
      id: 1,
      activoId: 1,
      timestamp: "2025-10-14T10:00:00",
      anomalyScore: 85,
      anomalyLikelihood: 90,
      severidad: "Alta",
      descripcion: "Vibración fuera de rango",
      threshold: 70,
      is_anomaly: true,
    },
    {
      id: 1,
      activoId: 1,
      timestamp: "2025-10-14T10:01:00",
      anomalyScore: 80,
      anomalyLikelihood: 75,
      severidad: "Alta",
      descripcion: "Vibración fuera de rango",
      threshold: 70,
      is_anomaly: true,
    },
    {
      id: 2,
      activoId: 2,
      timestamp: "2025-10-14T11:00:00",
      anomalyScore: 40,
      anomalyLikelihood: 30,
      severidad: "Media",
      descripcion: "Temperatura alta intermitente",
      threshold: 50,
      is_anomaly: false,
    },
  ];

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
      {/* Filtros */}
      <Grid container spacing={2} mb={2}>
        <Grid item xs={12} sm={4}>
          <FormControl fullWidth>
            <InputLabel>Activo</InputLabel>
            <Select
              value={filtroActivo}
              label="Activo"
              onChange={(e) => setFiltroActivo(e.target.value)}
            >
              <MenuItem value="Todos">Todos</MenuItem>
              {activos.map((a) => (
                <MenuItem key={a.id} value={a.nombre}>
                  {a.nombre}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sm={4}>
          <FormControl fullWidth>
            <InputLabel>Severidad</InputLabel>
            <Select
              value={filtroSeveridad}
              label="Severidad"
              onChange={(e) => setFiltroSeveridad(e.target.value)}
            >
              <MenuItem value="Todos">Todos</MenuItem>
              <MenuItem value="Baja">Baja</MenuItem>
              <MenuItem value="Media">Media</MenuItem>
              <MenuItem value="Alta">Alta</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sm={4}>
          <TextField
            fullWidth
            label="Buscar por palabra clave"
            value={filtroKeyword}
            onChange={(e) => setFiltroKeyword(e.target.value)}
          />
        </Grid>
      </Grid>

      {/* Resumen de alertas */}
      <AlertasResumen predicciones={prediccionesFiltradas} />

      {/* Dashboard por activo */}
      {activos.map((activo) => {
        const predActivo = prediccionesFiltradas.filter((p) => p.activoId === activo.id);
        if (predActivo.length === 0) return null;

        const ultimaPred = predActivo[predActivo.length - 1];
        const color = getColorBySeverity(ultimaPred.severidad);

        return (
          <Box key={activo.id} mb={4}>
            <Typography variant="h6" mb={1}>
              {activo.nombre}
            </Typography>

            {/* Métricas principales */}
            <Grid container spacing={3} mb={3}>
              <Grid item>
                <MetricCard
                  title="Último Anomaly Score"
                  value={ultimaPred.anomalyScore}
                  color={color}
                  icon={<ActivityIcon />}
                  suffix=""
                />
              </Grid>
              <Grid item>
                <MetricCard
                  title="Último Anomaly Likelihood"
                  value={ultimaPred.anomalyLikelihood}
                  color={color}
                  icon={<AlertIcon />}
                  suffix=""
                />
              </Grid>
              <Grid item>
                  <PrediccionCards predicciones={predActivo} />
              </Grid>
            </Grid>

            {/* Gráfico y cards */}
            <AnomalyChart predicciones={predActivo} activoNombre={activo.nombre} />
            
          </Box>
        );
      })}
    </Box>
  );
}
