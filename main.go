package main

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Invalid argument supplied")
		return
	}
	imgPath := os.Args[1]

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	basePath := dir + "/"
	if imgPath == "." {
		walk(basePath, "")
	} else {
		walk(basePath, imgPath)
	}
}

func walk(base, dir string) {
	thisDir := filepath.Join(base, dir)

	fmt.Printf("Walking through %s\n", thisDir)

	files, err := ioutil.ReadDir(thisDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() && (f.Mode()&os.ModeSymlink) != os.ModeSymlink {
			convert(thisDir, f)
		} else {
			walk(thisDir, f.Name())
		}
	}
}

func convert(path string, img os.FileInfo) {
	var buf bytes.Buffer

	imgPath := filepath.Join(path, img.Name())

	fmt.Printf("Converting %s\n", imgPath)

	m := decodeImg(imgPath)

	if m == nil {
		log.Println("Skipped")
		log.Println()
		return
	}

	fmt.Printf("Encoding %s\n", imgPath)

	if err := webp.Encode(&buf, m, &webp.Options{Quality:70.0}); err != nil {
		log.Println(err)
	}

	fmt.Printf("Saving %s\n", imgPath)

	if err := ioutil.WriteFile(filepath.Join(path, img.Name()[0:len(img.Name())-len(filepath.Ext(imgPath))]) + ".webp", buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	fmt.Println("Done")
	fmt.Println()
}

func decodeImg(imgPath string) image.Image {
	var m image.Image
	var data []byte

	fmt.Printf("Decoding %s\n", imgPath)

	data, err := ioutil.ReadFile(imgPath)

	ext := filepath.Ext(imgPath)

	switch ext {
	case ".jpg":
		m, err = jpeg.Decode(bytes.NewReader(data))
		break
	case ".jpeg":
		m, err = jpeg.Decode(bytes.NewReader(data))
		break
	case ".png":
		m, err = png.Decode(bytes.NewReader(data))
		break
	}
	if err != nil {
		log.Println(err)
	}
	return m
}