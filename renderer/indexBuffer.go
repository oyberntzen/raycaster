package renderer

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type IndexBuffer struct {
	indexBuffer uint32
	count       uint32
}

func NewIndexBuffer(data []uint32) *IndexBuffer {
	indexBuffer := &IndexBuffer{}
	indexBuffer.count = uint32(len(data))

	gl.GenBuffers(1, &indexBuffer.indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(indexBuffer.count)*4, gl.Ptr(data), gl.STATIC_DRAW)

	return indexBuffer
}

func (indexBuffer *IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer.indexBuffer)
}

func (indexBuffer *IndexBuffer) Delete() {
	gl.DeleteBuffers(1, &indexBuffer.indexBuffer)
}
