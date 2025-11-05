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
import Chip from '@mui/material/Chip';
import Grid from '@mui/material/Grid';
import { QrCode as QrCodeIcon } from '@phosphor-icons/react/dist/ssr/QrCode';
import { Warning as WarningIcon } from '@phosphor-icons/react/dist/ssr/Warning';
import { CheckCircle as CheckCircleIcon } from '@phosphor-icons/react/dist/ssr/CheckCircle';

const API_URL = process.env.NEXT_PUBLIC_API_GATEWAY_URL || '/api';

interface ActivoInfo {
  codigo: string;
  nombre: string;
  tipo: string;
  descripcion: string;
  ubicacion: string;
}

interface Sensor {
  id: number;
  nombre: string;
  tipo: string;
  unidad: string;
  valor?: number;
  estado?: string;
  ultima_lectura?: string;
}

export default function QRInfoPage(): React.JSX.Element {
  const params = useParams();
  const codigoActivo = params.codigo_activo as string;

  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [activoInfo, setActivoInfo] = React.useState<ActivoInfo | null>(null);
  const [sensores, setSensores] = React.useState<Sensor[]>([]);

  React.useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);

        // Obtener información básica del activo mediante código QR
        const response = await fetch(`${API_URL}/codigo_qr/${codigoActivo}`, {
          headers: {
            'ngrok-skip-browser-warning': 'true',
          },
        });
        
        if (!response.ok) {
          throw new Error('No se pudo obtener la información del activo');
        }

        const data = await response.json();
        setActivoInfo(data);

        // Aquí podrías hacer llamadas adicionales si necesitas datos de sensores
        // Por ejemplo, si el endpoint devuelve el ID del activo,
        // podrías hacer otra llamada para obtener los sensores:
        // if (data.activo_id) {
        //   const sensoresResponse = await fetch(`${API_URL}/sensores/${data.activo_id}`);
        //   const sensoresData = await sensoresResponse.json();
        //   setSensores(sensoresData.sensores || []);
        // }

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
                <Chip
                  label={activoInfo.tipo}
                  color="primary"
                  variant="outlined"
                  size="small"
                />
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

          {/* Sección de sensores (si hay) */}
          {sensores.length > 0 && (
            <>
              <Typography variant="h6" fontWeight="bold">
                Sensores
              </Typography>
              
              <Grid container spacing={2}>
                {sensores.map((sensor) => (
                  <Grid size={{ xs: 12, sm: 6, md: 4 }} key={sensor.id}>
                    <Card variant="outlined">
                      <CardContent>
                        <Stack spacing={1}>
                          <Stack direction="row" justifyContent="space-between" alignItems="center">
                            <Typography variant="subtitle2" fontWeight="bold">
                              {sensor.nombre}
                            </Typography>
                            {sensor.estado === 'activo' ? (
                              <CheckCircleIcon color="var(--mui-palette-success-main)" />
                            ) : (
                              <WarningIcon color="var(--mui-palette-warning-main)" />
                            )}
                          </Stack>
                          
                          <Typography variant="caption" color="text.secondary">
                            {sensor.tipo}
                          </Typography>
                          
                          {sensor.valor !== undefined && (
                            <Typography variant="h6" color="primary">
                              {sensor.valor} {sensor.unidad}
                            </Typography>
                          )}
                          
                          {sensor.ultima_lectura && (
                            <Typography variant="caption" color="text.secondary">
                              Última lectura: {new Date(sensor.ultima_lectura).toLocaleString('es-CL')}
                            </Typography>
                          )}
                        </Stack>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            </>
          )}

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
