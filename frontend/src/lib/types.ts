export interface MeshData {
  vertices: { id: number; x: number; y: number; z: number }[];
  faces: { id: number; vertices: number[]; type: 'triangle' | 'pentagon' }[];
}
