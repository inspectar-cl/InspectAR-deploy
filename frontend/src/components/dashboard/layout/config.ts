import type { NavItemConfig } from '@/types/nav';
import { paths } from '@/paths';

// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment -- NavItemConfig requires string icons
export const navItems = [
  { key: 'inicio', title: 'Inicio', href: paths.dashboard.inicio, icon: 'chart-pie' },
  { key: 'formulario', title: 'Formulario', href: paths.dashboard.formulario, icon: 'clipboard-text' },
  { key: 'solicitudes', title: 'Solicitudes', href: paths.dashboard.solicitudes, icon: 'clipboard-text' },
  { key: 'agregar-datos', title: 'Agregar Datos', href: paths.dashboard.agregardatos, icon: 'cloud-arrow-up' },
  { key: 'editar-datos', title: 'Gestionar Datos', href: paths.dashboard.gestion, icon: 'file-cloud' },
  { key: 'alertas', title: 'Alertas', href: paths.dashboard.alertas, icon: 'warning-circle' },
  { key: 'activos', title: 'Lista de activos', href: paths.dashboard.activos, icon: 'building-apartment' },
  { key: 'sensores', title: 'Sensores', href: paths.dashboard.sensores, icon: 'broadcast' },
  { key: 'graficos', title: 'Gráficos', href: paths.dashboard.graficos, icon: 'chart-line' },
  { key: 'reportes-tecnicos', title: 'Reportes técnicos', href: paths.dashboard.reportesTecnicos, icon: 'file-text' },
  { key: 'reporte-fallas', title: 'Reporte de fallas', href: paths.dashboard.reporteFallas, icon: 'siren' },
<<<<<<< HEAD
=======
  //{ key: 'planos', title: 'Planos', href: paths.dashboard.planos, icon: 'blueprint' },
>>>>>>> origin/feature-front
  { key: 'contactos', title: 'Contactos', href: paths.dashboard.contactos, icon: 'adressBook' },
  { key: 'settings', title: 'Settings', href: paths.dashboard.settings, icon: 'gear-six' },
  { key: 'account', title: 'Account', href: paths.dashboard.account, icon: 'user' },
  { key: 'ml', title: 'Predicción', href: paths.dashboard.ml, icon: 'brain' },
  { key: 'tendencias', title: 'Tendencias', href: paths.dashboard.tendencias, icon: 'trend-up' },

] satisfies NavItemConfig[];
