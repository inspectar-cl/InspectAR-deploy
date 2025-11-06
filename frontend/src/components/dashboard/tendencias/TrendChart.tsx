'use client';

import { Box, Typography, ToggleButton, ToggleButtonGroup } from '@mui/material';
import { useState } from 'react';
import type React from 'react';
import {
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Area,
  AreaChart
} from 'recharts';
import type { Variable } from '@/types/tendencias';

interface TrendChartProps {
  variable: Variable;
}

export default function TrendChart({ variable }: TrendChartProps) {
  const [timeRange, setTimeRange] = useState<'1h' | '6h' | '24h' | '7d'>('24h');

  const handleTimeRangeChange = (_event: React.MouseEvent<HTMLElement>, newRange: '1h' | '6h' | '24h' | '7d' | null) => {
    if (newRange !== null) {
      setTimeRange(newRange);
    }
  };

  const getFilteredData = () => {
    const now = new Date();
    const ranges = {
      '1h': 1,
      '6h': 6,
      '24h': 24,
      '7d': 168
    };
    
    const hoursToFilter = ranges[timeRange];
    const cutoffTime = new Date(now.getTime() - hoursToFilter * 60 * 60 * 1000);
    
    return variable.data.filter(point => new Date(point.timestamp) >= cutoffTime);
  };

  const filteredData = getFilteredData();

  const formatXAxis = (timestamp: string) => {
    const date = new Date(timestamp);
    if (timeRange === '7d') {
      return date.toLocaleDateString('es-ES', { day: '2-digit', month: '2-digit' });
    }
    return date.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <Box
      sx={{
        p: 3,
        borderRadius: 1,
        boxShadow: 'var(--mui-shadows-1)',
        backgroundColor: 'var(--mui-palette-background-paper)'
      }}
    >
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box>
          <Typography variant="h6" fontWeight="bold">
            Tendencia de {variable.nombre}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Análisis predictivo con Machine Learning
          </Typography>
        </Box>
        <ToggleButtonGroup
          value={timeRange}
          exclusive
          onChange={handleTimeRangeChange}
          size="small"
        >
          <ToggleButton value="1h">1H</ToggleButton>
          <ToggleButton value="6h">6H</ToggleButton>
          <ToggleButton value="24h">24H</ToggleButton>
          <ToggleButton value="7d">7D</ToggleButton>
        </ToggleButtonGroup>
      </Box>

      <ResponsiveContainer width="100%" height={400}>
        <AreaChart data={filteredData}>
          <defs>
            <linearGradient id="colorValor" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#2196f3" stopOpacity={0.3}/>
              <stop offset="95%" stopColor="#2196f3" stopOpacity={0}/>
            </linearGradient>
            <linearGradient id="colorPrediccion" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#4caf50" stopOpacity={0.2}/>
              <stop offset="95%" stopColor="#4caf50" stopOpacity={0}/>
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--mui-palette-divider)" />
          <XAxis 
            dataKey="timestamp" 
            tickFormatter={formatXAxis}
            stroke="var(--mui-palette-text-secondary)"
            style={{ fontSize: '12px' }}
          />
          <YAxis 
            label={{ value: variable.unidad, angle: -90, position: 'insideLeft' }}
            stroke="var(--mui-palette-text-secondary)"
            style={{ fontSize: '12px' }}
          />
          <Tooltip 
            formatter={(value: number, name: string) => {
              const nombres: Record<string, string> = {
                'valor': 'Valor Real',
                'prediccion': 'Predicción IA',
                'limite_superior': 'Límite Superior',
                'limite_inferior': 'Límite Inferior'
              };
              return [`${value.toFixed(2)} ${variable.unidad}`, nombres[name] || name];
            }}
            labelFormatter={(label: string | number) => new Date(label).toLocaleString('es-ES')}
            contentStyle={{ 
              backgroundColor: 'var(--mui-palette-background-paper)', 
              borderRadius: '8px',
              border: '1px solid var(--mui-palette-divider)'
            }}
          />
          <Legend />
          
          {/* Área entre límites */}
          <Area
            type="monotone"
            dataKey="limite_superior"
            stroke="transparent"
            fill="#ff9800"
            fillOpacity={0.1}
            name="Rango de Seguridad"
          />
          
          {/* Límites de predicción */}
          <Line
            type="monotone"
            dataKey="limite_superior"
            stroke="#ff9800"
            strokeDasharray="5 5"
            dot={false}
            name="Límite Superior"
            strokeWidth={2}
          />
          <Line
            type="monotone"
            dataKey="limite_inferior"
            stroke="#ff9800"
            strokeDasharray="5 5"
            dot={false}
            name="Límite Inferior"
            strokeWidth={2}
          />
          
          {/* Valor real */}
          <Line
            type="monotone"
            dataKey="valor"
            stroke="#2196f3"
            strokeWidth={3}
            dot={{ fill: '#2196f3', r: 4 }}
            activeDot={{ r: 6 }}
            name="Valor Real"
          />
          
          {/* Predicción */}
          <Line
            type="monotone"
            dataKey="prediccion"
            stroke="#4caf50"
            strokeWidth={3}
            strokeDasharray="5 5"
            dot={{ fill: '#4caf50', r: 4 }}
            name="Predicción IA"
          />
        </AreaChart>
      </ResponsiveContainer>

      {/* Estadísticas adicionales */}
      <Box display="flex" gap={3} mt={2} justifyContent="center">
        <Box textAlign="center">
          <Typography variant="body2" color="text.secondary">Último Valor</Typography>
          <Typography variant="h6" fontWeight="bold">
            {filteredData[filteredData.length - 1]?.valor.toFixed(2)} {variable.unidad}
          </Typography>
        </Box>
        <Box textAlign="center">
          <Typography variant="body2" color="text.secondary">Predicción</Typography>
          <Typography variant="h6" fontWeight="bold" color="#4caf50">
            {filteredData[filteredData.length - 1]?.prediccion?.toFixed(2)} {variable.unidad}
          </Typography>
        </Box>
        <Box textAlign="center">
          <Typography variant="body2" color="text.secondary">Tendencia</Typography>
          <Typography variant="h6" fontWeight="bold">
            {variable.cambio > 0 ? '+' : ''}{variable.cambio.toFixed(1)}%
          </Typography>
        </Box>
      </Box>
    </Box>
  );
}