package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// WINDOW width and height
const winWidth, winHeight int = 800, 600

// enum
type gameState int

const (
	start gameState = iota
	play
)

var state = start

// Font
var nums = [][]byte{
	{1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{1, 1, 0,
		0, 0, 1,
		0, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1}}

// Time
var frameStart time.Time

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
	score int
	color color
}

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
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

	if ball.y-ball.radius < 0 {
		ball.y = ball.radius
		ball.yv = -ball.yv
	} else if ball.y+ball.radius > float32(winHeight) {
		ball.y = float32(winHeight) - ball.radius
		ball.yv = -ball.yv
	}

	if ball.x < paddle1.x+paddle1.w/2 {
		if ball.y > paddle1.y-paddle1.h/2 && ball.y < paddle1.y+paddle1.h/2 {
			ball.xv = -ball.xv
			ball.x = paddle1.x + paddle1.w/2 + ball.radius
		}
	}

	if ball.x > paddle2.x-paddle2.w/2 {
		if ball.y > paddle2.y-paddle2.h/2 && ball.y < paddle2.y+paddle2.h/2 {
			ball.xv = -ball.xv
			ball.x = paddle2.x - paddle2.w/2 - ball.radius
		}
	}

	if ball.x < 0 {
		paddle2.score++
		ball.pos = getCenter()
		state = start
	} else if ball.x > float32(winWidth) {
		paddle1.score++
		ball.pos = getCenter()
		state = start
	}

	// if ball.x-ball.radius < 0 || ball.x+ball.radius > float32(winWidth) {
	// 	ball.pos = getCenter()
	// 	sdl.Delay(1000)
	// }
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8, controllerAxis int16, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}
	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
		paddle.y += paddle.speed * pct * elapsedTime
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

	// Joystick
	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{50, float32(winHeight) / 2}, 20, 100, 400, 0, color{255, 255, 255}}
	player2 := paddle{pos{float32(winWidth) - 50, float32(winHeight) / 2}, 20, 100, 400, 0, color{255, 255, 255}}
	ball := ball{pos{float32(winWidth) / 2, float32(winHeight) / 2}, 20, 300, 300, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var elapsedTime float32
	var controllerAxis int16

	running := true // not neceessary

	for running {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}

		if state == play {
			drawNumber(getCenter(), color{255, 255, 255}, 20, 2, pixels)
			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				state = play
			}
		}

		clear(pixels)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, unsafe.Pointer(&pixels[0]), winWidth*4)
		runtime.KeepAlive(pixels)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < 1.0/60.0 {
			time.Sleep(time.Duration((1.0/60.0 - elapsedTime) * float32(time.Second)))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
		// fmt.Printf("Ball position: (%.2f, %.2f), Player1: %.2f, Player2: %.2f\n", ball.x, ball.y, player1.y, player2.y)

	}
}
