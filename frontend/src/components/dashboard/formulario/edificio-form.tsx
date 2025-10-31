import * as React from 'react';
import { TextField, Grid } from '@mui/material';
import { useFormContext } from 'react-hook-form'; // Importante para la conexión
import type { SolicitudFormData, EdificioData } from '@/types/formulario';

export function EdificioForm(): React.JSX.Element {
    const {
        register,
        formState: { errors },
    } = useFormContext<SolicitudFormData>();

    const edificioErrors = errors.datosEspecificos as Partial<Record<keyof EdificioData, any>>;

    return (
        <Grid container spacing={2}>
            <Grid size={{xs:12}}>
                <TextField
                label="Nombre del Edificio"
                fullWidth
                required
                {...register("datosEspecificos.nombre" as const)}
                error={!!edificioErrors?.nombre}
                helperText={edificioErrors?.nombre?.message ?? ""}
                />
            </Grid>

            <Grid size={{xs:12}}>
                <TextField
                label="Dirección del Edificio"
                fullWidth
                required
                {...register("datosEspecificos.direccion" as const)}
                error={!!edificioErrors?.direccion}
                helperText={edificioErrors?.direccion?.message ?? ""}
                />
            </Grid>

            <Grid size={{xs:12}}>
                <TextField
                label="Latitud (Aprox.)"
                fullWidth
                required
                type="number"
                {...register("datosEspecificos.latitud" as const, {
                    valueAsNumber: true,
                })}
                error={!!edificioErrors?.latitud}
                helperText={edificioErrors?.latitud?.message ?? ""}
                slotProps={{
                    htmlInput: {
                        step: "any",
                    },
                    inputLabel: {
                        shrink: true,
                    }
                }}
                />
            </Grid>

            <Grid size={{xs:12}}>
                <TextField
                label="Longitud (Aprox.)"
                fullWidth
                required
                type="number"
                {...register("datosEspecificos.longitud" as const, {
                    valueAsNumber: true,
                })}
                error={!!edificioErrors?.longitud}
                helperText={edificioErrors?.longitud?.message ?? ""}
                slotProps={{
                    htmlInput: {
                        step: "any",
                    },
                    inputLabel: {
                        shrink: true,
                    }
                }}
                />
            </Grid>
        </Grid>
    );
}