export interface Activo {
  id: number
  tipoActivo: string
  estado: 'OK' | 'Medio' | 'Crítico' | 'NN'
  descripcion: string
  ubicacion: string
  img: string
  id_edificio: string
  id_ficha_tecnica: number
}