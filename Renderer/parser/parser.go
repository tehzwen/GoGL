package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Mesh struct {
	Vertices []float32
	Normals  []float32
	UVs      []float32
	Sparse   Sparse
}

type Sparse struct {
	Vertices []float32
	Normals  []float32
	UVs      []float32
}

func addFace(a *int, b *int, c *int, d *int, ua *int, ub *int, uc *int, ud *int, na *int, nb *int, nc *int, nd *int, mesh *Mesh) {
	vLen := len(mesh.Sparse.Vertices)

	ia := parseVertexIndex((*a), vLen)
	ib := parseVertexIndex((*b), vLen)
	ic := parseVertexIndex((*c), vLen)

	if d == nil {
		addVertex(ia, ib, ic, mesh)

	} else {

		id := parseVertexIndex((*d), vLen)
		addVertex(ia, ib, id, mesh)
		addVertex(ib, ic, id, mesh)
	}

	if ua != nil {

		uvLen := len(mesh.Sparse.UVs)

		ia := parseUVIndex((*ua), uvLen)
		ib := parseUVIndex((*ub), uvLen)
		ic := parseUVIndex((*uc), uvLen)

		if &d == nil {
			//add uv call
			addUV(ia, ib, ic, mesh)
		} else {
			id := parseUVIndex((*ud), uvLen)

			addUV(ia, ib, id, mesh)
			addUV(ib, ic, id, mesh)
		}
	}

	if na != nil {

		nLen := len(mesh.Normals)
		ia := parseVertexIndex((*na), nLen)

		if na == nb {
			ib = ia
		} else {
			ib = parseVertexIndex((*nb), nLen)
		}

		if na == nc {
			ic = ia
		} else {
			ic = parseVertexIndex((*nc), nLen)
		}

		if d == nil {
			addNormal(ia, ib, ic, mesh)
		} else {
			id := parseVertexIndex((*nd), nLen)

			addNormal(ia, ib, id, mesh)
			addNormal(ib, ic, id, mesh)
		}
	}

}

func addVertex(a int, b int, c int, mesh *Mesh) {
	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[a+0])
	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[a+1])
	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[a+2])

	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[b+0])
	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[b+1])
	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[b+2])

	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[c+0])
	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[c+1])
	mesh.Vertices = append(mesh.Vertices, (*mesh).Sparse.Vertices[c+2])
}

func addNormal(a int, b int, c int, mesh *Mesh) {
	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[a+0])
	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[a+1])
	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[a+2])

	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[b+0])
	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[b+1])
	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[b+2])

	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[c+0])
	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[c+1])
	mesh.Normals = append(mesh.Normals, (*mesh).Sparse.Normals[c+2])
}

func addUV(a int, b int, c int, mesh *Mesh) {
	mesh.UVs = append(mesh.UVs, (*mesh).Sparse.UVs[a+0])
	mesh.UVs = append(mesh.UVs, (*mesh).Sparse.UVs[a+1])

	mesh.UVs = append(mesh.UVs, (*mesh).Sparse.UVs[b+0])
	mesh.UVs = append(mesh.UVs, (*mesh).Sparse.UVs[b+1])

	mesh.UVs = append(mesh.UVs, (*mesh).Sparse.UVs[c+0])
	mesh.UVs = append(mesh.UVs, (*mesh).Sparse.UVs[c+1])
}

func parseUVIndex(value int, len int) int {

	if value >= 0 {
		return int((value - 1) * 2)
	}
	return (int(value) + len/2) * 2
}

func parseVertexIndex(value int, len int) int {

	if value >= 0 {
		return int((value - 1) * 3)
	}
	return (int(value) + len/3) * 3
}

func Parse(filePath string) Mesh {
	objFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer objFile.Close()
	//count := 0
	newMesh := Mesh{}

	scanner := bufio.NewScanner(objFile)
	for scanner.Scan() {
		line := scanner.Text()

		if line[0] == 'v' {
			if line[1] == ' ' {
				result := strings.Fields(line)
				for i := 1; i < 4; i++ {
					if s, err := strconv.ParseFloat(result[i], 32); err == nil {
						newMesh.Sparse.Vertices = append(newMesh.Sparse.Vertices, float32(s))
					}
				}
			} else if line[1] == 'n' {
				result := strings.Fields(line)
				for i := 1; i < 4; i++ {
					if s, err := strconv.ParseFloat(result[i], 32); err == nil {
						newMesh.Sparse.Normals = append(newMesh.Sparse.Normals, float32(s))
					}
				}
			} else if line[1] == 't' {
				result := strings.Fields(line)
				for i := 1; i < 3; i++ {
					if s, err := strconv.ParseFloat(result[i], 32); err == nil {
						newMesh.Sparse.UVs = append(newMesh.Sparse.UVs, float32(s))
					}
				}
			}
		} else if line[0] == 'f' {
			//means we have v/u/n
			if result, err := regexp.MatchString(`^f\s+(-?\d+)\/(-?\d+)\/(-?\d+)\s+(-?\d+)\/(-?\d+)\/(-?\d+)\s+(-?\d+)\/(-?\d+)\/(-?\d+)(?:\s+(-?\d+)\/(-?\d+)\/(-?\d+))?`, line); err == nil && result {
				//split on whitespace first
				whiteSpaceSplit := strings.Split(line, " ")
				var values []int

				//now iterate through the 3 values and pull the values out
				for i := 1; i < len(whiteSpaceSplit); i++ {
					slashSplit := strings.Split(whiteSpaceSplit[i], "/")
					for j := 0; j < len(slashSplit); j++ {
						if s, err := strconv.ParseInt(slashSplit[j], 10, 32); err == nil {
							values = append(values, int(s))
						}
					}
				}
				addFace(
					&values[0], &values[3], &values[6], &values[9],
					&values[1], &values[4], &values[7], &values[10],
					&values[2], &values[5], &values[8], &values[11],
					&newMesh)

			} else if result, err := regexp.MatchString(`^f\s+(-?\d+)\/(-?\d+)\s+(-?\d+)\/(-?\d+)\s+(-?\d+)\/(-?\d+)(?:\s+(-?\d+)\/(-?\d+))?`, line); err == nil && result {
				//split on whitespace first
				whiteSpaceSplit := strings.Split(line, " ")
				var values []int

				//now iterate through the 3 values and pull the values out
				for i := 1; i < len(whiteSpaceSplit); i++ {
					slashSplit := strings.Split(whiteSpaceSplit[i], "/")
					for j := 0; j < len(slashSplit); j++ {
						if s, err := strconv.ParseInt(slashSplit[j], 10, 32); err == nil {
							values = append(values, int(s))
						}
					}
				}
				addFace(
					&values[0], &values[2], &values[4], &values[6],
					&values[1], &values[3], &values[5], &values[7],
					nil, nil, nil, nil,
					&newMesh)
			} else if result, err := regexp.MatchString(`^f\s+(-?\d+)\/\/(-?\d+)\s+(-?\d+)\/\/(-?\d+)\s+(-?\d+)\/\/(-?\d+)(?:\s+(-?\d+)\/\/(-?\d+))?`, line); err == nil && result {
				//split on whitespace first
				whiteSpaceSplit := strings.Split(line, " ")
				var values []int

				//now iterate through the 3 values and pull the values out
				for i := 1; i < len(whiteSpaceSplit); i++ {
					slashSplit := strings.Split(whiteSpaceSplit[i], "/")
					for j := 0; j < len(slashSplit); j++ {
						if s, err := strconv.ParseInt(slashSplit[j], 10, 32); err == nil {
							values = append(values, int(s))
						}
					}
				}
				addFace(
					&values[0], &values[2], &values[4], nil,
					nil, nil, nil, nil,
					&values[1], &values[3], &values[5], nil,
					&newMesh)
			}
		}
	}

	fmt.Println("Vertices:", len(newMesh.Normals), "normals:", len(newMesh.Normals), "uvs:", len(newMesh.UVs))
	return newMesh
}
