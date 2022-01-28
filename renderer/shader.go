package renderer

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/oyberntzen/gogm"
)

type Shader struct {
	shader           uint32
	uniformLocations map[string]int32
}

func NewShader(vertexShader, fragmentShader string) *Shader {
	vs, err := compileShader(gl.VERTEX_SHADER, vertexShader)
	if err != nil {
		log.Panic(err)
	}

	fs, err := compileShader(gl.FRAGMENT_SHADER, fragmentShader)
	if err != nil {
		log.Panic(err)
	}

	shader := &Shader{}
	shader.shader = gl.CreateProgram()

	gl.AttachShader(shader.shader, vs)
	gl.AttachShader(shader.shader, fs)
	gl.LinkProgram(shader.shader)
	gl.ValidateProgram(shader.shader)

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	return shader
}

func (shader *Shader) Bind() {
	gl.UseProgram(shader.shader)
}

func (shader *Shader) Delete() {
	gl.DeleteProgram(shader.shader)
}

func (shader *Shader) UploadUniformInt32(name string, value int32) {
	location := shader.uniformLocation(name)
	gl.Uniform1i(location, value)
}

func (shader *Shader) UploadUniformArrayInt32(name string, values []int32) {
	location := shader.uniformLocation(name)
	gl.Uniform1iv(location, int32(len(values)), &values[0])
}

func (shader *Shader) UploadUniformFloat32(name string, value float32) {
	location := shader.uniformLocation(name)
	gl.Uniform1f(location, value)
}

func (shader *Shader) UploadUniformVec2Float32(name string, value *gogm.Vec2[float32]) {
	location := shader.uniformLocation(name)
	gl.Uniform2f(location, value[0], value[1])
}

func (shader *Shader) UploadUniformVec3Float32(name string, value *gogm.Vec3[float32]) {
	location := shader.uniformLocation(name)
	gl.Uniform3f(location, value[0], value[1], value[2])
}

func (shader *Shader) UploadUniformVec4Float32(name string, value *gogm.Vec4[float32]) {
	location := shader.uniformLocation(name)
	gl.Uniform4f(location, value[0], value[1], value[2], value[3])
}

func (shader *Shader) uniformLocation(name string) int32 {
	if location, ok := shader.uniformLocations[name]; ok {
		return location
	}
	location := gl.GetUniformLocation(shader.shader, gl.Str(name+"\x00"))
	shader.uniformLocations[name] = location
	return location
}

func compileShader(shaderType uint32, source string) (uint32, error) {
	id := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(id, 1, csources, nil)
	free()
	gl.CompileShader(id)

	var result int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)

		message := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(id, length, &length, gl.Str(message))
		gl.DeleteShader(id)
		if shaderType == gl.VERTEX_SHADER {
			return 0, fmt.Errorf("failed to compile vertex shader: %v", message)
		} else if shaderType == gl.FRAGMENT_SHADER {
			return 0, fmt.Errorf("failed to compile fragment shader: %v", message)
		} else {
			return 0, fmt.Errorf("failed to compile unknown shader: %v", message)
		}

	}

	return id, nil
}
