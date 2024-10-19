package main

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// WINDOW width and height
const winWidth, winHeight int = 800, 600

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius int
	xv     float32
	yv     float32
	color  color
}

type paddle struct {
	pos
	w     int
	h     int
	color color
}

func (ball *ball) draw(pixels []byte) {
	//YAGNI - Ya Aint Gonna need it

	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(paddle1 *paddle, paddle2 *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv

	if int(ball.y)-ball.radius < 0 || int(ball.y)+ball.radius > winHeight {
		ball.yv = -ball.yv
	}

	if int(ball.x)-ball.radius < 0 || int(ball.x)+ball.radius > winWidth {
		ball.x = float32(winWidth / 2)
		ball.y = float32(winHeight / 2)

		// Delay before it repositions
		sdl.Delay(1000)
	}

	if int(ball.x) < int(paddle1.x)+paddle1.w/2 {
		if int(ball.y) > int(paddle1.y)-paddle1.h/2 && int(ball.y) < int(paddle1.y)+paddle1.h/2 {
			ball.xv = -ball.xv
		}
	}

	if int(ball.x) > int(paddle2.x)+paddle2.w/2 {
		if int(ball.y) > int(paddle2.y)-paddle2.h/2 && int(ball.y) < int(paddle2.y)+paddle2.h/2 {
			ball.xv = -ball.xv
		}
	}
}

// Create Pong paddle
func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

// Create new paddle location
func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= 10
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += 10
	}
}

// AI player
func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}

// / Clear old pixel
func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func main() {
	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}

	// defer calls when function exits
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer calls when function exits
	defer renderer.Destroy()

	// Texture
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer calls when function exits
	defer tex.Destroy()

	// Draw pixels
	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{50, 100}, 20, 100, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 10, 5, color{255, 255, 255}}

	// Keyboard Input
	keyState := sdl.GetKeyboardState()

	// ensure continuous running
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
		clear(pixels)
		player1.draw(pixels)
		player1.update(keyState)
		player2.draw(pixels)
		player2.aiUpdate(&ball)
		ball.draw(pixels)
		ball.update(&player1, &player2)

		// have to put tex in game loop
		tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		// Delay 1 frame
		sdl.Delay(16)
	}
}
