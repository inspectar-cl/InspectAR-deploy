//Este archivo type es un test de como se OBTENDRIAN los campos del backend, considerar igual para guiarse "formulario.ts", este ultimo como principal para el guardado de datos en backend

export type TipoSolicitud = 'Edificio' | 'Activo' | 'Técnico';
export type TipoActivo = 'Ascensor' | 'BombaDeAgua' | 'PanelElectrico';
export type EspecialidadTecnico = 'Climatización' | 'Eléctrico' | 'Mecánico';
export type TipoSensor = 'Temperatura' | 'Presión' | 'Vibración';

export interface EdificioData {
  nombre: string;
  direccion: string;
  latitud: number;
  longitud: number;
}
export interface ActivoDataAPI {
  tipoActivo: TipoActivo;
  edificioId: number;
  ubicacion: string;
  descripcion?: string;
  imagen: string | null; //Ojo con este campo, no se como se hará para las imagenes en backend
}
export interface TecnicoData {
  nombre: string;
  activosAsociados: number[];
  correo: string;
  telefono: string;
  especialidad: string;
}
export interface SensorData {
  nombre: string;
  tipoSnsor: TipoSensor;
  activoAsociadoId: number;
}

// El tipo principal de solicitud a recibir
export interface SolicitudAPI {
  id: number;
  estado: 'Pendiente' | 'EnProgreso' | 'Resuelta';
  fechaCreacion: string;
  tipoSolicitud: TipoSolicitud;
  asunto: string;
  detalles?: string;
  datosEspecificos: EdificioData | ActivoDataAPI | TecnicoData | SensorData;
  // solicitanteId: string;
}