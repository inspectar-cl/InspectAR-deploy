import * as React from 'react';
import { 
    Box, 
    Container, 
    Typography, 
    Tabs, 
    Tab, 
    Grid, 
    CircularProgress, 
    Alert,
    Paper,
    FormControl,
    InputLabel,
    Select,
    MenuItem
} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material';
import type { 
  SolicitudAPI, 
  TipoSolicitud, 
  ApiTicket,
  EdificioData,
  ActivoDataAPI,
  TecnicoData
} from '@/types/form-solicitud';
import { SolicitudCard } from './solicitud-card';

import { MOCK_SOLICITUDES } from '@/mocks/solicitudes';
import { useUserToken } from '@/hooks/use-usertoken';

const TIPOS_FILTRO: (TipoSolicitud | 'Todas')[] = ['Todas', 'Edificio', 'Activo', 'Técnico'];
type EstadoSolicitud = 'Pendiente' | 'EnProgreso' | 'Resuelta';
const TIPOS_ESTADO: ('Todos' | EstadoSolicitud)[] = ['Todos', 'Pendiente', 'EnProgreso', 'Resuelta'];

function transformarApiATipoFrontend(ticket: ApiTicket): SolicitudAPI {
  let tipoSolicitud: TipoSolicitud;
  let datosEspecificos: any = {};
  
  // Mapea el estado (Backend "resuelto" -> Frontend "Resuelta")
  let estadoFrontend: EstadoSolicitud = 'Pendiente';
  if (ticket.estado === 'resuelto') estadoFrontend = 'Resuelta';
  if (ticket.estado === 'en progreso') estadoFrontend = 'EnProgreso';
  
  // Mapea la entidad y extrae los datos específicos
  switch (ticket.tipo_entidad) {
    case 'edificio':
      tipoSolicitud = 'Edificio';
      datosEspecificos = {
        nombre: ticket.edificio_nombre,
        direccion: ticket.edificio_direccion,
        latitud: ticket.edificio_latitud,
        longitud: ticket.edificio_longitud,
      } as EdificioData;
      break;
      
    case 'activo':
      tipoSolicitud = 'Activo';
      // Mapea el tipo de activo (ej. 'bomba de agua' -> 'BombaDeAgua')
      let tipoActivoForm: string = ticket.activo_tipo || '';
      if (tipoActivoForm === 'bomba de agua') tipoActivoForm = 'BombaDeAgua';
      if (tipoActivoForm === 'ascensor') tipoActivoForm = 'Ascensor';
      if (tipoActivoForm === 'panel electrico') tipoActivoForm = 'PanelElectrico';
      
      datosEspecificos = {
        tipoActivo: tipoActivoForm,
        edificioId: ticket.activo_edificio_id,
        ubicacion: ticket.activo_ubicacion,
        descripcion: ticket.activo_descripcion,
        imagen: null, // La API no parece enviar imagen
      } as ActivoDataAPI;
      break;
      
    case 'tecnico':
      tipoSolicitud = 'Técnico';
      datosEspecificos = {
        nombre: ticket.tecnico_nombre,
        correo: ticket.tecnico_email,
        telefono: ticket.tecnico_telefono,
        especialidad: ticket.tecnico_especialidad,
        activosAsociados: [], // La API no parece enviar esto
      } as TecnicoData;
      break;
      
    default:
      // Fallback por si llega un tipo no esperado
      tipoSolicitud = 'Edificio'; 
      datosEspecificos = { nombre: 'Error: Tipo no reconocido' };
  }

  return {
    id: ticket.id,
    estado: estadoFrontend,
    fechaCreacion: ticket.created_at,
    tipoSolicitud: tipoSolicitud,
    asunto: `${ticket.tipo_operacion.toUpperCase()} ${ticket.tipo_entidad}`,
    detalles: ticket.justificacion,
    datosEspecificos: datosEspecificos,
  };
}

export function ListaSolicitudes(): React.JSX.Element {
  const { user } = useUserToken();

  const [allSolicitudes, setAllSolicitudes] = React.useState<SolicitudAPI[]>([]);
  const [filtroActual, setFiltroActual] = React.useState<TipoSolicitud | 'Todas'>('Todas');
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [filtroEstado, setFiltroEstado] = React.useState<'Todos' | EstadoSolicitud>('Todos');

  React.useEffect(() => {
    if (!user?.token) return;

    const fetchSolicitudes = async (): Promise<void> => {
      setIsLoading(true);
      setError(null);

      try {
        const response = await fetch('/api/notification/tickets', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${user.token}`,
          },
        });
        if (!response.ok) {
          const errData = await response.json().catch(() => null);
          throw new Error(errData?.error || 'Error al obtener las solicitudes');
        }

        const data = await response.json();
        let ticketsArray: ApiTicket[];

        if (Array.isArray(data)) {
          ticketsArray = data as ApiTicket[];
        } 
        else if (data && typeof data === 'object' && Array.isArray(data.tickets)) {
          ticketsArray = data.tickets as ApiTicket[];
        } 
        else {
          throw new Error("No hay tickets o solicitudes");
        }

        const solicitudesTransformadas = ticketsArray.map(transformarApiATipoFrontend);
        setAllSolicitudes(solicitudesTransformadas);

      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ocurrió un error desconocido');
      } finally {
        setIsLoading(false);
      }
    };

    void fetchSolicitudes();
  }, [user]);

  // Filtros
  const filteredSolicitudes = React.useMemo(() => {
    let solicitudesFiltradas = allSolicitudes;

    if (filtroActual !== 'Todas') {
      solicitudesFiltradas = solicitudesFiltradas.filter(
        s => s.tipoSolicitud === filtroActual
      );
    }

    if (filtroEstado !== 'Todos') {
      solicitudesFiltradas = solicitudesFiltradas.filter(
        // Compara sin importar mayúsculas/minúsculas
        s => s.estado.toLowerCase() === filtroEstado.toLowerCase()
      );
    }

    return solicitudesFiltradas;
  }, [allSolicitudes, filtroActual, filtroEstado]);

  // Los tabs de navegacion
  const handleFiltroChange = (
    event: React.SyntheticEvent,
    newValue: TipoSolicitud | 'Todas',
  ): void => {
    setFiltroActual(newValue);
  };

  const handleEstadoChange = (
    event: SelectChangeEvent<'Todos' | EstadoSolicitud>,
  ): void => {
    setFiltroEstado(event.target.value as 'Todos' | EstadoSolicitud);
  };

  const handleTicketUpdated = (updatedTicket: SolicitudAPI) => {
    setAllSolicitudes((currentList) =>
      currentList.map((solicitud) =>
        solicitud.id === updatedTicket.id ? updatedTicket : solicitud
      )
    );
  };

  const renderContent = (): React.JSX.Element => {
    if (isLoading) {
      return (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 5 }}>
          <CircularProgress />
        </Box>
      );
    }
    if (error) {
      return <Alert severity="error">{error}</Alert>;
    }
    if (filteredSolicitudes.length === 0) {
      return <Typography sx={{ p: 3 }}>No se encontraron solicitudes para este filtro.</Typography>;
    }

    // Grid de Tarjetas
    return (
      <Grid container spacing={3} sx={{ pt: 3 }}>
        {filteredSolicitudes.map((solicitud) => (
          <Grid key={solicitud.id} size={{ xs: 12, sm: 6, md: 4 }}>
            <SolicitudCard solicitud={solicitud} onTicketUpdated={handleTicketUpdated}/>
          </Grid>
        ))}
      </Grid>
    );
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Visor de Solicitudes
      </Typography>

      <Paper elevation={1} sx={{ p: 2, mb: 3 }}>
        <Grid container spacing={2}>
          <Grid size={{ xs: 12, sm:6, md:4 }}>
            <FormControl fullWidth size="small">
              <InputLabel id="filtro-estado-label">Filtrar por Estado</InputLabel>
              <Select
                labelId="filtro-estado-label"
                label="Filtrar por Estado"
                value={filtroEstado}
                onChange={handleEstadoChange}
              >
                {TIPOS_ESTADO.map((estado) => (
                  <MenuItem key={estado} value={estado}>
                    {estado}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
        </Grid>
      </Paper>

      {/* Tabs */}
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={filtroActual} onChange={handleFiltroChange} aria-label="Filtros de solicitud" variant="scrollable" allowScrollButtonsMobile>
          {TIPOS_FILTRO.map((tipo) => (
            <Tab key={tipo} label={tipo} value={tipo} />
          ))}
        </Tabs>
      </Box>

      {renderContent()}
    </Container>
  );
}