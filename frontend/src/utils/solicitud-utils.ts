import type { TipoSolicitud } from '@/types/formulario';

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