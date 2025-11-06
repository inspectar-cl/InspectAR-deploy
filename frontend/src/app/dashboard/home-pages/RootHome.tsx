'use client';

import * as React from 'react';
import { Box, Typography, Stack, Button } from '@mui/material';
import { ActivosGrid } from './_common';
import { useRouter } from 'next/navigation';

export default function RootHome(): React.JSX.Element {
  const router = useRouter();

  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 1 }}>
        Panel Root – InspectAR
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Este es tu panel principal como usuario Root. Desde aquí puedes gestionar la plataforma completa:
        agregar edificios, activos, sensores y administrar accesos globales.
      </Typography>

      <Stack direction="row" spacing={2} sx={{ mb: 3 }}>
        <Button
          variant="contained"
          onClick={() => router.push('/dashboard/gestion')}
        >
          Gestión General
        </Button>
        <Button
          variant="outlined"
          onClick={() => router.push('/dashboard/agregardatos')}
        >
          Añadir Datos
        </Button>
      </Stack>

      <ActivosGrid title="Resumen global de activos" />
    </Box>
  );
}
