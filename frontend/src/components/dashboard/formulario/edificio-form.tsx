'use client';

import * as React from 'react';
import { TextField, Grid } from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, EdificioData } from '@/types/formulario';

export function EdificioForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const edificioErrors = errors.datosEspecificos as
    | Partial<Record<keyof EdificioData, { message?: string }>>
    | undefined;

  return (
    <Grid container spacing={2}>
      <Grid size={{xs:12, sm: 6}}>
        <TextField
          label="Nombre del Edificio"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const, {
            required: 'El nombre es obligatorio',
          })}
          error={Boolean(edificioErrors?.nombre)}
          helperText={edificioErrors?.nombre?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12, sm: 6}}>
        <TextField
          label="Dirección del Edificio"
          fullWidth
          required
          {...register('datosEspecificos.direccion' as const, {
            required: 'La dirección es obligatoria',
          })}
          error={Boolean(edificioErrors?.direccion)}
          helperText={edificioErrors?.direccion?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12, sm: 6}}>
        <TextField
          label="Latitud (Aprox.)"
          fullWidth
          required
          type="number"
          {...register('datosEspecificos.latitud' as const, {
            required: 'La latitud es obligatoria',
            valueAsNumber: true,
          })}
          error={Boolean(edificioErrors?.latitud)}
          helperText={edificioErrors?.latitud?.message ?? ''}
          slotProps={{
            htmlInput: { step: 'any' },
            inputLabel: { shrink: true },
          }}
        />
      </Grid>

      <Grid size={{xs:12, sm: 6}}>
        <TextField
          label="Longitud (Aprox.)"
          fullWidth
          required
          type="number"
          {...register('datosEspecificos.longitud' as const, {
            required: 'La longitud es obligatoria',
            valueAsNumber: true,
          })}
          error={Boolean(edificioErrors?.longitud)}
          helperText={edificioErrors?.longitud?.message ?? ''}
          slotProps={{
            htmlInput: { step: 'any' },
            inputLabel: { shrink: true },
          }}
        />
      </Grid>
    </Grid>
  );
}