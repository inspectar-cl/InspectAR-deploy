export type TipoSolicitud = 'Edificio' | 'Activo' | 'Técnico';
export type TipoActivo = 'Ascensor' | 'BombaDeAgua' | 'PanelElectrico';
export type EspecialidadTecnico = 'Climatización' | 'Eléctrico' | 'Mecánico';

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
  imagen: string | null;
}
export interface TecnicoData {
  nombre: string;
  activosAsociados: number[];
  correo: string;
  telefono: string;
  especialidad: string;
}

export interface SolicitudAPI {
  id: number;
  estado: 'Pendiente' | 'EnProgreso' | 'Resuelta' | string;
  fechaCreacion: string;
  tipoSolicitud: TipoSolicitud;
  tipoOperacion: string;
  asunto: string;
  detalles?: string;
  datosEspecificos: EdificioData | ActivoDataAPI | TecnicoData;
}

export interface ApiTicket {
  id: number;
  tipo_entidad: 'edificio' | 'activo' | 'tecnico' | string;
  tipo_operacion: 'ingreso' | 'modificacion' | 'eliminacion' | string;
  estado: 'resuelto' | 'pendiente' | 'en progreso' | string;
  usuario_email: string;
  justificacion: string;
  comentario_admin?: string;
  created_at: string;
  fecha_resolucion?: string;
  resuelto_por_email?: string;
  
  edificio_id?: number;
  edificio_nombre?: string;
  edificio_direccion?: string;
  edificio_latitud?: number;
  edificio_longitud?: number;

  activo_id?: number;
  activo_nombre?: string;
  activo_tipo?: string;
  activo_descripcion?: string;
  activo_ubicacion?: string;
  activo_edificio_id?: number;
  
  tecnico_id?: number;
  tecnico_nombre?: string;
  tecnico_email?: string;
  tecnico_telefono?: string;
  tecnico_especialidad?: string;
  tecnico_autorizado?: boolean;
}

export interface ApiSolicitudResponse {
  tickets: ApiTicket[];
  total_tickets: number;
  total_paginas: number;
  pagina_actual: number;
  items_por_pagina: number;
}