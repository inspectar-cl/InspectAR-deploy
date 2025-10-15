/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import React from "react";
import { Grid, Card, CardContent, Typography, Chip } from "@mui/material";
import { severidadColor } from "@/utils/severityUtils";
import { Prediccion } from "@/types/prediccion";

// Cards de predicciones
const PrediccionCards: React.FC<{ predicciones: Prediccion[] }> = ({ predicciones }) => (
  <Grid container spacing={2} mt={1}>
    {predicciones.map((p) => (
      <Grid item xs={12} sm={6} md={4} key={p.id}>
        <Card>
          <CardContent>
            <Typography variant="subtitle2" color="textSecondary">
              {p.timestamp} - Activo ID: {p.activoId}
            </Typography>
            <Typography variant="body1">{p.descripcion}</Typography>
            <Chip label={p.severidad} color={severidadColor(p.severidad)} />
            <Typography variant="body2">
              Probabilidad: {(p.anomalyLikelihood).toFixed(1)}%
            </Typography>
          </CardContent>
        </Card>
      </Grid>
    ))}
  </Grid>
);
export default PrediccionCards;