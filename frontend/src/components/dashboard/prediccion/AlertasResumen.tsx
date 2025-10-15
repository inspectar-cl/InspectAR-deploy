/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import React from "react";
import { Grid, Card, CardContent, Typography, Chip } from "@mui/material";
import { severidadColor } from "@/utils/severityUtils";
import { Prediccion } from "@/types/prediccion";
import { color } from "@mui/system/palette/palette";
import { minWidth } from "@mui/system";

export const AlertasResumen: React.FC<{ predicciones: Prediccion[] }> = ({ predicciones }) => {
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
      {["Alta", "Media", "Baja"].map((sev) => (
        <Grid item xs={12} sm={6} md={4} key={sev}>
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
              <Typography variant="h5">{contarPorSeveridad(sev as any)}</Typography>
              <Chip label={sev} color={severidadColor(sev)} />
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
};


export default AlertasResumen;