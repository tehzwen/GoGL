package common

import (
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
)

func WriteB64(filename, value string) {
	file, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(value))
	w.Close()
}

func ReadB64(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fz, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		panic(err)
	}

	return string(s)
}
