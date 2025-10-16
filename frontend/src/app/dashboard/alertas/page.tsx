import * as React from 'react';
import type { Metadata } from 'next';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

import { config } from '@/config';

import { AlertasTable } from '@/components/dashboard/alertas/AlertasTable';

export const metadata = { title: `Alertas | Dashboard | ${config.site.name}` } satisfies Metadata;

export default function Page(): React.JSX.Element {
  return (
    <Stack spacing={3}>
      <div>
        <Typography variant="h4">Alertas</Typography>
      </div>
      <AlertasTable />
    </Stack>
  )
}
