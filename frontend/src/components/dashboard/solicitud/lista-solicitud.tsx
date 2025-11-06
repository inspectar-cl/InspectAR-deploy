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
import type { SolicitudAPI, TipoSolicitud } from '@/types/form-solicitud';
import { SolicitudCard } from './solicitud-card';

import { MOCK_SOLICITUDES } from '@/mocks/solicitudes';
import { useUserToken } from '@/hooks/use-usertoken';

const TIPOS_FILTRO: (TipoSolicitud | 'Todas')[] = ['Todas', 'Edificio', 'Activo', 'Técnico'];
type EstadoSolicitud = 'Pendiente' | 'EnProgreso' | 'Resuelta';
const TIPOS_ESTADO: ('Todos' | EstadoSolicitud)[] = ['Todos', 'Pendiente', 'EnProgreso', 'Resuelta'];

export function ListaSolicitudes(): React.JSX.Element {
  const { user } = useUserToken();
  
  // --- Estados ---
  const [allSolicitudes, setAllSolicitudes] = React.useState<SolicitudAPI[]>([]);
  const [filtroActual, setFiltroActual] = React.useState<TipoSolicitud | 'Todas'>('Todas');
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [filtroEstado, setFiltroEstado] = React.useState<'Todos' | EstadoSolicitud>('Todos');

  // --- 1. Carga de Datos (API Call) ---
  React.useEffect(() => {
    if (!user?.token) return;

    const fetchSolicitudes = async (): Promise<void> => {
      setIsLoading(true);
      setError(null);

      /* //Aqui deberia estar el llamado de la api para obtener todas las solicitudes del back
      try {
        const response = await fetch('/api/obtener-solicitudes', {
          headers: {
            'Authorization': `Bearer ${user.token}`,
          },
        });
        if (!response.ok) {
          throw new Error('Error al obtener las solicitudes');
        }
        const data: SolicitudAPI[] = await response.json();
        setAllSolicitudes(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ocurrió un error desconocido');
      } finally {
        setIsLoading(false);
      }
        */

      try {
          //console.warn("Usando datos MOCK para simular la API");
          await new Promise((resolve) => {
            setTimeout(resolve, 1000);
          });
          setAllSolicitudes(MOCK_SOLICITUDES);
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

    // filtro de TIPO (Tabs)
    if (filtroActual !== 'Todas') {
      solicitudesFiltradas = solicitudesFiltradas.filter(
        s => s.tipoSolicitud === filtroActual
      );
    }

    // filtro de ESTADO (Select)
    if (filtroEstado !== 'Todos') {
      solicitudesFiltradas = solicitudesFiltradas.filter(
        s => s.estado === filtroEstado
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
            <SolicitudCard solicitud={solicitud} />
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
          {/* Aqui luego podrian agregarse más filtros */}
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

      {/* contenido o tarjetas*/}
      {renderContent()}
    </Container>
  );
}