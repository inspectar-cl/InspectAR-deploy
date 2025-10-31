// components/sensor-form.tsx
import * as React from 'react';
import Grid from '@mui/material/Grid';
import {
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Typography,
} from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, SensorData, TipoSensor } from '@/types/formulario';
import { useAuthUser } from '@/contexts/user-context'; // para obtener el usuario y token

const TIPOS_SENSOR: TipoSensor[] = ['Temperatura', 'Presión', 'Vibración'];

interface ActivoResumen {
  id: number;
  nombre: string;
  tipoActivo?: string;
}

export function SensorForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const { user } = useAuthUser();
  const sensorErrors = errors.datosEspecificos as Partial<Record<keyof SensorData, any>>;

  const [activos, setActivos] = React.useState<ActivoResumen[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  // Cargar activos asociados al usuario
  React.useEffect(() => {
    if (!user?.token) return;

    const fetchActivos = async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await fetch('/api/activos/usuario', {
          headers: { Authorization: `Bearer ${user.token}` },
        });

        if (!res.ok) throw new Error('Error al obtener los activos del usuario.');

        const data = await res.json();
        setActivos(data.activos || []);
      } catch (err: any) {
        console.error(err);
        setError('No se pudieron cargar los activos.');
      } finally {
        setLoading(false);
      }
    };

    fetchActivos();
  }, [user?.token]);

  return (
    <Grid container spacing={2}>
      <Grid size={{xs:12}}>
        <TextField
          label="Nombre del Sensor"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const)}
          error={!!sensorErrors?.nombre}
          helperText={sensorErrors?.nombre?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12}}>
        <FormControl fullWidth required error={!!sensorErrors?.tipoSnsor}>
          <InputLabel id="tipo-sensor-label">Tipo de Sensor</InputLabel>
          <Select
            labelId="tipo-sensor-label"
            label="Tipo de Sensor"
            defaultValue=""
            {...register('datosEspecificos.tipoSnsor' as const)}
          >
            {TIPOS_SENSOR.map((tipo) => (
              <MenuItem key={tipo} value={tipo}>
                {tipo}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>

      <Grid size={{xs:12}}>
        <FormControl fullWidth required error={!!sensorErrors?.activoAsociadoId}>
          <InputLabel id="activo-asociado-label">Activo Asociado</InputLabel>
          {loading ? (
            <CircularProgress size={24} sx={{ mt: 2, mb: 1 }} />
          ) : error ? (
            <Typography color="error" variant="body2">
              {error}
            </Typography>
          ) : (
            <Select
              labelId="activo-asociado-label"
              label="Activo Asociado"
              defaultValue=""
              {...register('datosEspecificos.activoAsociadoId' as const, {
                valueAsNumber: true,
                required: true,
              })}
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
          {sensorErrors?.activoAsociadoId && (
            <Typography color="error" variant="caption">
              {sensorErrors.activoAsociadoId.message ?? ''}
            </Typography>
          )}
        </FormControl>
      </Grid>
    </Grid>
  );
}