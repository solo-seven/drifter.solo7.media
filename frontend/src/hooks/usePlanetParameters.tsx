// hooks/usePlanetParameters.ts
import { useEffect, useState } from 'react';
import { PlanetParameters } from '@/types/planet';

const STORAGE_KEY = 'planet-generator-params';

export const usePlanetParameters = (defaultParams: PlanetParameters) => {
  const [parameters, setParameters] = useState<PlanetParameters>(() => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(STORAGE_KEY);
      return saved ? JSON.parse(saved) : defaultParams;
    }
    return defaultParams;
  });

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(parameters));
  }, [parameters]);

  return [parameters, setParameters] as const;
};