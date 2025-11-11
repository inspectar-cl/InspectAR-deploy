// src/components/dashboard/overview/SensorChartDialog.tsx

'use client';

import * as React from 'react';
import { useState, useMemo } from 'react';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import Alert from '@mui/material/Alert';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
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
import type { SensorRow } from '@/hooks/use-activos-with-sensors';

interface SensorChartDialogProps {
  open: boolean;
  onClose: () => void;
  sensor: SensorRow | null;
}

// Definimos el tipo para los datos filtrados
interface ChartData {
  time: string;
  value: number;
}

export function SensorChartDialog({
  open,
  onClose,
  sensor,
}: SensorChartDialogProps) {

  const [rangoMinutos, setRangoMinutos] = useState(30);

  const filteredData: ChartData[] = useMemo(() => {
    if (!sensor?.history24h) return [];

    if (rangoMinutos === 0) {
      return sensor.history24h.map((d) => ({
        time: dayjs(d.ts).format('DD/MM HH:mm'),
        value: d.value,
      }));
    }

    const limite = dayjs().subtract(rangoMinutos, 'minute');
    return sensor.history24h
      .filter((d) => dayjs(d.ts).isAfter(limite))
      .map((d) => ({
        time: dayjs(d.ts).format('HH:mm:ss'),
        value: d.value,
      }));
  }, [sensor, rangoMinutos]);

  // Reseteamos el rango de tiempo a 30 min cada vez que se cierra el diálogo
  // para que la próxima vez que se abra, tenga el valor por defecto.
  const handleClose = () => {
    onClose();
    // Opcional: resetear el estado al cerrar
    // setRangoMinutos(30); 
  };

  return (
    <Dialog
      fullWidth
      maxWidth="md"
      open={open}
      onClose={handleClose}
    >
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        {sensor ? `Histórico de ${sensor.name}` : 'Cargando...'}
        <FormControl size="small" sx={{ minWidth: 200 }}>
          <InputLabel id="rango-label">Rango de tiempo</InputLabel>
          <Select
            labelId="rango-label"
            value={rangoMinutos}
            label="Rango de tiempo"
            onChange={(e) => {
              setRangoMinutos(Number(e.target.value));
            }}
          >
            <MenuItem value={0}>Todo</MenuItem>
            <MenuItem value={3}>Últimos 3 minutos</MenuItem>
            <MenuItem value={30}>Últimos 30 minutos</MenuItem>
            <MenuItem value={60}>Última hora</MenuItem>
            <MenuItem value={240}>Últimas 4 horas</MenuItem>
            <MenuItem value={720}>Últimas 12 horas</MenuItem>
          </Select>
        </FormControl>
      </DialogTitle>
      <DialogContent dividers>
        {filteredData.length > 0 ? (
          <Box sx={{ width: '100%', height: 400 }}>
            <ResponsiveContainer>
              <LineChart
                data={filteredData}
                margin={{ top: 20, right: 30, left: 10, bottom: 20 }}
              >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time" />
                <YAxis
                  label={{
                    value: sensor?.unit ?? 'Unidad',
                    angle: -90,
                    position: 'insideLeft',
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
          </Box>
        ) : (
          <Alert severity="info">No hay datos en el rango seleccionado.</Alert>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cerrar</Button>
      </DialogActions>
    </Dialog>
  );
}