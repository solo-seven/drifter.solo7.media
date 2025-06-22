## Epic 1: Basic Planet Generation

### User Story 1.1: Generate Basic Icosidodecahedron Mesh

**As a user, I want to generate a basic icosidodecahedron mesh so that I can have a starting point for my planet**

#### User Workflow:
1. User opens the planet generator application
2. User clicks "New Planet" button
3. System generates a default icosidodecahedron mesh
4. User sees the wireframe mesh rendered in the 3D viewport
5. User can rotate/zoom the camera to inspect the mesh

#### Acceptance Criteria:
- [ ] Icosidodecahedron has 12 pentagonal and 20 triangular faces
- [ ] Mesh is normalized to unit sphere radius
- [ ] Vertices are evenly distributed on the sphere surface
- [ ] Mesh data includes vertices, faces, and edges
- [ ] Rendering shows clean wireframe without artifacts
- [ ] Camera controls allow 360° rotation and zoom

#### API Design:

**Request:** `POST /api/planet/generate`
```json
{
  "type": "icosidodecahedron",
  "radius": 1.0,
  "seed": 12345
}
```

**Response:** `200 OK`
```json
{
  "planetId": "uuid-here",
  "mesh": {
    "vertices": [
      {"id": 0, "x": 1.0, "y": 0.0, "z": 0.0},
      {"id": 1, "x": 0.809, "y": 0.588, "z": 0.0}
    ],
    "faces": [
      {"id": 0, "vertices": [0, 1, 2], "type": "triangle"},
      {"id": 1, "vertices": [0, 2, 3, 4, 5], "type": "pentagon"}
    ],
    "edges": [
      {"id": 0, "vertices": [0, 1]}
    ]
  },
  "metadata": {
    "vertexCount": 30,
    "faceCount": 32,
    "edgeCount": 60,
    "genus": 0
  }
}
```

#### JSON Schema:
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["planetId", "mesh", "metadata"],
  "properties": {
    "planetId": {
      "type": "string",
      "format": "uuid"
    },
    "mesh": {
      "$ref": "#/definitions/Mesh"
    },
    "metadata": {
      "$ref": "#/definitions/MeshMetadata"
    }
  },
  "definitions": {
    "Mesh": {
      "type": "object",
      "required": ["vertices", "faces", "edges"],
      "properties": {
        "vertices": {
          "type": "array",
          "items": {"$ref": "#/definitions/Vertex"}
        },
        "faces": {
          "type": "array",
          "items": {"$ref": "#/definitions/Face"}
        },
        "edges": {
          "type": "array",
          "items": {"$ref": "#/definitions/Edge"}
        }
      }
    },
    "Vertex": {
      "type": "object",
      "required": ["id", "x", "y", "z"],
      "properties": {
        "id": {"type": "integer"},
        "x": {"type": "number"},
        "y": {"type": "number"},
        "z": {"type": "number"}
      }
    },
    "Face": {
      "type": "object",
      "required": ["id", "vertices", "type"],
      "properties": {
        "id": {"type": "integer"},
        "vertices": {
          "type": "array",
          "items": {"type": "integer"}
        },
        "type": {
          "type": "string",
          "enum": ["triangle", "pentagon"]
        }
      }
    }
  }
}
```

### User Story 1.2: Subdivide Mesh for Increased Detail

**As a user, I want to subdivide the mesh to increase detail so that I can create more realistic planet surfaces**

#### User Workflow:
1. User has generated base icosidodecahedron
2. User adjusts subdivision level slider (0-5)
3. User clicks "Apply Subdivision"
4. System shows loading indicator
5. System updates mesh with subdivided geometry
6. User sees updated mesh with higher polygon count
7. UI displays vertex/face count statistics

#### Acceptance Criteria:
- [ ] Subdivision levels 0-5 are supported
- [ ] Each subdivision level approximately quadruples face count
- [ ] Subdivision maintains spherical shape
- [ ] Performance: Subdivision completes in <3 seconds for level 5
- [ ] UI shows real-time preview of polygon count before applying
- [ ] User can undo/redo subdivision operations
- [ ] Mesh remains manifold (no holes or non-connected vertices)

#### API Design:

**Request:** `PUT /api/planet/{planetId}/subdivide`
```json
{
  "level": 2,
  "algorithm": "loop",
  "smoothing": true
}
```

**Response:** `200 OK`
```json
{
  "planetId": "uuid-here",
  "mesh": {
    "vertices": [...],
    "faces": [...],
    "edges": [...]
  },
  "metadata": {
    "subdivisionLevel": 2,
    "vertexCount": 482,
    "faceCount": 960,
    "edgeCount": 1440
  },
  "performance": {
    "computationTimeMs": 245,
    "meshSizeBytes": 58240
  }
}
```

### User Story 1.3: 3D Visualization and Camera Controls

**As a user, I want to see the planet rendered in 3D with intuitive camera controls so that I can inspect my creation from all angles**

#### User Workflow:
1. User sees 3D viewport on page load
2. User can click and drag to rotate camera around planet
3. User can scroll/pinch to zoom in/out
4. User can right-click drag to pan camera
5. User can press 'R' to reset camera to default view
6. User can toggle between wireframe/solid/textured view modes
7. User can toggle UI overlay showing mesh statistics

#### Acceptance Criteria:
- [ ] Camera orbits around planet center smoothly (60 FPS)
- [ ] Zoom has min/max limits to prevent clipping
- [ ] Camera controls feel responsive (<16ms input lag)
- [ ] Wireframe mode clearly shows mesh topology
- [ ] Solid mode uses basic lighting (ambient + directional)
- [ ] Stats overlay shows: FPS, vertex count, face count
- [ ] Mobile touch controls work intuitively
- [ ] Camera position persists during mesh updates

#### Frontend Components Structure:

```typescript
// types/Planet.ts
export interface Planet {
  id: string;
  mesh: Mesh;
  metadata: MeshMetadata;
}

export interface Mesh {
  vertices: Vertex[];
  faces: Face[];
  edges: Edge[];
}

// components/PlanetViewer.tsx
interface PlanetViewerProps {
  planet: Planet;
  viewMode: 'wireframe' | 'solid' | 'textured';
  showStats: boolean;
  onCameraChange?: (position: Vector3, target: Vector3) => void;
}

// hooks/usePlanetGeneration.ts
export const usePlanetGeneration = () => {
  const generatePlanet = async (options: GenerateOptions) => {
    // API call to generate planet
  };
  
  const subdividePlanet = async (planetId: string, level: number) => {
    // API call to subdivide
  };
  
  return { generatePlanet, subdividePlanet };
};
```

## Implementation Roadmap

### Sprint 1 (Week 1-2): Foundation
1. **Backend**: Implement icosidodecahedron generation algorithm
2. **Backend**: Create API endpoints with JSON schema validation
3. **Frontend**: Set up Three.js scene with basic lighting
4. **Frontend**: Implement orbit camera controls

### Sprint 2 (Week 3-4): Mesh Operations
1. **Backend**: Implement Loop subdivision algorithm
2. **Backend**: Add mesh validation and error handling
3. **Frontend**: Create UI controls for subdivision
4. **Frontend**: Add loading states and progress indicators

### Sprint 3 (Week 5-6): Polish & Performance
1. **Backend**: Optimize mesh data structures
2. **Backend**: Add caching for subdivided meshes
3. **Frontend**: Implement view modes (wireframe/solid)
4. **Frontend**: Add statistics overlay
5. **Both**: Performance testing and optimization

## Technical Considerations

### Backend (Golang):
- Use a geometry library like `go-geo` or implement custom
- Consider using goroutines for parallel vertex processing
- Implement efficient mesh data structures (half-edge or winged-edge)

### Frontend (Next.js + Three.js):
- Use React Three Fiber for declarative Three.js
- Implement Web Workers for heavy computations
- Use IndexedDB for client-side mesh caching
- Consider LOD system for large meshes

### Data Transfer Optimization:
- Use binary formats (MessagePack/Protobuf) for large meshes
- Implement mesh compression (Draco)
- Add WebSocket support for real-time updates
- Stream mesh data for progressive loading
