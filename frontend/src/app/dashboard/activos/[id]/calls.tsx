/* eslint-disable @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/explicit-function-return-type -- Complex data operations require flexibility */
'use client';

import * as React from 'react';
import { useEffect, useState } from 'react';
import Card from '@mui/material/Card';
import CardMedia from '@mui/material/CardMedia';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import Alert from '@mui/material/Alert';
import Select from '@mui/material/Select'
import MenuItem from '@mui/material/MenuItem'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from 'recharts';
import dayjs from 'dayjs';

import { type Activo, type Prediccion } from '@/types/';
import { DownloadSimple } from '@phosphor-icons/react';
import { useUserToken } from '@/hooks/use-usertoken';
import ScoreChart from '@/components/dashboard/overview/score-chart';
import { Caudal } from '@/components/dashboard/overview/caudal';
import { Presion } from '@/components/dashboard/overview/presion';
import { TemperatureProgress } from '@/components/dashboard/overview/temperature';
import { LatestAlerts } from '@/components/dashboard/overview/latest-alerts';
import { WarningIcon } from '@phosphor-icons/react/dist/ssr/Warning';
import { ScatterWithArgs } from '@/components/dashboard/overview/MedicionTiempoReal';
import { ChatBotCard } from '@/components/dashboard/overview/chatbot';
import { useActivosWithSensors } from '@/hooks/use-activos-with-sensors';
import Services from '@/modules/Services';
import type { SensorRow } from '@/hooks/use-activos-with-sensors';
import { Cpu } from 'lucide-react';

// Constantes globales
const ESTADOS = ['OK', 'Medio', 'Crítico', 'NN'] as const;
type Estado = (typeof ESTADOS)[number];
type Trend = 'up' | 'down';
function normalizeTrend(t: unknown): Trend {
  if (t === 'up' || t === 'down') return t;
  // valor por defecto — ajusta según tu lógica
  return 'down';
}
const gs = new Services()

export default function ActivoDetailClient({ id }: { id: number}) {
  const [rangoMinutos, setRangoMinutos] = useState(30); 
  const { activos} = useActivosWithSensors();
  const { user, isLoading} = useUserToken();
  const [activo, setActivo] = React.useState<Activo | null>(null)
  const [openChart, setOpenChart] = useState(false);
  const [selectedSensor, setSelectedSensor] = useState<SensorRow | null>(null);
  const [alertas, setAlertas] = useState<Prediccion[]>([]);

  const filteredData = React.useMemo(() => {
    if (!selectedSensor?.history24h) return [];

    if (rangoMinutos === 0) {
      return selectedSensor.history24h.map((d) => ({
        time: dayjs(d.ts).format('DD/MM HH:mm'),
        value: d.value,
      }));
    }

    const limite = dayjs().subtract(rangoMinutos, 'minute');
    return selectedSensor.history24h
      .filter((d) => dayjs(d.ts).isAfter(limite))
      .map((d) => ({
        time: dayjs(d.ts).format('HH:mm:ss'),
        value: d.value,
      }));
  }, [selectedSensor, rangoMinutos]);

  interface SensorDato {
    valor: number;
    tiempo: string;
    [key: string]: unknown;
  }

  interface Sensor {
    sensor_id: string;
    datos: SensorDato[];
    [key: string]: unknown;
  }

  const [sensores, setSensores] = React.useState<Sensor[]>([])

  // Actualización del método a penas se recarga la página
  //const hasFetchedRef = React.useRef(false)
  
  useEffect(() => {
    //ta raro esto, a veces funciona con este codigo o a veces no
    //quitar el if y volver a colocarlo si no funca
    if (isLoading || !user) {return;}
    
    const fetchData = async () => {
      try {
        interface ActivoResponse {
          activo_id?: number;
          nombre?: string;
          estado?: string;
          ubicacion?: string;
          img?: string;
          id_edificio?: string;
          id_ficha_tecnica?: number;
          sensores?: Sensor[];
        }
        // falta solucion parche pal user
        const response = await gs.authorizedGet(`/obtener-activo-id/${id}`, user.token) as ActivoResponse;

        // Aqui deberian de cargarse la data de los activos (Ojala desde una llamada a API)
        const estado: Estado = ESTADOS.includes(response.estado as Estado)
          ? (response.estado as Estado)
          : 'NN';

        const activoTransformado: Activo = {
          id: response.activo_id ?? 0,
          tipoActivo: response.nombre ?? 'NN',
          estado,
          descripcion: 'Descripción personalizada :)',
          ubicacion: response.ubicacion ?? 'NN',
          img: 'https://www.sondagua.cl/blog/wp-content/uploads/2021/10/bomba-para-extraccion-de-agua.jpg',
          id_edificio: response.id_edificio ?? 'NN',
          id_ficha_tecnica: response.id_ficha_tecnica ?? 0,
        };

        setActivo(activoTransformado);
        // setSensores(response.sensores ?? []);
      } catch (err) {
        //console.error('Error al obtener el activo', err);
      }
    };
    
    // Nueva función para obtener datos de sensores
    const fetchDatos = async () => {
      try {
        const datos = await gs.get('/parser/lectura/1/datos') as { sensores?: Sensor[] };
        setSensores(datos.sensores ?? []);
      } catch (err) {
        //console.error('Error al obtener los datos de sensores', err);
      } 
    };

    const fetchAlertas = async () => {
      try {
        
        interface AnomaliasPaginadas {
          activo_id?: number;
          data?: {
            id?: number;
            activo_id?: number;
            timestamp?: string;
            anomaly_score?: number;
            anomaly_likelihood?: number;
            severidad?: "Baja" | "Media" | "Alta";
            descripcion?: string;
            threshold?: number;
            is_anomaly?: number | boolean;
            most_influential_variable?: string;
            contribution_magnitude?: number;
          }[];
          count?: number;
          message?: string;
          limit?: number;
          offset?: number;
        }

        const response = await gs.authorizedGet(
          `/anomalia-activo/${id}`, 
          user.token
        ) as AnomaliasPaginadas;


        // Transformar y filtrar solo las anomalías
        const alertasTransformadas = (response.data ?? [])
          .map((p) => ({
            id: p.id ?? 0,
            activoId: p.activo_id ?? id,
            timestamp: p.timestamp ?? new Date().toISOString(),
            anomalyScore: Math.round(p.anomaly_score ?? 0),
            anomalyLikelihood: Math.round(p.anomaly_likelihood ?? 0),
            severidad: p.severidad ?? "Baja",
            descripcion: p.descripcion ?? "Sin descripción",
            threshold: Math.round(p.threshold ?? 50),
            is_anomaly: p.is_anomaly === 1 || p.is_anomaly === true,
            most_influential_variable: p.most_influential_variable ?? "N/A",
            contribution_magnitude: p.contribution_magnitude ?? 0,
          }))
          .filter(a => a.is_anomaly);


        setAlertas(alertasTransformadas);
      } catch (err) {
        setAlertas([]);
      }
    };

    const fetchAll = async () => {
      await fetchData();
      await fetchDatos();
      await fetchAlertas();
    };

    // Llamado inicial inmediato
    void fetchAll();

    // Intervalo de actualización cada 5 minutos
    const interval = setInterval(() => {
      void fetchAll();
    }, 300000); // 300,000 ms = 5 minutos

    // Limpieza del intervalo al desmontar componente
    return () => { clearInterval(interval); };
    }, [id, isLoading, user]);

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
    const trend = (diff >= 0 ? 'up' : 'down')

    return {
      valor: ultimo,
      diff: Math.abs(diff),
      trend
    }
  }

  const presionInfo = getDiffInfo('pres1')
  const caudalInfo = getDiffInfo('caud1')
  if (!activo) {
    return <div style={{ padding: '1rem' }}>No se encontró el activo con ID: {`${id}`}</div>;
  }
  if (!sensorCaud || !sensorPres || !sensorTemp) {
    return <div>Loading...</div>; // o skeleton / placeholder
}
  return (
    
    <Box sx={{ p: 2 }}>
      {/* FILA SUPERIOR */}
      <Grid container spacing={2}>
        {/* Columna izquierda: Tarjeta del activo */}
        <Grid size={{md:8, xs:12}}>
          <Card sx={{ display: 'flex', height: '100%', minHeight: 400 }}>
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
                  <strong>ID:</strong> {`${activo.id}`}
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
                  href="/documentos/ficha_tecnica_bomba.pdf"
                  download="ficha_tecnica_bomba.pdf"
                  target
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
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, height: '100%' }}>
            {/* Caudal */}
            <Caudal diff={parseFloat(caudalInfo.diff.toFixed(2))}
              trend={normalizeTrend(caudalInfo.trend)}
              sx={{ height: '50%', minHeight: 196 }}
              value={`${caudalInfo.valor.toFixed(2)} m³/h`}
            />
            {/* Presión */}
            <Presion diff={parseFloat(presionInfo.diff.toFixed(2))} trend={normalizeTrend(presionInfo.trend)} sx={{ height: '50%', minHeight: 196 }} value={`${presionInfo.valor.toFixed(2)} Psi`} />
          </Box>
        </Grid>
      </Grid>

      {/* FILA INFERIOR */}
      <Grid container spacing={2} sx={{ mt: 2 }}>
        {/* Scatter: 83.33% (10/12) */}
        <Grid size={{md:10, xs:12}}>
          <ScatterWithArgs sx={{ height: 500 }} dataCaudal={sensorCaud} dataPresion={sensorPres} dataTemp={sensorTemp} />
        </Grid>
        {/* Temperatura: 16.67% (2/12) */}
        <Grid size={{md:2, xs:12}}>
          <Box sx={{ height: 500, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <TemperatureProgress value={temperatura} />
          </Box>
        </Grid>
        
        {/* ScoreChart: 25% (3/12) */}
        <Grid size={{md:3, xs:12}}>
          <ScoreChart sx={{ height: 450 }} />
        </Grid>
        {/* LatestAlerts: 75% (9/12) */}
        <Grid size={{md:9, xs:12}}>
          <LatestAlerts
            products={alertas.map(alerta => ({
              id: alerta.id,
              name: alerta.descripcion ?? `Anomalía - Severidad: ${alerta.severidad}`,
              icon: <WarningIcon size={32} weight="fill" color="#ff0000" />,
              updatedAt: new Date(alerta.timestamp)
            }))}
            sx={{ height: 450 }}
          />
        </Grid>
        
        {/* ChatBot: 100% (12/12) */}
        <Grid size={{md:12, xs:12}}>
          <ChatBotCard id={activo.id_ficha_tecnica}/>
        </Grid>
      </Grid>
        {activos
          .filter((a) => a.id === activo.id)
          .map((act) => (
            <Box key={act.assetName} sx={{ mt: 4 }}>
              <Typography variant="h5" gutterBottom>
                Sensores del activo
              </Typography>

              {/* Si el activo no tiene sensores */}
              {(!act.sensores || act.sensores.length === 0) ? (
                <Alert severity="info">
                  No se encontraron sensores para este activo.
                </Alert>
              ) : (
                <Box
                  sx={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fill, minmax(210px, 1fr))',
                    gap: 2,
                  }}
                >
                  {/* Mapeo de los sensores asociados al activo */}
                  {act.sensores.map((s) => {
                    const hora = s.lastSeen
                      ? dayjs(s.lastSeen).format('HH:mm:ss')
                      : 'Sin datos';
                    
                    return (
                      <Card
                        key={s.id}
                        sx={{
                          display: 'flex',
                          flexDirection: 'column',
                          alignItems: 'center',
                          justifyContent: 'space-between',
                          height: 230,
                          p: 2,
                          borderRadius: 3,
                          boxShadow: 3,
                          cursor: 'pointer',
                          '&:hover': { boxShadow: 6 },
                        }}
                        onClick={() => {
                          setSelectedSensor(s);
                          setOpenChart(true);
                        }}
                      >
                        {/* Ícono del sensor */}
                        <Box
                          sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            bgcolor: 'primary.main',
                            color: 'white',
                            borderRadius: '50%',
                            width: 48,
                            height: 48,
                          }}
                        >
                          <Cpu size={24} />
                        </Box>

                        {/* Nombre del sensor */}
                        <Typography 
                          variant="subtitle2" 
                          sx={{ 
                            mt: 1, 
                            color: 'text.secondary', 
                            textAlign: 'center',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            width: '100%',
                          }}
                        >
                          {s.name}
                        </Typography>

                        {/* Valor actual */}
                        <Typography 
                          variant="h4" 
                          sx={{ 
                            mb: 2,
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            width: '100%',
                            textAlign: 'center',
                          }}
                        >
                          {s.lastValue.toFixed(1)} {s.unit}
                        </Typography>

                        <Box sx={{ flexGrow: 1 }} />

                        {/* Hora */}
                        <Typography 
                          variant="caption" 
                          sx={{ color: 'text.secondary' }}
                        >
                          {hora}
                        </Typography>
                      </Card>
                    );
                  })}
                </Box>
              )}
            </Box>
          ))}

          <Dialog
            fullWidth
            maxWidth="md"
            open={openChart}
            onClose={() => { setOpenChart(false); }}
          >
            <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              {selectedSensor
                ? `Histórico de ${selectedSensor.name}`
                : 'Cargando...'}

              {/* Selector de rango de tiempo */}
              <FormControl size="small" sx={{ minWidth: 200 }}>
                <InputLabel id="rango-label">Rango de tiempo</InputLabel>
                <Select
                  labelId="rango-label"
                  value={rangoMinutos}
                  label="Rango de tiempo"
                  onChange={(e) => {
                  setRangoMinutos(Number(e.target.value));
                }}
                >
                  <MenuItem value={0}>Todo</MenuItem>
                  <MenuItem value={3}>Últimos 3 minutos</MenuItem>
                  <MenuItem value={30}>Últimos 30 minutos</MenuItem>
                  <MenuItem value={60}>Última hora</MenuItem>
                  <MenuItem value={240}>Últimas 4 horas</MenuItem>
                  <MenuItem value={720}>Últimas 12 horas</MenuItem>
                </Select>
              </FormControl>
            </DialogTitle>

            <DialogContent dividers>
              {filteredData.length > 0 ? (
                <Box sx={{ width: '100%', height: 400 }}>
                  <ResponsiveContainer>
                    <LineChart
                      data={filteredData}
                      margin={{ top: 20, right: 30, left: 10, bottom: 20 }}
                    >
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="time" />
                      <YAxis
                        label={{
                          value: selectedSensor?.unit ?? 'Unidad',
                          angle: -90,
                          position: 'insideLeft',
                        }}
                      />
                      <Tooltip />
                      <Line
                        type="monotone"
                        dataKey="value"
                        stroke="#1976d2"
                        dot={false}
                        strokeWidth={2}
                        isAnimationActive={false}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </Box>
              ) : (
                <Alert severity="info">
                  No hay datos en el rango seleccionado.
                </Alert>
              )}
            </DialogContent>

            <DialogActions>
              <Button onClick={() => { setOpenChart(false); }}>Cerrar</Button>
            </DialogActions>
          </Dialog>
    </Box>
  );
}