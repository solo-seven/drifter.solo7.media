// hooks/usePlanetAPI.ts
import { useState } from 'react';
import { PlanetParameters } from '@/types/planet';

export const usePlanetAPI = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const generatePlanet = async (params: PlanetParameters) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('/api/planet/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(params),
      });
      
      if (!response.ok) throw new Error('Generation failed');
      
      const data = await response.json();
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { generatePlanet, loading, error };
};