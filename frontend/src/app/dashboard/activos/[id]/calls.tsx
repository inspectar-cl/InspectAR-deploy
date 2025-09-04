'use client'

import * as React from 'react';
import { useEffect } from 'react';
import Card from '@mui/material/Card';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box'
import dayjs from 'dayjs';
import { Activo } from '@/types/'
import { DownloadSimple } from '@phosphor-icons/react';

import ScoreChart from '@/components/dashboard/overview/score-chart';
import { Caudal } from '@/components/dashboard/overview/caudal';
import { Presion } from '@/components/dashboard/overview/presion';
import { TemperatureProgress } from '@/components/dashboard/overview/temperature';
import { LatestAlerts } from '@/components/dashboard/overview/latest-alerts';
import { WarningIcon } from '@phosphor-icons/react/dist/ssr/Warning';
import Grid from '@mui/material/Grid';
import { ScatterWithArgs} from '@/components/dashboard/overview/medicionTiempoReal'
import { ChatBotCard } from '@/components/dashboard/overview/chatbot'

// Configuración rutas de obtención de datos desde db.
import Services from '@/modules/Services'

import dataAlertas from '@/mocks/alerts.json'

const gs = new Services()

export default function ActivoDetailClient({ id }: { id: string}) {

  const [activo, setActivo] = React.useState<Activo | null>(null)
  const [sensores, setSensores] = React.useState<any[]>([])
  const [documentId, setDocumentId] = React.useState<string | null>(null);

  // Actualización del método a penas se recarga la página
  const hasFetchedRef = React.useRef(false)
  
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await gs.get(`parser/lectura/${id}/datos`);
        console.log("response", response)

        // Aqui deberian de cargarse la data de los activos (Ojala desde una llamada a API)
        const activoTransformado: Activo = {
          id: response.activo_id ?? 'NN',
          tipoActivo: response.nombre ?? 'NN',
          estado: response.estado ?? 'NN',
          descripcion: 'NN',
          ubicacion: response.ubicacion ?? 'NN',
          img: 'https://www.sondagua.cl/blog/wp-content/uploads/2021/10/bomba-para-extraccion-de-agua.jpg',
          id_edificio: 'NN',
        };

        setActivo(activoTransformado);
        setSensores(response.sensores ?? []);
        
        // Aqui hago llamado a API para obtener id de documento del activo
        const docResponse = await gs.get(`/documentos/activo/${id}`);
        console.log("doc Response: ", docResponse)
        setDocumentId(docResponse.documento_id ?? 'none00'); //none00 como id vacia
      } catch (err) {
        console.error('Error al obtener el activo', err);
      }
    };

    // Llamado inicial inmediato
    fetchData();

    // Intervalo de actualización cada 5 segundos
    const interval = setInterval(() => {
      fetchData();
    }, 5000);

    // Limpieza del intervalo al desmontar componente
    return () => clearInterval(interval);
  }, [id]);

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
        {/* Columna izquierda: Caudal y Presion */}


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
                <Button
                  size="small"
                  variant="outlined"
                  startIcon={<DownloadSimple size={18} />}
                  component="a"
                  href="/archivos/ficha_tecnica_bomba.pdf"
                  download="ficha_tecnica_bomba.pdf"
                  target="_blank"
                  rel="noopener"
                  sx={{
                    textTransform: 'none',
                    borderRadius: 2,
                    px: 1.5
                  }}
                  aria-label="Descargar ficha técnica en PDF"
                >
                  Ficha técnica
                </Button>
              </CardActions>
            </Box>
          </Card>
        </Grid>
        <Grid size={{md:4, xs:12}}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {/* Caudal */}
            <Caudal diff={parseFloat(caudalInfo.diff.toFixed(2))}
              trend={caudalInfo.trend} 
              sx={{ height: 192 }} 
              value={`${caudalInfo.valor.toFixed(2)} m³/h`}
            />
            {/* Presión */}
            <Presion diff={parseFloat(presionInfo.diff.toFixed(2))} trend={presionInfo.trend} sx={{ height: 192}} value={`${presionInfo.valor.toFixed(2)} Psi`} />
          </Box>
        </Grid>
      </Grid>

      {/* FILA INFERIOR */}
      <Grid container spacing={2} sx={{ mt: 2 }}>
        {/* Scatter: 70% */}
        <Grid size={{md:10, xs:12}}>
          <ScatterWithArgs sx={{ height: 500 }} dataCaudal={sensorCaud} dataPresion={sensorPres} dataTemp={sensorTemp} />
        </Grid>
        {/* ScoreChart: 30% */}
        <Grid size={{md:2, xs:12}}>
          {/* Temperatura */}
          <TemperatureProgress value={temperatura} />
        </Grid>
        <Grid size={{md:3, xs:12}}>
          <ScoreChart sx={{ height: 450 }} />
        </Grid>
      
        <Grid size={{lg:9, md:6, xs:12}}>
          <LatestAlerts

            /* Aqui hay que modificar el como llegan las alertas */
            products={
              dataAlertas.alerts.map(alerta => ({
                id: alerta.id,
                name: alerta.name,
                icon: <WarningIcon size={32} weight="fill" color="#ff0000" />,
                updatedAt: dayjs(alerta.updatedAt).toDate()
              }))
            }
            sx={{ height: 450 }}
          />
        </Grid>
        <Grid size={{md:4, xs:12}}>
          <ChatBotCard id={documentId}/>
        </Grid>
      </Grid>
    </Box>
  );
}