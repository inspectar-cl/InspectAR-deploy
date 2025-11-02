import type { DriveStep } from 'driver.js';

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
    { 
        element: '#tour-metricas-principales', 
        popover: { 
            title: 'Metricas Principales', 
            description: 'Puedes observar diferentes puntajes otorgados por la IA para resumir su estado actual y futuro', 
            side: 'left',
            align: 'center',
        } 
    },
    { 
        element: '#tour-anomaly-score', 
        popover: { 
            title: 'Puntaje de Anomalia', 
            description: 'Mientras más alto sea el valor, significa un peor estado del activo. Procura revisar tu activo si este puntaje supera los 60 puntos', 
            side: 'bottom' 
        } 
    },
    { 
        element: '#tour-anomaly-likelihood', 
        popover: { 
            title: 'Probabilidad de anomalía', 
            description: 'Puntaje que describe el estado futuro del activo. Predice como serán los proximos valores de sensores del activo.', 
            side: 'bottom' 
        } 
    },
    { 
        element: '#tour-anomaly', 
        popover: { 
            title: 'Total de Anomalías', 
            description: 'Muestra cuantas anomalías se han detectado según la última información proporcionada de los sensores.', 
            side: 'bottom' 
        } 
    },
    { 
        element: '#tour-grafico-predicciones', 
        popover: { 
            title: 'Predicciones de tus activos', 
            description: 'Observa mayor detalle de fluctuaciones futuras del "Anomaly Score" y "Anomaly Likelihood"', 
            side: 'bottom' 
        } 
    },
    { 
        element: '#tour-ultimas-anomalias', 
        popover: { 
            title: 'Últimas anomalías detectadas', 
            description: 'Visualiza un resumen detallado de las últimas anomalías detectadas, la prioridad del problema y la probabilidad de ocurrencia. ¡Prioriza acciones y obten el control sobre tus activos!', 
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