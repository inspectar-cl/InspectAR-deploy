// types/formulario.ts

export type TipoSolicitud = 'Edificio' | 'Activo' | 'Técnico';
export type TipoActivo = 'Ascensor' | 'BombaDeAgua' | 'PanelElectrico';
export type EspecialidadTecnico = 'Climatización' | 'Eléctrico' | 'Mecánico';
export type TipoSensor = 'Temperatura' | 'Presión' | 'Vibración';

// --- 1. AÑADE ESTE NUEVO TIPO ---
export type TipoOperacion = 'Ingreso' | 'Modificacion' | 'Eliminacion';

export interface EdificioData {
  nombre: string;
  direccion: string;
  latitud: number;
  longitud: number;
}

export interface ActivoData {
  nombre: string;
  tipoActivo: TipoActivo;
  edificioId: number;
  ubicacion: string;
  descripcion?: string;
  imagen: FileList | null;
}

export interface TecnicoData {
  nombre: string;
  activosAsociados: number[]; // Array de IDs de Activos asociados
  correo: string;
  telefono: string;
  especialidad: string;
}

export interface SensorData {
  nombre: string;
  tipoSnsor: TipoSensor;
  activoAsociadoId: number;
}

export interface SolicitudFormData {
  tipoSolicitud: TipoSolicitud;
  tipoOperacion: TipoOperacion; 
  asunto: string;
  detalles?: string;
  datosEspecificos: EdificioData | ActivoData | TecnicoData | SensorData;
}