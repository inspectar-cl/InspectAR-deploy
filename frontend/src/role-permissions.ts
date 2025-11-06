import { paths } from './paths';

export const rolePermissions: Record<string, string[]> = {
  analista: [
    paths.dashboard.inicio,
    paths.dashboard.activos,
    paths.dashboard.reportesTecnicos,
    paths.dashboard.ml,
    paths.dashboard.account,
    paths.dashboard.settings,
    paths.dashboard.formulario,
    paths.dashboard.solicitudes,
    paths.dashboard.gestion,
  ],
  tecnico: [
    paths.dashboard.inicio,
    paths.dashboard.activos,
    paths.dashboard.sensores,
    paths.dashboard.alertas,
    paths.dashboard.reportesTecnicos,
    paths.dashboard.account,
    paths.dashboard.settings,
  ],
  residente: [
    paths.dashboard.inicio,
    paths.dashboard.contactos,
    paths.dashboard.account,
    paths.dashboard.settings,
    paths.dashboard.reporteFallas,
    paths.dashboard.formulario,
    paths.dashboard.solicitudes,
    paths.dashboard.gestion,
  ],           
  root: [
    paths.dashboard.inicio,
    paths.dashboard.account,
    paths.dashboard.alertas,
    paths.dashboard.activos,
    paths.dashboard.sensores,
    paths.dashboard.reportesTecnicos,
    paths.dashboard.reporteFallas,
    paths.dashboard.contactos,
    paths.dashboard.graficos,
    paths.dashboard.settings,
    paths.dashboard.activoDetail(':id'),
    paths.dashboard.formulario,
    paths.dashboard.solicitudes,
    paths.dashboard.agregardatos,
    paths.dashboard.gestion,
    paths.errors.notFound
  ]
};
