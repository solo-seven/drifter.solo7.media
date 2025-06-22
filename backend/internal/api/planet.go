package api

import (
	"encoding/json"
	"net/http"

	"github.com/solo-seven/drifter.solo7.media/internal/geometry"
	"github.com/solo-seven/drifter.solo7.media/internal/planet"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// RegisterPlanetRoutes sets up the routes for planet generation.
func RegisterPlanetRoutes(router *mux.Router) {
	router.HandleFunc("/api/planet/generate", GeneratePlanetHandler).Methods("POST")
}

// GeneratePlanetHandler handles the request to generate a new planet mesh.
func GeneratePlanetHandler(w http.ResponseWriter, r *http.Request) {
	// For now, we'll generate a default icosidodecahedron with radius 1.0
	mesh := geometry.GenerateIcosidodecahedron(1.0)

	planetID := uuid.New().String()

	response := planet.PlanetSchemaJson{
		PlanetId: planetID,
		Mesh:     *mesh,
		Metadata: planet.MeshMetadata{
			VertexCount: len(mesh.Vertices),
			FaceCount:   len(mesh.Faces),
			EdgeCount:   len(mesh.Edges),
			Genus:       0, // An icosidodecahedron is topologically a sphere
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
