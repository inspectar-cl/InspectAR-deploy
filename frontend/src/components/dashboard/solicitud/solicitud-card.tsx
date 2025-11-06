'use client'

import * as React from 'react';
import { Card, CircularProgress, CardContent, CardHeader, Typography, Chip, Divider, Box, CardActions, Button, type ChipProps } from '@mui/material';
import type { SolicitudAPI, ApiTicket } from '@/types/form-solicitud';
import { DetalleDatosEspecificos } from './detalles-datos';
import { TIPO_TO_SLUG, transformarApiATipoFrontend } from '@/utils/solicitud-utils';
import { useRouter } from 'next/navigation';
import { useUserToken } from '@/hooks/use-usertoken';

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
  const { id, asunto, tipoSolicitud, detalles, datosEspecificos, fechaCreacion, estado } = solicitud;
  const router = useRouter(); // Lo guardamos por si quieres añadir otro botón
  const { user } = useUserToken(); // Hook para obtener el token
  const [isResolving, setIsResolving] = React.useState(false);

  // --- 3. Lógica del botón "Resolver" (REEMPLAZADA) ---
  const handleResolver = async (): Promise<void> => {
    // Pedimos el comentario al usuario
    const comentario = prompt(
      'Por favor, ingresa un comentario de resolución:',
      'Resuelto con éxito.'
    );

    // Si el usuario cancela el prompt, no hacemos nada
    if (!comentario) {
      return;
    }

    if (!user?.token) {
      alert('Error: Token de usuario no encontrado.');
      return;
    }

    setIsResolving(true);
    try {
      const response = await fetch(`/api/resolver-ticket/${id}`, {
        method: 'PUT', // O 'POST'/'PATCH' según tu API
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${user.token}`,
        },
        body: JSON.stringify({ comentario_admin: comentario }),
      });

      if (!response.ok) {
        const errData = await response.json().catch(() => null);
        throw new Error(errData?.error || 'Error al resolver el ticket');
      }

      const responseData: { message: string; ticket: ApiTicket } = await response.json();

      // Transforma el ticket de la API al formato del frontend
      const updatedFrontendTicket = transformarApiATipoFrontend(responseData.ticket);

      // Pasa el ticket actualizado al componente padre (ListaSolicitudes)
      onTicketUpdated(updatedFrontendTicket);
      
    } catch (err) {
      console.error(err);
      alert(err instanceof Error ? err.message : 'Un error ocurrió.');
    } finally {
      setIsResolving(false);
    }
  };

  const handleNavigateToForm = (): void => {
    const slug = TIPO_TO_SLUG[solicitud.tipoSolicitud];
    if (!slug) return;
    const targetPath = `/dashboard/agregar-datos?tab=${slug}`;
    const dataString = JSON.stringify(solicitud);
    const encodedData = encodeURIComponent(dataString);
    const fullUrl = `${targetPath}?data=${encodedData}`;
    router.push(fullUrl);
  };

  const isResolved = estado.toLowerCase() === 'resuelta';

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
        
        {/* Botón para NAVEGAR AL FORMULARIO (antes llamado "Resolver") */}
        <Button 
          variant="outlined" 
          color="secondary" 
          onClick={handleNavigateToForm}
          disabled={isResolved} // No puedes procesar algo resuelto
        >
          Procesar
        </Button>

        {/* Botón para LLAMAR A LA API Y MARCAR COMO RESUELTO */}
        <Button
          variant="contained"
          color="primary"
          onClick={handleResolver}
          disabled={isResolving || isResolved} // Deshabilitado si está cargando o ya está resuelto
          sx={{ minWidth: 120 }} // Evita que el botón cambie de tamaño
        >
          {isResolving ? (
            <CircularProgress size={24} color="inherit" />
          ) : (
            'Marcar Resuelto'
          )}
        </Button>
      </CardActions>

    </Card>
  );
}