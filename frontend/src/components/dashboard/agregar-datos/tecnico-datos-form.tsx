'use client';

import * as React from 'react';
import Grid from '@mui/material/Grid';
import {
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  FormHelperText,
  Chip,
  Box,
} from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, TecnicoData, EspecialidadTecnico } from '@/types/formulario';
import { useUserToken } from '@/hooks/use-usertoken';

const ESPECIALIDADES: EspecialidadTecnico[] = ['Climatización', 'Eléctrico', 'Mecánico'];

interface ActivoSimple {
  id: number;
  nombre: string;
}

export function TecnicoForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const formErrors = errors.datosEspecificos as Partial<Record<keyof TecnicoData, any>>;

  const { user: userContext } = useUserToken();
  const [activos, setActivos] = React.useState<ActivoSimple[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [fetchError, setFetchError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const fetchActivos = async () => {
      if (!userContext?.token) return;
      setIsLoading(true);
      setFetchError(null);
      try {
        //Aqui lo mismo que con activo form
        const response = await fetch('/api/activos/del-usuario', {
          headers: { 'Authorization': `Bearer ${userContext.token}` },
        });
        if (!response.ok) throw new Error('No se pudieron cargar los activos');
        const data: ActivoSimple[] = await response.json();
        setActivos(data);
      } catch (err) {
        setFetchError(err instanceof Error ? err.message : 'Error desconocido');
      } finally {
        setIsLoading(false);
      }
    };
    fetchActivos();
  }, [userContext]);

  return (
    <Grid container spacing={2}>
      <Grid size={{ xs: 12 }}>
        <TextField
          label="Nombre del Técnico"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const)}
          error={!!formErrors?.nombre}
          helperText={formErrors?.nombre?.message ?? ''}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth required error={!!formErrors?.especialidad}>
          <InputLabel id="especialidad-label">Especialidad</InputLabel>
          <Select
            labelId="especialidad-label"
            label="Especialidad"
            defaultValue=""
            {...register('datosEspecificos.especialidad' as const)}
          >
            {ESPECIALIDADES.map((esp) => (
              <MenuItem key={esp} value={esp}>
                {esp}
              </MenuItem>
            ))}
          </Select>
          <FormHelperText>{formErrors?.especialidad?.message ?? ''}</FormHelperText>
        </FormControl>
      </Grid>
      
      <Grid size={{ xs: 12 }}>
        <TextField
          label="Correo Electrónico"
          type="email"
          fullWidth
          required
          {...register('datosEspecificos.correo' as const)}
          error={!!formErrors?.correo}
          helperText={formErrors?.correo?.message ?? ''}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <TextField
          label="Teléfono"
          fullWidth
          required
          {...register('datosEspecificos.telefono' as const)}
          error={!!formErrors?.telefono}
          helperText={formErrors?.telefono?.message ?? ''}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth error={!!formErrors?.activosAsociados || !!fetchError}>
          <InputLabel id="activos-asociados-label">Activos Asociados</InputLabel>
          <Select
            labelId="activos-asociados-label"
            label="Activos Asociados"
            multiple // <-- Permite selección múltiple
            defaultValue={[]} // <-- Importante para 'multiple'
            {...register('datosEspecificos.activosAsociados' as const)}
            disabled={isLoading || !!fetchError || activos.length === 0}
            // Renderiza Chips para las selecciones
            renderValue={(selected) => (
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                {(selected as number[]).map((value) => {
                  const activo = activos.find(a => a.id === value);
                  return <Chip key={value} label={activo ? activo.nombre : value} />;
                })}
              </Box>
            )}
          >
            {isLoading && <MenuItem disabled value=""><em>Cargando activos...</em></MenuItem>}
            {fetchError && <MenuItem disabled value=""><em>Error al cargar</em></MenuItem>}
            {!isLoading && activos.map((activo) => (
              <MenuItem key={activo.id} value={activo.id}>
                {activo.nombre}
              </MenuItem>
            ))}
          </Select>
          <FormHelperText>{formErrors?.activosAsociados?.message ?? fetchError ?? 'Selecciona uno o más activos'}</FormHelperText>
        </FormControl>
      </Grid>
    </Grid>
  );
}