export type TipoSolicitud = 'Edificio' | 'Activo' | 'Técnico';

export interface EdificioData {
  nombre: string;
  direccion: string;
  latitud: number;
  longitud: number;
}

export type TipoActivo = 'Ascensor' | 'BombaDeAgua' | 'PanelElectrico';
export interface ActivoData {
  nombre: string | null;
  tipoActivo: TipoActivo;
  edificioId: number;
  ubicacion: string;
  descripcion?: string;
}

export type EspecialidadTecnico = 'Climatización' | 'Eléctrico' | 'Mecánico';
export interface TecnicoData {
  nombre: string;
  correo: string;
  telefono: string;
  especialidad: string;
}

export type TipoSensor = 'Temperatura' | 'Presión' | 'Vibración'; // Cambiar si es necesario
export interface SensorData {
  nombre: string;
  tipoSnsor: TipoSensor;
  activoAsociadoId: number;
}

export interface SolicitudFormData {
  tipoSolicitud: TipoSolicitud;
  asunto: string;
  detalles?: string;
  datosEspecificos: EdificioData | ActivoData | TecnicoData | SensorData;
}