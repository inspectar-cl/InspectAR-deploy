// src/components/dashboard/overview/SensorChartCarousel.tsx

'use client';

import * as React from 'react';
import { useState, useMemo } from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import Box from '@mui/material/Box';
import Alert from '@mui/material/Alert';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from 'recharts';
import dayjs from 'dayjs';
import type { SensorRow } from '@/hooks/use-activos-with-sensors'; // Asegúrate que esta ruta es correcta
import { ArrowBackIosNew, ArrowForwardIos } from '@mui/icons-material';
import type { SxProps } from '@mui/material/styles';

// --- Props que recibirá el componente ---
interface SensorChartCarouselProps {
  sensors: SensorRow[];
  sx?: SxProps; // Para pasar estilos como height, etc.
}

// --- Tipo para los datos del gráfico ---
interface ChartData {
  time: string;
  value: number;
}

export function SensorChartCarousel({ sensors, sx }: SensorChartCarouselProps) {
  // --- Estado para el carrusel ---
  const [currentIndex, setCurrentIndex] = useState(0);

  // --- Estado para el rango de tiempo (igual que en el Dialog) ---
  const [rangoMinutos, setRangoMinutos] = useState(30);

  // --- Lógica de navegación ---
  const handlePrev = () => {
    setCurrentIndex((prev) => (prev - 1 + sensors.length) % sensors.length);
  };

  const handleNext = () => {
    setCurrentIndex((prev) => (prev + 1) % sensors.length);
  };

  // --- Sensor actual basado en el índice ---
  const currentSensor = sensors.length > 0 ? sensors[currentIndex] : null;

  // --- Memo para los datos (igual que en el Dialog) ---
  const filteredData: ChartData[] = useMemo(() => {
    if (!currentSensor?.history24h) return [];

    if (rangoMinutos === 0) {
      return currentSensor.history24h.map((d) => ({
        time: dayjs(d.ts).format('DD/MM HH:mm'),
        value: d.value,
      }));
    }

    const limite = dayjs().subtract(rangoMinutos, 'minute');
    return currentSensor.history24h
      .filter((d) => dayjs(d.ts).isAfter(limite))
      .map((d) => ({
        time: dayjs(d.ts).format('HH:mm:ss'),
        value: d.value,
      }));
  }, [currentSensor, rangoMinutos]); // Se recalcula si cambia el sensor o el rango

  // --- Renderizado ---
  // Caso 1: No hay sensores que mostrar
  if (!currentSensor || sensors.length === 0) {
    return (
      <Card sx={{ ...sx, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Alert severity="info" sx={{ width: '80%' }}>
          No hay sensores con datos históricos para mostrar en este activo.
        </Alert>
      </Card>
    );
  }

  // Caso 2: Hay sensores, mostramos el carrusel
  return (
    <Card sx={{ ...sx, display: 'flex', flexDirection: 'column' }}>
      {/* --- Cabecera con Navegación y Selector de Rango --- */}
      <CardHeader
        // Título con flechas de navegación
        title={
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <IconButton onClick={handlePrev} aria-label="Sensor anterior">
              <ArrowBackIosNew />
            </IconButton>
            <Typography variant="h6" component="div" sx={{ textAlign: 'center', flexGrow: 1 }}>
              {`Histórico: ${currentSensor.name}`}
            </Typography>
            <IconButton onClick={handleNext} aria-label="Siguiente sensor">
              <ArrowForwardIos />
            </IconButton>
          </Box>
        }
        // Selector de rango a la derecha
        action={
          <FormControl size="small" sx={{ minWidth: 200, mr: 1 }}>
            <InputLabel id="rango-label-carousel">Rango</InputLabel>
            <Select
              labelId="rango-label-carousel"
              value={rangoMinutos}
              label="Rango"
              onChange={(e) => {
                setRangoMinutos(Number(e.target.value));
              }}
            >
              <MenuItem value={0}>Todo</MenuItem>
              <MenuItem value={3}>3 min</MenuItem>
              <MenuItem value={30}>30 min</MenuItem>
              <MenuItem value={60}>1 hora</MenuItem>
              <MenuItem value={240}>4 horas</MenuItem>
              <MenuItem value={720}>12 horas</MenuItem>
            </Select>
          </FormControl>
        }
        sx={{ pb: 0 }}
      />
      
      {/* --- Contenido del Gráfico --- */}
      <CardContent sx={{ flexGrow: 1, height: 'calc(100% - 72px)' }}>
        {filteredData.length > 0 ? (
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={filteredData}
              margin={{ top: 5, right: 30, left: 10, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis
                label={{
                  value: currentSensor?.unit ?? 'Unidad',
                  angle: -90,
                  position: 'insideLeft',
                  offset: -5,
                }}
              />
              <Tooltip />
              <Line
                type="monotone"
                dataKey="value"
                stroke="#1976d2"
                dot={false}
                strokeWidth={2}
                isAnimationActive={false}
              />
            </LineChart>
          </ResponsiveContainer>
        ) : (
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%' }}>
            <Alert severity="info">
              No hay datos en el rango seleccionado para este sensor.
            </Alert>
          </Box>
        )}
      </CardContent>
    </Card>
  );
}