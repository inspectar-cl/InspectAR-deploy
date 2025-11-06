import * as React from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import type { SolicitudAPI, TipoSolicitud, ActivoDataAPI,  TecnicoData, EdificioData } from '@/types/form-solicitud';

interface DetalleProps {
  tipo: TipoSolicitud;
  datos: SolicitudAPI['datosEspecificos'];
}

// Componente helper para mostrar "Clave: Valor"
function DatoItem({ label, value }: { label: string, value: React.ReactNode }): React.JSX.Element {
  return (
    <Typography variant="body2" component="div">
      <Box component="span" sx={{ fontWeight: 'bold' }}>{label}:</Box> {value}
    </Typography>
  );
};

export function DetalleDatosEspecificos({ tipo, datos }: DetalleProps): React.JSX.Element {
  switch (tipo) {
    case 'Edificio':{
      const dEdificio = datos as EdificioData;
      return (
        <Box>
          <DatoItem label="Nombre" value={dEdificio.nombre} />
          <DatoItem label="Dirección" value={dEdificio.direccion} />
        </Box>
      );
    }
    case 'Activo': {
      const dActivo = datos as ActivoDataAPI;
      return (
        <Box>
          <DatoItem label="Tipo" value={dActivo.tipoActivo} />
          <DatoItem label="Ubicación" value={dActivo.ubicacion} />
          <DatoItem label="ID Edificio" value={dActivo.edificioId} />
        </Box>
      );
    }
    case 'Técnico': {
      const dTecnico = datos as TecnicoData;
      return (
        <Box>
          <DatoItem label="Nombre" value={dTecnico.nombre} />
          <DatoItem label="Especialidad" value={dTecnico.especialidad} />
          <DatoItem label="Correo" value={dTecnico.correo} />
        </Box>
      );
    }
    default: {
      return <Typography variant="body2" color="error">Tipo de datos no reconocido.</Typography>;
    }
  }
}