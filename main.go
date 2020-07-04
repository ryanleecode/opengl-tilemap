package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"opengl-tilemap/pkg"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/lafriks/go-tiled"
	"github.com/pkg/errors"
)

const (
	wWidth  = 800
	wHeight = 600
	title   = "OpenGL Tilemap"
)

func main() {
	runtime.LockOSThread()

	window, err := pkg.InitGlfw(
		pkg.InitGlfwOptions{Width: wWidth, Height: wHeight, Title: title})

	defer glfw.Terminate()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	prog, err := pkg.BuildProgram([]pkg.ShaderOptions{
		pkg.ShaderOptions{
			ShaderType: gl.VERTEX_SHADER,
			Shader:     "tile-vs.glsl",
			InputType:  pkg.ShaderInputFile,
		},
		pkg.ShaderOptions{
			ShaderType: gl.FRAGMENT_SHADER,
			Shader:     "tile-fs.glsl",
			InputType:  pkg.ShaderInputFile,
		},
	})
	if err != nil {
		panic(err)
	}

	tilemap, err := loadTilemap("tilemap.png")
	if err != nil {
		panic(err)
	}

	gameMap, err := tiled.LoadFromFile("tilemap.tmx")
	if err != nil {
		panic(err)
	}

	var vertices []float32
	for i := 0; i < gameMap.Width+1; i++ {
		for j := 0; j < gameMap.Height+1; j++ {
			vertices = append(vertices, float32(i), float32(j))
		}
	}

	var indices []uint
	for i := 0; i < gameMap.Width; i++ {
		for j := 0; j < gameMap.Height; j++ {
			topLeft := uint(i)
			topRight := uint(i + 1)
			bottomRight := uint(((j + 1) * (gameMap.Width + 1)) + i + 1)
			bottomLeft := uint(((j + 1) * (gameMap.Width + 1)) + i)

			indices = append(
				indices,
				topLeft, bottomLeft, bottomRight,
				topLeft, topRight, bottomRight)
		}
	}

	var vao, vbo, ebo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(vertices),
		gl.Ptr(vertices),
		gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,
		len(indices),
		gl.Ptr(indices),
		gl.STATIC_DRAW)

	gl.VertexAttribPointer(
		0, 3, gl.FLOAT, false,
		int32(unsafe.Sizeof(float32(0)))*3,
		gl.PtrOffset(0),
	)
	gl.EnableVertexAttribArray(0)

	//	fmt.Printf("%v", vertices)
	fmt.Printf("%v", indices)

	for _, layer := range gameMap.Layers {
		for _, tile := range layer.Tiles {
			if tile.Nil {
			}
		}
	}

	// tilemapWidth := (int32)(tilemap.Bounds().Dx())
	// tilemapHeight := (int32)(tilemap.Bounds().Dy())
	// const stride = 16

	texture := bindTexture(tilemap)

	tileWidth := gameMap.TileWidth
	log.Println(tileWidth)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(prog)
		gl.BindVertexArray(vao)

		gl.DrawElements(
			gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

		gl.BindTexture(gl.TEXTURE_2D, texture)

		window.SwapBuffers()
		glfw.PollEvents()
	}

}

func loadTilemap(fileName string) (*image.RGBA, error) {
	tilemapFile, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load tilemap %s", fileName)
	}

	tilemapPNG, err := png.Decode(tilemapFile)
	if err != nil {
		return nil, errors.Wrapf(err, "tilemap %s is not a png file", fileName)
	}

	return image.NewRGBA(tilemapPNG.Bounds()), nil
}

func bindTexture(img *image.RGBA) uint32 {
	bounds := img.Bounds()

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		(int32)(bounds.Dx()), (int32)(bounds.Dy()),
		0,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(image.NewRGBA(bounds).Pix))

	return texture
}
