'use client';

import * as React from 'react';
import { Box, Typography, Stack, Button } from '@mui/material';
import { ActivosGrid } from './_common';

export default function TecnicoHome() {
  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 1 }}>
        Panel del Técnico 
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Acceso rápido a activos por criticidad.
      </Typography>

      {/* Acciones rápidas opcionales */}
      <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
        <Button size="small" variant="outlined">Órdenes de trabajo</Button>
        <Button size="small" variant="outlined">Alertas abiertas</Button>
        <Button size="small" variant="outlined">Checklist de mantenimiento</Button>
      </Stack>

      <ActivosGrid title="Activos a tu cargo" />
    </Box>
  );
}
