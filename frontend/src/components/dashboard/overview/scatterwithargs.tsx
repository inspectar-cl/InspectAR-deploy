'use client';

import * as React from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import Divider from '@mui/material/Divider';
import { useTheme } from '@mui/material/styles';
import { ScatterChart } from '@mui/x-charts/ScatterChart';
import { useScatterSeries, useXScale, useYScale } from '@mui/x-charts/hooks';

type SensorData = {
  sensor_id: string
  datos: { tiempo: string; valor: number }[]
}

type Props = {
  sx?: any
  dataCaudal: SensorData
  dataPresion: SensorData
  dataTemp: SensorData
}

export function ScatterWithArgs({ sx, dataCaudal, dataPresion, dataTemp }: Props) {
  const theme = useTheme()

  const convert = (sensor: SensorData) =>
    sensor.datos.map((item, index) => ({
      x: new Date(item.tiempo).getTime(),
      y: item.valor,
      id: index,
    }))

  const series = React.useMemo(
    () => [
      { id: 'caudal', data: convert(dataCaudal), label: 'Caudal' },
      { id: 'presion', data: convert(dataPresion), label: 'Presión' },
      { id: 'temperatura', data: convert(dataTemp), label: 'Temperatura' },
    ],
    [dataCaudal, dataPresion, dataTemp]
  )

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
