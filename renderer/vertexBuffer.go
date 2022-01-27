package renderer

import (
	"reflect"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type VertexBuffer struct {
	vertexBuffer, vertexArray uint32
}

func NewVertexBuffer(data interface{}) *VertexBuffer {
	vertexBuffer := &VertexBuffer{}

	gl.GenBuffers(1, &vertexBuffer.vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, reflect.ValueOf(data).Len()*int(reflect.TypeOf(data).Elem().Size()), gl.Ptr(data), gl.STATIC_DRAW)

	return vertexBuffer
}

func (vertexBuffer *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.vertexBuffer)
}

func (vertexBuffer *VertexBuffer) Delete() {
	gl.DeleteBuffers(1, &vertexBuffer.vertexBuffer)
}
