---
trigger: always_on
---

# Drifter - Planet Generator Architecture Guide

## Overview

Drifter is a web application that creates procedurally generated planets using geometric subdivision and geophysical simulations. The system follows a client-server architecture with a clear separation of concerns.

## Tech Stack

- **Frontend**: Next.js (React) + TypeScript + Three.js
- **Backend**: Golang
- **Communication**: REST API with JSON/JSON Schema
- **3D Rendering**: Three.js with React Three Fiber

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                        │
├─────────────────────────────────────────────────────────────────┤
│  UI Layer           │  3D Visualization    │  State Management  │
│  - Parameter Forms  │  - Three.js Scene   │  - React Hooks     │
│  - Status Display   │  - Camera Controls  │  - Local Storage   │
│  - Export/Import    │  - Mesh Rendering   │  - API Integration │
└─────────────────────────────────────────────────────────────────┘
                                   │
                                   │ JSON via REST API
                                   ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Backend (Golang)                          │
├─────────────────────────────────────────────────────────────────┤
│  API Layer          │  Domain Logic       │  Data Layer        │
│  - HTTP Handlers    │  - Mesh Generation  │  - Mesh Storage    │
│  - JSON Validation  │  - Subdivision      │  - Caching         │
│  - WebSocket        │  - Terrain Gen      │  - Serialization   │
│                     │  - Physics Sim      │                    │
└─────────────────────────────────────────────────────────────────┘
```

## Core Components

### Frontend Components

1. **PlanetGenerator** (Main Component)
   - Manages overall application state
   - Coordinates between parameter controls and 3D viewport

2. **Parameter Controls**
   - `SliderInput`, `SelectInput`, `ToggleInput`: Reusable form controls
   - `CollapsibleSection`: Organizes parameters into logical groups
   - Parameter sections for each generation aspect (Basic, Subdivision, Terrain, etc.)

3. **3D Viewport** (To be implemented)
   - Three.js scene management
   - Camera controls (orbit, zoom, pan)
   - Mesh rendering with different view modes

4. **Status Display**
   - Real-time generation status
   - Mesh statistics (vertex/face count)
   - Performance metrics

### Backend Domains (Domain-Driven Design)

1. **Geometry Domain**
   - Icosidodecahedron generation
   - Mesh subdivision algorithms
   - Vertex/Face/Edge management

2. **Terrain Domain**
   - Noise-based height generation
   - Biome classification
   - Surface properties

3. **Geophysics Domain**
   - Plate tectonics simulation
   - Erosion algorithms
   - Climate modeling

4. **Visualization Domain**
   - Mesh optimization for rendering
   - LOD (Level of Detail) generation
   - Export formats

## Data Flow

1. **Parameter Update Flow**
   ```
   User Input → React State → Validation → API Request → Backend Processing → Response → 3D Update
   ```

2. **Generation Pipeline**
   ```
   Parameters → Base Mesh → Subdivision → Terrain → Tectonics → Erosion → Final Mesh
   ```

3. **API Communication**
   - Request: Planet parameters (JSON)
   - Response: Mesh data + metadata (JSON)
   - Large meshes: Consider binary formats or streaming

## API Communication

### JSON Schema Contract
The project uses JSON Schema for API contract definition:

1. **Schema Definition**: Define in `schemas/environment.schema.json`
2. **Code Generation**: Run `make generate-go-schema` to create Go structs
3. **Type Safety**: Both frontend TypeScript and backend Go share the same contract
4. **Validation**: Automatic request/response validation

### API Endpoints

- `POST /api/planet/generate` - Create new planet
- `PUT /api/planet/{id}/subdivide` - Apply subdivision
- `GET /api/planet/{id}` - Retrieve planet data
- `WS /ws/planet/{id}` - Real-time generation updates (if implemented)

### Data Format
- Request: Planet parameters (JSON)
- Response: Mesh data + metadata (JSON)
- Large meshes: Consider binary formats or streaming

## Key Design Decisions

1. **Separation of Concerns**: UI logic, 3D rendering, and planet generation are clearly separated
2. **Incremental Complexity**: Start with basic mesh generation, add features progressively
3. **Performance Considerations**: 
   - Client-side caching of generated meshes
   - Progressive mesh streaming for large planets
   - Web Workers for heavy computations
4. **Extensibility**: Parameter-driven design allows easy addition of new generation features

## Development Workflow

### 1. Frontend Development
- Component development in isolation
- Mock API responses for testing
- Storybook for UI component gallery (if configured)
- Hot reload with Next.js dev server

### 2. Backend Development
- Domain-driven modules in `internal/` directory
- Unit tests with minimum 80% coverage requirement
- Hot reload with Air for rapid development
- Benchmarks for performance-critical code

### 3. Schema-Driven Development
- Define data models in `schemas/` directory (JSON Schema)
- Generate Go structs: `make generate-go-schema`
- Ensures frontend/backend contract consistency
- Type safety across the stack

### 4. Testing Strategy
- **Backend**: Minimum 80% test coverage
- **Frontend**: Jest tests with `npm test` with minimum 80% test coverage
- **Integration**: Test API contracts with JSON Schema validation
- **E2E**: Full application testing (to be implemented)

### 5. CI/CD Workflow
```bash
# Local CI simulation
make lint          # Lint both frontend and backend
make test-backend  # Run backend tests with coverage
make test-frontend # Run frontend tests
make build-backend # Build production binary
make build-frontend # Build production frontend
```

## Getting Started

### Quick Start (Local Development)

1. **Install dependencies**
   ```bash
   make deps  # Installs both Go modules and npm packages
   ```

2. **Start everything**
   ```bash
   make all   # Installs deps and starts both services
   # OR
   make start # Start services in background
   ```

3. **Access the application**
   - Frontend: `http://localhost:3000`
   - Backend API: `http://localhost:8080`

### Docker Development

```bash
# Build and run with Docker Compose
make docker-up

# Stop services
make docker-down

# Build individual images
make docker-build-frontend
make docker-build-backend
```

## Development Commands

### Backend Development

```bash
# Run with hot-reload (using Air)
make dev-backend

# Run normally
make run-backend

# Build binary
make build-backend

# Run tests (requires 80% coverage)
make test-backend

# Lint code
make lint-backend
```

### Frontend Development

```bash
# Run development server
make dev-frontend

# Build for production
make build-frontend

# Run tests
make test-frontend

# Lint and format code
make lint-frontend
make format-frontend
```

### Utility Commands

```bash
# View logs (when running with 'make start')
make logs

# Stop all services
make stop

# Clean all build artifacts and Docker images
make clean

# Generate Go structs from JSON schemas
make generate-go-schema

# Full regeneration and build
make regen
```

## Project Structure

```
drifter/
├── backend/              # Golang backend
│   ├── main.go          # Entry point
│   ├── internal/        # Internal packages
│   │   └── world/       # Planet generation domain
│   └── bin/             # Built binaries
├── frontend/            # Next.js frontend
│   ├── app/            # Next.js app directory
│   ├── components/     # React components
│   └── hooks/          # Custom React hooks
├── schemas/            # JSON schemas
│   └── environment.schema.json
├── logs/               # Application logs
├── docker-compose.yml  # Docker orchestration
└── Makefile           # Build automation
```

## Future Enhancements

- **Real-time Collaboration**: Multiple users editing the same planet
- **GPU Acceleration**: Compute shaders for erosion simulation
- **Procedural Textures**: Material generation based on terrain
- **Export Formats**: glTF, OBJ, heightmaps
- **Planet Templates**: Pre-configured Earth-like, Mars-like presets

## Performance Targets

- Base mesh generation: < 100ms
- Subdivision (level 5): < 3 seconds
- Full pipeline with erosion: < 30 seconds
- 60 FPS rendering for meshes up to 1M vertices

## Troubleshooting

### Viewing Logs
When running services with `make start`:
```bash
make logs  # View backend logs
# Logs are stored in logs/backend.log and logs/frontend.log
```

### Common Issues

1. **Port conflicts**: Ensure ports 3000 (frontend) and 8080 (backend) are available
2. **Coverage failures**: Backend tests require 80% coverage (excluding internal/world)
3. **Process cleanup**: Use `make stop` to properly terminate all services
4. **Clean start**: Run `make clean` to remove all artifacts and start fresh

### Development Tips

- Use `make dev-backend` for hot-reload during backend development
- Backend binary is built to `backend/bin/drifter`
- PID files track running services for proper cleanup
- Docker development isolates dependencies if local setup has issues