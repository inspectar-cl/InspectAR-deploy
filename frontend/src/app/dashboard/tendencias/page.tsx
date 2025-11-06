'use client';

import { useState, useMemo } from 'react';
import {
  Box,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Divider,
  Chip
} from '@mui/material';
import Grid from '@mui/material/Grid';
import { 
  Activity as ActivityIcon, 
  TrendingUp as TrendingUpIcon,
  AlertTriangle as AlertIcon,
  Gauge as GaugeIcon
} from 'lucide-react';
import MetricCard from '@/components/dashboard/tendencias/MetricCard';
import TendenciasCards from '@/components/dashboard/tendencias/TendenciasCards';
import TrendChart from '@/components/dashboard/tendencias/TrendChart';
import type { Activo } from '@/types/activos';
import type { Variable, TendenciaDataPoint } from '@/types/tendencias';

// Función para generar datos de tendencia simulados
const generarDatosTendencia = (
  baseValue: number, 
  variabilidad: number, 
  puntos = 288
): TendenciaDataPoint[] => {
  const now = new Date();
  const data: TendenciaDataPoint[] = [];
  
  for (let i = puntos; i >= 0; i--) {
    const timestamp = new Date(now.getTime() - i * 5 * 60 * 1000); // Cada 5 minutos
    const ruido = (Math.random() - 0.5) * variabilidad;
    const tendencia = -i * 0.01; // Tendencia suave
    const valor = baseValue + ruido + tendencia;
    const prediccion = valor + (Math.random() - 0.5) * (variabilidad * 0.5);
    
    data.push({
      timestamp: timestamp.toISOString(),
      valor: parseFloat(valor.toFixed(2)),
      prediccion: parseFloat(prediccion.toFixed(2)),
      limite_superior: parseFloat((baseValue + variabilidad * 1.5).toFixed(2)),
      limite_inferior: parseFloat((baseValue - variabilidad * 1.5).toFixed(2))
    });
  }
  
  return data;
};

// Datos simulados de activos
const activosSimulados: Activo[] = [
  {
    id: 1,
    tipoActivo: 'Bomba de Agua',
    estado: 'Medio',
    descripcion: 'Bomba Centrífuga A',
    ubicacion: 'Sala de Bombas',
    img: '/activos/bomba.jpg',
    id_edificio: 'ED-001',
    id_ficha_tecnica: 1
  },
  {
    id: 2,
    tipoActivo: 'Transformador',
    estado: 'OK',
    descripcion: 'Transformador Principal',
    ubicacion: 'Subestación Eléctrica',
    img: '/activos/transformador.jpg',
    id_edificio: 'ED-001',
    id_ficha_tecnica: 2
  },
  {
    id: 3,
    tipoActivo: 'Ascensor',
    estado: 'Crítico',
    descripcion: 'Ascensor Central',
    ubicacion: 'Edificio Central - Hall',
    img: '/activos/ascensor.jpg',
    id_edificio: 'ED-002',
    id_ficha_tecnica: 3
  },
  {
    id: 4,
    tipoActivo: 'Bomba de Agua',
    estado: 'OK',
    descripcion: 'Bomba de Emergencia',
    ubicacion: 'Planta de Emergencia',
    img: '/activos/bomba.jpg',
    id_edificio: 'ED-001',
    id_ficha_tecnica: 4
  }
];

// Generar variables simuladas para cada activo
const generarVariablesPorActivo = (activoId: number): Variable[] => {
  const configuraciones = [
    { nombre: 'Temperatura', baseValue: 75, variabilidad: 8, unidad: '°C', severidad: 'advertencia' },
    { nombre: 'Vibración', baseValue: 3.8, variabilidad: 0.8, unidad: 'mm/s', severidad: 'normal' },
    { nombre: 'Presión', baseValue: 102.5, variabilidad: 5, unidad: 'PSI', severidad: 'normal' },
    { nombre: 'RPM', baseValue: 1850, variabilidad: 80, unidad: 'RPM', severidad: activoId === 3 ? 'critico' : 'normal' },
    { nombre: 'Consumo Energético', baseValue: 12.4, variabilidad: 2, unidad: 'kW', severidad: 'advertencia' },
    { nombre: 'Humedad', baseValue: 65, variabilidad: 8, unidad: '%', severidad: 'normal' }
  ];

  return configuraciones.map((config, index) => {
    const cambio = (Math.random() - 0.5) * 10;
    const valorActual = config.baseValue + (Math.random() - 0.5) * config.variabilidad;
    
    let tendencia: 'ascendente' | 'descendente' | 'estable';
    if (Math.abs(cambio) < 1) {
      tendencia = 'estable';
    } else if (cambio > 0) {
      tendencia = 'ascendente';
    } else {
      tendencia = 'descendente';
    }

    let estado: 'normal' | 'advertencia' | 'critico';
    if (config.severidad === 'critico') {
      estado = 'critico';
    } else if (config.severidad === 'advertencia' || Math.abs(cambio) > 5) {
      estado = 'advertencia';
    } else {
      estado = 'normal';
    }

    return {
      id: activoId * 10 + index + 1,
      nombre: config.nombre,
      valor: parseFloat(valorActual.toFixed(2)),
      unidad: config.unidad,
      tendencia,
      cambio: parseFloat(cambio.toFixed(1)),
      estado,
      activoId,
      data: generarDatosTendencia(config.baseValue, config.variabilidad)
    };
  });
};

export default function TendenciasPage() {
  const [activoSeleccionado, setActivoSeleccionado] = useState<number>(1);
  const [variableSeleccionada, setVariableSeleccionada] = useState<number>(11);

  // Generar variables para todos los activos (memoizado)
  const todasLasVariables = useMemo(() => {
    return activosSimulados.flatMap(activo => generarVariablesPorActivo(activo.id));
  }, []);

  // Filtrar variables del activo seleccionado
  const variablesActivo = useMemo(() => {
    return todasLasVariables.filter(v => v.activoId === activoSeleccionado);
  }, [todasLasVariables, activoSeleccionado]);

  const activoActual = activosSimulados.find(a => a.id === activoSeleccionado);
  const variableActual = variablesActivo.find(v => v.id === variableSeleccionada);

  // Calcular métricas
  const totalVariables = variablesActivo.length;
  const variablesEnAdvertencia = variablesActivo.filter(v => v.estado === 'advertencia').length;
  const variablesEnCritico = variablesActivo.filter(v => v.estado === 'critico').length;
  const promedioTendencia = variablesActivo.reduce((acc, v) => acc + Math.abs(v.cambio), 0) / totalVariables;

  const handleActivoChange = (nuevoActivoId: number) => {
    setActivoSeleccionado(nuevoActivoId);
    // Seleccionar la primera variable del nuevo activo
    const primeraVariable = todasLasVariables.find(v => v.activoId === nuevoActivoId);
    if (primeraVariable) {
      setVariableSeleccionada(primeraVariable.id);
    }
  };

  const handleVariableClick = (variableId: number) => {
    setVariableSeleccionada(variableId);
  };

  return (
    <Box p={3}>
      {/* Header */}
      <Box display="flex" alignItems="center" gap={2} mb={2}>
        <TrendingUpIcon size={32} color="#2196f3" />
        <Typography variant="h4" fontWeight="bold">
          Tendencias y Predicciones
        </Typography>
      </Box>

      <Typography variant="body1" color="text.secondary" mb={3}>
        Monitoreo en tiempo real con predicciones basadas en Machine Learning
      </Typography>

      <Divider sx={{ mb: 3 }} />

      {/* Filtros */}
      <Grid container spacing={2} mb={3}>
        <Grid size={{ xs: 12, md: 6 }}>
          <FormControl fullWidth>
            <InputLabel>Activo</InputLabel>
            <Select
              value={activoSeleccionado}
              label="Activo"
              onChange={(e) => { handleActivoChange(Number(e.target.value)); }}
            >
              {activosSimulados.map((activo) => (
                <MenuItem key={activo.id} value={activo.id}>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Chip 
                      label={activo.estado} 
                      size="small"
                      color={
                        activo.estado === 'OK' ? 'success' : 
                        activo.estado === 'Medio' ? 'warning' : 
                        'error'
                      }
                    />
                    {activo.descripcion}
                  </Box>
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid size={{ xs: 12, md: 6 }}>
          <FormControl fullWidth>
            <InputLabel>Variable a Analizar</InputLabel>
            <Select
              value={variableSeleccionada}
              label="Variable a Analizar"
              onChange={(e) => { setVariableSeleccionada(Number(e.target.value)); }}
            >
              {variablesActivo.map((variable) => (
                <MenuItem key={variable.id} value={variable.id}>
                  <Box display="flex" alignItems="center" gap={1}>
                    <Chip 
                      label={variable.estado} 
                      size="small"
                      color={
                        variable.estado === 'normal' ? 'success' :
                        variable.estado === 'advertencia' ? 'warning' :
                        'error'
                      }
                    />
                    {variable.nombre} ({variable.unidad})
                  </Box>
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
      </Grid>

      {/* Métricas principales */}
      <Grid container spacing={3} mb={3}>
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <MetricCard
            title="Variables Monitoreadas"
            value={totalVariables}
            color="var(--mui-palette-primary-main)"
            icon={<GaugeIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <MetricCard
            title="En Advertencia"
            value={variablesEnAdvertencia}
            color="var(--mui-palette-warning-main)"
            icon={<AlertIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <MetricCard
            title="Estado Crítico"
            value={variablesEnCritico}
            color="var(--mui-palette-error-main)"
            icon={<AlertIcon />}
            suffix=""
          />
        </Grid>
        
        <Grid size={{ xs: 12, sm: 6, md: 3 }}>
          <MetricCard
            title="Cambio Promedio"
            value={promedioTendencia.toFixed(1)}
            color="var(--mui-palette-success-main)"
            icon={<ActivityIcon />}
            suffix="%"
            trend={{ value: 1.2, isPositive: promedioTendencia > 0 }}
          />
        </Grid>
      </Grid>

      {/* Información del activo actual */}
      {activoActual && (
        <Box 
          sx={{ 
            p: 2, 
            mb: 3, 
            backgroundColor: 'var(--mui-palette-background-paper)',
            borderRadius: 1,
            boxShadow: 'var(--mui-shadows-1)'
          }}
        >
          <Typography variant="h6" fontWeight="bold" mb={1}>
            {activoActual.descripcion}
          </Typography>
          <Box display="flex" gap={2} flexWrap="wrap">
            <Typography variant="body2" color="text.secondary">
              📍 {activoActual.ubicacion}
            </Typography>
          </Box>
        </Box>
      )}

      {/* Gráfico de tendencias */}
      {variableActual && (
        <Box mb={3}>
          <TrendChart variable={variableActual} />
        </Box>
      )}

      {/* Tarjetas de variables */}
      <Box>
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
          <Typography variant="h6" fontWeight="bold">
            Estado de Variables - {activoActual?.descripcion}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Última actualización: {new Date().toLocaleTimeString('es-ES')}
          </Typography>
        </Box>
        <TendenciasCards 
          variables={variablesActivo}
          onVariableClick={handleVariableClick}
        />
      </Box>
    </Box>
  );
}