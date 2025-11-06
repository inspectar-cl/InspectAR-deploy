import * as React from 'react';
import { 
    Stack, Button, FormControl, InputLabel, Select, MenuItem, TextField, Typography, Alert, Paper , Divider, CircularProgress
} from '@mui/material';
import { useForm, FormProvider } from 'react-hook-form';
import { EdificioForm } from './edificio-form';
import { ActivoForm } from './activo-form';
import { TecnicoForm } from './tecnico-form';
//import { SensorForm } from './sensor-form';

import type { EdificioData, SolicitudFormData } from '@/types/formulario';

import { useUserToken } from '@/hooks/use-usertoken';

type TipoSolicitud = 'Edificio' | 'Activo' | 'Técnico';

const TIPOS_SOLICITUD: TipoSolicitud[] = ['Edificio', 'Activo', 'Técnico'];

export function FormularioSolicitud(): React.JSX.Element {
    const methods = useForm<SolicitudFormData>({
        defaultValues: {
            tipoSolicitud: 'Edificio',
            asunto: '',
            detalles: '',
            datosEspecificos: {} as EdificioData
        }
    });

    const { register, handleSubmit, watch, formState: { errors, isSubmitting } } = methods;

    const tipoActual = watch("tipoSolicitud"); // Observar el campo de tipo de solicitud
    const { user} = useUserToken();

    const onSubmit = async (data: SolicitudFormData): Promise<void> => {
        const URL_ENDPOINT = '/api/solicitudes/crear'; // URL del backend
        const token = user?.token;

        if (!token) {
            //console.error("Token no disponible.");
            return;
        }

        try {
            const payload = {
                ...data,
            };

            const response = await fetch(URL_ENDPOINT, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify(payload)
            });

            if (!response.ok) {
                throw new Error('Error en el backend al crear la solicitud.');
            }

            //console.log("Solicitud creada con éxito:", payload);
            //alert("Solicitud enviada con éxito.");
            methods.reset();
        } catch (error) {
            //console.error("Fallo al enviar el formulario:", error);
            //alert("Fallo al enviar el formulario. Vea la consola para más detalles.");
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
            <Typography variant="h5" component="h1" gutterBottom>
                Crear Nueva Solicitud
            </Typography>
            <Divider sx={{ my: 2 }} />
            <FormProvider {...methods}>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Stack spacing={3}>
                        
                        <FormControl fullWidth required error={Boolean(errors.tipoSolicitud)}>
                            <InputLabel id="tipo-solicitud-label">Tipo de Solicitud</InputLabel>
                            <Select
                                labelId="tipo-solicitud-label"
                                label="Tipo de Solicitud"
                                defaultValue="Edificio"
                                {...register("tipoSolicitud", { required: "Debe seleccionar un tipo" })}
                            >
                                {TIPOS_SOLICITUD.map(tipo => (
                                    <MenuItem key={tipo} value={tipo}>{tipo}</MenuItem>
                                ))}
                            </Select>
                            {errors.tipoSolicitud && <Typography color="error" variant="caption">{errors.tipoSolicitud.message}</Typography>}
                        </FormControl>

                        <TextField
                            label="Asunto / Título"
                            fullWidth
                            required
                            {...register("asunto", { required: "El Asunto es obligatorio" })}
                            error={Boolean(errors.asunto)}
                            helperText={errors.asunto?.message}
                        />

                        <TextField
                            label="Detalles de la Solicitud"
                            fullWidth
                            multiline
                            rows={3}
                            {...register("detalles")}
                        />
                        
                        <Divider sx={{ my: 2 }} />

                        <Typography variant="h6" gutterBottom>
                            Datos Específicos ({tipoActual})
                        </Typography>

                        {renderFormularioEspecifico()}

                        <Button 
                            type="submit" 
                            variant="contained" 
                            color="primary" 
                            fullWidth
                            disabled={isSubmitting}
                            sx={{ mt: 3 }}
                        >
                            {isSubmitting ? <CircularProgress  /> : 'Enviar Solicitud'}
                        </Button>
                    </Stack>
                </form>
            </FormProvider>
        </Paper>
    );
}