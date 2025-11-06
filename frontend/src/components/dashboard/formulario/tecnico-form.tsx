import * as React from 'react';
import Grid from '@mui/material/Grid';
import {
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, TecnicoData, EspecialidadTecnico } from '@/types/formulario';
import { useUserToken } from '@/hooks/use-usertoken';

const ESPECIALIDADES: EspecialidadTecnico[] = ['Climatización', 'Eléctrico', 'Mecánico'];

export function TecnicoForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const tecnicoErrors = errors.datosEspecificos as
    | Partial<Record<keyof TecnicoData, { message?: string }>>
    | undefined;
  const { user } = useUserToken();

  React.useEffect(() => {
    if (!user?.token) return;

    const fetchActivos = async (): Promise<void> => {
      try {
        const res = await fetch('/api/obtener-todos-activos?sensores=true', {
          headers: { Authorization: `Bearer ${user.token}` },
        });
        if (!res.ok) throw new Error('Error al obtener los activos del usuario.');
      } catch (_err: unknown) {
        // Error silenciado
      }
    };

    void fetchActivos();
  }, [user?.token]);

  return (
    <Grid container spacing={2}>
      <Grid size={{xs:12}}> {/* Corregido */}
        <TextField
          label="Nombre del Técnico"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const, {
            required: 'El nombre es obligatorio', // Añadida validación
          })}
          error={Boolean(tecnicoErrors?.nombre)}
          helperText={tecnicoErrors?.nombre?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12}}> {/* Corregido */}
        <TextField
          label="Correo Electrónico"
          fullWidth
          required
          type="email"
          {...register('datosEspecificos.correo' as const, {
            required: 'El correo es obligatorio', // Añadida validación
          })}
          error={Boolean(tecnicoErrors?.correo)}
          helperText={tecnicoErrors?.correo?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12}}> {/* Corregido */}
        <TextField
          label="Teléfono"
          fullWidth
          required
          type="tel"
          {...register('datosEspecificos.telefono' as const, {
            required: 'El teléfono es obligatorio', // Añadida validación
          })}
          error={Boolean(tecnicoErrors?.telefono)}
          helperText={tecnicoErrors?.telefono?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12}}> {/* Corregido */}
        <FormControl fullWidth required error={Boolean(tecnicoErrors?.especialidad)}>
          <InputLabel id="especialidad-label">Especialidad</InputLabel>
          <Select
            labelId="especialidad-label"
            label="Especialidad"
            defaultValue=""
            {...register('datosEspecificos.especialidad' as const, {
              required: 'La especialidad es obligatoria', // Añadida validación
            })}
          >
            {ESPECIALIDADES.map((e) => (
              <MenuItem key={e} value={e}>{e}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>

    </Grid>
  );
}
