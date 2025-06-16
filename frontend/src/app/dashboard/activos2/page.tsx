import React from 'react';
import type { Metadata } from 'next';
import { config } from '@/config';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';

import DataGridDemo from '@/components/dashboard/alertas/alert-grid';

export const metadata = { title: `Activos | Dashboard | ${config.site.name}` } satisfies Metadata;

export default function ActivosPage(): React.JSX.Element {
  return (
    <Stack spacing={2}>
      <div>
        <Typography variant="h4">Lista de Activos</Typography>
      </div>
    <Grid size={{lg:20, md:10, xs:20}}>
        <DataGridDemo sx={{ height: '100%' }}/>
    </Grid>
    </Stack>
  );
}