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
    paths.dashboard.contactos,
    paths.dashboard.account,
    paths.dashboard.settings,
    paths.dashboard.reporteFallas,
    paths.dashboard.formulario,
    paths.dashboard.solicitudes,
    paths.dashboard.agregarDatos('edificio'),
    paths.dashboard.agregarDatos('sensor'),
    paths.dashboard.agregarDatos('tecnico'),
    paths.dashboard.agregarDatos('activo'),
  ],           
  admin: [
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
    paths.errors.notFound
  ]
};
