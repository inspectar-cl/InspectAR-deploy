// components/dashboard/gestion/tecnico-edit-form.tsx
'use client';

import * as React from 'react';
import Grid from '@mui/material/Grid';
import {
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormHelperText,
} from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, EspecialidadTecnico } from '@/types/formulario';

// Asegúrate de que estos tipos coincidan con tu formulario de creación
const ESPECIALIDADES: EspecialidadTecnico[] = ['Climatización', 'Eléctrico', 'Mecánico'];

// Definimos los campos que este formulario SÍ edita
interface TecnicoEditData {
  nombre: string;
  correo: string;
  telefono: string;
  especialidad: EspecialidadTecnico;
}

export function TecnicoEditForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const formErrors = errors.datosEspecificos as
    | Partial<Record<keyof TecnicoEditData, { message?: string }>>
    | undefined;

  return (
    <Grid container spacing={2}>

      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth required error={Boolean(formErrors?.especialidad)}>
          <InputLabel id="especialidad-label">Especialidad</InputLabel>
          <Select
            labelId="especialidad-label"
            label="Especialidad"
            defaultValue=""
            {...register('datosEspecificos.especialidad' as const, {
              required: 'La especialidad es obligatoria',
            })}
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
          {...register('datosEspecificos.correo' as const, {
            required: 'El correo es obligatorio',
          })}
          error={Boolean(formErrors?.correo)}
          helperText={formErrors?.correo?.message ?? ''}
          InputLabelProps={{ shrink: true }}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <TextField
          label="Teléfono"
          fullWidth
          required
          {...register('datosEspecificos.telefono' as const, {
            required: 'El teléfono es obligatorio',
          })}
          error={Boolean(formErrors?.telefono)}
          helperText={formErrors?.telefono?.message ?? ''}
          InputLabelProps={{ shrink: true }}
        />
      </Grid>

      {/* Omitimos 'activosAsociados' ya que no es editable aquí */}
    </Grid>
  );
}