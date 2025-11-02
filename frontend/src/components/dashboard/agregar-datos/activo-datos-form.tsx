// components/activo-form.tsx
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
import type { SolicitudFormData, ActivoData, TipoActivo } from '@/types/formulario';
import { useUserToken } from '@/hooks/use-usertoken';

const TIPOS_ACTIVO: TipoActivo[] = ['Ascensor', 'BombaDeAgua', 'PanelElectrico'];

interface EdificioSimple {
  id: number;
  nombre: string;
}

export function ActivoForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const formErrors = errors.datosEspecificos as
    | Partial<Record<keyof ActivoData, { message?: string }>>
    | undefined;
  
  const { user: userContext } = useUserToken();
  const [edificios, setEdificios] = React.useState<EdificioSimple[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [fetchError, setFetchError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const fetchEdificios = async (): Promise<void> => {
      if (!userContext?.token) return;
      setIsLoading(true);
      setFetchError(null);
      try {
        // Aqui debe ir una ruta para obtener todos los edificios
        const response = await fetch('/api/edificios/listar', {
          headers: { 'Authorization': `Bearer ${userContext.token}` },
        });
        if (!response.ok) throw new Error('No se pudieron cargar los edificios');
        const data = (await response.json()) as EdificioSimple[];
        setEdificios(data);
      } catch (err) {
        setFetchError(err instanceof Error ? err.message : 'Error desconocido');
      } finally {
        setIsLoading(false);
      }
    };
    void fetchEdificios();
  }, [userContext]);

  return (
    <Grid container spacing={2}>
      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth required>
          <InputLabel id="tipo-activo-label">Tipo de Activo</InputLabel>
          <Select
            labelId="tipo-activo-label"
            label="Tipo de Activo"
            defaultValue=""
            {...register('datosEspecificos.tipoActivo' as const)}
          >
            {TIPOS_ACTIVO.map((tipo) => (
              <MenuItem key={tipo} value={tipo}>
                {tipo}
              </MenuItem>
            ))}
          </Select>
          <FormHelperText>{formErrors?.tipoActivo?.message ?? ''}</FormHelperText>
        </FormControl>
      </Grid>
      
      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth required error={Boolean(formErrors?.edificioId) || Boolean(fetchError)}>
          <InputLabel id="edificio-id-label">Edificio Asociado</InputLabel>
          <Select
            labelId="edificio-id-label"
            label="Edificio Asociado"
            defaultValue=""
            {...register('datosEspecificos.edificioId' as const, { valueAsNumber: true })}
            disabled={
              isLoading || Boolean(fetchError) || edificios.length === 0
            }
          >
            {Boolean(isLoading) && (
              <MenuItem disabled value="">
                <em>Cargando edificios...</em>
              </MenuItem>
            )}
            {Boolean(fetchError) && (
              <MenuItem disabled value="">
                <em>Error al cargar</em>
              </MenuItem>
            )}
            {!isLoading && edificios.map((edificio) => (
              <MenuItem key={edificio.id} value={edificio.id}>
                {edificio.nombre}
              </MenuItem>
            ))}
          </Select>
          <FormHelperText>{formErrors?.edificioId?.message ?? fetchError}</FormHelperText>
        </FormControl>
      </Grid>

      <Grid size={{ xs: 12 }}>
        <TextField
          label="Ubicación (Ej: Piso 5, Sala de máquinas)"
          fullWidth
          required
          {...register('datosEspecificos.ubicacion' as const)}
          error={Boolean(formErrors?.ubicacion)}
          helperText={formErrors?.ubicacion?.message ?? ''}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <TextField
          label="Descripción"
          fullWidth
          multiline
          rows={3}
          {...register('datosEspecificos.descripcion' as const)}
        />
      </Grid>

      <Grid size={{ xs: 12 }}>
        <TextField
          label="Imagen"
          type="file"
          fullWidth
          slotProps={{
            inputLabel: {
              shrink: true,
            }
          }}
          {...register('datosEspecificos.imagen' as const)}
          error={Boolean(formErrors?.imagen)}
          helperText={formErrors?.imagen?.message ?? 'Sube una imagen del activo (opcional)'}
        />
      </Grid>
    </Grid>
  );
}