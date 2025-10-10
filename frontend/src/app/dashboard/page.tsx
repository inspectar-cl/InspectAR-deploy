'use client';

import * as React from 'react';
import { useEffect } from 'react';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import Chip from '@mui/material/Chip';
import CircularProgress from '@mui/material/CircularProgress';

import { TemperatureProgress } from '@/components/dashboard/overview/temperature';
import { SensorScatter } from '@/components/dashboard/overview/sensor-scatter';

import Services from '@/modules/Services';
import { useUserToken } from '@/hooks/use-usertoken';
import { activosMock } from '@/mocks/ActivosEdificiosMocks';

const gs = new Services();

interface SensorData {
  sensor_id: string;
  datos: { tiempo: string; valor: number }[];
}

interface Activo {
  id: number
  edificio_id: number
  creado_en: string
  nombre: string
  estado: 'OK' | 'Medio' | 'Crítico' | 'NN'
  tipo: string
  ubicacion: string
}

const estadoColor: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
  OK: 'success',
  Medio: 'warning',
  Crítico: 'error',
} as const;

const fallbackImg = 'https://via.placeholder.com/640x360?text=Activo';

export default function Page(): React.JSX.Element {
  const { user, isLoading, error } = useUserToken();

  const [sensores, setSensores] = React.useState<SensorData[]>([]);

  const [activosDelUsuario, setActivosDelUsuario] = React.useState<Activo[] | null>(null);
  const [activosError, setActivosError] = React.useState<string | null>(null);
  const [loadingActivos, setLoadingActivos] = React.useState(true);

  type Estado = 'OK' | 'Medio' | 'Crítico' | 'NN';

  useEffect(() => {
    let alive = true;

    const fetchData = async (): Promise<void> => {
      try {
        const response = await gs.get('/lectura/Caldera1/datos') as { sensores?: SensorData[] };
        if (!alive) return;
        setSensores(response.sensores ?? []);
      } catch {
        if (!alive) return;
        setSensores([]);
      }
    };

    void fetchData();
    const interval = setInterval(() => { void fetchData(); }, 5000);

    return () => { alive = false; clearInterval(interval); };
  }, []); // 👈 importante

  useEffect( () => {
    // Si el usuario está cargando o no está definido, salimos.
    if (isLoading || !user) {
      setLoadingActivos(false);
      return;
    }

    // Función asíncrona para obtener y filtrar los activos
    const fetchActivos = async () => {
      setLoadingActivos(true);
      setActivosError(null);

      try {
        // Corregido: 'edificios' es el array de edificios del user
        const edificiosIds = new Set((user.edificio || []).map(e => e.id));
        const severity: Record<Estado, number> = { Crítico: 3, Medio: 2, OK: 1, NN: 0 };
        
        // LLAMADA ASÍNCRONA CORREGIDA (usando await)
        const responseActivo = await gs.authorizedGet('/obtener-activos', user.token) as { activos?: Activo[], error?: any};
        console.log('Respuesta de activos:', responseActivo);
        if (responseActivo.error) {
          setActivosError(responseActivo.error.mensaje || 'Error al obtener activos.');
          setActivosDelUsuario(null);
          return;
        }

        // Filtrado de activos
        const activosFiltrados = (responseActivo.activos ?? [])
          .filter(a => edificiosIds.has(a.edificio_id))
          .sort((a, b) => severity[b.estado] - severity[a.estado]);

        setActivosDelUsuario(activosFiltrados);
      } catch (err) {
        setActivosError('Error de red o token inválido/expirado.');
        setActivosDelUsuario(null);
      } finally {
        setLoadingActivos(false);
      }
    };

    void fetchActivos();
  }, [user, isLoading]);

  //Vista “Residente”: tarjetas 2x2 de activos críticos
  const isResidente = user?.role === 'residente';

  if (isLoading || loadingActivos) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress />
      </Box>
    );
  }
  if (error || activosError) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <Typography color="error">{error}</Typography>
      </Box>
    );
  }

  if (isResidente) {
    const activosFinales = activosDelUsuario || [];
    return (
      <Box>
        <Typography variant="h4" sx={{ mb: 2 }}>
          Activos de tu edificio
        </Typography>

        {activosFinales.length === 0 ? (
          <Box sx={{ textAlign: 'center', py: 6 }}>
            <Typography variant="h6" color="text.secondary">
              No hay activos críticos en tu edificio
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Todo en orden por ahora. Te avisaremos si algo cambia.
            </Typography>
          </Box>
        ) : (
          <Grid container spacing={3}>
            {activosFinales.map((a) => (
              <Grid key={a.id} size={{ xs: 12, md: 6 }}>
                <Card sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, height: { md: 220 } }}>
                  <CardMedia
                    component="img"
                    src={fallbackImg}
                    alt={a.tipo}
                    sx={{ width: { md: 260 }, height: { xs: 200, md: '100%' }, objectFit: 'cover' }}
                  />
                  <Box sx={{ display: 'flex', flexDirection: 'column', flex: 1 }}>
                    <CardContent sx={{ pb: 1.5 }}>
                      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                        <Typography variant="h6">{a.tipo}</Typography>
                        <Chip
                          label={a.estado}
                          color={estadoColor[a.estado] ?? 'default'}
                          size="small"
                          sx={{ fontWeight: 600 }}
                          onClick={() => {
                            /* Funcion vacia por ahora para evitar issue de que boton no tiene onClick*/
                          }}
                        />
                      </Stack>
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                        <strong>Ubicación:</strong> {a.ubicacion}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" noWrap title={a.creado_en}>
                        {a.creado_en || 'Sin descripción'}
                      </Typography>
                    </CardContent>
                  </Box>
                </Card>
              </Grid>
            ))}
          </Grid>
        )}
      </Box>
    );
  }

  // ------ Vista “no residente”: lo que ya tenías ------
  const sensorTemp = sensores.find(s => s.sensor_id === 'sensortemp');
  const temperaturaInfo = sensorTemp?.datos?.at(-1)?.valor ?? 0;

  return (
    <Grid container spacing={3}>
      <Grid size={{ lg: 20, md: 12, xs: 12 }}>
        {sensorTemp ? <SensorScatter sx={{ height: 480 }} dataTemp={sensorTemp} /> : null}
      </Grid>
      <Grid size={{ md: 2, xs: 12 }}>
        <TemperatureProgress value={temperaturaInfo} />
      </Grid>
    </Grid>
  );
}
