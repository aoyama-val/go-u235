package main

import (
	"time"

	m "github.com/aoyama-val/go-u235/model"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEEN_WIDTH = 640
	SCREEN_HEIGHT = 400
	CELL_SIZE_PX  = 20
	FPS           = 30
)

type Resources struct {
	textures map[string]*sdl.Texture
	chunks   map[string]*mix.Chunk
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("go-u235", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, SCREEEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	err = renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	// Initialize SDL2 mixer
	if err := mix.Init(int(mix.WAV)); err != nil {
		panic("cannot init mixer")
	}
	defer mix.Quit()

	// Open default playback device
	if err := mix.OpenAudio(int(mix.DEFAULT_FREQUENCY), uint16(mix.DEFAULT_FORMAT), int(mix.DEFAULT_CHANNELS), mix.DEFAULT_CHUNKSIZE); err != nil {
		panic("cannot open audio")
	}
	defer mix.CloseAudio()

	resources := loadResources(renderer)

	running := true
	game := m.NewGame()

	for running {
		command := ""
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.State == sdl.PRESSED && t.Repeat == 0 {
					keyCode := t.Keysym.Sym
					switch keyCode {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_LEFT:
						command = "left"
					case sdl.K_RIGHT:
						command = "right"
					case sdl.K_LSHIFT:
						command = "shoot"
						chunk := resources.chunks["hit.wav"]
						chunk.Play(-1, 0)
					case sdl.K_RSHIFT:
						command = "shoot"
					}
				}
			}
		}
		game.Update(command)
		render(renderer, window, game, resources)
		time.Sleep((1000 / FPS) * time.Millisecond)
	}
}

func loadResources(renderer *sdl.Renderer) *Resources {
	var resources Resources
	var err error
	var image *sdl.Surface

	var imagePaths []string = []string{
		"back.bmp",
		"down.bmp",
		"dust.bmp",
		"left.bmp",
		"player.bmp",
		"right.bmp",
		"target.bmp",
		"title.bmp",
		"up.bmp",
		"wall.bmp",
	}

	resources.textures = make(map[string]*sdl.Texture)
	for _, path := range imagePaths {
		fullPath := "resources/image/" + path
		if image, err = sdl.LoadBMP(fullPath); err != nil {
			panic("cannot load image: " + path)
		}
		texture, err := renderer.CreateTextureFromSurface(image)
		if err != nil {
			panic("cannot convert to texture: " + path)
		}
		resources.textures[path] = texture
	}

	// Load WAV file with short duration as *mix.Chunk
	var soundPaths []string = []string{
		"crash.wav",
		"hit.wav",
	}

	resources.chunks = make(map[string]*mix.Chunk)
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

	if game.IsOver {
		renderer.SetDrawColor(0, 0, 0, 128)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: SCREEEN_WIDTH, H: SCREEN_HEIGHT})
	}

	renderTexture(renderer, resources, "player.bmp", 100, 150)
	renderTexture(renderer, resources, "title.bmp", 100, 350)
	renderer.Present()
}

func renderTexture(renderer *sdl.Renderer, resources *Resources, textureKey string, x int32, y int32) {
	texture := resources.textures[textureKey]
	_, _, w, h, err := texture.Query()
	if err != nil {
		panic("cannot query texture: " + textureKey)
	}
	renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
}
