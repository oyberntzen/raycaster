package renderer

import (
	"reflect"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type ShaderStorageBufferObject struct {
	ssbo uint32
}

func NewShaderStorageBufferObject(data interface{}, index uint32) *ShaderStorageBufferObject {
	ssbo := &ShaderStorageBufferObject{}

	gl.GenBuffers(1, &ssbo.ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo.ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, reflect.ValueOf(data).Len()*int(reflect.TypeOf(data).Elem().Size()), gl.Ptr(data), gl.DYNAMIC_COPY)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, index, ssbo.ssbo)

	return ssbo
}

func (ssbo *ShaderStorageBufferObject) Bind() {
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo.ssbo)
}

func (ssbo *ShaderStorageBufferObject) Delete() {
	gl.DeleteBuffers(1, &ssbo.ssbo)
}
