import { Box, Typography, Chip } from '@mui/material';
import { TrendingUp as TrendingUpIcon, TrendingDown as TrendingDownIcon, Minus as MinusIcon } from 'lucide-react';
import type { Variable } from '@/types/tendencias';

interface VariableCardProps {
  variable: Variable;
  onClick?: (variableId: number) => void;
}

export default function VariableCard({ variable, onClick }: VariableCardProps) {
  const getEstadoColor = (estado: string) => {
    switch (estado) {
      case 'normal':
        return '#4caf50';
      case 'advertencia':
        return '#ff9800';
      case 'critico':
        return '#f44336';
      default:
        return '#2196f3';
    }
  };

  const getTendenciaIcon = () => {
    if (variable.tendencia === 'ascendente') {
      return <TrendingUpIcon size={20} color="#4caf50" />;
    }
    if (variable.tendencia === 'descendente') {
      return <TrendingDownIcon size={20} color="#f44336" />;
    }
    return <MinusIcon size={20} color="#9e9e9e" />;
  };

  const getTendenciaTexto = () => {
    const signo = variable.cambio > 0 ? '+' : '';
    return `${signo}${variable.cambio.toFixed(1)}% vs última hora`;
  };

  return (
    <Box
      onClick={() => { onClick?.(variable.id); }}
      sx={{
        p: 2,
        borderRadius: 2,
        boxShadow: 1,
        display: 'flex',
        flexDirection: 'column',
        gap: 2,
        backgroundColor: 'white',
        border: `2px solid ${getEstadoColor(variable.estado)}`,
        cursor: onClick ? 'pointer' : 'default',
        transition: 'transform 0.2s',
        '&:hover': {
          transform: onClick ? 'translateY(-4px)' : 'none',
          boxShadow: onClick ? 3 : 1
        }
      }}
    >
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Typography variant="h6" fontWeight="bold">
          {variable.nombre}
        </Typography>
        <Chip 
          label={variable.estado.toUpperCase()} 
          size="small"
          sx={{ 
            backgroundColor: getEstadoColor(variable.estado),
            color: 'white',
            fontWeight: 'bold'
          }}
        />
      </Box>

      <Box display="flex" alignItems="baseline" gap={1}>
        <Typography variant="h3" fontWeight="bold" color={getEstadoColor(variable.estado)}>
          {variable.valor.toFixed(2)}
        </Typography>
        <Typography variant="body1" color="text.secondary">
          {variable.unidad}
        </Typography>
      </Box>

      <Box display="flex" alignItems="center" gap={1}>
        {getTendenciaIcon()}
        <Typography 
          variant="body2" 
          color={
            variable.tendencia === 'ascendente' 
              ? '#4caf50' 
              : variable.tendencia === 'descendente' 
                ? '#f44336' 
                : 'text.secondary'
          }
        >
          {getTendenciaTexto()}
        </Typography>
      </Box>
    </Box>
  );
}