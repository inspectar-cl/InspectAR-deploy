import { paths } from './paths';

export const rolePermissions: Record<string, string[]> = {
  analista: [
    paths.dashboard.inicio,
    paths.dashboard.activos,
    paths.dashboard.reportesTecnicos,
    paths.dashboard.ml,
    paths.dashboard.account,
    paths.dashboard.settings,
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
    paths.dashboard.account,
    paths.dashboard.settings,
    paths.dashboard.reporteFallas,
  ],
  administrador: [
    paths.dashboard.inicio,
    paths.dashboard.account,
    paths.dashboard.settings,
    paths.dashboard.formulario,
    paths.dashboard.activoDetail(':id'),
    paths.dashboard.contactos,
    paths.dashboard.reportesTecnicos,
    paths.dashboard.sensores,
    paths.dashboard.ml,
  ],         
  root: [
    paths.dashboard.inicio,
    paths.dashboard.activos,
    paths.dashboard.sensores,
    paths.dashboard.reportesTecnicos,
    paths.dashboard.contactos,
    paths.dashboard.activoDetail(':id'),
    paths.dashboard.solicitudes,
    paths.dashboard.agregardatos,
    paths.dashboard.gestion,
    paths.dashboard.ml,
  ]
};
