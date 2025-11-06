'use client';

import * as React from 'react';
import { useParams } from 'next/navigation';
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CircularProgress from '@mui/material/CircularProgress';
import Container from '@mui/material/Container';
import Divider from '@mui/material/Divider';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import Alert from '@mui/material/Alert';
import Grid from '@mui/material/Grid';
import { QrCodeIcon } from '@phosphor-icons/react/dist/csr/QrCode';
import { WarningIcon } from '@phosphor-icons/react/dist/csr/Warning';
import Services from '@/modules/Services';

interface ActivoInfo {
  codigo: string;
  nombre: string;
  tipo: string;
  descripcion: string;
  ubicacion: string;
}

export default function QRInfoPage(): React.JSX.Element {
  const params = useParams();
  const codigoActivo = params.codigo_activo as string;

  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [activoInfo, setActivoInfo] = React.useState<ActivoInfo | null>(null);

  React.useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);

        const services = new Services();
        const data = await services.get(`/codigo_qr/${codigoActivo}`);

        // Verificar si hubo error
        if (data.error || data.mensaje === 'error inesperado') {
          throw new Error('No se pudo obtener la información del activo');
        }

        setActivoInfo(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error desconocido');
      } finally {
        setLoading(false);
      }
    };

    if (codigoActivo) {
      void fetchData();
    }
  }, [codigoActivo]);

  if (loading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
          bgcolor: 'var(--mui-palette-background-default)',
        }}
      >
        <Stack spacing={2} alignItems="center">
          <CircularProgress />
          <Typography color="text.secondary">Cargando información del activo...</Typography>
        </Stack>
      </Box>
    );
  }

  if (error || !activoInfo) {
    return (
      <Container maxWidth="md" sx={{ py: 8 }}>
        <Alert severity="error" icon={<WarningIcon fontSize="var(--Icon-fontSize)" />}>
          <Typography variant="h6" gutterBottom>
            Error al cargar información
          </Typography>
          <Typography variant="body2">{error || 'No se encontró el activo'}</Typography>
        </Alert>
      </Container>
    );
  }

  return (
    <Box
      sx={{
        minHeight: '100vh',
        bgcolor: 'var(--mui-palette-background-default)',
        py: 4,
      }}
    >
      <Container maxWidth="lg">
        <Stack spacing={3}>
          {/* Header con código QR */}
          <Card>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 2 }}>
                <QrCodeIcon fontSize="var(--icon-fontSize-lg)" />
                <Box sx={{ flex: 1 }}>
                  <Typography variant="h5" fontWeight="bold">
                    {activoInfo.nombre}
                  </Typography>
                  <Typography color="text.secondary" variant="body2">
                    Código: {activoInfo.codigo}
                  </Typography>
                </Box>
              </Stack>

              <Divider sx={{ my: 2 }} />

              <Grid container spacing={2}>
                <Grid size={{ xs: 12, sm: 6 }}>
                  <Typography variant="caption" color="text.secondary">
                    Ubicación
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    {activoInfo.ubicacion}
                  </Typography>
                </Grid>
                <Grid size={{ xs: 12, sm: 6 }}>
                  <Typography variant="caption" color="text.secondary">
                    Tipo
                  </Typography>
                  <Typography variant="body1" fontWeight="medium">
                    {activoInfo.tipo}
                  </Typography>
                </Grid>
                <Grid size={{ xs: 12 }}>
                  <Typography variant="caption" color="text.secondary">
                    Descripción
                  </Typography>
                  <Typography variant="body1">
                    {activoInfo.descripcion}
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>

          {/* Información adicional */}
          <Card variant="outlined">
            <CardContent>
              <Typography variant="body2" color="text.secondary" align="center">
                Esta información es de solo lectura y está disponible públicamente mediante código QR.
              </Typography>
              <Typography variant="caption" color="text.secondary" align="center" display="block" sx={{ mt: 1 }}>
                Para acceder a más funcionalidades, inicia sesión en la plataforma InspectAR.
              </Typography>
            </CardContent>
          </Card>
        </Stack>
      </Container>
    </Box>
  );
}
