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
import type { SolicitudFormData, SensorData, TipoSensor } from '@/types/formulario';
import { useUserToken } from '@/hooks/use-usertoken';

const TIPOS_SENSOR: TipoSensor[] = ['Temperatura', 'Presión', 'Vibración'];

interface ActivoSimple {
  id: number;
  nombre: string;
}

export function SensorForm(): React.JSX.Element {
  const {
    register,
    formState: { errors },
  } = useFormContext<SolicitudFormData>();

  const formErrors = errors.datosEspecificos as
    | Partial<Record<keyof SensorData, { message?: string }>>
    | undefined;

  // Logica para cargar Activos
  const { user: userContext } = useUserToken();
  const [activos, setActivos] = React.useState<ActivoSimple[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [fetchError, setFetchError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const fetchActivos = async (): Promise<void> => {
      if (!userContext?.token) return;
      setIsLoading(true);
      setFetchError(null);
      try {
        // Aqui debe ir una ruta para obtener los activos de los usuarios (A menos que compliquemos un poco más el form agregandole primero la selccion de usuario)
        const response = await fetch('/api/activos/del-usuario', {
          headers: { 'Authorization': `Bearer ${userContext.token}` },
        });
        if (!response.ok) throw new Error('No se pudieron cargar los activos');
        const data = (await response.json()) as ActivoSimple[];
        setActivos(data);
      } catch (err) {
        setFetchError(err instanceof Error ? err.message : 'Error desconocido');
      } finally {
        setIsLoading(false);
      }
    };
    void fetchActivos();
  }, [userContext]);

  return (
    <Grid container spacing={2}>
      <Grid size={{ xs: 12 }}>
        <TextField
          label="Nombre del Sensor"
          fullWidth
          required
          {...register('datosEspecificos.nombre' as const)}
          error={Boolean(formErrors?.nombre)}
          helperText={formErrors?.nombre?.message ?? ''}
        />
      </Grid>
      
      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth required error={Boolean(formErrors?.tipoSnsor)}>
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
          <FormHelperText>{formErrors?.tipoSnsor?.message ?? ''}</FormHelperText>
        </FormControl>
      </Grid>

      <Grid size={{ xs: 12 }}>
        <FormControl fullWidth required error={Boolean(formErrors?.activoAsociadoId) || Boolean(fetchError)}>
          <InputLabel id="activo-asociado-label">Activo Asociado</InputLabel>
          <Select
            labelId="activo-asociado-label"
            label="Activo Asociado"
            defaultValue=""
            {...register('datosEspecificos.activoAsociadoId' as const, { valueAsNumber: true })}
            disabled={isLoading || Boolean(fetchError) || activos.length === 0}
          >
            {Boolean(isLoading) && (
              <MenuItem disabled value="">
                <em>Cargando activos...</em>
              </MenuItem>
            )}
            {Boolean(fetchError) && (
              <MenuItem disabled value="">
                <em>Error al cargar</em>
              </MenuItem>
            )}
            {!isLoading && activos.map((activo) => (
              <MenuItem key={activo.id} value={activo.id}>
                {activo.nombre} (ID: {activo.id})
              </MenuItem>
            ))}
          </Select>
          <FormHelperText>{formErrors?.activoAsociadoId?.message ?? fetchError}</FormHelperText>
        </FormControl>
      </Grid>
    </Grid>
  );
}