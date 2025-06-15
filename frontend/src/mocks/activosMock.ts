import { Activo } from '@/types'

export const activosMock: Activo[] = [
  { id: 'A1', tipoActivo: 'Ascensor #1', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago' },
  { id: 'A2', tipoActivo: 'Ascensor #2', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago' },
  { id: 'A3', tipoActivo: 'Ascensor #3', estado: 'Medio', descripcion: 'Mantenimiento programado no realizado', ubicacion: 'Edificio A, Santiago' },
  { id: 'A4', tipoActivo: 'Ascensor #4', estado: 'Medio', descripcion: 'Puerta atascada, requiere revisión', ubicacion: 'Edificio B, Santiago' },
  { id: 'C1', tipoActivo: 'Caldera #1', estado: 'Crítico', descripcion: 'Temperatura fuera de rango, riesgo de daño', ubicacion: 'Edificio A, Santiago' },
  { id: 'C2', tipoActivo: 'Caldera #2', estado: 'Crítico', descripcion: 'Anomalía detectada – revisar urgentemente', ubicacion: 'Edificio B, Santiago' },
  { id: 'C3', tipoActivo: 'Caldera #3', estado: 'Crítico', descripcion: 'Fuga de gas detectada', ubicacion: 'Edificio C, Santiago' },
  { id: 'C4', tipoActivo: 'Caldera #4', estado: 'Crítico', descripcion: 'Presión excesiva, riesgo de explosión', ubicacion: 'Edificio C, Santiago' },
  { id: 'B1', tipoActivo: 'Bomba de Agua #1', estado: 'Crítico', descripcion: 'Motor sobrecalentado, riesgo de falla', ubicacion: 'Edificio A, Santiago' },
  { id: 'B2', tipoActivo: 'Bomba de Agua #2', estado: 'Crítico', descripcion: 'Fuga de presión, riesgo de parada', ubicacion: 'Edificio A, Santiago' },
  { id: 'B3', tipoActivo: 'Bomba de Agua #3', estado: 'Crítico', descripcion: 'Riesgo crítico de falla en 7 días', ubicacion: 'Edificio B, Santiago' },
  { id: 'B4', tipoActivo: 'Bomba de Agua #4', estado: 'Crítico', descripcion: 'Nivel de agua crítico', ubicacion: 'Edificio B, Santiago' },
  { id: 'E1', tipoActivo: 'Sistema Eléctrico #1', estado: 'Medio', descripcion: 'Sobrecarga detectada, monitorear consumo', ubicacion: 'Edificio A, Santiago' },
  { id: 'E2', tipoActivo: 'Sistema Eléctrico #2', estado: 'Medio', descripcion: 'Pico de voltaje registrado', ubicacion: 'Edificio B, Santiago' },
  { id: 'E3', tipoActivo: 'Sistema Eléctrico #3', estado: 'Medio', descripcion: 'Variación de frecuencia, revisar panel', ubicacion: 'Edificio C, Santiago' },
]