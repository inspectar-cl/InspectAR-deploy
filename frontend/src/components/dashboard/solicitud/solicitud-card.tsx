'use client'

import * as React from 'react';
import { Card, CircularProgress, CardContent, CardHeader, Typography, Chip, Divider, Box, CardActions, Button, type ChipProps } from '@mui/material';
import type { SolicitudAPI, ApiTicket } from '@/types/form-solicitud';
import { DetalleDatosEspecificos } from './detalles-datos';
import { TIPO_TO_SLUG, transformarApiATipoFrontend } from '@/utils/solicitud-utils';
import { useRouter } from 'next/navigation';
import { useUserToken } from '@/hooks/use-usertoken';
import { paths } from '@/paths'; // Importa tus paths
import Services from '@/modules/Services';

const services = new Services();

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

interface SolicitudCardProps {
  solicitud: SolicitudAPI;
  onTicketUpdated: (solicitud: SolicitudAPI) => void;
}

export function SolicitudCard({ solicitud, onTicketUpdated }: SolicitudCardProps): React.JSX.Element {
  const { id, asunto, tipoSolicitud, detalles, datosEspecificos, fechaCreacion, estado, tipoOperacion } = solicitud;
  const router = useRouter(); // Lo guardamos por si quieres añadir otro botón
  const { user } = useUserToken(); // Hook para obtener el token
  const [isResolving, setIsResolving] = React.useState(false);

  // --- 3. Lógica del botón "Resolver" (REEMPLAZADA) ---
  const handleProcesar = (): void => {
    const slug = TIPO_TO_SLUG[tipoSolicitud];
    if (!slug) return;


    // --- Lógica Condicional ---
    if (tipoOperacion === 'ingreso') {
      // 1. Ir a "Agregar Datos" y pre-llenar el formulario
      const targetPath = paths.dashboard.agregardatos;
      
      // Creamos un 'data' que SÍ coincide con SolicitudFormData
      // (asunto, detalles, tipoSolicitud, datosEspecificos)
      const dataParaFormulario = {
        tipoSolicitud: solicitud.tipoSolicitud,
        asunto: solicitud.asunto,
        detalles: solicitud.detalles,
        datosEspecificos: solicitud.datosEspecificos,
      };
      
      const dataString = JSON.stringify(dataParaFormulario);
      const encodedData = encodeURIComponent(dataString);
      const fullUrl = `${targetPath}?tab=${solicitud.tipoSolicitud}&data=${encodedData}`;
      router.push(fullUrl);
      
    } else {
      const targetPath = paths.dashboard.gestion;
      const tabValue = `${slug}s`; 
      
      const fullUrl = `${targetPath}?tab=${tabValue}`; 
      router.push(fullUrl);
    }
  };

  const handleResolverApi = async (): Promise<void> => {
    const comentario = '';
    if (!comentario || !user?.token) return;

    setIsResolving(true);
    try {
      const responseData = await services.authorizedPut(
        `/resolver-ticket/${id}`,
        { comentario_admin: comentario },
        user.token
      ) as { message: string; ticket: ApiTicket };

      const updatedFrontendTicket = transformarApiATipoFrontend(responseData.ticket);
      
      // Llama a la función del padre para actualizar el estado
      onTicketUpdated(updatedFrontendTicket);

    } catch (err) {
      // alert(err instanceof Error ? err.message : 'Un error ocurrió.');
    } finally {
      setIsResolving(false);
    }
  };

  // NOTE: navigation helper removed (handled by `handleProcesar`).
  // Ensure estado is treated safely when checking resolved state.
  const isResolved = (estado || '').toLowerCase() === 'resuelta';

  return (
    <Card sx={{ height: '100%' }}>
      <CardHeader
        action={
          <Chip label={tipoSolicitud} color={getChipColor(tipoSolicitud)} size="small" />
        }
        title={asunto}
        titleTypographyProps={{ variant: 'h6' }}
        subheader={
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 0.5 }}>
            <Typography variant="caption" color="text.secondary">
              {new Date(fechaCreacion).toLocaleDateString()}
            </Typography>
            {/* El Chip de estado ahora se actualizará automáticamente */}
            <Chip 
              label={estado} 
              size="small" 
              variant="outlined" 
              color={isResolved ? 'success' : 'default'}
            />
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
        {/* Botón de NAVEGACIÓN */}
        <Button 
          variant="outlined" 
          color="secondary" 
          onClick={handleProcesar}
          disabled={isResolved}
        >
          Procesar
        </Button>
        
        {/* Botón de API */}
        <Button 
          variant="contained" 
          color="primary" 
          onClick={handleResolverApi}
          disabled={isResolving || isResolved}
          sx={{ minWidth: 120 }}
        >
          {isResolving ? <CircularProgress size={24} color="inherit" /> : 'Resolver'}
        </Button>
      </CardActions>

    </Card>
  );
}