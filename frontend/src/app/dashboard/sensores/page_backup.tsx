'use client';
import * as React from 'react';
import Grid from '@mui/material/Grid';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Alert from '@mui/material/Alert';
import RefreshIcon from '@mui/icons-material/Refresh';
import { LatestSensors } from '@/components/dashboard/overview/latest-sensors';
import { useActivosWithSensors } from '@/hooks/useActivosWithSensors';
import dayjs from 'dayjs';

export default function Page(): React.JSX.Element {

  // Ejemplo de dos activos críticos distintos con sus sensores
  const [activos, setActivos] = React.useState<Activo[]>([
    {
      assetName: "Bomba de Agua — Edificio A",
      imageUrl: "https://www.iprecom.com/wp-content/uploads/2019/07/bomba-de-agua.jpg", 
      sensores: [
        {
          id: 'SEN-C-001',
          name: 'Caudalímetro #1',
          type: 'caudal',
          unit: 'L/s',
          lastValue: 42.3,
          lastSeen: dayjs(now).subtract(3, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 42, 5, 'L/s'),
        },
        {
          id: 'SEN-P-017',
          name: 'Transmisor de Presión #2',
          type: 'presión',
          unit: 'bar',
          lastValue: 2.1,
          lastSeen: dayjs(now).subtract(4, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 2.0, 0.2, 'bar'),
        },
        {
          id: 'SEN-P-018',
          name: 'Transmisor de Presión #3',
          type: 'presión',
          unit: 'bar',
          lastValue: 2.1,
          lastSeen: dayjs(now).subtract(2, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 2.0, 0.2, 'bar'),
        },
        {
          id: 'SEN-P-019',
          name: 'Transmisor de Presión #4',
          type: 'presión',
          unit: 'bar',
          lastValue: 2.1,
          lastSeen: dayjs(now).subtract(4, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 2.0, 0.2, 'bar'),
        },
        {
          id: 'SEN-P-020',
          name: 'Transmisor de Presión #5',
          type: 'presión',
          unit: 'bar',
          lastValue: 2.1,
          lastSeen: dayjs(now).subtract(3, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 2.0, 0.2, 'bar'),
        },
        {
          id: 'SEN-P-021',
          name: 'Transmisor de Presión #6',
          type: 'presión',
          unit: 'bar',
          lastValue: 2.1,
          lastSeen: dayjs(now).subtract(1, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 2.0, 0.2, 'bar'),
        },
      ],
    },
    {
      assetName: "Panel Eléctrico — Subestación B",
      imageUrl: "https://www.eabel.com/wp-content/uploads/2024/04/Electrical-Control-Panel-0401.webp", 
      sensores: [
        {
          id: 'SEN-T-003',
          name: 'Sensor de Temperatura',
          type: 'temperatura',
          unit: '°C',
          lastValue: 85.2,
          lastSeen: dayjs(now).subtract(1, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 85, 3, '°C'),
        },
        {
          id: 'SEN-V-009',
          name: 'Sensor de Vibración',
          type: 'vibración',
          unit: 'mm/s',
          lastValue: 3.4,
          lastSeen: dayjs(now).subtract(8, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 3.2, 0.6, 'mm/s'),
        },
      ],
    },
    {
      assetName: "Ascensor Torre C",
      imageUrl: "https://elevabalear.com/wp-content/uploads/2024/11/caracteristicas-de-los-ascensores-electricos.jpg",
      sensores: [
        {
          id: 'SEN-M-021',
          name: 'Motor Eléctrico RPM',
          type: 'velocidad',
          unit: 'rpm',
          lastValue: 1200,
          lastSeen: dayjs(now).subtract(2, 'minute').toDate(),
          history24h: mockHistory24h(now, 12, 1200, 50, 'rpm'),
        },
      ],
    },
  ]);

  // Helper para setear lastSeen en un sensor del activo i
  function updateSensorLastSeen(actIdx: number, sensorId: string, minutesAgo: number) {
    setActivos(prev => {
      const clone = structuredClone(prev) as Activo[];
      const sensores = clone[actIdx].sensores;
      const s = sensores.find(x => x.id === sensorId);
      if (s) s.lastSeen = dayjs().subtract(minutesAgo, 'minute').toDate();
      return clone;
    });
  }

  // Simula caída: pone lastSeen a 8 minutos (=> inactivo)
  function simulateDown(actIdx: number, sensorId: string) {
    updateSensorLastSeen(actIdx, sensorId, 8);
  }

  // Simula recovery: pone lastSeen a 1 minuto (=> activo)
  function simulateUp(actIdx: number, sensorId: string) {
    updateSensorLastSeen(actIdx, sensorId, 1);
  }

  return (
    <Grid
      container
      spacing={3}
      justifyContent="center"
      alignItems="flex-start"
    >
      {activos.map((activo, idx) => (
        <Grid key={idx} size={{ lg: 8, md: 12, xs: 12 }}>
          <Stack direction="row" spacing={1} sx={{ mb: 1, justifyContent: 'flex-end' }}>
            <Button size="small" variant="outlined" onClick={() => simulateDown(idx, activo.sensores[0].id)}>
              Simular caída ({activo.sensores[0].name})
            </Button>
            <Button size="small" variant="outlined" onClick={() => simulateUp(idx, activo.sensores[0].id)}>
              Simular recovery
            </Button>
          </Stack>

          <LatestSensors
            assetName={activo.assetName}
            imageUrl={activo.imageUrl}
            sensors={activo.sensores}
            sx={{ height: '100%' }}
          />
        </Grid>
      ))}
    </Grid>
  );
}
