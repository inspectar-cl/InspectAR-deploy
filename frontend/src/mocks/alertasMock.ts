import {Alerta, Activo} from '@/types'
import { activosMock } from '@/mocks/'

// Fecha estimada de ocurrencia: aleatoria entre 1 y 10 días desde hoy
const generarFechaProxima = (): string => {
  const hoy = new Date()
  hoy.setDate(hoy.getDate() + Math.floor(Math.random() * 10) + 1)
  return hoy.toLocaleDateString('es-CL')
}

const generarAlertas = (activos: Activo[]): Alerta[] =>
  activos
    .filter((a) => a.estado !== 'OK')
    .map((a) => ({
      id: a.id,
      tipoActivo: a.tipoActivo,
      descripcion: a.descripcion,
      estado: a.estado,
      ubicacion: a.ubicacion,
      fecha: generarFechaProxima(),
      atendida: false,
    }))

export const alertasMock: Alerta[] = generarAlertas(activosMock)
