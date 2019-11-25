package geometry

import (
	"fmt"
	"io"
	"os"

	"github.com/go-gl/mathgl/mgl32"
)

// Model is a renderable collection of vecs.
type Mesh struct {
	// For the v, vt and vn in the obj file.
	Normals, Vecs []mgl32.Vec3
	Uvs           []mgl32.Vec2

	// For the fun "f" in the obj file.
	VecIndices, NormalIndices, UvIndices []float32
}

// NewModel will read an OBJ model file and create a Model from its contents
func NewModel(file string) Mesh {
	// Open the file for reading and check for errors.
	objFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	// Don't forget to close the file reader.
	defer objFile.Close()

	// Create a model to store stuff.
	model := Mesh{}

	// Read the file and get it's contents.
	for {
		var lineType string

		// Scan the type field.
		_, err := fmt.Fscanf(objFile, "%s", &lineType)

		// Check if it's the end of the file
		// and break out of the loop.
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Check the type.
		switch lineType {
		// VERTICES.
		case "v":
			// Create a vec to assign digits to.
			vec := mgl32.Vec3{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])

			// Add the vector to the model.
			model.Vecs = append(model.Vecs, vec)

		// NORMALS.
		case "vn":
			// Create a vec to assign digits to.
			vec := mgl32.Vec3{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f %f\n", &vec[0], &vec[1], &vec[2])

			// Add the vector to the model.
			model.Normals = append(model.Normals, vec)

		// TEXTURE VERTICES.
		case "vt":
			// Create a Uv pair.
			vec := mgl32.Vec2{}

			// Get the digits from the file.
			fmt.Fscanf(objFile, "%f %f\n", &vec[0], &vec[1])

			// Add the uv to the model.
			model.Uvs = append(model.Uvs, vec)

		// INDICES.
		case "f":
			// Create a vec to assign digits to.
			norm := make([]float32, 3)
			vec := make([]float32, 3)
			uv := make([]float32, 3)

			// Get the digits from the file.
			matches, _ := fmt.Fscanf(objFile, "%f/%f/%f %f/%f/%f %f/%f/%f\n", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])

			if matches == 0 {
				panic("Cannot read your file")
			}

			// Add the numbers to the model.
			model.NormalIndices = append(model.NormalIndices, norm[0])
			model.NormalIndices = append(model.NormalIndices, norm[1])
			model.NormalIndices = append(model.NormalIndices, norm[2])

			model.VecIndices = append(model.VecIndices, vec[0])
			model.VecIndices = append(model.VecIndices, vec[1])
			model.VecIndices = append(model.VecIndices, vec[2])

			model.UvIndices = append(model.UvIndices, uv[0])
			model.UvIndices = append(model.UvIndices, uv[1])
			model.UvIndices = append(model.UvIndices, uv[2])
		}
	}

	// Return the newly created Model.
	return model
}

// GetRenderableVertices returns a slice of float32s
// formatted in X, Y, Z, U, V. That is, XYZ of the
// vertex and the texture position.
func (model Mesh) GetRenderableVertices() []float32 {
	// Create a slice for the outward float32s.
	var out []float32

	// Loop over each vec3 in the indices property.
	for _, position := range model.VecIndices {
		index := int(position) - 1
		vec := model.Vecs[index]
		uv := model.Uvs[int(model.UvIndices[index])-1]

		out = append(out, vec.X(), vec.Y(), vec.Z(), uv.X(), uv.Y())
	}

	// Return the array.
	return out
}
