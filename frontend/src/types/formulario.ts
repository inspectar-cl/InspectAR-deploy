export type TipoSolicitud = 'Edificio' | 'Activo' | 'Sensor' | 'Técnico';

export interface EdificioData {
  nombre: string;
  direccion: string;
  latitud: number;
  longitud: number;
}

export type TipoActivo = 'Ascensor' | 'BombaDeAgua' | 'PanelElectrico';
export interface ActivoData {
  tipoActivo: TipoActivo;
  edificioId: number;
  ubicacion: string;
  descripcion?: string;
  imagen: FileList | null;
}

export type EspecialidadTecnico = 'Climatización' | 'Eléctrico' | 'Mecánico';
export interface TecnicoData {
  nombre: string;
  activosAsociados: number[]; // Array de IDs de Activos asociados
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