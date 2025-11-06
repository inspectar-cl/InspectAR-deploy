'use client';

import * as React from 'react';
import { Box, Typography, Stack, Button } from '@mui/material';
import { ActivosGrid } from './_common';
import { useRouter } from 'next/navigation';

export default function AdministradorHome(): React.JSX.Element {
  const router = useRouter();

  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 1 }}>
        Bienvenido/a Administrador
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        Aquí puedes supervisar el estado general de los activos, gestionar formularios,
        revisar reportes técnicos y coordinar el mantenimiento del edificio.
      </Typography>

      <Stack direction="row" spacing={2} sx={{ mb: 3 }}>
        <Button
          variant="contained"
          onClick={() => {
            router.push('/dashboard/formulario');
          }}
        >
          Ver Formularios
        </Button>
        <Button
          variant="outlined"
          onClick={() => {
            router.push('/dashboard/reportes-tecnicos');
          }}
        >
          Reportes Técnicos
        </Button>
      </Stack>

      <ActivosGrid title="Activos del edificio" />
    </Box>
  );
}
