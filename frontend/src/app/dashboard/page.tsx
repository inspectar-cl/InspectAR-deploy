'use client'

import * as React from 'react';
import { useEffect } from 'react';
import type { Metadata } from 'next';
import Grid from '@mui/material/Grid';
import { TemperatureProgress } from '@/components/dashboard/overview/temperature';

import { config } from '@/config';
import { SensorScatter } from '@/components/dashboard/overview/sensor-scatter';
import { type Activo } from '@/types/'

import Services from '@/modules/Services'

const gs = new Services()


export default function Page(): React.JSX.Element {
  const [activo, setActivo] = React.useState<Activo | null>(null)
  const [sensores, setSensores] = React.useState<any[]>([])

  const hasFetchedRef = React.useRef(false)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await gs.get("/lectura/Caldera1/datos");
        console.log("response", response);
        setSensores(response.sensores ?? []);
      } catch (err) {
        console.error("Error al obtener el activo", err);
      }
    };

    // fetchData();

    const interval = setInterval (() => {
      fetchData();
    }, 5000);

    return () => { clearInterval(interval); };
  })

  const sensorTemp = sensores.find(s => s.sensor_id === 'sensortemp')
  
  const getDiffInfo = (sensorId: string) => {
    const datos = sensores.find(s => s.sensor_id === sensorId)?.datos ?? []
    const ultimo = datos.at(-1)?.valor ?? 0
    const penultimo = datos.at(-2)?.valor ?? 0
    const diff = ultimo - penultimo
    const trend = (diff >= 0 ? 'up' : 'down')
    
    return {
      valor: ultimo,
      diff: Math.abs(diff),
      trend
    }
  }
  
  const temperaturaInfo = sensorTemp?.datos?.at(-1)?.valor ?? 0


  return (
    <Grid container spacing={3}>

      <Grid size={{lg:20, md:12, xs:12}}>
        <SensorScatter sx={{ height: 480 }} dataTemp={sensorTemp}/>
      </Grid>

      <Grid size={{md:2, xs:12}}>
        <TemperatureProgress value={temperaturaInfo} />
      </Grid>

    </Grid>
  );
}
