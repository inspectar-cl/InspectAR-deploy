 
'use client';
import * as React from 'react';
import Grid from '@mui/material/Grid';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import CircularProgress from '@mui/material/CircularProgress';
import Alert from '@mui/material/Alert';
import RefreshIcon from '@mui/icons-material/Refresh';
import { LatestSensors } from '@/components/dashboard/overview/latest-sensors';
import { useActivosWithSensors } from '@/hooks/use-activos-with-sensors';

export default function Page(): React.JSX.Element {
  const { activos, loading, error, refresh } = useActivosWithSensors();

  if (loading) {
    return (
      <Box 
        display="flex" 
        flexDirection="column" 
        alignItems="center" 
        justifyContent="center" 
        minHeight="60vh"
        gap={2}
      >
        <CircularProgress size={60} />
        <Typography variant="h6" color="text.secondary">
          Cargando datos de sensores...
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Obteniendo información de activos y estado de sensores
        </Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Box 
        display="flex" 
        flexDirection="column" 
        alignItems="center" 
        justifyContent="center" 
        minHeight="60vh"
        gap={2}
      >
        <Alert severity="error" sx={{ maxWidth: 600 }}>
          <Typography variant="h6" gutterBottom>
            Error al cargar los datos
          </Typography>
          <Typography variant="body2" gutterBottom>
            {error}
          </Typography>
          <Button 
            variant="outlined" 
            startIcon={<RefreshIcon />} 
            onClick={refresh}
            sx={{ mt: 1 }}
          >
            Reintentar
          </Button>
        </Alert>
      </Box>
    );
  }

  if (activos.length === 0) {
    return (
      <Box 
        display="flex" 
        flexDirection="column" 
        alignItems="center" 
        justifyContent="center" 
        minHeight="60vh"
        gap={2}
      >
        <Typography variant="h5" color="text.secondary">
          No se encontraron activos
        </Typography>
        <Typography variant="body1" color="text.secondary" textAlign="center">
          No hay activos registrados en el sistema o no están disponibles en este momento.
        </Typography>
        <Button 
          variant="outlined" 
          startIcon={<RefreshIcon />} 
          onClick={refresh}
        >
          Actualizar
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      {/* Header con botón de refresh */}
      <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 3 }}>
        <Typography variant="h4" component="h1">
          Monitoreo de Sensores
        </Typography>
        <Button 
          variant="outlined" 
          startIcon={<RefreshIcon />} 
          onClick={refresh}
          disabled={loading}
        >
          Actualizar
        </Button>
      </Stack>

      {/* Grid de activos */}
      <Grid
        container
        spacing={3}
        justifyContent="center"
        alignItems="flex-start"
      >
        {activos.map((activo) => (
          <Grid key={activo.assetName} size={{ lg: 8, md: 12, xs: 12 }}>
            
            {/* Botones de simulación (solo para desarrollo) */}
            {/* {process.env.NODE_ENV === 'development' && activo.sensores.length > 0 && (
              <Stack direction="row" spacing={1} sx={{ mb: 1, justifyContent: 'flex-end' }}>
                <Button 
                  size="small" 
                  variant="outlined" 
                  onClick={() => { simulateDown(idx, activo.sensores[0].id); }}
                >
                  Simular caída ({activo.sensores[0].name})
                </Button>
                <Button 
                  size="small" 
                  variant="outlined" 
                  onClick={() => { simulateUp(idx, activo.sensores[0].id); }}
                >
                  Simular recovery
                </Button>
              </Stack>
            )} */}

            <LatestSensors
              assetName={activo.assetName}
              imageUrl={activo.imageUrl}
              sensors={activo.sensores}
              sx={{ height: '100%' }}
            />
          </Grid>
        ))}
      </Grid>

      {/* Información de última actualización */}
      {/* <Box sx={{ mt: 3, textAlign: 'center' }}>
        <Typography variant="caption" color="text.secondary">
          Última actualización: {dayjs().format('DD/MM/YYYY HH:mm:ss')} • 
          Actualización automática cada 30 segundos
        </Typography>
      </Box> */}
    </Box>
  );
}
