package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eliiasg/trifont"

	"github.com/eliiasg/fonttriangulator/ttf2mesh"
)

func main() {
	if len(os.Args) != 4 {
		panic("font path, output path and quality must be specified")
	}
	qual, err := strconv.ParseInt(os.Args[3], 10, 32)
	if err != nil || qual < 5 || qual > 100 {
		panic("Invalid quality, must be between 5 and 100")
	}
	font := ttf2mesh.LoadTTF(os.Args[1])
	defer font.Free()
	triFont := triangulate(font, uint8(qual))
	os.Remove(os.Args[2])
	f, err := os.Create(os.Args[2])
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}
	triFont.ToBinary(f)
	fmt.Println("Done")
}

func triangulate(font *ttf2mesh.File, quality uint8) *trifont.Font {
	// setup
	resFont := &trifont.Font{Chars: make(map[rune]trifont.Char)}
	mx := font.Ascender()
	mn := font.Descender()
	slideX := -font.MinLSideBearing()
	slideY := -mn
	scale := 1 / (mx - mn)
	// glyphs
	fmt.Println("Triangulating glyphs: ")
	start := time.Now().UnixMilli()
	for _, gly := range font.Glyphs() {
		fmt.Print(string(rune(gly.Symbol())))
		mesh := gly.ToMesh(quality)
		if mesh == nil {
			continue
		}
		// vertices
		verts := make([][2]float32, len(mesh.Vertices()))
		for i, vert := range mesh.Vertices() {
			verts[i] = transform(slideX, slideY, scale, vert)
		}
		// indices
		inds := make([]uint16, len(mesh.Indices()))
		for i, ind := range mesh.Indices() {
			inds[i] = uint16(ind)
		}
		resFont.Chars[rune(gly.Symbol())] = trifont.Char{
			Vertices: verts,
			Indices:  inds,
			Advance:  gly.Advance() * scale,
		}
	}
	fmt.Printf("\nFont triangulated in %v ms, saving...\n", time.Now().UnixMilli()-start)
	return resFont
}

func transform(slideX, slideY, scale float32, vert ttf2mesh.Vertex) [2]float32 {
	return [2]float32{(vert.X + slideX) * scale, (vert.Y + slideY) * scale}
}
