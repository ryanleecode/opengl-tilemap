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
	"github.com/go-gl/gl/v3.3-core/gl"
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	runtime.LockOSThread()

	window, err := pkg.InitGlfw(
		pkg.InitGlfwOptions{Width: wWidth, Height: wHeight, Title: title})

	defer glfw.Terminate()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

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

	gl.UseProgram(prog)

	gameMap, err := tiled.LoadFromFile("tilemap.tmx")
	if err != nil {
		panic(err)
	}

	gameMapWidth := gameMap.Width
	gameMapHeight := gameMap.Height
	colSize := gameMapWidth + 1
	rowSize := gameMapHeight + 1

	var vertices []float32
	for j := 0; j < rowSize; j++ {
		for i := 0; i < colSize; i++ {
			//	vertices = append(vertices,
			//		float32(i)-float32(gameMapWidth)/2,
			//		float32(j)-float32(gameMapHeight)/2,
			//	)
		}
	}

	fmt.Printf("%v\n", vertices)

	var indices []uint32

	// var uvs []float32
	i := 0
	for _, layer := range gameMap.Layers {
		for _, tile := range layer.Tiles {

			if !tile.Nil {
				x := i % gameMapWidth
				y := -(i / gameMapHeight)

				xOffset := float32(gameMapWidth) / 2
				yOffset := float32(gameMapHeight)/2 - 1

				tileSetWidth := tile.Tileset.Columns
				tileSetHeight := tile.Tileset.TileCount / tileSetWidth

				asdf := float32((tile.ID + 1) % uint32(tileSetWidth))
				if float32(tile.ID+1)/float32(tileSetWidth) == float32((tile.ID+1)/uint32(tileSetWidth)) {
					asdf = float32(tileSetWidth)
				}

				vertices = append(vertices,
					float32(x)+1-xOffset, float32(y)+1+yOffset, // Top Right V
					asdf/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth))/float32(tileSetHeight)),

					float32(x)-xOffset, float32(y)+1+yOffset, // Top Left V
					float32(tile.ID%uint32(tileSetWidth))/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth))/float32(tileSetHeight)),

					float32(x)-xOffset, float32(y)+yOffset, // Bottom Left V
					float32(tile.ID%uint32(tileSetWidth))/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth)+1)/float32(tileSetHeight)),

					float32(x)+1-xOffset, float32(y)+yOffset, // Bottom Right V
					asdf/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth)+1)/float32(tileSetHeight)),
				)

				indices = append(
					indices,
					uint32(i*4), uint32(i*4+1), uint32(i*4+2),
					uint32(i*4), uint32(i*4+2), uint32(i*4+3))

				fmt.Printf("%d, %d, %d, %d, %d, [(%f, %f), (%f, %f), (%f, %f), (%f, %f)]\n", x, y, tileSetWidth, tileSetHeight, tile.ID,
					asdf/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth))/float32(tileSetHeight)),

					float32(tile.ID%uint32(tileSetWidth))/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth))/float32(tileSetHeight)),
					float32(tile.ID%uint32(tileSetWidth))/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth)+1)/float32(tileSetHeight)),
					asdf/float32(tileSetWidth), 1-(float32(tile.ID/uint32(tileSetWidth)+1)/float32(tileSetHeight)),
				)
			}
			i++
		}
	}

	// fmt.Printf("%v\n", vertices)

	for j := 0; j < gameMapHeight; j++ {
		for i := 0; i < gameMapWidth; i++ {
			/* 			bottomLeft := uint32(i + (j * colSize))
			   			bottomRight := uint32(i + (j * colSize) + 1)
			   			topLeft := uint32(((j + 1) * colSize) + i)
			   			topRight := uint32(((j + 1) * colSize) + i + 1)

			   			indices = append(
			   				indices,
			   				topRight, topLeft, bottomLeft,
			   				topRight, bottomLeft, bottomRight) */
		}
	}

	/*
		tilemap, err := loadTilemap("tilemap.png")
		if err != nil {
			panic(err)
		} */

	/* 	fmt.Printf("%v\n", vertices)
	   	fmt.Printf("%v\n", indices) */

	/* 	indices = []uint32{
	0, 1, 4, 0, 3, 4,
	1, 2, 5, 1, 4, 5}
	*/
	/* 	indices = []uint32{
	   		1 2 5 1 4 5}
	   	vertices = []float32{
	   		0.5, 0.5, 0.0,
	   		0.5, -0.5, 0.0,
	   		-0.5, -0.5, 0.0,
	   		-0.5, 0.5, 0.0,
	   	}
	*/

	//	var vertices []float32
	//	i := 0
	//	for _, layer := range gameMap.Layers {
	//		for _, tile := range layer.Tiles {
	//			if !tile.Nil {
	//				x := i % gameMapWidth
	//				y := i / gameMapHeight
	//				fmt.Printf("%d, %d\n", x, y)
	//				vertices = append(vertices, float32(x)-float32(gameMapWidth)/2, float32(y)-float32(gameMapHeight)/2)
	//				vertices = append(vertices,
	//					float32(0.0499219968799), float32(0.0499219968799),
	//					float32(0), float32(0.0499219968799),
	//					float32(0), float32(0),
	//					float32(0.0499219968799), float32(0.0499219968799),
	//					float32(0), float32(0),
	//					float32(0.0499219968799), float32(0),
	//				)
	//			}
	//
	//			i += 1
	//
	//		}
	//	}

	tileset := gameMap.Tilesets[0]
	tileSetImg, err := loadTileset(tileset.Image.Source)
	if err != nil {
		panic(err)
	}
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT) // set texture wrapping to GL_REPEAT (default wrapping method)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set texture figl.ring parametgl.
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST_MIPMAP_LINEAR)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
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

	textureUniform := gl.GetUniformLocation(prog, gl.Str("texture1\x00"))
	gl.Uniform1i(textureUniform, 0)
	//
	// 	gl.BindFragDataLocation(prog, 0, gl.Str("outputColor\x00"))

	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// tilemapWidth := (int32)(tilemap.Bounds().Dx())
	// tilemapHeight := (int32)(tilemap.Bounds().Dy())
	// const stride = 16

	// texture := bindTexture(tilemap)

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
