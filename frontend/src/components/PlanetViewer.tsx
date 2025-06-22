"use client";

import React from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Edges } from '@react-three/drei';
import * as THREE from 'three';

interface MeshData {
  vertices: { id: number; x: number; y: number; z: number }[];
  faces: { id: number; vertices: number[]; type: 'triangle' | 'pentagon' }[];
}

interface PlanetViewerProps {
  meshData: MeshData;
}

const PlanetViewer: React.FC<PlanetViewerProps> = ({ meshData }) => {
  if (!meshData) return null;

  const positions = new Float32Array(meshData.vertices.length * 3);
  meshData.vertices.forEach((v, i) => {
    positions[i * 3] = v.x;
    positions[i * 3 + 1] = v.y;
    positions[i * 3 + 2] = v.z;
  });

  const indices = new Uint32Array(meshData.faces.flatMap(f => {
    if (f.vertices.length === 3) {
      return f.vertices;
    } else if (f.vertices.length === 5) {
      return [
        f.vertices[0], f.vertices[1], f.vertices[4],
        f.vertices[1], f.vertices[2], f.vertices[3],
        f.vertices[1], f.vertices[3], f.vertices[4],
      ];
    }
    return [];
  }));

  const geometry = new THREE.BufferGeometry();
  geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
  geometry.setIndex(new THREE.BufferAttribute(indices, 1));
  geometry.computeVertexNormals();

  return (
    <Canvas style={{ background: '#111' }} camera={{ position: [0, 0, 3], fov: 75 }}>
      <ambientLight intensity={0.5} />
      <pointLight position={[10, 10, 10]} />
      <mesh geometry={geometry}>
        <meshStandardMaterial color="#444" wireframe={true} />
        <Edges>
          <lineBasicMaterial color="#fff" linewidth={2} />
        </Edges>
      </mesh>
      <OrbitControls />
    </Canvas>
  );
};

export default PlanetViewer;
