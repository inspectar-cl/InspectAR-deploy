/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import * as React from "react";
import { Grid, Card, CardContent, Typography, Chip } from "@mui/material";
import { severidadColor } from "@/utils/SeverityUtils";
import { type Prediccion } from "@/types/prediccion";

function AlertasResumen({ predicciones }: { predicciones: Prediccion[] }) {
  const contarPorSeveridad = (sev: "Baja" | "Media" | "Alta") =>
    predicciones.filter((p) => p.severidad === sev).length;

  return (
     <Grid
      container
      spacing={3}
      justifyContent="center"
      alignItems="stretch"
      sx={{ mb: 2 }}
    >
      {(["Alta", "Media", "Baja"] as ("Alta" | "Media" | "Baja")[]).map((sev) => (
        <Grid size={{ xs: 12, sm: 4 }} key={sev}>
          <Card
            sx={{
              display: "flex",
              flexDirection: "column",
              justifyContent: "center",
              alignItems: "center",
              height: 150,
              minWidth: 150,
              borderRadius: 4,
              boxShadow: 4,
            }}
          >
            <CardContent
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                textAlign: "center",
                gap: 1.5,
              }}
            >
              <Typography variant="subtitle1">{sev}</Typography>
              <Typography variant="h5">{contarPorSeveridad(sev)}</Typography>
              <Chip label={sev} color={severidadColor(sev) as "default" | "primary" | "secondary" | "error" | "info" | "success" | "warning"} />
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
};


export default AlertasResumen;