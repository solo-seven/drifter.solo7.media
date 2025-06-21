* * *

## 🧱 Project Overview

**Goal**: Build a server-client simulation platform where a Golang backend generates and runs a procedural simulation, and a Next.js frontend serves as a real-time viewer/observer for the simulation world.

* * *

## 🗂️ High-Level Architecture

### 1. **Environment Definition**

- YAML or JSON-based schema with:
    - **Map data**: terrain layout, sectors, spatial boundaries.
    - **Object models**: 3D models or procedural rulesets (e.g., trees, buildings).
    - **Agents**: physical models + behavior rules (FSMs or behavior trees).

### 2. **Golang Backend**

Responsible for:

- Parsing environment definitions
- Procedural generation of world state
- Simulating agent behavior and environment changes
- WebSocket API for real-time updates
- REST/GraphQL API for CRUD operations and world inspection

### 3. **Next.js Frontend**

Responsible for:

- Loading and visualizing world state (likely using WebGL/Three.js)
- Showing agent behaviors and state changes in real time
- User controls for camera, playback, and filters
- Debug panels (FPS, state overlay, agent logs)

* * *

## 🧩 Core Components

### 📦 Environment Package (Go)

- `MapLoader`: loads spatial data into grid or mesh
- `ModelRegistry`: manages assets for objects and agents
- `Spawner`: procedural logic for placing models and agents
- `Serializer`: exports world state to a WebSocket/REST-compatible format

### 🧠 Simulation Engine (Go)

- `AgentEngine`: schedules and runs agent logic each tick
- `WorldState`: in-memory state container of current simulation
- `Physics/RulesEngine`: handles time, motion, collision, causality
- `EventBus`: decoupled message bus for triggering updates/events

### 🌐 API Layer (Go)

- REST for:
    - Starting/stopping/resetting simulations
    - Uploading definitions
    - Inspecting entity state
- WebSocket for:
    - Streaming world state and simulation events
    - User command input (teleport camera, inspect agent, etc.)

### 🖥️ Viewer Client (Next.js)

- Three.js or similar for rendering
- `useWorldState()` React hook using WebSocket
- Camera/navigation controls
- UI panel for logs, stats, object/agent inspection

* * *

## 🚀 Getting Started Steps

### Phase 1: Core Definition & Bootstrapping

1. Define a schema for the environment definition (YAML/JSON)
    - Use JSON Schema to enforce structure if desired
2. Scaffold Golang project with `go mod init` and basic layout
    - `/cmd/server`
    - `/internal/{env,sim,api}`
3. Scaffold Next.js viewer project in `/apps/viewer`
4. Set up shared protobuf or JSON formats for agent/object/map data
5. Implement map loader that can parse and print debug info

### Phase 2: Procedural Generation Prototype

6. Write basic model/agent registry (load from disk)
7. Build basic spawner logic for placing agents and objects
8. Implement server tick loop with agent movement simulation
9. Build WebSocket bridge that emits world state

### Phase 3: Viewer Integration

10. Build minimal viewer that connects to server and renders:
    - Basic floor plane
    - Cube agents
11. Add controls for camera movement and agent inspection
12. Show debug info (agent positions, simulation time)

* * *

## 🛠️ Tools & Libraries

### Golang

- `gorilla/websocket` or `nhooyr/websocket`
- `go-json`, `gopkg.in/yaml.v3`
- `go-fsm` or custom FSM implementation for behaviors

### Frontend

- `Three.js` or `react-three-fiber`
- `Tailwind` or `shadcn/ui` for UI
- `zustand` or `redux` for global state

* * *

## ✅ Next Steps (Today/This Week)

1. **Define your environment definition schema**
2. **Scaffold the Golang backend structure**
3. **Scaffold the Next.js app with Three.js or R3F**
4. **Decide on state format** (protobuf vs JSON)
5. **Start building the map loader in Go**
6. **Set up WebSocket connection from frontend to backend**
