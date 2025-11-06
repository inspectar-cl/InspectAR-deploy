'use client';

import * as React from 'react';
import Grid from '@mui/material/Grid';
import {
  TextField,
  MenuItem,
  InputLabel,
  FormControl,
  Select,
  FormHelperText,
  CircularProgress,
} from '@mui/material';
import { useFormContext } from 'react-hook-form';
import type { SolicitudFormData, ActivoData, TipoActivo } from '@/types/formulario';
import { useUserToken } from '@/hooks/use-usertoken'; // Importa tu hook

const TIPOS_ACTIVO: TipoActivo[] = ['Ascensor', 'BombaDeAgua', 'PanelElectrico'];

interface EdificioSimple {
  id: number;
  nombre: string;
}
interface EdificiosApiResponse {
  edificios: EdificioSimple[];
  total: number;
}

export function ActivoForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const activoErrors = errors.datosEspecificos as
    | Partial<Record<keyof ActivoData, { message?: string }>>
    | undefined;

  // --- Lógica para cargar edificios ---
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
        const response = await fetch('/api/gestion/edificios', { // Ruta de tu API
          headers: { 'Authorization': `Bearer ${userContext.token}` },
        });
        if (!response.ok) throw new Error('No se pudieron cargar los edificios');
        const data = (await response.json()) as EdificiosApiResponse;
        setEdificios(data.edificios);
      } catch (err) {
        setFetchError(err instanceof Error ? err.message : 'Error desconocido');
      } finally {
        setIsLoading(false);
      }
    };
    void fetchEdificios();
  }, [userContext]);

  const ITEM_HEIGHT = 35;

  return (
    <Grid container spacing={2}>
      <Grid size={{xs:12}}>
        <TextField
          label="Nombre del Activo"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const, { // <-- 'nombre' añadido
            required: 'El nombre es obligatorio',
          })}
          error={Boolean(activoErrors?.nombre)}
          helperText={activoErrors?.nombre?.message ?? ''}
        />
      </Grid>
      <Grid size={{xs:12}}>
        <FormControl fullWidth required error={Boolean(activoErrors?.tipoActivo)}>
          <InputLabel id="tipo-activo-label">Tipo de Activo</InputLabel>
          <Select
            labelId="tipo-activo-label"
            label="Tipo de Activo"
            defaultValue=""
            {...register('datosEspecificos.tipoActivo' as const, {
              required: 'El tipo es obligatorio',
            })}
          >
            {TIPOS_ACTIVO.map((tipo) => (
              <MenuItem key={tipo} value={tipo}>{tipo}</MenuItem>
            ))}
          </Select>
          <FormHelperText>{activoErrors?.tipoActivo?.message}</FormHelperText>
        </FormControl>
      </Grid>

      <Grid size={{xs:12}}>
        <FormControl fullWidth required error={Boolean(activoErrors?.edificioId) || Boolean(fetchError)}>
          <InputLabel id="edificio-id-label">Edificio Asociado</InputLabel>
          <Select
            labelId="edificio-id-label"
            label="Edificio Asociado"
            defaultValue=""
            {...register('datosEspecificos.edificioId' as const, {
              required: 'El edificio es obligatorio',
              valueAsNumber: true,
            })}
            disabled={isLoading || Boolean(fetchError) || edificios.length === 0}
            MenuProps={{
              slotProps: {
                paper: {
                  sx: { maxHeight: ITEM_HEIGHT * 4.5 }, // 4.5 ítems visibles
                },
              },
            }}
          >
            {isLoading ? <MenuItem disabled value="">
                <CircularProgress size={20} sx={{ mr: 1, verticalAlign: 'middle' }} />
                <em>Cargando edificios...</em>
              </MenuItem> : null}
            {fetchError ? <MenuItem disabled value="">
                <em>Error al cargar edificios</em>
              </MenuItem> : null}
            {!isLoading && !fetchError && edificios.length === 0 && (
              <MenuItem disabled value="">
                <em>No hay edificios disponibles</em>
              </MenuItem>
            )}
            {!isLoading && !fetchError && edificios.map((edificio) => (
              <MenuItem key={edificio.id} value={edificio.id}>
                {edificio.nombre}
              </MenuItem>
            ))}
          </Select>
          <FormHelperText>{activoErrors?.edificioId?.message ?? fetchError}</FormHelperText>
        </FormControl>
      </Grid>

      <Grid size={{xs:12}}>
        <TextField
          label="Ubicación dentro del edificio"
          fullWidth
          required
          {...register('datosEspecificos.ubicacion' as const, {
            required: 'La ubicación es obligatoria',
          })}
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
    </Grid>
  );
}
