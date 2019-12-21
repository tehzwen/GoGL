package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ParsedMaterial struct {
	Ns    float32
	Ka    []float32
	Kd    []float32
	Ks    []float32
	D     float32
	MapKD string
}

func ParseMTLFile(filename string, materialName string) ParsedMaterial {
	mtlFile, err := os.Open("../Editor/models/" + filename)
	if err != nil {
		panic(err)
	}

	defer mtlFile.Close()

	mtlDetails := ParsedMaterial{}
	materialFound := false

	scanner := bufio.NewScanner(mtlFile)
	for scanner.Scan() {
		line := scanner.Text()

		whiteSpaceSplit := strings.Split(line, " ")

		if whiteSpaceSplit[0] == "newmtl" && whiteSpaceSplit[1] != materialName {
			materialFound = false
		}

		if whiteSpaceSplit[0] == "newmtl" && whiteSpaceSplit[1] == materialName {
			materialFound = true
		}

		if materialFound {
			if whiteSpaceSplit[0] == "Ns" {
				nS, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
				if err != nil {
					panic(err)
				}
				mtlDetails.Ns = float32(nS)
			} else if whiteSpaceSplit[0] == "Ka" {
				ka1, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
				if err != nil {
					panic(err)
				}
				ka2, err := strconv.ParseFloat(whiteSpaceSplit[2], 64)
				if err != nil {
					panic(err)
				}
				ka3, err := strconv.ParseFloat(whiteSpaceSplit[3], 64)
				if err != nil {
					panic(err)
				}

				mtlDetails.Ka = append(mtlDetails.Ka, float32(ka1))
				mtlDetails.Ka = append(mtlDetails.Ka, float32(ka2))
				mtlDetails.Ka = append(mtlDetails.Ka, float32(ka3))
			} else if whiteSpaceSplit[0] == "Kd" {
				kd1, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
				if err != nil {
					panic(err)
				}
				kd2, err := strconv.ParseFloat(whiteSpaceSplit[2], 64)
				if err != nil {
					panic(err)
				}
				kd3, err := strconv.ParseFloat(whiteSpaceSplit[3], 64)
				if err != nil {
					panic(err)
				}

				mtlDetails.Kd = append(mtlDetails.Kd, float32(kd1))
				mtlDetails.Kd = append(mtlDetails.Kd, float32(kd2))
				mtlDetails.Kd = append(mtlDetails.Kd, float32(kd3))
			} else if whiteSpaceSplit[0] == "Ks" {
				ks1, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
				if err != nil {
					panic(err)
				}
				ks2, err := strconv.ParseFloat(whiteSpaceSplit[2], 64)
				if err != nil {
					panic(err)
				}
				ks3, err := strconv.ParseFloat(whiteSpaceSplit[3], 64)
				if err != nil {
					panic(err)
				}

				mtlDetails.Ks = append(mtlDetails.Ks, float32(ks1))
				mtlDetails.Ks = append(mtlDetails.Ks, float32(ks2))
				mtlDetails.Ks = append(mtlDetails.Ks, float32(ks3))
			} else if whiteSpaceSplit[0] == "d" {
				d, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
				if err != nil {
					panic(err)
				}

				mtlDetails.D = float32(d)
			} else if whiteSpaceSplit[0] == "map_Kd" {
				mtlDetails.MapKD = whiteSpaceSplit[1]
			}
		}

		/*if whiteSpaceSplit[0] == "Ns" {
			nS, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
			if err != nil {
				panic(err)
			}
			mtlDetails.Ns = float32(nS)
		} else if whiteSpaceSplit[0] == "Ka" {
			ka1, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
			if err != nil {
				panic(err)
			}
			ka2, err := strconv.ParseFloat(whiteSpaceSplit[2], 64)
			if err != nil {
				panic(err)
			}
			ka3, err := strconv.ParseFloat(whiteSpaceSplit[3], 64)
			if err != nil {
				panic(err)
			}

			mtlDetails.Ka = append(mtlDetails.Ka, float32(ka1))
			mtlDetails.Ka = append(mtlDetails.Ka, float32(ka2))
			mtlDetails.Ka = append(mtlDetails.Ka, float32(ka3))
		} else if whiteSpaceSplit[0] == "Kd" {
			kd1, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
			if err != nil {
				panic(err)
			}
			kd2, err := strconv.ParseFloat(whiteSpaceSplit[2], 64)
			if err != nil {
				panic(err)
			}
			kd3, err := strconv.ParseFloat(whiteSpaceSplit[3], 64)
			if err != nil {
				panic(err)
			}

			mtlDetails.Kd = append(mtlDetails.Kd, float32(kd1))
			mtlDetails.Kd = append(mtlDetails.Kd, float32(kd2))
			mtlDetails.Kd = append(mtlDetails.Kd, float32(kd3))
		} else if whiteSpaceSplit[0] == "Ks" {
			ks1, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
			if err != nil {
				panic(err)
			}
			ks2, err := strconv.ParseFloat(whiteSpaceSplit[2], 64)
			if err != nil {
				panic(err)
			}
			ks3, err := strconv.ParseFloat(whiteSpaceSplit[3], 64)
			if err != nil {
				panic(err)
			}

			mtlDetails.Ks = append(mtlDetails.Ks, float32(ks1))
			mtlDetails.Ks = append(mtlDetails.Ks, float32(ks2))
			mtlDetails.Ks = append(mtlDetails.Ks, float32(ks3))
		} else if whiteSpaceSplit[0] == "d" {
			d, err := strconv.ParseFloat(whiteSpaceSplit[1], 64)
			if err != nil {
				panic(err)
			}

			mtlDetails.D = float32(d)
		} else if whiteSpaceSplit[0] == "map_Kd" {
			mtlDetails.MapKD = whiteSpaceSplit[1]
		} */
	}

	fmt.Println(mtlDetails.Ka)
	return mtlDetails
}
