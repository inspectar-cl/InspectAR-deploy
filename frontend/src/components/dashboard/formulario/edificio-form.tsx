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
            <Grid size={{xs:12}}>
                <TextField
                label="Nombre del Edificio"
                fullWidth
                required
                {...register("datosEspecificos.nombre" as const)}
                error={Boolean(edificioErrors?.nombre)}
                helperText={edificioErrors?.nombre?.message ?? ''}
                />
            </Grid>

            <Grid size={{xs:12}}>
                <TextField
                label="Dirección del Edificio"
                fullWidth
                required
                {...register("datosEspecificos.direccion" as const)}
                error={Boolean(edificioErrors?.direccion)}
                helperText={edificioErrors?.direccion?.message ?? ''}
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
                error={Boolean(edificioErrors?.latitud)}
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
                error={Boolean(edificioErrors?.longitud)}
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