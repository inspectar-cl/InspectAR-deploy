/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client'

import * as React from 'react';
import { useEffect } from 'react';
import Grid from '@mui/material/Grid';
import { TemperatureProgress } from '@/components/dashboard/overview/temperature';
import { SensorScatter } from '@/components/dashboard/overview/sensor-scatter';

import Services from '@/modules/Services'

const gs = new Services()

interface SensorData {
  sensor_id: string;
  datos: {
    tiempo: string;
    valor: number;
  }[];
}

export default function Page(): React.JSX.Element {
  const [sensores, setSensores] = React.useState<SensorData[]>([])

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await gs.get("/lectura/Caldera1/datos") as { sensores?: SensorData[] };
        setSensores(response.sensores ?? []);
      } catch (err) {
        // Error handling for failed sensor data fetch
        setSensores([]);
      }
    };

    void fetchData(); // Initial fetch

    const interval = setInterval (() => {
      void fetchData(); // Periodic fetch
    }, 5000);

    return () => { clearInterval(interval); };
  })

  const sensorTemp = sensores.find(s => s.sensor_id === 'sensortemp')
  const temperaturaInfo = sensorTemp?.datos?.at(-1)?.valor ?? 0


  return (
    <Grid container spacing={3}>

      <Grid size={{lg:20, md:12, xs:12}}>
        {sensorTemp ? <SensorScatter sx={{ height: 480 }} dataTemp={sensorTemp}/> : null}
      </Grid>

      <Grid size={{md:2, xs:12}}>
        <TemperatureProgress value={temperaturaInfo} />
      </Grid>

    </Grid>
  );
}
