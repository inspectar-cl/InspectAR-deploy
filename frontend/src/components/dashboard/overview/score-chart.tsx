'use client';

import * as React from 'react';
import Stack from '@mui/material/Stack';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import { useTheme } from '@mui/material/styles';
import { Gauge } from '@mui/x-charts/Gauge';

export default function ScoreChart({ sx }: { sx?: any }) {
  const theme = useTheme();

  const h1 = theme.typography.h1;

  return (
    <Card sx={sx}>
      <CardHeader title="Puntaje de Anomalias" />
      <CardContent>
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={{ xs: 1, md: 3 }} sx={{ alignItems: 'center', justifyContent: 'center' }}>
          <Gauge
            width={300}
            height={300}
            value={75}
            startAngle={0}
            endAngle={360}
            innerRadius="80%"
            outerRadius="100%"
            fontSize={h1.fontSize}
          />
        </Stack>
      </CardContent>
    </Card>
  );
}