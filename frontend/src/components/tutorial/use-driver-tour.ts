'use client'; 

import { useCallback, useRef } from 'react';
import { driver } from 'driver.js';
import type { Driver, Config } from 'driver.js';
import type { TourKey, TourSteps } from './tour-config';
import { getTourSteps } from './tour-config';
import 'driver.js/dist/driver.css';

const DEFAULT_DRIVER_CONFIG: Config = {
  showProgress: true,
  allowClose: true,
  nextBtnText: 'Siguiente',
  prevBtnText: 'Anterior',
  doneBtnText: 'Cerrar',
};

export const useDriverTour = (tourKey: TourKey): { startTour: () => void } => {
    const driverRef = useRef<Driver | null>(null);
    
    const startTour = useCallback(() => {
        const steps: TourSteps = getTourSteps(tourKey);

        if (steps.length === 0) {
            return;
        }

        if (driverRef.current) {
            driverRef.current.destroy();
        }
        
        // Se crea una nueva instancia de Driver con los pasos específicos
        const driverObj = driver({
            ...DEFAULT_DRIVER_CONFIG,
            steps,
        });

        driverRef.current = driverObj;
        driverObj.drive(); 
        
    }, [tourKey]);
    
    return { startTour };
};