'use client'; 

import { Stack } from '@mui/material';
import Typography from '@mui/material/Typography';
import { FormularioSolicitud } from '@/components/dashboard/formulario/formulario-solicitud'; 

export default function SolicitudesPage(): React.JSX.Element {
  return (
    <Stack spacing={6}>
        <div>
            <Typography variant="h4">Formulario de Solicitud de Datos</Typography>
        </div>
        
      <FormularioSolicitud />
    </Stack>
  );
}