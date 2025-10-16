 
import React from "react";
import { Grid, Card, CardContent, Typography, Chip } from "@mui/material";
import { severidadColor } from "@/utils/SeverityUtils";
import { type Prediccion } from "@/types/prediccion";

// Cards de predicciones
function PrediccionCards({ predicciones }: { predicciones: Prediccion[] }): React.ReactElement {
  return (
    <Grid container spacing={2} mt={1}>
      {predicciones.map((p) => (
        <Grid size={{ xs: 12, sm: 4 }} key={p.id}>
          <Card>
            <CardContent>
              <Typography variant="subtitle2" color="textSecondary">
                {p.timestamp} - Activo ID: {p.activoId}
              </Typography>
              <Typography variant="body1">{p.descripcion}</Typography>
              <Chip
                label={p.severidad}
                color={severidadColor(p.severidad) as
                  | "primary"
                  | "secondary"
                  | "success"
                  | "error"
                  | "info"
                  | "warning"
                  | "default"}
              />
              <Typography variant="body2">
                Probabilidad: {(p.anomalyLikelihood).toFixed(1)}%
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
}
export default PrediccionCards;