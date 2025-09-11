/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/

'use client';

import * as React from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import Divider from '@mui/material/Divider';
import dayjs from 'dayjs';
import Select from '@mui/material/Select'
import MenuItem from '@mui/material/MenuItem'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'


import { ScatterChart } from '@mui/x-charts/ScatterChart';
import { useScatterSeries, useXScale, useYScale } from '@mui/x-charts/hooks';
import {type SxProps, type Theme } from '@mui/material';


interface SensorData {
  sensor_id: string
  datos: { tiempo: string; valor: number }[]
}

interface SensorScatterProps {
  sx?: SxProps<Theme>;
  //dataTemp: SensorData
  dataTemp: SensorData
}

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

export function SensorScatter({ sx, dataTemp }: SensorScatterProps) {
  const [rangoMinutos, setrangoMinutos] = React.useState(360); // por defecto: ultimos 30 min

  const convert = (sensor: SensorData) => {
    if (!sensor?.datos || sensor.datos.length === 0) {
      // router.push('/errors');
      return []; //Para que no falle el map, nos vamos a error page
    }

    const oneHourAgo = dayjs().subtract(rangoMinutos, 'minute');

    return sensor.datos
      .filter(item => dayjs(item.tiempo).isAfter(oneHourAgo))
      .map((item, index) => ({
        x: dayjs(item.tiempo).valueOf(),
        y: item.valor,
        id: index,
      }));
  };
  
  const series = React.useMemo(
    () => [
      { id: 'temperatura', data: convert(dataTemp), label: 'Temperatura' },
      // { id: 'temperatura', data: data1, label: 'Temperatura', color: 'red'},
    ],
    [ dataTemp, rangoMinutos, convert]
  );

  return (
    <Card sx={sx}>
      <CardHeader title="Mediciones Sensor en Tiempo Real" 
        action={
          <FormControl size="small" sx={{ minWidth: 160, mb: 2 }}>
            <InputLabel id="rango-label">Rango de tiempo</InputLabel>
            <Select
              labelId="rango-label"
              value={rangoMinutos}
              label="Rango de tiempo"
              onChange={(e) => { setrangoMinutos(Number(e.target.value)); }}
            >
              <MenuItem value={30}>Últimos 30 minutos</MenuItem>
              <MenuItem value={60}>Última hora</MenuItem>
              <MenuItem value={240}>Últimas 4 horas</MenuItem>
              <MenuItem value={240}>Últimas 8 horas</MenuItem>
              <MenuItem value={240}>Últimas 24 horas</MenuItem>
            </Select>
          </FormControl>
        }
      />
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
          {/* <LinkPoints seriesId="caudal" />
          <LinkPoints seriesId="presion" /> */}
          <LinkPoints seriesId="temperatura" />
        </ScatterChart>
      </CardContent>
      <Divider />
    </Card>
  );
}
