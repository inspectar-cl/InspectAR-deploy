'use client';

import * as React from 'react';
import {
  Box, Grid, Typography, Card, CardContent, CardMedia, Chip, CircularProgress, Stack
} from '@mui/material';
import Services from '@/modules/Services';
import { useUserToken } from '@/hooks/use-usertoken';

const gs = new Services();

export type Estado = 'OK' | 'Medio' | 'Crítico' | 'NN';

export interface Activo {
  id: number
  img?: string
  edificio_id: number
  creado_en: string
  nombre: string
  estado: Estado
  tipo: string
  ubicacion: string
}

export interface ApiActivosResponse {
  activos?: Activo[];
  error?: { mensaje: string };
}

const estadoColor: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
  OK: 'success',
  Medio: 'warning',
  Crítico: 'error',
} as const;

const fallbackImg = 'https://via.placeholder.com/640x360?text=Activo';

/** Hook que trae activos del usuario y los ordena por criticidad*/
export function useActivosDelUsuario(): {
  activos: Activo[] | null;
  loading: boolean;
  err: string | null;
} {
  const { user, isLoading } = useUserToken();
  const [activos, setActivos] = React.useState<Activo[] | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [err, setErr] = React.useState<string | null>(null);

  React.useEffect(() => {
    const run = async () => {
      if (isLoading || !user) {
        setLoading(false);
        return;
      }

      setLoading(true);
      setErr(null);
      try {
        const edificiosIds = new Set(user.edificio?.map((e) => e.id) ?? []);
        const severity: Record<Estado, number> = { Crítico: 3, Medio: 2, OK: 1, NN: 0 };

        const resp = await gs.authorizedGet('/obtener-activos-id/1', user.token) as ApiActivosResponse;
        if (resp?.error) {
          setErr(resp.error.mensaje || 'Error al obtener activos.');
          setActivos(null);
          return;
        }

        const filtrados = (resp?.activos ?? [])
          .filter(a => edificiosIds.has(a.edificio_id))
          .sort((a, b) => severity[b.estado] - severity[a.estado]);

        setActivos(filtrados);
      } catch {
        setErr('Error de red o token inválido/expirado.');
        setActivos(null);
      } finally {
        setLoading(false);
      }
    };

    void run();
  }, [isLoading, user]);

  return { activos, loading, err };
}

/** Sección de tarjetas 2x2 para mostrar activos (reutilizable para todos los roles) */
export function ActivosGrid({ title }: { title: string }) {
  const { activos, loading, err } = useActivosDelUsuario();

  if (loading) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (err) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <Typography color="error">{err}</Typography>
      </Box>
    );
  }

  const lista = activos || [];

  return (
    <Box>
      <Typography variant="h4" sx={{ mb: 2 }}>{title}</Typography>

      {lista.length === 0 ? (
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
        {lista.map((a) => {
            const getImageByTipo = (tipo: string): string => {
            const type = tipo.toLowerCase();
            if (type === 'ascensor') {
                return 'https://www.schindler.cl/content/dam/website/lac/images/modernizacion/cabinas/ascensor-interior-4.jpg/_jcr_content/renditions/original./ascensor-interior-4.jpg'; 
            }
            if (type.includes('Caldera')) {
                return 'https://gstingenieria.cl/wp-content/uploads/2020/10/sala-de-calderas.jpg'; 
            }
            if (type === 'bomba de agua') {
                return 'https://www.tstservicios.com/wp-content/uploads/2021/06/bomba-agua-necesito.jpg'; 
            }
            if (type.includes('Sistema Eléctrico') || type.includes('electrico')) {
                return 'https://previews.123rf.com/images/greentellect/greentellect1903/greentellect190300075/120087711-open-circuit-board-connection-or-eletrical-panel-in-modern-building.jpg'; 
            }
                return fallbackImg; // fallback
            };

            // Usar img del activo si existe, sino una por tipo
            const imageSrc =
            a?.img && a.img.trim() !== '' ? a.img : getImageByTipo(a.tipo);
            return (
            <Grid key={a.id} size={{ xs: 12, md: 6 }}>
                <Card
                sx={{
                    display: 'flex',
                    flexDirection: { xs: 'column', md: 'row' },
                    height: { md: 220 },
                }}
                >
                <CardMedia
                    component="img"
                    src={imageSrc}
                    alt={a.tipo}
                    sx={{
                    width: { md: 260 },
                    height: { xs: 200, md: '100%' },
                    objectFit: 'cover',
                    backgroundColor: '#f5f5f5',
                    }}
                />
                <Box sx={{ display: 'flex', flexDirection: 'column', flex: 1 }}>
                    <CardContent sx={{ pb: 1.5 }}>
                    <Stack
                        direction="row"
                        justifyContent="space-between"
                        alignItems="center"
                        sx={{ mb: 1 }}
                    >
                        <Typography variant="h6">{a.tipo}</Typography>
                        <Chip
                        label={a.estado}
                        color={estadoColor[a.estado] ?? 'default'}
                        size="small"
                        sx={{ fontWeight: 600 }}
                        />
                    </Stack>
                    <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ mb: 1 }}
                    >
                        <strong>Ubicación:</strong> {a.ubicacion}
                    </Typography>
                    <Typography
                        variant="body2"
                        color="text.secondary"
                        noWrap
                        title={a.creado_en}
                    >
                        {a.creado_en || 'Sin descripción'}
                    </Typography>
                    </CardContent>
                </Box>
                </Card>
            </Grid>
            );
        })}
        </Grid>
      )}
    </Box>
  );
}
