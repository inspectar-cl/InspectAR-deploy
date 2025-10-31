'use client'; 

import { Stack } from '@mui/material';
import Typography from '@mui/material/Typography';
import { ListaSolicitudes } from '@/components/dashboard/solicitud/lista-solicitud'; 

export default function SolicitudesPage() {
  return (
    <Stack spacing={3}>
        <div>
            <Typography variant="h4">Solicitudes</Typography>
        </div>
        
      <ListaSolicitudes />
    </Stack>
  );
}