// components/edificio-form.tsx
'use client';

import * as React from 'react';
import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, EdificioData } from '@/types/formulario';

export function EdificioDataForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  // Helper para acceder a los errores de este sub-formulario
  const formErrors = errors.datosEspecificos as
    | Partial<Record<keyof EdificioData, { message?: string }>>
    | undefined;

  return (
    <Grid container spacing={2}>
      <Grid size={{ xs: 12 }}>
        <TextField
          label="Nombre del Edificio"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const, {
            required: 'El nombre es obligatorio',
          })}
          error={Boolean(formErrors?.nombre)}
          helperText={formErrors?.nombre?.message ?? ''}
        />
      </Grid>
      <Grid size={{ xs: 12 }}>
        <TextField
          label="Dirección"
          fullWidth
          required
          // --- AÑADIDA VALIDACIÓN DE RHF ---
          {...register('datosEspecificos.direccion' as const, {
            required: 'La dirección es obligatoria',
          })}
          error={Boolean(formErrors?.direccion)}
          helperText={formErrors?.direccion?.message ?? ''}
        />
      </Grid>
      <Grid size={{ xs: 12, sm:6}}>
        <TextField
          label="Latitud"
          type="number"
          fullWidth
          required
          slotProps={{
            htmlInput: { step: 'any' },
            inputLabel: { shrink: true },
          }}
          {...register('datosEspecificos.latitud' as const, {
            required: 'La latitud es obligatoria',
            valueAsNumber: true,
          })}
          error={Boolean(formErrors?.latitud)}
          helperText={formErrors?.latitud?.message ?? ''}
        />
      </Grid>
      <Grid size={{ xs: 12, sm:6}}>
        <TextField
          label="Longitud"
          type="number"
          fullWidth
          required
          slotProps={{
            htmlInput: { step: 'any' },
            inputLabel: { shrink: true },
          }}
          {...register('datosEspecificos.longitud' as const, {
            required: 'La longitud es obligatoria',
            valueAsNumber: true,
          })}
          error={Boolean(formErrors?.longitud)}
          helperText={formErrors?.longitud?.message ?? ''}
        />
      </Grid>
    </Grid>
  );
}