package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"opengl-tilemap/pkg"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lafriks/go-tiled"
	"github.com/pkg/errors"
)

const (
	wWidth  = 800
	wHeight = 600
	title   = "OpenGL Tilemap"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if err := glfw.Init(); err != nil {
		panic(errors.Wrap(err, "failed to initialize glfw"))
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(wWidth, wHeight, title, nil, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to create glfw window"))
	}

	window.MakeContextCurrent()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	if err = gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	prog, err := pkg.BuildProgram([]pkg.ShaderOptions{
		{
			ShaderType: gl.VERTEX_SHADER,
			Shader:     "tile-vs.glsl",
			InputType:  pkg.ShaderInputFile,
		},
		{
			ShaderType: gl.FRAGMENT_SHADER,
			Shader:     "tile-fs.glsl",
			InputType:  pkg.ShaderInputFile,
		},
	})
	if err != nil {
		panic(err)
	}

	gl.UseProgram(prog)

	gameMap, err := tiled.LoadFromFile("tilemap.tmx")
	if err != nil {
		panic(err)
	}

	gameMapTileWidth := gameMap.Width
	gameMapTileHeight := gameMap.Height

	var vertices []float32

	var indices []uint32

	// var uvs []float32
	i := 0
	for _, layer := range gameMap.Layers {
		for _, tile := range layer.Tiles {

			if !tile.Nil {
				x := i % gameMapTileWidth
				y := -(i / gameMapTileHeight)

				xOffset := float32(gameMapTileWidth) / 2
				yOffset := float32(gameMapTileHeight)/2 - 1

				tileSetColCount := tile.Tileset.Columns
				tileSetRowCount := tile.Tileset.TileCount / tileSetColCount

				// uOffset := float32(tile.Tileset.Spacing*x) / float32(gameMapWidth*tileSetWidth)
				// vOffset := float32(tile.Tileset.Spacing*y) / float32(gameMapHeight*tileSetHeight)

				initialUVOffset := uint32(tile.Tileset.Margin)

				//	pixelWidth := uint32(tileSetWidth) * uint32(tile.Tileset.TileWidth+tile.Tileset.Spacing)

				col := tile.ID % uint32(tileSetColCount)
				uDenom := float32(uint32(tileSetColCount) * uint32(tile.Tileset.TileWidth+tile.Tileset.Spacing))
				u := initialUVOffset + col*uint32(tile.Tileset.TileWidth) + col*uint32(tile.Tileset.Spacing)

				row := tile.ID / uint32(tileSetColCount)
				vDenom := float32(uint32(tileSetRowCount) * uint32(tile.Tileset.TileHeight+tile.Tileset.Spacing))
				v := initialUVOffset + row*uint32(tile.Tileset.TileHeight) + row*uint32(tile.Tileset.Spacing)

				vertices = append(vertices,
					float32(x)+1-xOffset, float32(y)+1+yOffset, // Top Right V
					float32(u+uint32(tile.Tileset.TileWidth))/uDenom, 1-(float32(v)/vDenom),

					float32(x)-xOffset, float32(y)+1+yOffset, // Top Left V
					float32(u)/uDenom, 1-(float32(v)/vDenom),

					float32(x)-xOffset, float32(y)+yOffset, // Bottom Left V
					float32(u)/uDenom, 1-(float32(v+uint32(tile.Tileset.TileHeight))/vDenom),

					float32(x)+1-xOffset, float32(y)+yOffset, // Bottom Right V
					float32(u+uint32(tile.Tileset.TileWidth))/uDenom, 1-(float32(v+uint32(tile.Tileset.TileHeight))/vDenom),
				)

				indices = append(
					indices,
					uint32(i*4), uint32(i*4+1), uint32(i*4+2),
					uint32(i*4), uint32(i*4+2), uint32(i*4+3))

			}
			i++
		}
	}

	tileset := gameMap.Tilesets[0]
	tileSetImg, err := loadTileset(tileset.Image.Source)
	if err != nil {
		panic(err)
	}
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE) // set texture wrapping to GL_REPEAT (default wrapping method)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// set texture figl.ring parametgl.
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_LINEAR)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(tileSetImg.Rect.Size().X), int32(tileSetImg.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(tileSetImg.Pix))

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(wWidth)/wHeight, 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(prog, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 20.5}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(prog, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(prog, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(prog, gl.Str("tileAtlas\x00"))
	gl.Uniform1i(textureUniform, 0)

	var vao, vbo, ebo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		int(unsafe.Sizeof(float32(0)))*len(vertices),
		gl.Ptr(vertices),
		gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,
		int(unsafe.Sizeof(uint32(0)))*len(indices),
		gl.Ptr(indices),
		gl.STATIC_DRAW)

	gl.VertexAttribPointer(
		0, 2, gl.FLOAT, false,
		int32(unsafe.Sizeof(float32(0)))*4,
		gl.PtrOffset(0),
	)
	gl.EnableVertexAttribArray(0)

	println(2 * int(unsafe.Sizeof(float32(0))))
	gl.VertexAttribPointer(
		1, 2, gl.FLOAT, false,
		int32(unsafe.Sizeof(float32(0)))*4,
		gl.PtrOffset(2*int(unsafe.Sizeof(float32(0)))),
	)
	gl.EnableVertexAttribArray(1)

	tileWidth := gameMap.TileWidth
	log.Println(tileWidth)

	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.UseProgram(prog)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
		gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)

		/* 		gl.DrawElements(
		gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))  */
		gl.DrawElements(
			gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

		// gl.BindTexture(gl.TEXTURE_2D, texture)

		window.SwapBuffers()
		glfw.PollEvents()
	}

}

func loadTileset(fileName string) (*image.RGBA, error) {
	tilesetFile, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load tileset %s", fileName)
	}

	tilesetPNG, err := png.Decode(tilesetFile)
	if err != nil {
		return nil, errors.Wrapf(err, "tileset %s is not a png file", fileName)
	}

	rgba := image.NewRGBA(tilesetPNG.Bounds())

	draw.Draw(rgba, rgba.Bounds(), imaging.FlipV(tilesetPNG), image.Point{0, 0}, draw.Src)

	return rgba, nil
}

func loadTexture(img *image.RGBA) uint32 {

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(img.Rect.Size().X), int32(img.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	return texture
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
