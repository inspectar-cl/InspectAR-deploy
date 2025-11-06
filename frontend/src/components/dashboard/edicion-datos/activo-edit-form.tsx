'use client';

import * as React from 'react';
import Grid from '@mui/material/Grid';
import { TextField } from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData } from '@/types/formulario';

// Definimos los campos que este formulario SÍ edita
interface ActivoEditData {
  nombre: string;
  ubicacion: string;
  descripcion?: string;
}

export function ActivoEditForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>(); // Aún usamos la estructura padre

  // Apuntamos a los errores de los campos específicos
  const formErrors = errors.datosEspecificos as
    | Partial<Record<keyof ActivoEditData, { message?: string }>>
    | undefined;

  return (
    // Usamos Grid item xs={12} para un layout correcto
    <Grid container spacing={2}>
      
      {/* 1. CAMPO NOMBRE */}
      <Grid size={{ xs:12 }}>
        <TextField
          label="Nombre del Activo"
          fullWidth
          required
          // Usamos 'datosEspecificos.nombre' para que coincida con el 'reset'
          {...register('datosEspecificos.nombre' as const, {
            required: 'El nombre es obligatorio',
          })}
          error={Boolean(formErrors?.nombre)}
          helperText={formErrors?.nombre?.message ?? ''}
          InputLabelProps={{ shrink: true }} // Para que el label no se solape
        />
      </Grid>

      {/* 2. CAMPO UBICACIÓN */}
      <Grid size={{ xs:12 }}>
        <TextField
          label="Ubicación"
          fullWidth
          required
          {...register('datosEspecificos.ubicacion' as const, {
            required: 'La ubicación es obligatoria',
          })}
          error={Boolean(formErrors?.ubicacion)}
          helperText={formErrors?.ubicacion?.message ?? ''}
          InputLabelProps={{ shrink: true }}
        />
      </Grid>

      {/* 3. CAMPO DESCRIPCIÓN */}
      <Grid size={{ xs:12 }}>
        <TextField
          label="Descripción"
          fullWidth
          multiline
          rows={3}
          {...register('datosEspecificos.descripcion' as const)}
          error={Boolean(formErrors?.descripcion)}
          helperText={formErrors?.descripcion?.message ?? ''}
          InputLabelProps={{ shrink: true }}
        />
      </Grid>
    </Grid>
  );
}