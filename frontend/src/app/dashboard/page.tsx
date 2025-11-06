'use client';

import * as React from 'react';
import { CircularProgress, Box, Typography } from '@mui/material';
import { useUserToken } from '@/hooks/use-usertoken';

// Vistas por rol
import ResidenteHome from './home-pages/ResidenteHome';
import TecnicoHome from './home-pages/TecnicoHome';
import AnalistaHome from './home-pages/AnalistaHome';
import AdministradorHome from './home-pages/AdministradorHome';
import RootHome from './home-pages/RootHome';

export default function Page(): React.JSX.Element {
  const { user, isLoading, error } = useUserToken();

  if (isLoading) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error || !user) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <Typography color="error">{error || 'No autenticado'}</Typography>
      </Box>
    );
  }

  const role = (user.role || '').toLowerCase();

  // 👉 Dirige por rol
  if (role === 'residente') return <ResidenteHome />;
  if (role === 'tecnico') return <TecnicoHome />;
  if (role === 'analista') return <AnalistaHome />;
  if (role === 'administrador') return <AdministradorHome />;
  if (role === 'root') return <RootHome />;

  // Fallback: si llega un rol desconocido
  return (
    <Box sx={{ py: 8, textAlign: 'center' }}>
      <Typography variant="h6">Rol no soportado: {String(user.role)}</Typography>
    </Box>
  );
}
