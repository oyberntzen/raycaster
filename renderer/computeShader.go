package renderer

import (
	"log"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/oyberntzen/gogm"
)

type ComputeShader struct {
	shader           uint32
	uniformLocations map[string]int32
}

func NewComputeShader(computeShader string) *ComputeShader {
	cs, err := compileShader(gl.COMPUTE_SHADER, computeShader)
	if err != nil {
		log.Panic(err)
	}

	shader := &ComputeShader{}
	shader.uniformLocations = make(map[string]int32)
	shader.shader = gl.CreateProgram()

	gl.AttachShader(shader.shader, cs)
	gl.LinkProgram(shader.shader)

	var result int32
	gl.GetProgramiv(shader.shader, gl.LINK_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetProgramiv(shader.shader, gl.INFO_LOG_LENGTH, &length)

		message := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(shader.shader, length, &length, gl.Str(message))
		gl.DeleteProgram(shader.shader)
		log.Panic(message)

	}

	gl.ValidateProgram(shader.shader)

	gl.DeleteShader(cs)

	return shader
}

func (shader *ComputeShader) Bind() {
	gl.UseProgram(shader.shader)
}

func (shader *ComputeShader) Delete() {
	gl.DeleteProgram(shader.shader)
}

func (shader *ComputeShader) UploadUniformInt32(name string, value int32) {
	location := shader.uniformLocation(name)
	gl.Uniform1i(location, value)
}

func (shader *ComputeShader) UploadUniformArrayInt32(name string, values []int32) {
	location := shader.uniformLocation(name)
	gl.Uniform1iv(location, int32(len(values)), &values[0])
}

func (shader *ComputeShader) UploadUniformFloat32(name string, value float32) {
	location := shader.uniformLocation(name)
	gl.Uniform1f(location, value)
}

func (shader *ComputeShader) UploadUniformVec2Float32(name string, value *gogm.Vec2[float32]) {
	location := shader.uniformLocation(name)
	gl.Uniform2f(location, value[0], value[1])
}

func (shader *ComputeShader) UploadUniformVec3Float32(name string, value *gogm.Vec3[float32]) {
	location := shader.uniformLocation(name)
	gl.Uniform3f(location, value[0], value[1], value[2])
}

func (shader *ComputeShader) UploadUniformVec4Float32(name string, value *gogm.Vec4[float32]) {
	location := shader.uniformLocation(name)
	gl.Uniform4f(location, value[0], value[1], value[2], value[3])
}

func (shader *ComputeShader) uniformLocation(name string) int32 {
	if location, ok := shader.uniformLocations[name]; ok {
		return location
	}
	location := gl.GetUniformLocation(shader.shader, gl.Str(name+"\x00"))
	shader.uniformLocations[name] = location
	return location
}

func (shader *ComputeShader) Run(numGroupsX, numGroupsY, numGroupsZ uint32) {
	gl.DispatchCompute(numGroupsX, numGroupsY, numGroupsZ)
}
