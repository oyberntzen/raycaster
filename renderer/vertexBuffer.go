package renderer

import (
	"reflect"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type shaderDataType int

const (
	ShaderDataTypeFloat32 shaderDataType = iota
	ShaderDataTypeInt32
	ShaderDataTypeUint32
	ShaderDataTypeBool
)

func (dataType shaderDataType) size() uint32 {
	switch dataType {
	case ShaderDataTypeFloat32:
		return 4
	case ShaderDataTypeInt32:
		return 4
	case ShaderDataTypeUint32:
		return 4
	case ShaderDataTypeBool:
		return 1
	}
	return 0
}

func (dataType shaderDataType) openGLType() uint32 {
	switch dataType {
	case ShaderDataTypeFloat32:
		return gl.FLOAT
	case ShaderDataTypeInt32:
		return gl.INT
	case ShaderDataTypeUint32:
		return gl.UNSIGNED_INT
	case ShaderDataTypeBool:
		return gl.BOOL
	}
	panic("unknown shader data type")
}

type VertexBuffer struct {
	vertexBuffer, vertexArray uint32
}

type LayoutElement struct {
	ShaderDataType shaderDataType
	Count          uint32
}

func NewVertexBuffer(data interface{}, layout []LayoutElement) *VertexBuffer {
	vertexBuffer := &VertexBuffer{}

	gl.GenBuffers(1, &vertexBuffer.vertexBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, reflect.ValueOf(data).Len()*int(reflect.TypeOf(data).Elem().Size()), gl.Ptr(data), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &vertexBuffer.vertexArray)
	gl.BindVertexArray(vertexBuffer.vertexArray)
	stride := int32(0)
	for _, element := range layout {
		stride += int32(element.ShaderDataType.size() * element.Count)
	}
	offset := uint32(0)
	for i, element := range layout {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointerWithOffset(uint32(i), int32(element.Count), element.ShaderDataType.openGLType(), false, stride, uintptr(offset))
		offset += element.Count * element.ShaderDataType.size()
	}

	return vertexBuffer
}

func (vertexBuffer *VertexBuffer) Bind() {
	// gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.vertexBuffer)
	gl.BindVertexArray(vertexBuffer.vertexArray)
}

func (vertexBuffer *VertexBuffer) Delete() {
	gl.DeleteBuffers(1, &vertexBuffer.vertexBuffer)
	gl.DeleteVertexArrays(1, &vertexBuffer.vertexArray)
}
