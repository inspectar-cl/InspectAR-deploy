import { Card, CardContent, Typography, Box, Chip } from '@mui/material';
import Grid from '@mui/material/Grid';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import type { Variable } from '@/types/tendencias';

interface TendenciasCardsProps {
  variables: Variable[];
  onVariableClick?: (variableId: number) => void;
}

export default function TendenciasCards({ variables, onVariableClick }: TendenciasCardsProps) {
  const getTrendIcon = (tendencia: 'ascendente' | 'descendente' | 'estable') => {
    switch (tendencia) {
      case 'ascendente':
        return <TrendingUp size={20} color="var(--mui-palette-success-main)" />;
      case 'descendente':
        return <TrendingDown size={20} color="var(--mui-palette-error-main)" />;
      case 'estable':
        return <Minus size={20} color="var(--mui-palette-text-secondary)" />;
    }
  };

  return (
    <Grid container spacing={2}>
      {variables.map((variable) => (
        <Grid key={variable.id} size={{ xs: 12, sm: 6, md: 4 }}>
          <Card
            sx={{
              height: '100%',
              cursor: 'pointer',
              transition: 'all 0.2s',
              borderLeft: `4px solid ${
                variable.estado === 'normal'
                  ? 'var(--mui-palette-success-main)'
                  : variable.estado === 'advertencia'
                  ? 'var(--mui-palette-warning-main)'
                  : 'var(--mui-palette-error-main)'
              }`,
              '&:hover': {
                transform: 'translateY(-4px)',
                boxShadow: 'var(--mui-shadows-8)'
              }
            }}
            onClick={() => { onVariableClick?.(variable.id); }}
          >
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={2}>
                <Typography variant="h6" fontWeight="bold">
                  {variable.nombre}
                </Typography>
                <Chip
                  label={variable.estado}
                  size="small"
                  color={
                    variable.estado === 'normal'
                      ? 'success'
                      : variable.estado === 'advertencia'
                      ? 'warning'
                      : 'error'
                  }
                />
              </Box>

              <Typography variant="h4" fontWeight="bold" mb={1}>
                {variable.valor} {variable.unidad}
              </Typography>

              <Box display="flex" alignItems="center" gap={1}>
                {getTrendIcon(variable.tendencia)}
                <Typography
                  variant="body2"
                  sx={{
                    color:
                      variable.cambio > 0
                        ? 'var(--mui-palette-success-main)'
                        : variable.cambio < 0
                        ? 'var(--mui-palette-error-main)'
                        : 'var(--mui-palette-text-secondary)',
                    fontWeight: 500
                  }}
                >
                  {variable.cambio > 0 ? '+' : ''}
                  {variable.cambio}% en las últimas 24h
                </Typography>
              </Box>

              <Typography variant="caption" color="text.secondary" mt={1} display="block">
                Tendencia: {variable.tendencia}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
}