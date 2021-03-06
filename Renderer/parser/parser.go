package parser

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MTLMaterial struct {
	MTLLib string
	Name   string
	Start  int
	End    int
}

type OBJObject struct {
	Geometry  Mesh
	Materials []MTLMaterial
	Name      string
	Smooth    bool
}

func SerializeOBJ(vals []OBJObject) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)

	err := e.Encode(vals)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func DeserializeOBJ(value string) []OBJObject {
	var objects []OBJObject
	by, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Println(value)
		panic(err)
	}
	b := bytes.Buffer{}

	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&objects)
	if err != nil {
		panic(err)
	}

	return objects
}

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

		if d == nil {
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

func addLineGeometry(vertices []float32, uvs []float32, mesh *Mesh) {
	vLen := len(mesh.Sparse.Vertices)
	uvLen := len(mesh.Sparse.UVs)

	fmt.Println("LINE GEOMETRY! ", vLen, " ", uvLen)

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
	//fmt.Println(a, b, c, len((*mesh).Sparse.Normals))
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

func Parse(filePath string) []OBJObject {
	objFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer objFile.Close()
	objectCount := 0
	materialCount := 0

	objects := []OBJObject{}
	objects = append(objects, OBJObject{})
	mtlLib := ""

	scanner := bufio.NewScanner(objFile)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		if line[0] == 'v' {
			if line[1] == ' ' {
				result := strings.Fields(line)
				for i := 1; i < 4; i++ {
					if s, err := strconv.ParseFloat(result[i], 32); err == nil {
						//fmt.Println(len(objects), objectCount)
						objects[objectCount].Geometry.Sparse.Vertices = append(objects[objectCount].Geometry.Sparse.Vertices, float32(s))
					}
				}
			} else if line[1] == 'n' {
				result := strings.Fields(line)
				for i := 1; i < 4; i++ {
					if s, err := strconv.ParseFloat(result[i], 32); err == nil {
						objects[objectCount].Geometry.Sparse.Normals = append(objects[objectCount].Geometry.Sparse.Normals, float32(s))
					}
				}
			} else if line[1] == 't' {
				result := strings.Fields(line)
				for i := 1; i < 3; i++ {
					if s, err := strconv.ParseFloat(result[i], 32); err == nil {
						objects[objectCount].Geometry.Sparse.UVs = append(objects[objectCount].Geometry.Sparse.UVs, float32(s))
					}
				}
			}
		} else if line[0] == 'f' {
			//means we have f vertex/uv/normal vertex/uv/normal vertex/uv/normal
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

				//check if the last three values are nil
				if len(values) == 9 {
					addFace(
						&values[0], &values[3], &values[6], nil,
						&values[1], &values[4], &values[7], nil,
						&values[2], &values[5], &values[8], nil,
						&objects[objectCount].Geometry)

				} else {
					addFace(
						&values[0], &values[3], &values[6], &values[9],
						&values[1], &values[4], &values[7], &values[10],
						&values[2], &values[5], &values[8], &values[11],
						&objects[objectCount].Geometry)
				}
				//means we have f vertex/uv vertex/uv vertex/uv
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
					&objects[objectCount].Geometry)

				//means we have f vertex//normal vertex//normal vertex//normal
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
					&objects[objectCount].Geometry)

				//means we have f vertex vertex vertex
			} else if result, err := regexp.MatchString(`^f\s+(-?\d+)\s+(-?\d+)\s+(-?\d+)(?:\s+(-?\d+))?`, line); err == nil && result {

				whiteSpaceSplit := strings.Split(line, " ")
				var values []int

				for i := 1; i < len(whiteSpaceSplit); i++ {
					if s, err := strconv.ParseInt(whiteSpaceSplit[i], 10, 32); err == nil {
						values = append(values, int(s))
					}
				}
				addFace(
					&values[0], &values[1], &values[2], &values[3],
					nil, nil, nil, nil,
					&values[0], &values[1], &values[2], &values[3],
					&objects[objectCount].Geometry)

				//check for object
			}
		} else if result, err := regexp.MatchString(`^[og]\s*(.+)?`, line); err == nil && result {
			whiteSpaceSplit := strings.Split(line, " ")

			//check if this is the initial object or not
			if len(objects) == 0 {
				tempObject := OBJObject{}
				tempObject.Geometry = Mesh{}
				tempObject.Materials = []MTLMaterial{}
				objects = append(objects, tempObject)
				objectCount = len(objects) - 1
				objects[objectCount].Name = whiteSpaceSplit[1]
			} else {
				if len(objects[objectCount].Materials) > 0 {
					objects[objectCount].Materials[materialCount-1].End = len(objects[objectCount].Geometry.Vertices) / 3
				}

				materialCount = 0
				//create a new object and increment counters
				tempObject := OBJObject{}
				tempObject.Name = whiteSpaceSplit[1]
				tempObject.Geometry = Mesh{
					Sparse: objects[objectCount].Geometry.Sparse,
				}
				objects = append(objects, tempObject)
				objectCount = len(objects) - 1
			}

		} else if result, err := regexp.MatchString(`^mtllib `, line); err == nil && result {
			whiteSpaceSplit := strings.Split(line, " ")
			//mtllib shows up first so we will make an object here
			mtlLib = whiteSpaceSplit[1]

		} else if result, err := regexp.MatchString(`^usemtl `, line); err == nil && result {
			whiteSpaceSplit := strings.Split(line, " ")

			if objectCount == 0 && materialCount == 0 {
				//create the first material
				tempMaterial := MTLMaterial{
					MTLLib: mtlLib,
					Name:   whiteSpaceSplit[1],
					Start:  0,
				}
				objects[objectCount].Materials = append(objects[objectCount].Materials, tempMaterial)
				materialCount++

			} else {

				tempMaterial := MTLMaterial{
					MTLLib: mtlLib,
					Name:   whiteSpaceSplit[1],
					Start:  len(objects[objectCount].Geometry.Vertices) / 3,
				}

				objects[objectCount].Materials = append(objects[objectCount].Materials, tempMaterial)
				if materialCount > 0 {
					objects[objectCount].Materials[materialCount-1].End = len(objects[objectCount].Geometry.Vertices) / 3
				}
				materialCount++
			}
		} else if line[0] == 'l' {
			whiteSpaceSplit := strings.Split(line, " ")
			UVs := []float32{}
			Vertices := []float32{}

			if !strings.Contains(line, "/") {
				//do some shit
				//fmt.Println("slash not present?")
				//fmt.Println(whiteSpaceSplit[1:len(whiteSpaceSplit)])

				for i := 1; i < len(whiteSpaceSplit); i++ {
					fmt.Println(whiteSpaceSplit[i])
					val, err := strconv.ParseFloat(whiteSpaceSplit[i], 32)
					if err != nil {
						panic(err)
					}

					Vertices = append(Vertices, float32(val))
				}

				//addLineGeometry(whiteSpaceSplit[1:len(whiteSpaceSplit)], []float32{})
			} else {
				//TODO
				fmt.Println("fuck you, ", whiteSpaceSplit[1], UVs)
			}

			addLineGeometry(Vertices, UVs, &objects[objectCount].Geometry)
		}
	}
	objects[objectCount].Materials[materialCount-1].End = len(objects[objectCount].Geometry.Vertices) / 3

	return objects
}
