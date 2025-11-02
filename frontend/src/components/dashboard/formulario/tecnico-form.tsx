import * as React from 'react';
import Grid from '@mui/material/Grid';
import {
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Box,
  CircularProgress,
  Typography,
} from '@mui/material';
import { Controller, useFormContext } from 'react-hook-form';
import type { SolicitudFormData, TecnicoData, EspecialidadTecnico } from '@/types/formulario';
import { useUserToken } from '@/hooks/use-usertoken';

const ESPECIALIDADES: EspecialidadTecnico[] = ['Climatización', 'Eléctrico', 'Mecánico'];

interface ActivoResumen {
  id: number;
  nombre: string;
  tipoActivo?: string;
}

export function TecnicoForm(): React.JSX.Element {
  const {
    register,
    control,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const tecnicoErrors = errors.datosEspecificos as
    | Partial<Record<keyof TecnicoData, { message?: string }>>
    | undefined;
  const { user } = useUserToken();

  const [activos, setActivos] = React.useState<ActivoResumen[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (!user?.token) return;

    const fetchActivos = async (): Promise<void> => {
      setLoading(true);
      setError(null);
      try {
        const res = await fetch('/api/activos/usuario', {
          headers: { Authorization: `Bearer ${user.token}` },
        });
        if (!res.ok) throw new Error('Error al obtener los activos del usuario.');

        const data = (await res.json()) as { activos?: ActivoResumen[] };
        setActivos(data.activos ?? []);
      } catch (err: unknown) {
        //console.error(err);
        setError('No se pudieron cargar los activos.');
      } finally {
        setLoading(false);
      }
    };

    void fetchActivos();
  }, [user?.token]);

  return (
    <Grid container spacing={2}>
      <Grid size={{ xs: 12 }}>
        <TextField
          label="Nombre del Técnico"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const)}
          error={Boolean(tecnicoErrors?.nombre)}
          helperText={tecnicoErrors?.nombre?.message ?? ''}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <TextField
          label="Correo Electrónico"
          fullWidth
          required
          type="email"
          {...register('datosEspecificos.correo' as const)}
          error={Boolean(tecnicoErrors?.correo)}
          helperText={tecnicoErrors?.correo?.message ?? ''}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <TextField
          label="Teléfono"
          fullWidth
          required
          type="tel"
          {...register('datosEspecificos.telefono' as const)}
          error={Boolean(tecnicoErrors?.telefono)}
          helperText={tecnicoErrors?.telefono?.message ?? ''}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth required error={Boolean(tecnicoErrors?.especialidad)}>
          <InputLabel id="especialidad-label">Especialidad</InputLabel>
          <Select
            labelId="especialidad-label"
            label="Especialidad"
            defaultValue=""
            {...register('datosEspecificos.especialidad' as const)}
          >
            {ESPECIALIDADES.map((e) => (
              <MenuItem key={e} value={e}>
                {e}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>

      <Grid size={{ xs: 12 }}>
        <Controller
          name="datosEspecificos.activosAsociados"
          control={control}
          defaultValue={[]}
          render={({ field }) => (
            <FormControl fullWidth error={Boolean(tecnicoErrors?.activosAsociados)}>
              <InputLabel id="activos-label">Activos Asociados</InputLabel>

              {loading ? (
                <CircularProgress size={24} sx={{ mt: 2, mb: 1 }} />
              ) : error ? (
                <Typography color="error" variant="body2">
                  {error}
                </Typography>
              ) : (
                <Select
                  labelId="activos-label"
                  multiple
                  label="Activos Asociados"
                  value={field.value}
                  onChange={(e) => {field.onChange(e.target.value)}}
                  renderValue={(selected) => (
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                      {selected.map((id) => {
                        const activo = activos.find((a) => a.id === id);
                        return (
                          <Chip
                            key={id}
                            label={activo ? activo.nombre : `Activo ${id}`}
                          />
                        );
                      })}
                    </Box>
                  )}
                >
                  {activos.length > 0 ? (
                    activos.map((a) => (
                      <MenuItem key={a.id} value={a.id}>
                        {a.nombre} {a.tipoActivo ? `(${a.tipoActivo})` : ''}
                      </MenuItem>
                    ))
                  ) : (
                    <MenuItem disabled>No hay activos disponibles</MenuItem>
                  )}
                </Select>
              )}

              {tecnicoErrors?.activosAsociados && (
                <Typography color="error" variant="caption">
                  {tecnicoErrors.activosAsociados.message ?? ''}
                </Typography>
              )}
            </FormControl>
          )}
        />
      </Grid>
    </Grid>
  );
}
