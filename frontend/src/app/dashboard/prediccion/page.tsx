/* eslint-disable @typescript-eslint/explicit-function-return-type -- tipo inferido */
"use client";

import { useState, useMemo, useEffect } from "react";
import {
  Box,
  Typography,
  TextField,
  MenuItem,
  Select,
  InputLabel,
  FormControl,
  Grid,
  CircularProgress,
} from "@mui/material";
import { Activity as ActivityIcon, AlertTriangle as AlertIcon } from "lucide-react";
import { type Prediccion } from "@/types/prediccion";
import { getColorBySeverity } from "@/utils/SeverityUtils";
import MetricCard from "@/components/dashboard/prediccion/MetricCard";
import PrediccionCards from "@/components/dashboard/prediccion/PrediccionCards";
import AlertasResumen from "@/components/dashboard/prediccion/AlertasResumen";
import AnomalyChart from "@/components/dashboard/prediccion/AnomalyChart";
import Services from "@/modules/Services";
import { useUserToken } from "@/hooks/use-usertoken";

const gs = new Services();

// Tipos
interface Activo {
  id: number;
  nombre: string;
  tipo: string;
  ubicacion: string;
}

export default function Page() {
  const { user } = useUserToken();
  const [activos, setActivos] = useState<Activo[]>([]);
  const [predicciones, setPredicciones] = useState<Prediccion[]>([]);
  const [loading, setLoading] = useState(true);

  // 🔹 Cargar datos desde backend
  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!user?.token) return;

        // 1️⃣ Obtener activos
        const activosData = await gs.authorizedGet("/obtener-todos-activos", user.token) as { data?: Activo[] };
        const activosList = activosData?.data || [];
        setActivos(activosList);

        // 2️⃣ Obtener predicciones (en paralelo por activo)
        const results = await Promise.all(
          activosList.map(async (activo: Activo) => {
            try {
              const res = await gs.get(`/anomalies/activo/${activo.id}`, user.token) as { data?: Prediccion[] };
              return Array.isArray(res.data)
                ? res.data.map((p: Prediccion) => ({ ...p, activoId: activo.id }))
                : [];
            } catch {
              return [];
            }
          })
        );

        const allPreds = results.flat();
        setPredicciones(allPreds);
      } catch (error) { /* empty */ } finally {
        setLoading(false);
      }
    };

    void fetchData();
  }, [user]);

  // 🔹 Filtros
  const [filtroActivo, setFiltroActivo] = useState<string>("Todos");
  const [filtroSeveridad, setFiltroSeveridad] = useState<string>("Todos");
  const [filtroKeyword, setFiltroKeyword] = useState<string>("");

  const prediccionesFiltradas = useMemo(() => {
    return predicciones.filter((p) => {
      const activoMatch =
        filtroActivo === "Todos" || activos.find((a) => a.id === p.activoId)?.nombre === filtroActivo;
      const severidadMatch = filtroSeveridad === "Todos" || p.severidad === filtroSeveridad;
      const keywordMatch = p.descripcion?.toLowerCase().includes(filtroKeyword.toLowerCase());
      return activoMatch && severidadMatch && keywordMatch;
    });
  }, [predicciones, filtroActivo, filtroSeveridad, filtroKeyword, activos]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" height="50vh">
        <CircularProgress />
        <Typography ml={2}>Cargando datos...</Typography>
      </Box>
    );
  }

  return (
    <Box p={2}>
      {/* 🔹 Filtros superiores */}
      <Grid container spacing={2} mb={2}>
        <Grid size={{ xs: 12, sm: 4 }}>
          <FormControl fullWidth>
            <InputLabel>Activo</InputLabel>
            <Select
              value={filtroActivo}
              label="Activo"
              onChange={(e) => { setFiltroActivo(e.target.value); }}
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

        <Grid size={{ xs: 12, sm: 4 }}>
          <FormControl fullWidth>
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

      {/* 🔹 Resumen general de alertas */}
      <AlertasResumen predicciones={prediccionesFiltradas} />

      {/* 🔹 Dashboard por activo */}
      {activos.map((activo) => {
        // Todas las predicciones del activo (sin depender de filtro global)
        const predActivo = predicciones.filter((p) => p.activoId === activo.id);
        const ultimaPred = predActivo[predActivo.length - 1];

        return (
          <Box
            key={activo.id}
            mb={5}
            p={2}
            sx={{
              borderBottom: "1px solid #ddd",
              backgroundColor: "#fafafa",
              borderRadius: 2,
            }}
          >
            <Typography variant="h6" mb={2}>
              {activo.nombre}
            </Typography>

            {predActivo.length > 0 ? (
              <>
                {/* 🔹 Métricas y alertas lado a lado */}
                <Grid container spacing={3} mb={3}>
                  <Grid size={{ xs: 12, sm: 4 }}>
                    <MetricCard
                      title="Anomaly Score"
                      value={ultimaPred.anomalyScore}
                      color={getColorBySeverity(ultimaPred.severidad)}
                      icon={<ActivityIcon />}
                    />
                  </Grid>
                  <Grid size={{ xs: 12, sm: 4 }}>
                    <MetricCard
                      title="Anomaly Likelihood"
                      value={ultimaPred.anomalyLikelihood}
                      color={getColorBySeverity(ultimaPred.severidad)}
                      icon={<AlertIcon />}
                    />
                  </Grid>
                  <Grid size={{ xs: 12, sm: 12 }}>
                    <AlertasResumen predicciones={predActivo} />
                  </Grid>
                </Grid>

                {/* 🔹 Gráfico + tarjetas detalladas */}
                <AnomalyChart predicciones={predActivo} activoNombre={activo.nombre} />
                <Box mt={2}>
                  <PrediccionCards predicciones={predActivo} />
                </Box>
              </>
            ) : (
              <Typography variant="body2" color="text.secondary">
                No hay datos de predicción para este activo.
              </Typography>
            )}
          </Box>
        );
      })}
    </Box>
  );
}
