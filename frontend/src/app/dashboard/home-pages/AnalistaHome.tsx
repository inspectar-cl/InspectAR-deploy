'use client';

import * as React from 'react';
import { Box, Typography, Stack, Button } from '@mui/material';
import { ActivosGrid } from './_common';

export default function AnalistaHome() {
  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 1 }}>
        Panel del Analista 
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Vista general de activos por criticidad. Próximamente comparativas y tendencias.
      </Typography>

      {/* Atajos opcionales */}
      <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
        <Button size="small" variant="outlined">Reportes</Button>
        <Button size="small" variant="outlined">Comparativas</Button>
      </Stack>

      <ActivosGrid title="Activos monitoreados" />
    </Box>
  );
}
