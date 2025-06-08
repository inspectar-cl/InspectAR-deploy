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

  //     { x: 95,  y: 200, id: 1 },

  const data1 = [
    { x: 50, y: 300, id: 1 },
    { x: 90, y: 450, id: 2 },
    { x: 130, y: 500, id: 3 },
    { x: 200, y: 350, id: 4 },
    { x: 280, y: 280, id: 5 },
  ];
  const data2 = [
    { x: 50, y: 100, id: 2 },
    { x: 130, y: 300, id: 3 },
    { x: 190, y: 250, id: 4 },
    { x: 250, y: 200, id: 5 },
    { x: 290, y: 180, id: 6 },
  ];
  const data3 = [
    { x: 60, y: 60, id: 2 },
    { x: 140, y: 120, id: 3 },
    { x: 190, y: 75, id: 4 },
    { x: 250, y: 50, id: 5 },
    { x: 300, y: 40, id: 6 },
  ];

  // Series en el formato que MUI X Charts espera
  const series = React.useMemo(
    () => [
      { id: 's1', data: data1, label: 'Caudal' },
      { id: 's2', data: data2, label: 'Presion' },
      { id: 's3', data: data3, label: 'Temperatura' },
    ],
    []
  );

  // Componente que une los puntos de cada serie
  function LinkPoints({ seriesId, close }: { seriesId: string; close?: boolean }) {
    const scatter = useScatterSeries(seriesId);
    const xScale = useXScale();
    const yScale = useYScale();

    if (!scatter?.data) return null;

    const { color, data } = scatter;
    // Construye la ruta SVG M x1,y1 L x2,y2 …
    const pathD =
      `M ${data.map(({ x, y }) => `${xScale(x)},${yScale(y)}`).join(' L ')}`

    return <path fill="none" stroke={color} strokeWidth={2} d={pathD} />;
  }

  return (
    <Card sx={sx}>
      <CardHeader title="Scatter Plot" />
      <CardContent>
        <ScatterChart series={series} height={350}>
          {/* Estas dos capas dibujan las líneas entre puntos */}
          <LinkPoints seriesId="s1" />
          <LinkPoints seriesId="s2" />
          <LinkPoints seriesId="s3" />
        </ScatterChart>
      </CardContent>
      <Divider />
    </Card>
  );
}
