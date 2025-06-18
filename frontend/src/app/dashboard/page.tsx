import * as React from 'react';
import type { Metadata } from 'next';
import Grid from '@mui/material/Grid';
import { TemperatureProgress } from '@/components/dashboard/overview/temperature';

import { config } from '@/config';
//Import SensorScatter es temporal, luego eliminar
import { SensorScatter } from '@/components/dashboard/overview/sensor-scatter';

export const metadata = { title: `Overview | Dashboard | ${config.site.name}` } satisfies Metadata;

export default function Page(): React.JSX.Element {

  const dataTemp = 1;

  return (
    <Grid container spacing={3}>

      <Grid size={{lg:20, md:12, xs:12}}>
        <SensorScatter sx={{ height: 480 }} dataTemp={dataTemp}/>
      </Grid>

      <Grid size={{md:2, xs:12}}>
        {/* Temperatura */}
        <TemperatureProgress value={dataTemp} />
      </Grid>

    </Grid>
  );
}
