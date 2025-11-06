'use client';

import React from 'react';
import Typography from '@mui/material/Typography';
import { useDriverTour } from "@/components/tutorial/use-driver-tour";
import type { TourKey } from "@/components/tutorial/tour-config";
import { TourButton } from "@/components/tutorial/tour-button";
import ContentActivosPage from './contentpage';
import { Box } from '@mui/material';

export default function ActivosPage(): React.JSX.Element {
  const tourKey: TourKey = 'lista-activos';
  const { startTour } = useDriverTour(tourKey);

  return (
    <Box p={2}>
      <Box 
          display="flex" 
          justifyContent="space-between" 
          alignItems="center" 
          mb={3}
          id="tour-header-activos" 
      >
          <Typography variant="h4">
              Lista de Activos
          </Typography>

          <TourButton 
            onClick={startTour}
            tooltipTitle="Iniciar Tutorial de Activos"
            style="pulse 3s infinite"
          />
      </Box>
      <ContentActivosPage />
    </Box>
  );
}
