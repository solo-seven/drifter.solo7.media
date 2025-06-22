import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { MeshData } from '@/lib/types';
import PlanetGenerator from '../planet-generator';

const meta: Meta<typeof PlanetGenerator> = {
  title: 'Components/PlanetGenerator',
  component: PlanetGenerator,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
  },
};

export default meta;
type Story = StoryObj<typeof PlanetGenerator>;

// Interactive component that handles state
const PlanetGeneratorWithState = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [meshData, setMeshData] = useState<MeshData | null>(null);

  const handleGenerate = async (newMeshData: MeshData | null) => {
    setIsLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000));
    setMeshData(newMeshData);
    setIsLoading(false);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-lg overflow-hidden">
        <div className="p-6">
          <h1 className="text-2xl font-bold text-gray-800 mb-6">Planet Generator</h1>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="md:col-span-1">
              <PlanetGenerator 
                onGenerate={handleGenerate} 
                setIsLoading={setIsLoading} 
              />
            </div>
            <div className="md:col-span-2 bg-gray-100 rounded-lg p-4 flex items-center justify-center">
              {isLoading ? (
                <div className="text-gray-500">Generating planet...</div>
              ) : meshData ? (
                <div className="text-center">
                  <div className="w-64 h-64 mx-auto bg-indigo-100 rounded-full flex items-center justify-center mb-4">
                    <span className="text-indigo-600">3D Planet Preview</span>
                  </div>
                  <p className="text-sm text-gray-600">
                    Vertices: {meshData.vertices.length / 3}
                    <br />
                    Faces: {meshData.faces.length / 3}
                  </p>
                </div>
              ) : (
                <div className="text-gray-400">
                  Click "Generate" to create a new planet
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export const Default: Story = {
  render: () => <PlanetGeneratorWithState />,
};

// Mock data for the story
const mockMeshData: MeshData = {
  vertices: new Float32Array([
    0, 0, 0, 1, 0, 0, 0, 1, 0,
    1, 0, 0, 1, 1, 0, 0, 1, 0
  ]),
  normals: new Float32Array([
    0, 0, 1, 0, 0, 1, 0, 0, 1,
    0, 0, 1, 0, 0, 1, 0, 0, 1
  ]),
  faces: new Uint32Array([0, 1, 2, 3, 4, 5]),
  uvs: new Float32Array([
    0, 0, 1, 0, 0, 1,
    1, 0, 1, 1, 0, 1
  ]),
  colors: new Float32Array([
    1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1,
    1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1
  ])
};

export const WithMockData: Story = {
  render: () => {
    const [isLoading, setIsLoading] = useState(false);
    
    return (
      <div className="p-4">
        <PlanetGenerator 
          onGenerate={() => {}} 
          setIsLoading={setIsLoading} 
        />
      </div>
    );
  },
  parameters: {
    docs: {
      description: {
        story: 'This shows the PlanetGenerator component with mock data. In a real application, the `onGenerate` prop would handle the generated mesh data.'
      }
    }
  }
};
