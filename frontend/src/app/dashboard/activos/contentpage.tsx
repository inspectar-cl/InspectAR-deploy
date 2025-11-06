'use client';

import React, { useState } from 'react';
import type { Metadata } from 'next';
import { config } from '@/config';
import Stack from '@mui/material/Stack';
import Grid from '@mui/material/Grid';
import {MapaActivos} from '@/components/dashboard//overview/mapa-activos';

import DataGridDemo from '@/components/dashboard/alertas/alert-grid';

export const metadata = { title: `Activos | Dashboard | ${config.site.name}` } satisfies Metadata;

export default function ContentActivosPage(): React.JSX.Element {
  const [edificioSeleccionado, setEdificioSeleccionado] = useState<string | null>(null);

  return (
    <Stack spacing={2} id="tour-mapa-activos">
      <Grid size={{lg:8, md:6, xs:12}}>
        <MapaActivos onSeleccionarEdificio={(id) =>
            { setEdificioSeleccionado((prev) => (prev === id ? null : id)); }
        }/>
      </Grid>
      <Grid id="tour-lista-activos" size={{lg:20, md:10, xs:20}}>
          <DataGridDemo sx={{ height: '100%' }} edificioSeleccionado={edificioSeleccionado} />
      </Grid>
    </Stack>
  );
}