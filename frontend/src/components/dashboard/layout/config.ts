import type { NavItemConfig } from '@/types/nav';
import { paths } from '@/paths';

export const navItems = [
  { key: 'inicio', title: 'Inicio', href: paths.dashboard.inicio, icon: 'chart-pie' },
  { key: 'alertas', title: 'Alertas', href: paths.dashboard.alertas, icon: 'warning-circle' },
  { key: 'activos', title: 'Lista de activos', href: paths.dashboard.activos, icon: 'building-apartment' },
  { key: 'activos2', title: 'Lista de activos 2', href: paths.dashboard.activos2, icon: 'building-apartment' },
  { key: 'graficos', title: 'Gráficos', href: paths.dashboard.graficos, icon: 'chart-line' },
  { key: 'informe', title: 'Informes', href: paths.dashboard.informe, icon: 'file-text' },
  { key: 'planos', title: 'Planos', href: paths.dashboard.planos, icon: 'blueprint' },
  { key: 'settings', title: 'Settings', href: paths.dashboard.settings, icon: 'gear-six' },
  { key: 'account', title: 'Account', href: paths.dashboard.account, icon: 'user' },
  { key: 'error', title: 'Error', href: paths.errors.notFound, icon: 'x-square' },
] satisfies NavItemConfig[];
