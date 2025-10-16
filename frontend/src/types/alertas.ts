export interface Alerta {
  id: string
  tipoActivo: string
  descripcion: string
  estado: 'Medio' | 'Crítico'
  ubicacion: string
  fecha: string
  atendida: boolean
}