export interface TendenciaDataPoint {
  timestamp: string;
  valor: number;
  prediccion?: number;
  limite_superior?: number;
  limite_inferior?: number;
}

export interface Variable {
  id: number;
  nombre: string;
  valor: number;
  unidad: string;
  tendencia: 'ascendente' | 'descendente' | 'estable';
  cambio: number;
  estado: 'normal' | 'advertencia' | 'critico';
  activoId: number;
  data: TendenciaDataPoint[];
}

export interface TendenciaResumen {
  totalVariables: number;
  variablesAdvertencia: number;
  variablesCriticas: number;
  cambioPromedio: number;
  ultimaActualizacion: string;
}