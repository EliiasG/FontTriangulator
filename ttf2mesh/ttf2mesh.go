package ttf2mesh

//#cgo CFLAGS: -w
//#include "ttf2mesh.h"
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

const (
	QualityLow    uint8 = 10
	QualityNormal uint8 = 20
	Qualityhigh   uint8 = 50
)

type File C.ttf_t
type Glyph C.ttf_glyph_t
type Mesh C.ttf_mesh_t

type Vertex struct {
	X, Y float32
}

// File

func LoadTTF(path string) *File {
	var file *C.ttf_t
	str := C.CString(path)
	C.ttf_load_from_file(str, &file, false)
	C.free(unsafe.Pointer(str))
	return (*File)(file)
}

func (f *File) Glyphs() []Glyph {
	return unsafe.Slice((*Glyph)(f.glyphs), f.nglyphs)
}

func (f *File) Free() {
	C.ttf_free((*C.ttf_t)(f))
}

func (f *File) Ascender() float32 {
	return float32(f.hhea.ascender)
}

func (f *File) Descender() float32 {
	return float32(f.hhea.descender)
}

func (f *File) MinLSideBearing() float32 {
	return float32(f.hhea.minLSideBearing)
}

// Glyph

func (g *Glyph) ToMesh(quality uint8) *Mesh {
	var mesh *C.ttf_mesh_t
	C.ttf_glyph2mesh((*C.ttf_glyph_t)(g), &mesh, C.uchar(quality), 1)
	return (*Mesh)(mesh)
}

func (g *Glyph) Symbol() int32 {
	return int32(g.symbol)
}

func (g *Glyph) XBounds() [2]float32 {
	return *(*[2]float32)(unsafe.Pointer(&g.xbounds))
}

func (g *Glyph) YBounds() [2]float32 {
	return *(*[2]float32)(unsafe.Pointer(&g.ybounds))
}

func (g *Glyph) Advance() float32 {
	return float32(g.advance)
}

// Mesh

func (m *Mesh) Vertices() []Vertex {
	// convert C pointer to Vertex pointer, then convert pointer to slice
	return unsafe.Slice((*Vertex)(unsafe.Pointer(m.vert)), m.nvert)
}

func (m *Mesh) Indices() []int32 {
	// convert C pointer to int32 pointer, then convert pointer to slice, but 3 times longer since 3 indeices per face
	return unsafe.Slice((*int32)(unsafe.Pointer(m.faces)), m.nfaces*3)
}

func (m *Mesh) Free() {
	C.ttf_free_mesh((*C.ttf_mesh_t)(m))
}
