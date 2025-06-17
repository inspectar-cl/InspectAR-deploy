'use client';

import { activosMock } from '@/mocks'; // mock de activos
import * as React from 'react';
import Card from '@mui/material/Card';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box'
import dayjs from 'dayjs';
import type { Metadata } from 'next';
import { config } from '@/config';

import ScoreChart from '@/components/dashboard/overview/score-chart';
import { Scatter } from '@/components/dashboard/overview/scatter';
import { Budget } from '@/components/dashboard/overview/budget';
import { TotalCustomers } from '@/components/dashboard/overview/total-customers';
import { TemperatureProgress } from '@/components/dashboard/overview/temperature';
import { LatestProducts } from '@/components/dashboard/overview/latest-products';
import { WarningIcon } from '@phosphor-icons/react/dist/ssr/Warning';
import Grid from '@mui/material/Grid';
import { ScatterWithArgs} from '@/components/dashboard/overview/scatterwithargs'

// export const metadata = { title: `Activos | Dashboard | ${config.site.name}` } satisfies Metadata;

// Configuración rutas de obtención de datos desde db.
import Services from '@/modules/Services'

const gs = new Services()

const uris = {
  GET: `/lectura`
}

export default function ActivoDetailPage({ params }: { params: { id: string } }) {
  const { id } = params;

  console.log("id:", id)

  type Activo = {
    id: string
    tipoActivo: string
    estado: 'OK' | 'Medio' | 'Crítico' | 'NN'
    descripcion: string
    ubicacion: string
    img: string
  }

  const [activo, setActivo] = React.useState<Activo | null>(null)
  const [sensores, setSensores] = React.useState<any[]>([])

  const getDetalleActivo = async () => {
    try {
      const response = await gs.get(`/lectura/${id}/datos`)

      const activoTransformado: Activo = {
        id: response.activo_id ?? 'NN',
        tipoActivo: response.nombre ?? 'NN',
        estado: 'NN', // Por ahora no viene
        descripcion: 'NN', // Por ahora no viene
        ubicacion: response.ubicacion ?? 'NN',
        img: 'https://www.sondagua.cl/blog/wp-content/uploads/2021/10/bomba-para-extraccion-de-agua.jpg' // temporal
      }

      setActivo(activoTransformado)
      setSensores(response.sensores ?? [])

    } catch (error) {
      console.error('Error al obtener el activo', error)
    }
  }
  
  // Actualización del método a penas se recarga la página
  const hasFetchedRef = React.useRef(false)
  
  React.useEffect(() => {
    if (!hasFetchedRef.current) {
      hasFetchedRef.current = true
      getDetalleActivo()
    }
  }, [])

  // Estos nombres tendrían que ser dinámicos, de momento quedarán así.
  // Extracción de los valores de cada sensor:
  const sensorTemp = sensores.find(s => s.sensor_id === 'temp1')
  const sensorPres = sensores.find(s => s.sensor_id === 'pres1')
  const sensorCaud = sensores.find(s => s.sensor_id === 'caud1')

  // Obtención de valores actuales
  const temperatura = sensorTemp?.datos?.at(-1)?.valor ?? 0

  const getDiffInfo = (sensorId: string) => {
    const datos = sensores.find(s => s.sensor_id === sensorId)?.datos ?? []
    const ultimo = datos.at(-1)?.valor ?? 0
    const penultimo = datos.at(-2)?.valor ?? 0
    const diff = ultimo - penultimo
    const trend = (diff >= 0 ? 'up' : 'down') as 'up' | 'down'

    return {
      valor: ultimo,
      diff: Math.abs(diff),
      trend
    }
  }

  const presionInfo = getDiffInfo('pres1')
  const caudalInfo = getDiffInfo('caud1')

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
            {/* Caudal */}
            <Budget diff={parseFloat(caudalInfo.diff.toFixed(2))}
              trend={caudalInfo.trend} 
              sx={{ height: 192 }} 
              value={`${caudalInfo.valor.toFixed(2)} m³/h`}
            />
            {/* Presión */}
            <TotalCustomers diff={parseFloat(presionInfo.diff.toFixed(2))} trend={presionInfo.trend} sx={{ height: 192}} value={`${presionInfo.valor.toFixed(2)} Psi`} />
          </Box>
        </Grid>
      </Grid>

      {/* FILA INFERIOR */}
      <Grid container spacing={2} sx={{ mt: 2 }}>
        {/* Scatter: 70% */}
        <Grid size={{md:10, xs:12}}>
          <ScatterWithArgs sx={{ height: 480 }} dataCaudal={sensorCaud} dataPresion={sensorPres} dataTemp={sensorTemp} />
        </Grid>
        {/* ScoreChart: 30% */}
        <Grid size={{md:2, xs:12}}>
          {/* Temperatura */}
          <TemperatureProgress value={temperatura} /> {/* 🔥 Aquí lo usas */}
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
