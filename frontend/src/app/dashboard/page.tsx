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
import { useUser } from '@/hooks/use-user';
import { activosMock } from '@/mocks/ActivosEdificiosMocks';

const gs = new Services();

interface SensorData {
  sensor_id: string;
  datos: { tiempo: string; valor: number }[];
}

const estadoColor: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
  OK: 'success',
  Medio: 'warning',
  Crítico: 'error',
} as const;

const fallbackImg = 'https://via.placeholder.com/640x360?text=Activo';

export default function Page(): React.JSX.Element {
  const { user, isLoading: loadingUser, error: userError } = useUser();

  const [sensores, setSensores] = React.useState<SensorData[]>([]);

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

  //Vista “Residente”: tarjetas 2x2 de activos críticos
  const isResidente = user?.role === 'residente';
  const edificios = (user as any)?.edificioId as EdificioAPI[] | undefined;

  // construimos un Set con claves de edificio para comparar
  const buildingKeySet = React.useMemo<Set<string>>(() => {
    if (!Array.isArray(edificios) || !edificios.length) return new Set();
    return new Set(edificios.map(e => String(e.id).toUpperCase()));
  }, [edificios]);

  const activosDelUsuario = React.useMemo(() => {
    const severity: Record<Estado, number> = { Crítico: 3, Medio: 2, OK: 1, NN: 0 };

    if (buildingKeySet.size === 0) return [];
    return activosMock
      .filter(a => buildingKeySet.has(String(a.id_edificio).toUpperCase()))
      .sort((a, b) => (severity[b.estado as Estado] ?? 0) - (severity[a.estado as Estado] ?? 0));
  }, [buildingKeySet]);

  type Estado = 'OK' | 'Medio' | 'Crítico' | 'NN';

  type EdificioAPI = {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
  };

  if (loadingUser) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress />
      </Box>
    );
  }
  if (userError) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <Typography color="error">{userError}</Typography>
      </Box>
    );
  }

  if (isResidente) {
    return (
      <Box>
        <Typography variant="h4" sx={{ mb: 2 }}>
          Activos de tu edificio
        </Typography>

        {activosDelUsuario.length === 0 ? (
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
            {activosDelUsuario.map((a) => (
              <Grid key={a.id} size={{ xs: 12, md: 6 }}>
                <Card sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, height: { md: 220 } }}>
                  <CardMedia
                    component="img"
                    src={a.img || fallbackImg}
                    alt={a.tipoActivo}
                    sx={{ width: { md: 260 }, height: { xs: 200, md: '100%' }, objectFit: 'cover' }}
                  />
                  <Box sx={{ display: 'flex', flexDirection: 'column', flex: 1 }}>
                    <CardContent sx={{ pb: 1.5 }}>
                      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                        <Typography variant="h6">{a.tipoActivo}</Typography>
                        <Chip
                          label={a.estado}
                          color={estadoColor[a.estado] ?? 'default'}
                          size="small"
                          sx={{ fontWeight: 600 }}
                        />
                      </Stack>
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                        <strong>Ubicación:</strong> {a.ubicacion}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" noWrap title={a.descripcion}>
                        {a.descripcion || 'Sin descripción'}
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
