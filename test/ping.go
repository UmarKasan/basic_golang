package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// WINDOW width and height
const winWidth, winHeight int = 800, 600

// time variable
var frameStart time.Time
var elapsedTime float32

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color  color
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	color color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func getCenter() pos {
	return pos{float32(winWidth / 2), float32(winHeight) / 2}
}

func (ball *ball) update(paddle1 *paddle, paddle2 *paddle, elapsedTime float32) {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime
	// ball.x += ball.xv
	// ball.y += ball.yv

	if ball.y-ball.radius < 0 || ball.y+ball.radius > float32(winHeight) {
		ball.yv = -ball.yv
	}

	if ball.x-ball.radius < 0 || ball.x+ball.radius > float32(winWidth) {
		// ball.x = float32(winWidth / 2)
		// ball.y = float32(winHeight / 2)
		ball.pos = getCenter()
		// Delay before it repositions
		sdl.Delay(1000)

	}

	if ball.x < paddle1.x+paddle1.w/2 {
		if ball.y > paddle1.y-paddle1.h/2 && ball.y < paddle1.y+paddle1.h/2 {
			ball.xv = -ball.xv
		}
	}

	if ball.x > paddle2.x-paddle2.w/2 {
		if ball.y > paddle2.y-paddle2.h/2 && ball.y < paddle2.y+paddle2.h/2 {
			ball.xv = -ball.xv
		}
	}
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}

	// Ensure the paddle stays within the screen bounds
	if paddle.y < paddle.h/2 {
		paddle.y = paddle.h / 2
	} else if paddle.y > float32(winHeight)-paddle.h/2 {
		paddle.y = float32(winHeight) - paddle.h/2
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	if paddle.y < ball.y {
		paddle.y += paddle.speed * elapsedTime
	} else if paddle.y > ball.y {
		paddle.y -= paddle.speed * elapsedTime
	}
	// Ensure the paddle stays within the screen bounds
	if paddle.y < paddle.h/2 {
		paddle.y = paddle.h / 2
	} else if paddle.y > float32(winHeight)-paddle.h/2 {
		paddle.y = float32(winHeight) - paddle.h/2
	}
}

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
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{50, 100}, 20, 100, 300, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, 20, 100, 300, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 400, 400, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	running := true
	for running {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		clear(pixels)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		player1.update(keyState, elapsedTime)
		player2.aiUpdate(&ball, elapsedTime)
		ball.update(&player1, &player2, elapsedTime)

		tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		runtime.KeepAlive(pixels)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime := float32(time.Since(frameStart).Seconds())
		if elapsedTime < 1.0/60.0 {
			time.Sleep(time.Duration((1.0/60.0 - elapsedTime) * float32(time.Second)))
		}
		fmt.Println(elapsedTime)
		// sdl.Delay(16 - uint32(elapsedTime))
		// sdl.Delay(16)
	}
}
