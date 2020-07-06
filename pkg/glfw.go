package pkg

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pkg/errors"
)

type InitGlfwOptions struct {
	Width  int
	Height int
	Title  string
}

func InitGlfw(opts InitGlfwOptions) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize glfw")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(opts.Width, opts.Height, opts.Title, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create glfw window")
	}

	window.MakeContextCurrent()

	return window, nil
}
