/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import type { Metadata } from 'next';
import { config } from '@/config';
import ActivoDetailClient from './calls';
import AnomaliasDashboard from './anomalias-dashboard';
import { Box } from '@mui/material';

export const metadata = { title: `Activos | Dashboard | ${config.site.name}` } satisfies Metadata;

export default async function ActivoDetailPage({ params }: { params: Promise<{ id: number }> }) {
  const {id} = await params;
  
  return (
    <Box>
      {/* Dashboard principal del activo (sin modificar) */}
      <ActivoDetailClient id={id} />
      
      {/* Dashboard de anomalías ML (nuevo) */}
      <AnomaliasDashboard 
        activoId={id} 
        activoNombre={`Activo ${id}`} 
      />
    </Box>
  );
}