package main

import (
	"bytes"
	"encoding/binary"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type Vertex interface {
	PosX() float64
	PosY() float64
	PosZ() float64

	Write(buf *Buffer)
}

type BaseVertex struct {
	x, y, z float64
}

func (v *BaseVertex) PosX() float64 {
	return v.x
}

func (v *BaseVertex) PosY() float64 {
	return v.y
}

func (v *BaseVertex) PosZ() float64 {
	return v.z
}

func (v *BaseVertex) Write(buf *Buffer) {
	binary.Write(buf.buf, binary.LittleEndian, v.x)
	binary.Write(buf.buf, binary.LittleEndian, v.y)
	binary.Write(buf.buf, binary.LittleEndian, v.z)
}

func NewVertex(x, y, z float64) Vertex {
	return &BaseVertex{x, y, z}
}

type Buffer struct {
	Id  uint32
	buf *bytes.Buffer

	CurrentVertex Vertex

	Vertices []Vertex
}

func (buf *Buffer) Vertex(vertex Vertex) {
	vertex.Write(buf)
	buf.Vertices = append(buf.Vertices, vertex)
}

func (b *Buffer) Bind(target uint32) {
	gl.BindBuffer(target, b.Id)
}

func (b *Buffer) Upload(target uint32) {
	gl.BufferData(target, len(b.buf.Bytes()), gl.Ptr(b.buf.Bytes()), gl.DYNAMIC_DRAW)
}

func CreateBuffer() *Buffer {
	var id uint32

	gl.GenBuffers(1, &id)

	return &Buffer{
		Id:       id,
		buf:      bytes.NewBuffer([]byte{}),
		Vertices: make([]Vertex, 0),
	}
}
