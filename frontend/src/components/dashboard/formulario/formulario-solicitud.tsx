'use client';

import * as React from 'react';
import { 
    Grid, Stack, Button, FormControl, InputLabel, Select, MenuItem, TextField, Typography, Alert, Paper , Divider, CircularProgress
} from '@mui/material';
import { useForm, FormProvider } from 'react-hook-form';
import { EdificioForm } from './edificio-form';
import { ActivoForm } from './activo-form';
import { TecnicoForm } from './tecnico-form';

import type { EdificioData, ActivoData, TecnicoData, SolicitudFormData } from '@/types/formulario'; // Importa todos los tipos
import { useUserToken } from '@/hooks/use-usertoken';
import {decodeJwtToken} from '@/hooks/use-auth'
import { CustomAlert } from '@/components/dashboard/alert-popups/CustomAlertPopup';

type TipoSolicitud = 'Edificio' | 'Activo' | 'Técnico';
type TipoOperacion = 'Ingreso' | 'Modificacion' | 'Eliminacion'; // Tipos de operación

const TIPOS_SOLICITUD: TipoSolicitud[] = ['Edificio', 'Activo', 'Técnico'];
const TIPOS_OPERACION: TipoOperacion[] = ['Ingreso', 'Modificacion', 'Eliminacion'];

function transformarDatosParaApi(data: SolicitudFormData, userEmail: string): any {
  const { tipoSolicitud, asunto, detalles, datosEspecificos, tipoOperacion } = data;
  
  const tipo_entidad = tipoSolicitud.toLowerCase();
  const tipo_operacion = tipoOperacion.toLowerCase();

  const payload: any = {
    tipo_entidad,
    tipo_operacion,
    usuario_email: userEmail,
    justificacion: detalles || asunto,
  };

  // Añadir campos específicos
  switch (tipoSolicitud) {
    case 'Edificio':
      const edificio = datosEspecificos as EdificioData;
      payload.edificio_nombre = edificio.nombre;
      payload.edificio_direccion = edificio.direccion;
      payload.edificio_latitud = edificio.latitud;
      payload.edificio_longitud = edificio.longitud;
      // 'edificio_id' solo se añadiría si la operación es 'modificacion' o 'eliminacion'
      // Por ahora, asumimos que este formulario es solo para 'ingreso'
      break;
      
    case 'Activo':
      const activo = datosEspecificos as ActivoData;
      payload.activo_nombre = activo.nombre; // <-- Añadido
      payload.activo_tipo = activo.tipoActivo.toLowerCase().replace('deagua', ' de agua'); // 'BombaDeAgua' -> 'bomba de agua'
      payload.activo_descripcion = activo.descripcion;
      payload.activo_ubicacion = activo.ubicacion;
      payload.activo_edificio_id = activo.edificioId;
      break;

    case 'Técnico':
      const tecnico = datosEspecificos as TecnicoData;
      payload.tecnico_nombre = tecnico.nombre;
      payload.tecnico_email = tecnico.correo;
      payload.tecnico_telefono = tecnico.telefono;
      payload.tecnico_especialidad = tecnico.especialidad;
      payload.tecnico_autorizado = true; // Valor fijo, según tu ejemplo de API
      // 'activosAsociados' no parece ser parte del payload de solicitud
      break;
  }
  
  return payload;
}

export function FormularioSolicitud(): React.JSX.Element {

    const methods = useForm<SolicitudFormData>({
        defaultValues: {
            tipoSolicitud: 'Edificio',
            tipoOperacion: 'Ingreso',
            asunto: '',
            detalles: '',
            datosEspecificos: {} as EdificioData
        }
    });

    const { register, handleSubmit, watch, formState: { errors, isSubmitting }, reset } = methods; // <-- 2. Obtén 'reset'

    const tipoActual = watch("tipoSolicitud");
    const { user } = useUserToken();
    // const [apiError, setApiError] = React.useState<string | null>(null); // <-- Reemplazado por alertState

    // --- 3. Añade el estado para tu CustomAlert ---
    const [alertState, setAlertState] = React.useState({
      open: false,
      title: '',
      message: '',
      severity: 'success' as 'success' | 'info' | 'warning' | 'error',
    });

    // --- 4. useEffect para cerrar el pop-up automáticamente ---
    React.useEffect(() => {
      if (alertState.open) {
        const timer = setTimeout(() => {
          setAlertState(prev => ({ ...prev, open: false }));
        }, 4000); // Cierra después de 4 segundos
        return () => clearTimeout(timer); // Limpia el timer si el componente se desmonta
      }
    }, [alertState.open]);


    const onSubmit = async (data: SolicitudFormData): Promise<void> => {
        const decodedpayload = decodeJwtToken(user?.token);
        const URL_ENDPOINT = '/api/crear-ticket';
        const token = user?.token;
        const email = decodedpayload?.email ?? '';

        if (!token || !email) {
            setAlertState({ // <-- 5. Usa setAlertState para errores
              open: true,
              title: 'Error de Autenticación',
              message: 'Token o email de usuario no disponible. Por favor, inicie sesión.',
              severity: 'error',
            });
            return;
        }

        try {
            const payload = transformarDatosParaApi(data, email);
            const response = await fetch(URL_ENDPOINT, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify(payload)
            });

            if (!response.ok) {
                const errorData = await response.json().catch(() => null);
                throw new Error(errorData?.error || 'Error en el backend al crear la solicitud.');
            }

            // --- 6. Muestra el pop-up de ÉXITO ---
            setAlertState({
              open: true,
              title: '¡Éxito!',
              message: `Solicitud para ${data.tipoSolicitud} enviada correctamente.`,
              severity: 'success',
            });
            reset(); // Resetea el formulario

        } catch (error) {
            console.error("Fallo al enviar el formulario:", error);
            // --- 7. Muestra el pop-up de ERROR ---
            setAlertState({
              open: true,
              title: 'Error al Enviar',
              message: error instanceof Error ? error.message : 'Fallo al enviar el formulario.',
              severity: 'error',
            });
        }
    };

    const renderFormularioEspecifico = (): React.JSX.Element => {
        switch (tipoActual) {
            case 'Edificio':
                return <EdificioForm />;
            case 'Activo':
                return <ActivoForm />;
            case 'Técnico':
                return <TecnicoForm />;
            //case 'Sensor':
                //return <SensorForm />;
            default:
                return <Alert severity="error">Seleccione un tipo de solicitud válido.</Alert>;
        }
    };

    return (
        <Paper elevation={3} sx={{ p: 4, width: "80%", margin: '0 auto' }}>
            <CustomAlert
              title={alertState.title}
              message={alertState.message}
              severity={alertState.severity}
              open={alertState.open}
              onClose={() => setAlertState(prev => ({ ...prev, open: false }))}
            />
            <Typography variant="h5" component="h1" gutterBottom>
                Crear Nueva Solicitud
            </Typography>
            <Divider sx={{ my: 2 }} />
            <FormProvider {...methods}>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Stack spacing={3}>
                        
                        <Grid container spacing={3}> {/* 4. Usa Grid para layout de 2 columnas */}
                            <Grid size={{xs:12, sm: 6}}>
                                <FormControl fullWidth required error={Boolean(errors.tipoSolicitud)}>
                                    <InputLabel id="tipo-solicitud-label">Tipo de Entidad</InputLabel>
                                    <Select
                                        labelId="tipo-solicitud-label"
                                        label="Tipo de Entidad"
                                        defaultValue="Edificio"
                                        {...register("tipoSolicitud", { required: "Debe seleccionar un tipo" })}
                                    >
                                        {TIPOS_SOLICITUD.map(tipo => (
                                            <MenuItem key={tipo} value={tipo}>{tipo}</MenuItem>
                                        ))}
                                    </Select>
                                </FormControl>
                            </Grid>
                            
                            <Grid size={{xs:12, sm: 6}}>
                                {/* --- 5. NUEVO CAMPO 'TIPO OPERACIÓN' --- */}
                                <FormControl fullWidth required error={Boolean(errors.tipoOperacion)}>
                                    <InputLabel id="tipo-operacion-label">Tipo de Operación</InputLabel>
                                    <Select
                                        labelId="tipo-operacion-label"
                                        label="Tipo de Operación"
                                        defaultValue="Ingreso"
                                        {...register("tipoOperacion", { required: "Debe seleccionar una operación" })}
                                    >
                                        {TIPOS_OPERACION.map(tipo => (
                                            <MenuItem key={tipo} value={tipo}>{tipo}</MenuItem>
                                        ))}
                                    </Select>
                                </FormControl>
                            </Grid>
                        </Grid>

                        <TextField
                            label="Asunto / Título"
                            fullWidth
                            required
                            {...register("asunto", { required: "El Asunto es obligatorio" })}
                            error={Boolean(errors.asunto)}
                            helperText={errors.asunto?.message}
                        />

                        <TextField
                            label="Justificación / Detalles"
                            fullWidth
                            multiline
                            rows={3}
                            {...register("detalles")}
                            helperText="Esta información se usará como la 'justificación' de la solicitud."
                        />
                        
                        <Divider sx={{ my: 2 }} />

                        <Typography variant="h6" gutterBottom>
                            Datos Específicos ({tipoActual})
                        </Typography>

                        {renderFormularioEspecifico()}

                        <Button 
                            type="submit" 
                            variant="contained" 
                            fullWidth
                            disabled={isSubmitting}
                            sx={{ mt: 3 }}
                        >
                            {isSubmitting ? <CircularProgress size={24} color="inherit" /> : 'Enviar Solicitud'}
                        </Button>
                    </Stack>
                </form>
            </FormProvider>
        </Paper>
    );
}