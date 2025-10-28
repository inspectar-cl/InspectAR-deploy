import { DriveStep } from 'driver.js';

//Palabras clave a utilizar para identificar los diferentes tours disponibles en la aplicación.
export type TourKey = 'home' | 'activo' | 'predicciones-filtros';

export type TourSteps = DriveStep[];

// Pasos de cada tour

const HOME_TOUR_STEPS: TourSteps = [
  { 
    element: '#header-nav', 
    popover: { title: 'Inicio Rápido', description: 'Aquí accedes a todas las secciones principales.', side: 'bottom' } 
  },
  { 
    element: '.new-feature-badge', 
    popover: { title: 'Nueva Función', description: '¡Explora la última actualización aquí!', side: 'left' } 
  },
];

const ACTIVO_TOUR_STEPS: TourSteps = [
  { 
    element: '.activo-summary-card', 
    popover: { title: 'Resumen de Activos', description: 'Vista general de tu inventario.', side: 'top' } 
  },
  // ...
];

const PREDICCIONES_FILTROS_TOUR_STEPS: TourSteps = [
    { 
        element: '#tour-filtro-activo', 
        popover: { 
            title: 'Filtro 1: Selección de Activo', 
            description: 'Usa este filtro para aislar las anomalías de una sola máquina (ej: Bomba de Agua).', 
            side: 'bottom' 
        } 
    },
    { 
        element: '#tour-filtro-severidad', 
        popover: { 
            title: 'Filtro 2: Nivel de Severidad', 
            description: 'Concentra tu atención en alertas de "Alta" o "Media" severidad para priorizar acciones.', 
            side: 'bottom' 
        } 
    },
];

//Mapeo de todas las claves de tour a sus respectivos pasos.

const ALL_TOUR_STEPS: Record<TourKey, TourSteps> = {
    'home': HOME_TOUR_STEPS,
    'activo': ACTIVO_TOUR_STEPS,
    'predicciones-filtros': PREDICCIONES_FILTROS_TOUR_STEPS,
};

// Devuelve los pasos del tour segun la clave proporcionada.
export const getTourSteps = (key: TourKey): TourSteps => {
  return ALL_TOUR_STEPS[key] || [];
};