import type {
  TipoSolicitud,
  EdificioData,
  ActivoData,
  TecnicoData,
} from '@/types/formulario'; // Asegúrate que esta ruta sea correcta
import type {
  ApiTicket,
  SolicitudAPI,
  ActivoDataAPI, // Asumo que este tipo está en form-solicitud.ts
} from '@/types/form-solicitud';

export const TIPO_TO_SLUG: Record<TipoSolicitud, string> = {
  'Edificio': 'edificio',
  'Activo': 'activo',
  //'Sensor': 'sensor',
  'Técnico': 'tecnico',
};

export const SLUG_TO_TIPO: Record<string, TipoSolicitud> = {
  'edificio': 'Edificio',
  'activo': 'Activo',
  //'sensor': 'Sensor',
  'tecnico': 'Técnico',
};

export const TABS = [
  { slug: 'edificio', label: 'Edificio' },
  { slug: 'activo', label: 'Activo' },
  //{ slug: 'sensor', label: 'Sensor' },
  { slug: 'tecnico', label: 'Técnico' },
];

export function transformarApiATipoFrontend(ticket: ApiTicket): SolicitudAPI {
  let tipoSolicitud: TipoSolicitud;
  let datosEspecificos: any = {};

  // Mapea el estado (Backend "resuelto" -> Frontend "Resuelta")
  let estadoFrontend: string = 'Pendiente'; // Default
  if (ticket.estado === 'resuelto') estadoFrontend = 'Resuelta';
  if (ticket.estado === 'en progreso') estadoFrontend = 'EnProgreso';
  // Puedes añadir más mapeos si es necesario

  // Mapea la entidad y extrae los datos específicos
  switch (ticket.tipo_entidad) {
    case 'edificio':
      tipoSolicitud = 'Edificio';
      datosEspecificos = {
        nombre: ticket.edificio_nombre,
        direccion: ticket.edificio_direccion,
        latitud: ticket.edificio_latitud,
        longitud: ticket.edificio_longitud,
      } as EdificioData;
      break;

    case 'activo':
      tipoSolicitud = 'Activo';
      // Mapea el tipo de activo (ej. 'bomba de agua' -> 'BombaDeAgua')
      let tipoActivoForm: string = ticket.activo_tipo || '';
      if (tipoActivoForm === 'bomba de agua') tipoActivoForm = 'BombaDeAgua';
      if (tipoActivoForm === 'ascensor') tipoActivoForm = 'Ascensor';
      if (tipoActivoForm === 'panel electrico') tipoActivoForm = 'PanelElectrico';

      datosEspecificos = {
        // Ojo: ActivoDataAPI espera 'tipoActivo' y 'edificioId'
        // Tu API de ticket parece enviarlos como 'activo_tipo' y 'activo_edificio_id'
        nombre: ticket.activo_nombre, // Asumiendo que ActivoDataAPI tiene 'nombre'
        tipoActivo: tipoActivoForm,
        edificioId: ticket.activo_edificio_id,
        ubicacion: ticket.activo_ubicacion,
        descripcion: ticket.activo_descripcion,
        imagen: null, // La API de tickets no envía imagen
      } as ActivoDataAPI;
      break;

    case 'tecnico':
      tipoSolicitud = 'Técnico';
      datosEspecificos = {
        nombre: ticket.tecnico_nombre,
        correo: ticket.tecnico_email,
        telefono: ticket.tecnico_telefono,
        especialidad: ticket.tecnico_especialidad,
        activosAsociados: [], // La API de tickets no envía esto
      } as TecnicoData;
      break;

    default:
      // Fallback por si llega un tipo no esperado
      tipoSolicitud = 'Edificio';
      datosEspecificos = { nombre: 'Error: Tipo no reconocido' };
  }

  return {
    id: ticket.id,
    estado: estadoFrontend,
    fechaCreacion: ticket.created_at,
    tipoSolicitud: tipoSolicitud,
    tipoOperacion: ticket.tipo_operacion,
    // Asunto autogenerado para que la card tenga un título
    asunto: `${ticket.tipo_operacion.toUpperCase()} DE ${ticket.tipo_entidad.toUpperCase()}`,
    detalles: ticket.justificacion, // La justificación es el detalle
    datosEspecificos: datosEspecificos,
  };
}