'use client';

import * as React from 'react';
import { Box, Typography } from '@mui/material';
import { ActivosGrid } from './_common';

export default function ResidenteHome(): React.JSX.Element {
  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 1 }}>
        Bienvenido/a
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
        Aquí verás el estado de tus activos y una descripción breve.
      </Typography>

      <ActivosGrid title="Activos de tu edificio" />
    </Box>
  );
}
