import { activosMock } from '@/mocks'; // tu mock de activos
import React from 'react';
import Card from '@mui/material/Card';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box'
import dayjs from 'dayjs';

import ScoreChart from '@/components/dashboard/overview/score-chart';
import { Scatter } from '@/components/dashboard/overview/scatter';
import { Budget } from '@/components/dashboard/overview/budget';
import { TotalCustomers } from '@/components/dashboard/overview/total-customers';
import { TemperatureProgress } from '@/components/dashboard/overview/temperature';
import { LatestProducts } from '@/components/dashboard/overview/latest-products';
import { WarningIcon } from '@phosphor-icons/react/dist/ssr/Warning';
import Grid from '@mui/material/Grid';

export default async function ActivoDetailPage({ params }: { params: { id: string } }) {
  const { id } = params;

  const activo = activosMock.find((a) => a.id === id);

  if (!activo) {
    return <div style={{ padding: '1rem' }}>No se encontró el activo con ID: {id}</div>;
  }

  return (
    <Box sx={{ p: 2 }}>
      {/* FILA SUPERIOR */}
      <Grid container spacing={2}>
        {/* Columna izquierda: Budget y TotalCustomers */}


        {/* Columna derecha: Tarjeta del activo */}
        <Grid size={{md:8, xs:12}}>
          <Card sx={{ display: 'flex', height: 400 }}>
            {/* Imagen izquierda */}
            <CardMedia
              component="img"
              sx={{ width: 350, height: '100%', objectFit: 'cover' }}
              image={activo.img || '/static/images/cards/default-placeholder.png'}
              alt={activo.tipoActivo}
            />
            {/* Info derecha */}
            <Box sx={{ display: 'flex', flexDirection: 'column', width: '100%' }}>
              <CardContent sx={{ flex: '1 0 auto' }}>
                <Typography gutterBottom variant="h2" component="div">
                  {activo.tipoActivo}
                </Typography>
                <Typography variant="h6" sx={{ color: 'text.secondary', mb: 1 }}>
                  <strong>ID:</strong> {activo.id}
                </Typography>
                <Typography variant="h6" sx={{ color: 'text.secondary', mb: 1 }}>
                  <strong>Estado:</strong> {activo.estado}
                </Typography>
                <Typography variant="h6" sx={{ color: 'text.secondary', mb: 1 }}>
                  <strong>Descripción:</strong> {activo.descripcion}
                </Typography>
                <Typography variant="h6" sx={{ color: 'text.secondary' }}>
                  <strong>Ubicación:</strong> {activo.ubicacion}
                </Typography>
              </CardContent>
              <CardActions>
                <Button size="small">Ver historial de mantenimiento</Button>
                <Button size="small">Alertas recientes</Button>
              </CardActions>
            </Box>
          </Card>
        </Grid>
        <Grid size={{md:4, xs:12}}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <Budget diff={5} trend="up" sx={{ height: 192 }} value="48 m³/h" />
            <TotalCustomers diff={10} trend="down" sx={{ height: 192}} value="15 Psi" />
          </Box>
        </Grid>
      </Grid>

      {/* FILA INFERIOR */}
      <Grid container spacing={2} sx={{ mt: 2 }}>
        {/* Scatter: 70% */}
        <Grid size={{md:10, xs:12}}>
          <Scatter sx={{ height: 480 }} />
        </Grid>
        {/* ScoreChart: 30% */}
        <Grid size={{md:2, xs:12}}>
          <TemperatureProgress  value={75} /> {/* 🔥 Aquí lo usas */}
        </Grid>
        <Grid size={{md:3, xs:12}}>
          <ScoreChart sx={{ height: 450 }} />
        </Grid>
      <Grid size={{lg:9, md:6, xs:12}}>
        <LatestProducts
          products={[
            {
              id: 'SENS-005',
              name: 'Motor sobrecalentado, riesgo de falla',
              icon: <WarningIcon size={32} weight="fill" color="#ff0000"/>,
              updatedAt: dayjs().subtract(18, 'minutes').subtract(5, 'hour').toDate(),
            },
            {
              id: 'SENS-004',
              name: 'Subida del caudal',
              icon: <WarningIcon size={32} weight="fill" color="#ff0000"/>,
              updatedAt: dayjs().subtract(41, 'minutes').subtract(3, 'hour').toDate(),
            },
            {
              id: 'SENS-003',
              name: 'Presion muy baja',
              icon: <WarningIcon size={32} weight="fill" color="#ff9214" />,
              updatedAt: dayjs().subtract(5, 'minutes').subtract(3, 'hour').toDate(),
            },
            {
              id: 'SENS-002',
              name: 'Aumento inusual de temperatura',
              icon: <WarningIcon size={32} weight="fill" color="#ff9214" />,
              updatedAt: dayjs().subtract(23, 'minutes').subtract(2, 'hour').toDate(),
            },
          ]}
          sx={{ height: 450 }}
        />
      </Grid>
      </Grid>
    </Box>
  );
}
