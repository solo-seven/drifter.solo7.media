package geometry

import (
	"math"
	"strconv"

	"github.com/solo-seven/drifter.solo7.media/internal/planet"
)

// GenerateIcosidodecahedron creates a new icosidodecahedron mesh.
func GenerateIcosidodecahedron(radius float64) *planet.Mesh {
	phi := (1.0 + math.Sqrt(5.0)) / 2.0

	// 30 vertices of the icosidodecahedron
	verts := []planet.Vertex{
		{X: 0, Y: 0, Z: 2}, {X: 0, Y: 0, Z: -2},
		{X: 0, Y: 2, Z: 0}, {X: 0, Y: -2, Z: 0},
		{X: 2, Y: 0, Z: 0}, {X: -2, Y: 0, Z: 0},
		{X: 1, Y: phi, Z: 1 / phi}, {X: 1, Y: phi, Z: -1 / phi},
		{X: 1, Y: -phi, Z: 1 / phi}, {X: 1, Y: -phi, Z: -1 / phi},
		{X: -1, Y: phi, Z: 1 / phi}, {X: -1, Y: phi, Z: -1 / phi},
		{X: -1, Y: -phi, Z: 1 / phi}, {X: -1, Y: -phi, Z: -1 / phi},
		{X: 1 / phi, Y: 1, Z: phi}, {X: 1 / phi, Y: 1, Z: -phi},
		{X: 1 / phi, Y: -1, Z: phi}, {X: 1 / phi, Y: -1, Z: -phi},
		{X: -1 / phi, Y: 1, Z: phi}, {X: -1 / phi, Y: 1, Z: -phi},
		{X: -1 / phi, Y: -1, Z: phi}, {X: -1 / phi, Y: -1, Z: -phi},
		{X: phi, Y: 1 / phi, Z: 1}, {X: phi, Y: 1 / phi, Z: -1},
		{X: phi, Y: -1 / phi, Z: 1}, {X: phi, Y: -1 / phi, Z: -1},
		{X: -phi, Y: 1 / phi, Z: 1}, {X: -phi, Y: 1 / phi, Z: -1},
		{X: -phi, Y: -1 / phi, Z: 1}, {X: -phi, Y: -1 / phi, Z: -1},
	}

	for i := range verts {
		v := &verts[i]
		length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
		if length > 0 {
			v.X = (v.X / length) * radius
			v.Y = (v.Y / length) * radius
			v.Z = (v.Z / length) * radius
		}
		v.Id = i
	}

	faces := []planet.Face{
		// 20 triangular faces
		{Id: 0, Vertices: []int{0, 14, 20}, Type: "triangle"},
		{Id: 1, Vertices: []int{0, 20, 2}, Type: "triangle"},
		{Id: 2, Vertices: []int{1, 15, 21}, Type: "triangle"},
		{Id: 3, Vertices: []int{1, 21, 3}, Type: "triangle"},
		{Id: 4, Vertices: []int{2, 18, 10}, Type: "triangle"},
		{Id: 5, Vertices: []int{2, 10, 6}, Type: "triangle"},
		{Id: 6, Vertices: []int{3, 12, 28}, Type: "triangle"},
		{Id: 7, Vertices: []int{3, 28, 8}, Type: "triangle"},
		{Id: 8, Vertices: []int{4, 22, 24}, Type: "triangle"},
		{Id: 9, Vertices: []int{4, 24, 8}, Type: "triangle"},
		{Id: 10, Vertices: []int{5, 26, 27}, Type: "triangle"},
		{Id: 11, Vertices: []int{5, 27, 10}, Type: "triangle"},
		{Id: 12, Vertices: []int{6, 7, 16}, Type: "triangle"},
		{Id: 13, Vertices: []int{15, 17, 9}, Type: "triangle"},
		{Id: 14, Vertices: []int{11, 19, 13}, Type: "triangle"},
		{Id: 15, Vertices: []int{23, 25, 9}, Type: "triangle"},
		{Id: 16, Vertices: []int{29, 13, 7}, Type: "triangle"},
		{Id: 17, Vertices: []int{1, 19, 11}, Type: "triangle"},
		{Id: 18, Vertices: []int{5, 11, 27}, Type: "triangle"},
		{Id: 19, Vertices: []int{29, 7, 12}, Type: "triangle"},

		// 12 pentagonal faces
		{Id: 20, Vertices: []int{0, 2, 6, 16, 14}, Type: "pentagon"},
		{Id: 21, Vertices: []int{1, 3, 8, 25, 15}, Type: "pentagon"},
		{Id: 22, Vertices: []int{4, 8, 24, 22, 9}, Type: "pentagon"},
		{Id: 23, Vertices: []int{5, 10, 27, 26, 11}, Type: "pentagon"},
		{Id: 24, Vertices: []int{14, 16, 7, 18, 20}, Type: "pentagon"},
		{Id: 25, Vertices: []int{15, 25, 23, 17, 21}, Type: "pentagon"},
		{Id: 26, Vertices: []int{18, 6, 10, 11, 19}, Type: "pentagon"},
		{Id: 27, Vertices: []int{20, 2, 12, 28, 29}, Type: "pentagon"},
		{Id: 28, Vertices: []int{21, 3, 13, 29, 28}, Type: "pentagon"},
		{Id: 29, Vertices: []int{22, 24, 4, 9, 23}, Type: "pentagon"},
		{Id: 30, Vertices: []int{26, 5, 13, 19, 27}, Type: "pentagon"},
		{Id: 31, Vertices: []int{18, 7, 12, 2, 29}, Type: "pentagon"},
	}

	edges := deriveEdgesFromFaces(faces)

	return &planet.Mesh{
		Vertices: verts,
		Faces:    faces,
		Edges:    edges,
	}
}

// deriveEdgesFromFaces generates a list of unique edges from the mesh's faces.
func deriveEdgesFromFaces(faces []planet.Face) []planet.Edge {
	edgeMap := make(map[string]bool)
	var edges []planet.Edge

	for _, face := range faces {
		for i := 0; i < len(face.Vertices); i++ {
			v1 := face.Vertices[i]
			v2 := face.Vertices[(i+1)%len(face.Vertices)] // Wrap around

			if v1 > v2 {
				v1, v2 = v2, v1
			}

			edgeKey := strconv.Itoa(v1) + "-" + strconv.Itoa(v2)
			if !edgeMap[edgeKey] {
				edgeMap[edgeKey] = true
				edges = append(edges, planet.Edge{Id: len(edges), Vertices: []int{v1, v2}})
			}
		}
	}

	return edges
}
