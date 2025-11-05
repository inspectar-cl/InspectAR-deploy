interface Prediccion {
  id: number;
  activoId: number;
  timestamp: string;
  anomalyScore: number;
  anomalyLikelihood: number;
  severidad: "Baja" | "Media" | "Alta";
  descripcion: string;
  threshold: number;
  is_anomaly: boolean;
  most_influential_variable: string;
  contribution_magnitude: number;
}
export type { Prediccion };