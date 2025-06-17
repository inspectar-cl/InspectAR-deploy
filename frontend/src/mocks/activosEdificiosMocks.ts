import { Activo } from '@/types'

export const activosMock: Activo[] = [
  { id: 'A1', tipoActivo: 'Ascensor', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago', img: "https://www.schindler.cl/content/dam/website/lac/images/modernizacion/cabinas/ascensor-interior-4.jpg/_jcr_content/renditions/original./ascensor-interior-4.jpg", id_edificio: "EA"},
  { id: 'A2', tipoActivo: 'Ascensor', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago', img: "https://img.archiexpo.es/images_ae/photo-g/190735-18379375.webp", id_edificio: "EA"},
  { id: 'A3', tipoActivo: 'Ascensor', estado: 'Medio', descripcion: 'Mantenimiento programado no realizado', ubicacion: 'Edificio A, Santiago', img: "", id_edificio: "EA"},
  { id: 'A4', tipoActivo: 'Ascensor', estado: 'Medio', descripcion: 'Puerta atascada, requiere revisión', ubicacion: 'Edificio B, Santiago', img: "", id_edificio: "EB"},
  { id: 'C1', tipoActivo: 'Caldera', estado: 'Crítico', descripcion: 'Temperatura fuera de rango, riesgo de daño', ubicacion: 'Edificio A, Santiago', img: "", id_edificio: "EA"},
  { id: 'C2', tipoActivo: 'Caldera', estado: 'Crítico', descripcion: 'Anomalía detectada – revisar urgentemente', ubicacion: 'Edificio B, Santiago', img: "", id_edificio: "EB"},
  { id: 'C3', tipoActivo: 'Caldera', estado: 'Crítico', descripcion: 'Fuga de gas detectada', ubicacion: 'Edificio C, Santiago', img: "", id_edificio: "EC"},
  { id: 'C4', tipoActivo: 'Caldera', estado: 'Crítico', descripcion: 'Presión excesiva, riesgo de explosión', ubicacion: 'Edificio C, Santiago', img: "", id_edificio: "EA"},
  { id: 'B1', tipoActivo: 'Bomba de Agua', estado: 'Crítico', descripcion: 'Motor sobrecalentado, riesgo de falla', ubicacion: 'Edificio A, Santiago', img: "https://www.sondagua.cl/blog/wp-content/uploads/2021/10/bomba-para-extraccion-de-agua.jpg", id_edificio: "EA"},
  { id: 'B2', tipoActivo: 'Bomba de Agua', estado: 'Crítico', descripcion: 'Fuga de presión, riesgo de parada', ubicacion: 'Edificio A, Santiago', img: "https://www.sondagua.cl/blog/wp-content/uploads/2021/10/bomba-para-extraccion-de-agua.jpg", id_edificio: "EA"},
  { id: 'B3', tipoActivo: 'Bomba de Agua', estado: 'Crítico', descripcion: 'Riesgo crítico de falla en 7 días', ubicacion: 'Edificio B, Santiago', img: "https://www.sondagua.cl/blog/wp-content/uploads/2021/10/bomba-para-extraccion-de-agua.jpg", id_edificio: "EB"},
  { id: 'B4', tipoActivo: 'Bomba de Agua', estado: 'Crítico', descripcion: 'Nivel de agua crítico', ubicacion: 'Edificio B, Santiago', img: "https://www.sondagua.cl/blog/wp-content/uploads/2021/10/bomba-para-extraccion-de-agua.jpg", id_edificio: "EB"},
  { id: 'E1', tipoActivo: 'Sistema Eléctrico', estado: 'Medio', descripcion: 'Sobrecarga detectada, monitorear consumo', ubicacion: 'Edificio A, Santiago', img: "", id_edificio: "EA"},
  { id: 'E2', tipoActivo: 'Sistema Eléctrico', estado: 'Medio', descripcion: 'Pico de voltaje registrado', ubicacion: 'Edificio B, Santiago', img: "", id_edificio: "EB"},
  { id: 'E3', tipoActivo: 'Sistema Eléctrico', estado: 'Medio', descripcion: 'Variación de frecuencia, revisar panel', ubicacion: 'Edificio C, Santiago', img: "", id_edificio: "EC"},
]