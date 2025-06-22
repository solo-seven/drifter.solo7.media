"use client";

import React, { useState } from 'react';
import PlanetViewer from '@/components/PlanetViewer';
import PlanetGenerator from '@/components/planet-generator'; // Corrected import path
import type { MeshData } from '@/lib/types'; // Assuming types are defined here

const PlanetPage: React.FC = () => {
  const [meshData, setMeshData] = useState<MeshData | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  return (
    <div className="flex h-screen bg-gray-100">
      <div className="w-[400px] flex-shrink-0 bg-white overflow-y-auto shadow-lg">
        <PlanetGenerator 
          onGenerate={setMeshData} 
          setIsLoading={setIsLoading} 
        />
      </div>
      <div className="flex-1 relative">
        {isLoading && (
          <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center z-10">
            <p className="text-white text-2xl">Generating Planet...</p>
          </div>
        )}
        <PlanetViewer meshData={meshData} />
      </div>
    </div>
  );
};

export default PlanetPage;
