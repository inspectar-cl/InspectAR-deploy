'use client'

import * as React from 'react';
import { Card, CardContent, CardHeader, Typography, Chip, Divider, Box, CardActions, Button, type ChipProps } from '@mui/material';
import type { SolicitudAPI } from '@/types/form-solicitud';
import { DetalleDatosEspecificos } from './detalles-datos';
import { TIPO_TO_SLUG } from '@/utils/solicitud-utils';
import { useRouter } from 'next/navigation';

// Helper para dar color a los chips (opcional)
const getChipColor = (tipo: SolicitudAPI['tipoSolicitud']): ChipProps['color'] => {
  switch (tipo) {
    case 'Edificio': return 'primary';
    case 'Activo': return 'secondary';
    //case 'Sensor': return 'warning';
    case 'Técnico': return 'info';
    default: return 'default';
  }
};

export function SolicitudCard({ solicitud }: { solicitud: SolicitudAPI }): React.JSX.Element {
  const { asunto, tipoSolicitud, detalles, datosEspecificos, fechaCreacion, estado } = solicitud;
  const router = useRouter();

  const handleResolver = (): void => {
    const slug = TIPO_TO_SLUG[solicitud.tipoSolicitud];
    if (!slug) {
      return;
    }

    const targetPath = `/dashboard/agregar-datos/${slug}`;

    const dataString = JSON.stringify(solicitud);
    const encodedData = encodeURIComponent(dataString);

    const fullUrl = `${targetPath}?data=${encodedData}`;
    router.push(fullUrl);
  };

  return (
    <Card sx={{ height: '100%' }}>
      <CardHeader
        // Chip de Tipo de Solicitud
        action={
          <Chip label={tipoSolicitud} color={getChipColor(tipoSolicitud)} size="small" />
        }
        title={asunto}
        titleTypographyProps={{ variant: 'h6' }}
        // Sub-encabezado con fecha y estado
        subheader={
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 0.5 }}>
            <Typography variant="caption" color="text.secondary">
              {new Date(fechaCreacion).toLocaleDateString()}
            </Typography>
            <Chip label={estado} size="small" variant="outlined" />
          </Box>
        }
      />
      <CardContent sx={{ flexGrow: 1 }}>
        {detalles ? (
          <Typography variant="body2" color="text.secondary" component="p">
            {detalles}
          </Typography>
        ): null}
        
        <Divider sx={{ my: 2 }} />

        <Typography variant="subtitle2" gutterBottom>
          Datos Específicos
        </Typography>
        
        {/* Usamos el componente de detalles */}
        <DetalleDatosEspecificos tipo={tipoSolicitud} datos={datosEspecificos} />
        
      </CardContent>

      <CardActions sx={{ justifyContent: 'flex-end', p: 2 }}>
        <Button 
          variant="contained" 
          color="primary" 
          onClick={handleResolver}
          // Deshabilitar si el estado ya es "Resuelta" (opcional)
          // disabled={solicitud.estado === 'Resuelta'} 
        >
          Resolver
        </Button>
      </CardActions>

    </Card>
  );
}