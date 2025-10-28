'use client'; 

import { useCallback, useRef } from 'react';
import { driver, Driver, Config } from 'driver.js';
import { getTourSteps, TourKey, TourSteps } from './tour-config';
import 'driver.js/dist/driver.css';

const DEFAULT_DRIVER_CONFIG: Config = {
  showProgress: true,
  allowClose: true,
  nextBtnText: 'Siguiente',
  prevBtnText: 'Anterior',
  doneBtnText: 'Cerrar',
};

export const useDriverTour = (tourKey: TourKey) => {
    const driverRef = useRef<Driver | null>(null);
    
    const startTour = useCallback(() => {
        const steps: TourSteps = getTourSteps(tourKey);

        console.log('Iniciando tour con pasos:', steps);

        if (steps.length === 0) {
            console.warn(`No se encontraron pasos para la clave: ${tourKey}`);
            return;
        }

        if (driverRef.current) {
            driverRef.current.destroy();
        }
        
        // Se crea una nueva instancia de Driver con los pasos específicos
        const driverObj = driver({
            ...DEFAULT_DRIVER_CONFIG,
            steps: steps,
        });

        driverRef.current = driverObj;
        driverObj.drive(); 
        
    }, [tourKey]);
    
    return { startTour };
};