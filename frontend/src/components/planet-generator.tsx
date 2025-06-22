import React, { useState } from 'react';
import { Globe, Play, RotateCcw, Save, Download } from 'lucide-react';
import type { MeshData } from '@/lib/types';
import { CollapsibleSection } from './ui/collapsible-section';
import SliderInput from './ui/slider-input';
import SelectInput from './ui/select-input';
import ToggleInput from './ui/toggle-input';

// Type definitions
interface PlanetParameters {
  basic: {
    seed: number;
    radius: number;
    meshType: 'icosidodecahedron' | 'icosahedron' | 'cube';
  };
  subdivision: {
    level: number;
    algorithm: 'loop' | 'catmull-clark' | 'butterfly';
    smoothing: boolean;
  };
  terrain: {
    enabled: boolean;
    scale: number;
    octaves: number;
    persistence: number;
    lacunarity: number;
    noiseType: 'perlin' | 'simplex' | 'cellular';
  };
  tectonics: {
    enabled: boolean;
    plateCount: number;
    plateSpeed: number;
    convergenceStrength: number;
    divergenceStrength: number;
  };
  erosion: {
    enabled: boolean;
    hydraulicIterations: number;
    thermalIterations: number;
    rainAmount: number;
    evaporationRate: number;
  };
}

interface PlanetGeneratorProps {
  onGenerate: (meshData: MeshData | null) => void;
  setIsLoading: (isLoading: boolean) => void;
}

// UI Components have been moved to their own files in the components/ui directory

const PlanetGenerator: React.FC<PlanetGeneratorProps> = ({ onGenerate, setIsLoading }) => {
  const [parameters, setParameters] = useState<PlanetParameters>({
    basic: {
      seed: Math.floor(Math.random() * 10000),
      radius: 1,
      meshType: 'icosidodecahedron',
    },
    subdivision: {
      level: 2,
      algorithm: 'loop',
      smoothing: true,
    },
    terrain: {
      enabled: true,
      scale: 2.5,
      octaves: 4,
      persistence: 0.5,
      lacunarity: 2.0,
      noiseType: 'perlin',
    },
    tectonics: {
      enabled: false,
      plateCount: 12,
      plateSpeed: 0.1,
      convergenceStrength: 0.5,
      divergenceStrength: 0.5,
    },
    erosion: {
      enabled: false,
      hydraulicIterations: 50,
      thermalIterations: 50,
      rainAmount: 1.0,
      evaporationRate: 0.5,
    },
  });

  const [isGenerating, setIsGenerating] = useState(false);
  const [generationStatus, setGenerationStatus] = useState('Ready to generate');
  const [generationError, setGenerationError] = useState<string | null>(null);

  const handleGenerate = async () => {
    if (isGenerating) return;
    setIsGenerating(true);
    setIsLoading(true);
    onGenerate(null);
    setGenerationStatus('Generating planet...');
    setGenerationError(null);

    try {
      const response = await fetch('http://localhost:8080/api/planet/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(parameters),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`API Error: ${response.status} ${errorText}`);
      }

      const data = await response.json();
      onGenerate(data.mesh as MeshData);
      setGenerationStatus('Planet generated successfully.');
    } catch (error) {
      const message = error instanceof Error ? error.message : 'An unknown error occurred.';
      setGenerationError(message);
      setGenerationStatus('Generation failed.');
    } finally {
      setIsGenerating(false);
      setIsLoading(false);
    }
  };

  const randomizeSeed = () => {
    setParameters({
      ...parameters,
      basic: { ...parameters.basic, seed: Math.floor(Math.random() * 10000) }
    });
  };

  return (
    <div className="h-full flex flex-col">
        {/* Header */}
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold text-gray-800 flex items-center gap-2">
              <Globe size={24} />
              Planet Generator
            </h1>
            <button
              onClick={handleGenerate}
              disabled={isGenerating}
              className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 disabled:bg-indigo-300 flex items-center gap-2 transition-colors"
            >
              <Play size={16} />
              {isGenerating ? 'Generating...' : 'Generate'}
            </button>
          </div>
          <div className="mt-3 text-sm text-gray-500">
            <p>{generationStatus}</p>
            {generationError && <p className="text-red-500 mt-1">{generationError}</p>}
          </div>
        </div>

        {/* Parameter Controls */}
        <div className="flex-1 overflow-y-auto p-4">
          <CollapsibleSection title="Basic Settings" defaultOpen>
            <div className="flex items-center gap-2 mb-4">
              <div className="flex-1">
                <SliderInput
                  label="Seed"
                  value={parameters.basic.seed}
                  onChange={(seed) => setParameters({ ...parameters, basic: { ...parameters.basic, seed } })}
                  min={0}
                  max={9999}
                />
              </div>
              <button
                onClick={randomizeSeed}
                className="p-2 mt-4 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
                title="Randomize Seed"
              >
                <RotateCcw size={18} />
              </button>
            </div>
            <SliderInput
              label="Radius"
              value={parameters.basic.radius}
              onChange={(radius) => setParameters({ ...parameters, basic: { ...parameters.basic, radius } })}
              min={0.5}
              max={5}
              step={0.1}
              unit=" units"
            />
            <SelectInput
              label="Base Mesh Type"
              value={parameters.basic.meshType}
              onChange={(meshType) =>
                setParameters({
                  ...parameters,
                  basic: { ...parameters.basic, meshType: meshType as PlanetParameters['basic']['meshType'] }
                })
              }
              options={[
                { value: 'icosidodecahedron', label: 'Icosidodecahedron' },
                { value: 'icosahedron', label: 'Icosahedron' },
                { value: 'cube', label: 'Cube (Normalized)' },
              ]}
            />
          </CollapsibleSection>

          <CollapsibleSection title="Subdivision">
            <SliderInput
              label="Subdivision Level"
              value={parameters.subdivision.level}
              onChange={(level) =>
                setParameters({
                  ...parameters,
                  subdivision: { ...parameters.subdivision, level }
                })
              }
              min={0}
              max={5}
            />
            <SelectInput
              label="Algorithm"
              value={parameters.subdivision.algorithm}
              onChange={(algorithm) =>
                setParameters({
                  ...parameters,
                  subdivision: { ...parameters.subdivision, algorithm: algorithm as PlanetParameters['subdivision']['algorithm'] }
                })
              }
              options={[
                { value: 'loop', label: 'Loop' },
                { value: 'catmull-clark', label: 'Catmull-Clark' },
                { value: 'butterfly', label: 'Butterfly' },
              ]}
            />
            <ToggleInput
              label="Smoothing"
              checked={parameters.subdivision.smoothing}
              onChange={(smoothing) =>
                setParameters({
                  ...parameters,
                  subdivision: { ...parameters.subdivision, smoothing }
                })
              }
            />
          </CollapsibleSection>

          <CollapsibleSection title="Terrain Generation">
            <ToggleInput
              label="Enable Terrain"
              checked={parameters.terrain.enabled}
              onChange={(enabled) =>
                setParameters({
                  ...parameters,
                  terrain: { ...parameters.terrain, enabled }
                })
              }
            />
            {parameters.terrain.enabled && (
              <>
                <SelectInput
                  label="Noise Type"
                  value={parameters.terrain.noiseType}
                  onChange={(noiseType) =>
                    setParameters({
                      ...parameters,
                      terrain: { ...parameters.terrain, noiseType: noiseType as PlanetParameters['terrain']['noiseType'] }
                    })
                  }
                  options={[
                    { value: 'perlin', label: 'Perlin Noise' },
                    { value: 'simplex', label: 'Simplex Noise' },
                    { value: 'cellular', label: 'Cellular (Worley) Noise' },
                  ]}
                />
                <SliderInput
                  label="Terrain Scale"
                  value={parameters.terrain.scale}
                  onChange={(scale) =>
                    setParameters({
                      ...parameters,
                      terrain: { ...parameters.terrain, scale }
                    })
                  }
                  min={0.1}
                  max={10}
                  step={0.1}
                />
                <SliderInput
                  label="Octaves"
                  value={parameters.terrain.octaves}
                  onChange={(octaves) =>
                    setParameters({
                      ...parameters,
                      terrain: { ...parameters.terrain, octaves }
                    })
                  }
                  min={1}
                  max={8}
                />
              </>
            )}
          </CollapsibleSection>

          <CollapsibleSection title="Tectonic Plates">
            <ToggleInput
              label="Enable Tectonics"
              checked={parameters.tectonics.enabled}
              onChange={(enabled) =>
                setParameters({
                  ...parameters,
                  tectonics: { ...parameters.tectonics, enabled }
                })
              }
            />
            {parameters.tectonics.enabled && (
              <SliderInput
                label="Plate Count"
                value={parameters.tectonics.plateCount}
                onChange={(plateCount) =>
                  setParameters({
                    ...parameters,
                    tectonics: { ...parameters.tectonics, plateCount }
                  })
                }
                min={2}
                max={50}
              />
            )}
          </CollapsibleSection>

          <CollapsibleSection title="Erosion Simulation">
            <ToggleInput
              label="Enable Erosion"
              checked={parameters.erosion.enabled}
              onChange={(enabled) =>
                setParameters({
                  ...parameters,
                  erosion: { ...parameters.erosion, enabled }
                })
              }
            />
            {parameters.erosion.enabled && (
              <>
                <SliderInput
                  label="Hydraulic Iterations"
                  value={parameters.erosion.hydraulicIterations}
                  onChange={(hydraulicIterations) =>
                    setParameters({
                      ...parameters,
                      erosion: { ...parameters.erosion, hydraulicIterations }
                    })
                  }
                  min={0}
                  max={200}
                  step={10}
                />
                <SliderInput
                  label="Rain Amount"
                  value={parameters.erosion.rainAmount}
                  onChange={(rainAmount) =>
                    setParameters({
                      ...parameters,
                      erosion: { ...parameters.erosion, rainAmount }
                    })
                  }
                  min={0.1}
                  max={5.0}
                  step={0.1}
                />
              </>
            )}
          </CollapsibleSection>

          {/* Export/Import Actions */}
          <div className="mt-6 pt-6 border-t border-gray-200">
            <div className="flex gap-2">
              <button className="flex-1 px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center justify-center gap-2 text-sm transition-colors">
                <Save size={16} />
                Save Config
              </button>
              <button className="flex-1 px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 flex items-center justify-center gap-2 text-sm transition-colors">
                <Download size={16} />
                Export Planet
              </button>
            </div>
          </div>
        </div>
    </div>
  );
};

export default PlanetGenerator;
