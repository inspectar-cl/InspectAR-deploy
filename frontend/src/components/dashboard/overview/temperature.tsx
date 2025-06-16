import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import LinearProgress from '@mui/material/LinearProgress';
import Stack from '@mui/material/Stack';
import type { SxProps } from '@mui/material/styles';
import Typography from '@mui/material/Typography';
import { Thermometer as ThermometerIcon } from '@phosphor-icons/react/dist/ssr/Thermometer';
import Box from '@mui/material/Box'

export interface TemperatureProgressProps {
  sx?: SxProps;
  value: number;
}

export function TemperatureProgress({ value, sx }: TemperatureProgressProps): React.JSX.Element {
  return (
    <Card
      sx={{
        ...sx,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'flex-start',
        alignItems: 'center',
        height: 480,
        width: 210,
        paddingTop: 2,
        paddingBottom: 2,
      }}
    >
      <Avatar
        sx={{
          backgroundColor: 'var(--mui-palette-error-main)',
          height: 56,
          width: 56,
        }}
      >
        <ThermometerIcon fontSize="var(--icon-fontSize-lg)" />
      </Avatar>

      <Typography variant="subtitle2" sx={{ mt: 1, color: 'text.secondary' }}>
        Temperatura
      </Typography>

      <Typography variant="h4" sx={{ mb: 2 }}>
        {value}°C
      </Typography>

      {/* Espaciador que empuja la barra hacia abajo */}
      <Box sx={{ flexGrow: 1 }} />

      <Box sx={{ display: 'flex', justifyContent: 'center', mb: 20 }}>
        <LinearProgress
          value={value}
          variant="determinate"
          sx={{
            transform: 'rotate(-90deg)',
            height: 10,
            width: 300,
            borderRadius: 5,
            backgroundColor: 'var(--mui-palette-neutral-300)',
            '& .MuiLinearProgress-bar': {
              backgroundColor: 'var(--mui-palette-error-main)',
            },
          }}
        />
      </Box>
    </Card>
  );
}
