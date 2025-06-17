export type Activo = {
  id: string
  tipoActivo: string
  estado: 'OK' | 'Medio' | 'Crítico' | 'NN'
  descripcion: string
  ubicacion: string
  img: string
  id_edificio: string
}