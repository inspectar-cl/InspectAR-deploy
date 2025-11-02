import * as React from 'react';
import Grid from '@mui/material/Grid';
import { TextField, MenuItem, InputLabel, FormControl, Select } from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, ActivoData, TipoActivo } from '@/types/formulario';

const TIPOS_ACTIVO: TipoActivo[] = ['Ascensor', 'BombaDeAgua', 'PanelElectrico'];

export function ActivoForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const activoErrors = errors.datosEspecificos as
      | Partial<Record <keyof ActivoData, { message?: string }>>
      | undefined

  return (
    <Grid container spacing={2}>
      <Grid size={{xs:12}}>
        <FormControl fullWidth required error={Boolean(activoErrors?.tipoActivo)}>
          <InputLabel id="tipo-activo-label">Tipo de Activo</InputLabel>
          <Select
            labelId="tipo-activo-label"
            label="Tipo de Activo"
            defaultValue=""
            {...register('datosEspecificos.tipoActivo' as const, { required: true })}
          >
            {TIPOS_ACTIVO.map(tipo => (
              <MenuItem key={tipo} value={tipo}>{tipo}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>

      <Grid size={{xs:12}}>
        <TextField
          label="ID del Edificio Asociado"
          fullWidth
          required
          type="number"
          {...register('datosEspecificos.edificioId' as const, { valueAsNumber: true })}
          error={Boolean(activoErrors?.edificioId)}
          helperText={activoErrors?.edificioId?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12}}>
        <TextField
          label="Ubicación dentro del edificio"
          fullWidth
          required
          {...register('datosEspecificos.ubicacion' as const)}
          error={Boolean(activoErrors?.ubicacion)}
          helperText={activoErrors?.ubicacion?.message ?? ''}
        />
      </Grid>

      <Grid size={{xs:12}}>
        <TextField
          label="Descripción (opcional)"
          fullWidth
          multiline
          rows={3}
          {...register('datosEspecificos.descripcion' as const)}
        />
      </Grid>

      <Grid size={{xs:12}}>
        <TextField
          type="file"
          fullWidth
          inputProps={{ accept: 'image/*' }}
          {...register('datosEspecificos.imagen' as const)}
          error={Boolean(activoErrors?.imagen)}
          helperText={activoErrors?.imagen?.message ?? ''}
        />
      </Grid>
    </Grid>
  );
}
