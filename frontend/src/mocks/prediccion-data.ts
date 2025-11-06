import type { Prediccion } from "@/types/prediccion";

interface Activo {
  id: number;
  nombre: string;
  tipo: string;
  ubicacion: string;
}

export const mockActivos: Activo[] = [
  {
    id: 1,
    nombre: "Bomba Centrífuga B",
    tipo: "Bomba",
    ubicacion: "Planta Principal - Sector A"
  },
  {
    id: 2,
    nombre: "Ascensor Central",
    tipo: "Ascensor",
    ubicacion: "Planta Principal - Sector B"
  },
  {
    id: 3,
    nombre: "Circuito Eléctrico C",
    tipo: "Circuito Eléctrico",
    ubicacion: "Planta Secundaria"
  }
];

export const mockPredicciones: Prediccion[] = [
  // Activo 1 - Bomba Centrífuga A
  {
    id: 1,
    activoId: 1,
    timestamp: new Date(Date.now() - 1000 * 60 * 5).toISOString(),
    anomalyScore: 85,
    anomalyLikelihood: 92,
    severidad: "Alta",
    descripcion: "Vibración excesiva detectada",
    threshold: 75,
    is_anomaly: true,
    most_influential_variable: "Vibración",
    contribution_magnitude: 89
  },
  {
    id: 2,
    activoId: 1,
    timestamp: new Date(Date.now() - 1000 * 60 * 15).toISOString(),
    anomalyScore: 72,
    anomalyLikelihood: 78,
    severidad: "Media",
    descripcion: "Temperatura del motor por encima del rango normal",
    threshold: 75,
    is_anomaly: true,
    most_influential_variable: "Temperatura",
    contribution_magnitude: 72
  },
  {
    id: 3,
    activoId: 1,
    timestamp: new Date(Date.now() - 1000 * 60 * 30).toISOString(),
    anomalyScore: 45,
    anomalyLikelihood: 38,
    severidad: "Baja",
    descripcion: "Operación normal",
    threshold: 75,
    is_anomaly: false,
    most_influential_variable: "Presión",
    contribution_magnitude: 12
  },
  // Activo 2 - Compresor B-12
  {
    id: 4,
    activoId: 2,
    timestamp: new Date(Date.now() - 1000 * 60 * 8).toISOString(),
    anomalyScore: 91,
    anomalyLikelihood: 95,
    severidad: "Alta",
    descripcion: "Presión de descarga anormalmente alta",
    threshold: 70,
    is_anomaly: true,
    most_influential_variable: "Presión",
    contribution_magnitude: 94
  },
  {
    id: 5,
    activoId: 2,
    timestamp: new Date(Date.now() - 1000 * 60 * 20).toISOString(),
    anomalyScore: 68,
    anomalyLikelihood: 71,
    severidad: "Media",
    descripcion: "Consumo energético elevado",
    threshold: 70,
    is_anomaly: true,
    most_influential_variable: "Corriente",
    contribution_magnitude: 68
  },
  {
    id: 6,
    activoId: 2,
    timestamp: new Date(Date.now() - 1000 * 60 * 40).toISOString(),
    anomalyScore: 52,
    anomalyLikelihood: 48,
    severidad: "Baja",
    descripcion: "Parámetros dentro de rangos aceptables",
    threshold: 70,
    is_anomaly: false,
    most_influential_variable: "Temperatura",
    contribution_magnitude: 22
  },
  // Activo 3 - Motor Eléctrico C
  {
    id: 7,
    activoId: 3,
    timestamp: new Date(Date.now() - 1000 * 60 * 3).toISOString(),
    anomalyScore: 78,
    anomalyLikelihood: 82,
    severidad: "Alta",
    descripcion: "Desbalance de fases detectado",
    threshold: 65,
    is_anomaly: true,
    most_influential_variable: "Voltaje",
    contribution_magnitude: 81
  },
  {
    id: 8,
    activoId: 3,
    timestamp: new Date(Date.now() - 1000 * 60 * 12).toISOString(),
    anomalyScore: 43,
    anomalyLikelihood: 39,
    severidad: "Baja",
    descripcion: "Funcionamiento estable",
    threshold: 65,
    is_anomaly: false,
    most_influential_variable: "Corriente",
    contribution_magnitude: 18
  },
  {
    id: 9,
    activoId: 3,
    timestamp: new Date(Date.now() - 1000 * 60 * 25).toISOString(),
    anomalyScore: 61,
    anomalyLikelihood: 65,
    severidad: "Media",
    descripcion: "Ligero incremento en temperatura de bobinados",
    threshold: 65,
    is_anomaly: false,
    most_influential_variable: "Temperatura",
    contribution_magnitude: 45
  },
  {
    id: 10,
    activoId: 1,
    timestamp: new Date(Date.now() - 1000 * 60 * 45).toISOString(),
    anomalyScore: 88,
    anomalyLikelihood: 90,
    severidad: "Alta",
    descripcion: "Cavitación detectada en impulsor",
    threshold: 75,
    is_anomaly: true,
    most_influential_variable: "Presión de succión",
    contribution_magnitude: 87
  }
];

// Función para simular delay de red
export const simulateNetworkDelay = (ms = 1000): Promise<void> => {
  return new Promise(resolve => {
    setTimeout(resolve, ms);
  });
};