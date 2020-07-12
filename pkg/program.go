package pkg

import (
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/pkg/errors"
)

type ShaderInputType uint32

const (
	ShaderInputFile ShaderInputType = iota
	ShaderInputRaw
)

type ShaderOptions struct {
	ShaderType uint32
	Shader     string
	InputType  ShaderInputType
}

func CompileShader(options ShaderOptions) (uint32, error) {
	switch options.InputType {
	case ShaderInputFile:
		return compileShaderFromFile(options.Shader, uint32(options.ShaderType))
	case ShaderInputRaw:
		return compileShaderFromSource(options.Shader, uint32(options.ShaderType))
	default:
		return 0, errors.Errorf("unsupported input type: %d", options.InputType)
	}
}

func compileShaderFromFile(fileName string, shaderType uint32) (uint32, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, errors.Wrapf(err, "failed compile shader file %s", fileName)
	}

	return compileShaderFromSource(string(b)+"\x00", shaderType)
}

func compileShaderFromSource(src string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	glSrc, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, glSrc, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, errors.Wrapf(errors.New(log), "failed to compile shader source %v", src)
	}

	return shader, nil
}

func BuildProgram(shaderOptionsList []ShaderOptions) (uint32, error) {
	prog := gl.CreateProgram()

	for _, shaderOptions := range shaderOptionsList {
		shader, err := CompileShader(shaderOptions)
		if err != nil {
			return 0, errors.Wrap(err, "failed to build program")
		}
		gl.AttachShader(prog, shader)
		gl.DeleteShader(shader)
	}

	gl.LinkProgram(prog)

	return prog, nil
}
