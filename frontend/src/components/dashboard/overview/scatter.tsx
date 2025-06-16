'use client';

import * as React from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import Divider from '@mui/material/Divider';
import { useTheme } from '@mui/material/styles';
import { ScatterChart } from '@mui/x-charts/ScatterChart';
import { useScatterSeries, useXScale, useYScale } from '@mui/x-charts/hooks';

export function Scatter({ sx }: { sx?: any }) {
  const theme = useTheme();

  // Simula una base de timestamps reales para ejemplo
  const now = new Date().getTime();

  const data1 = [
    { x: now - 5 * 60 * 1000, y: 300, id: 1 },
    { x: now - 4 * 60 * 1000, y: 450, id: 2 },
    { x: now - 3 * 60 * 1000, y: 500, id: 3 },
    { x: now - 2 * 60 * 1000, y: 350, id: 4 },
    { x: now - 1 * 60 * 1000, y: 280, id: 5 },
  ];

  const data2 = [
    { x: now - 5 * 60 * 1000, y: 100, id: 1 },
    { x: now - 4 * 60 * 1000, y: 250, id: 2 },
    { x: now - 3 * 60 * 1000, y: 300, id: 3 },
    { x: now - 2 * 60 * 1000, y: 200, id: 4 },
    { x: now - 1 * 60 * 1000, y: 180, id: 5 },
  ];

  const data3 = [
    { x: now - 5 * 60 * 1000, y: 60, id: 1 },
    { x: now - 4 * 60 * 1000, y: 100, id: 2 },
    { x: now - 3 * 60 * 1000, y: 75, id: 3 },
    { x: now - 2 * 60 * 1000, y: 50, id: 4 },
    { x: now - 1 * 60 * 1000, y: 40, id: 5 },
  ];

  const series = React.useMemo(
    () => [
      { id: 'caudal', data: data1, label: 'Caudal' },
      { id: 'presion', data: data2, label: 'Presión' },
      { id: 'temperatura', data: data3, label: 'Temperatura' },
    ],
    []
  );

  // Dibuja las líneas conectando los puntos
  function LinkPoints({ seriesId }: { seriesId: string }) {
    const scatter = useScatterSeries(seriesId);
    const xScale = useXScale();
    const yScale = useYScale();

    if (!scatter?.data) return null;

    const { color, data } = scatter;
    const pathD = `M ${data.map(({ x, y }) => `${xScale(x)},${yScale(y)}`).join(' L ')}`;

    return <path fill="none" stroke={color} strokeWidth={2} d={pathD} />;
  }

  return (
    <Card sx={sx}>
      <CardHeader title="Mediciones en Tiempo Real" />
      <CardContent>
        <ScatterChart
          series={series}
          height={350}
          xAxis={[
            {
              scaleType: 'time',
              label: 'Hora',
              valueFormatter: (value) =>
                new Date(value as number).toLocaleTimeString('es-CL', {
                  hour: '2-digit',
                  minute: '2-digit',
                }),
            },
          ]}
          yAxis={[{label: 'Valor Medido' }]}
        >
          <LinkPoints seriesId="caudal" />
          <LinkPoints seriesId="presion" />
          <LinkPoints seriesId="temperatura" />
        </ScatterChart>
      </CardContent>
      <Divider />
    </Card>
  );
}
