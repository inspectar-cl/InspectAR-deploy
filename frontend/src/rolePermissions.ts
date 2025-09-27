import { paths } from './paths';

export const rolePermissions: Record<string, string[]> = {
  analista: [
    paths.dashboard.inicio,
    paths.dashboard.account,
    paths.dashboard.alertas,
    paths.dashboard.activos,
    paths.dashboard.graficos,
    paths.dashboard.settings,
    paths.errors.notFound
  ],
  tecnico: [
    paths.dashboard.inicio,
    paths.dashboard.activos,
    paths.dashboard.alertas,
    paths.dashboard.contactos,
    paths.dashboard.settings,
    paths.errors.notFound
  ],
  residente: [
    paths.dashboard.inicio,
    paths.dashboard.account,
    paths.dashboard.settings,
    paths.errors.notFound
  ]
};
