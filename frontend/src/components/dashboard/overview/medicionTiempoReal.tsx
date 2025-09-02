'use client';

/* Esta seccion del codigo corresponde a la grafica con las variables a sensorizar

Aqui se hace llamado a las API para la obtencion de la data correspondiente a cada sensor del activo

*/

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

import { useRouter } from 'next/navigation';

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
  const router = useRouter();
  const [rangoMinutos, setrangoMinutos] = React.useState(30); // por defecto: ultimos 30 min

  const convert = (sensor: SensorData) => {
    if (!sensor?.datos || sensor.datos.length === 0) {
      //router.push('/errors'); // Esta linea de aqui hace que vayamos a la page de errors si no hay data en la grafica de sensores
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
      { id: 'caudal', data: convert(dataCaudal), label: 'Caudal' },
      { id: 'presion', data: convert(dataPresion), label: 'Presión' },
      { id: 'temperatura', data: convert(dataTemp), label: 'Temperatura' },
    ],
    [dataCaudal, dataPresion, dataTemp, rangoMinutos]
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
      <CardHeader title="Mediciones en Tiempo Real" 
        action={
          <FormControl size="small" sx={{ minWidth: 220, mb: 2 }}>
            <InputLabel id="rango-label">Rango de tiempo</InputLabel>
            <Select
              labelId="rango-label"
              value={rangoMinutos}
              label="Rango de tiempo"
              onChange={(e) => setrangoMinutos(Number(e.target.value))}
            >
              <MenuItem value={30}>Últimos 30 minutos</MenuItem>
              <MenuItem value={60}>Última hora</MenuItem>
              <MenuItem value={240}>Últimas 4 horas</MenuItem>
              <MenuItem value={720}>Últimas 12 horas</MenuItem>
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
          <LinkPoints seriesId="caudal" />
          <LinkPoints seriesId="presion" />
          <LinkPoints seriesId="temperatura" />
        </ScatterChart>
      </CardContent>
      <Divider />
    </Card>
  );
}
