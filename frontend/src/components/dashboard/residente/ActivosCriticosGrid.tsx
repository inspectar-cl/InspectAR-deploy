'use client';

import * as React from 'react';
import Grid from '@mui/material/Grid';
import Card from '@mui/material/Card';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import Chip from '@mui/material/Chip';
import Typography from '@mui/material/Typography';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';

import { activosMock } from '@/mocks/ActivosEdificiosMocks'; // <-- si lo tienes en archivo; si no, pega ahí tu array


type Props = {
  edificioId: string | undefined;
};

const estadoColor: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
  'OK': 'success',
  'Medio': 'warning',
  'Crítico': 'error'
} as const;

const fallbackImg = 'https://via.placeholder.com/640x360?text=Activo';

type Estado = 'OK' | 'Medio' | 'Crítico' | 'NN';

// severidad completa para todos los estados posibles
const severity: Record<Estado, number> = {
  'Crítico': 3,
  'Medio': 2,
  'OK': 1,
  'NN': 0
} as const;

export function ActivosCriticosGrid({ edificioId }: Props) {
  // Mostrar TODOS los activos del edificio (y ordenar por severidad opcional)
    const activos = React.useMemo(() => {
    return activosMock
        .filter(a => !edificioId || a.id_edificio === edificioId)
        .sort((a, b) => severity[b.estado as Estado] - severity[a.estado as Estado]);
    }, [edificioId]);

    if (!activos.length) {
        return (
        <Box sx={{ textAlign: 'center', py: 6 }}>
            <Typography variant="h5" color="text.secondary">
            No hay activos críticos en tu edificio
            </Typography>
            <Typography variant="body2" color="text.secondary">
            Todo en orden por ahora. Te avisaremos si algo cambia.
            </Typography>
        </Box>
        );
    }

    return (
        <Grid container spacing={3}>
        {activos.map((a) => (
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
    );
}
