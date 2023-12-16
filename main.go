package main

import (
	"fmt"
	"time"

	m "github.com/aoyama-val/go-u235/model"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEEN_WIDTH = 640
	SCREEN_HEIGHT = 400
	CELL_SIZE_PX  = 16
	FPS           = 30
)

type Texture struct {
	texture *sdl.Texture
	w       int32
	h       int32
}

type Resources struct {
	textures map[string]*Texture
	chunks   map[string]*mix.Chunk
}

func main() {
	// [SDL2の使い方 - mirichiの日記](https://mirichi.hatenadiary.org/entry/20141018/p1)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("go-u235", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, SCREEEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	sdl.ShowCursor(sdl.DISABLE)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	err = renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	initMixer()
	defer mix.Quit()
	defer mix.CloseAudio()

	resources := loadResources(renderer)

	running := true
	game := m.NewGame()

	printUsage()

	for running {
		commands := GetCommands()
		if len(commands) > 0 && commands[0] == m.COMMAND_QUIT {
			running = false
			break
		}
		if game.IsOver && len(commands) > 0 && commands[0] == m.COMMAND_RESTART {
			game = game.Restart()
			continue
		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
		game.Update(commands)
		render(renderer, window, game, resources)
		playSounds(game, resources)
		time.Sleep((1000 / FPS) * time.Millisecond)
	}
}

func printUsage() {
	fmt.Printf("Keys:\n")
	fmt.Printf("    Left   : Move player left\n")
	fmt.Printf("    Right  : Move player right\n")
	fmt.Printf("    Shift  : Shoot\n")
	fmt.Printf("    Escape : Quit game\n")
	fmt.Printf("    Space  : Restart when game over\n")
}

func GetCommands() []m.Command {
	commands := make([]m.Command, 0, 3)
	keyState := sdl.GetKeyboardState()
	if keyState[sdl.SCANCODE_LEFT] != 0 {
		commands = append(commands, m.COMMAND_LEFT)
	}
	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		commands = append(commands, m.COMMAND_RIGHT)
	}
	if keyState[sdl.SCANCODE_LSHIFT] != 0 || keyState[sdl.SCANCODE_RSHIFT] != 0 {
		commands = append(commands, m.COMMAND_SHOOT)
	}
	if keyState[sdl.SCANCODE_ESCAPE] != 0 {
		commands = append(commands, m.COMMAND_QUIT)
	}
	if keyState[sdl.SCANCODE_SPACE] != 0 {
		commands = append(commands, m.COMMAND_RESTART)
	}
	return commands
}

func initMixer() {
	if err := mix.Init(int(mix.WAV)); err != nil {
		panic("cannot init mixer")
	}

	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		panic("cannot open audio")
	}
}

func loadResources(renderer *sdl.Renderer) *Resources {
	var resources Resources
	resources.textures = make(map[string]*Texture)
	resources.chunks = make(map[string]*mix.Chunk)

	var imagePaths []string = []string{
		"back.bmp",
		"down.bmp",
		"dust.bmp",
		"left.bmp",
		"numbers.bmp",
		"player.bmp",
		"right.bmp",
		"target.bmp",
		"title.bmp",
		"up.bmp",
		"wall.bmp",
	}
	for _, path := range imagePaths {
		fullPath := "resources/image/" + path
		image, err := sdl.LoadBMP(fullPath)
		if err != nil {
			panic("cannot load image: " + path)
		}
		texture, err := renderer.CreateTextureFromSurface(image)
		if err != nil {
			panic("cannot convert to texture: " + path)
		}

		_, _, w, h, err := texture.Query()
		if err != nil {
			panic("cannot query texture: " + path)
		}

		resources.textures[path] = &Texture{
			texture: texture,
			w:       w,
			h:       h,
		}
	}

	var soundPaths []string = []string{
		"crash.wav",
		"hit.wav",
	}
	for _, path := range soundPaths {
		fullPath := "resources/sound/" + path
		chunk, err := mix.LoadWAV(fullPath)
		if err != nil {
			panic("cannot load wav: " + path)
		}
		resources.chunks[path] = chunk
	}

	return &resources
}

func render(renderer *sdl.Renderer, window *sdl.Window, game *m.Game, resources *Resources) {
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()

	renderTexture(renderer, resources, "title.bmp", 1, 0)

	for i := 1; i <= 22; i++ {
		renderTexture(renderer, resources, "wall.bmp", m.X_MIN-1, i)
		renderTexture(renderer, resources, "wall.bmp", m.X_MAX+1, i)
	}
	for i := m.X_MIN; i <= m.X_MAX; i++ {
		renderTexture(renderer, resources, "wall.bmp", i, 1)
	}

	for i := 0; i < SCREEEN_WIDTH/CELL_SIZE_PX; i++ {
		renderTexture(renderer, resources, "back.bmp", i, m.Y_MAX+1)
	}

	renderNumber(renderer, resources, 18, 0, fmt.Sprintf("%8d", game.HighScore))
	renderNumber(renderer, resources, 32, 0, fmt.Sprintf("%8d", game.Score))

	renderTexture(renderer, resources, "player.bmp", game.Player.X, game.Player.Y)

	for _, bullet := range game.Bullets {
		switch bullet.Direction {
		case m.DIRECTION_LEFT:
			renderTexture(renderer, resources, "left.bmp", bullet.X, bullet.Y)
		case m.DIRECTION_RIGHT:
			renderTexture(renderer, resources, "right.bmp", bullet.X, bullet.Y)
		case m.DIRECTION_UP:
			renderTexture(renderer, resources, "up.bmp", bullet.X, bullet.Y)
		case m.DIRECTION_DOWN:
			renderTexture(renderer, resources, "down.bmp", bullet.X, bullet.Y)
		}
	}

	for _, target := range game.Targets {
		renderTexture(renderer, resources, "target.bmp", target.X, target.Y)
	}

	if game.IsOver {
		renderer.SetDrawColor(255, 0, 0, 128)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: SCREEEN_WIDTH, H: SCREEN_HEIGHT})
	}

	renderer.Present()
}

func renderTexture(renderer *sdl.Renderer, resources *Resources, textureKey string, x int, y int) {
	texture := resources.textures[textureKey]
	renderer.Copy(texture.texture, nil, &sdl.Rect{X: int32(CELL_SIZE_PX * x), Y: int32(CELL_SIZE_PX * y), W: texture.w, H: texture.h})
}

func renderNumber(renderer *sdl.Renderer, resources *Resources, x int, y int, numstr string) {
	texture := resources.textures["numbers.bmp"]
	digitWidthInPx := 8
	xInPx := int32(CELL_SIZE_PX * x)
	yInPx := int32(CELL_SIZE_PX * y)
	for i := 0; i < len(numstr); i++ {
		digit := numstr[i]
		if 0x30 <= digit && digit <= 0x39 {
			renderer.Copy(
				texture.texture,
				&sdl.Rect{X: int32(digitWidthInPx * int(digit-0x30)), Y: 0, W: int32(digitWidthInPx), H: texture.h},
				&sdl.Rect{X: xInPx, Y: yInPx, W: int32(digitWidthInPx), H: texture.h},
			)
		}
		xInPx += int32(digitWidthInPx)
	}
}

func playSounds(game *m.Game, resources *Resources) {
	for _, requestedSound := range game.RequestedSounds {
		chunk := resources.chunks[requestedSound]
		chunk.Play(-1, 0)
	}
	game.RequestedSounds = nil
}
