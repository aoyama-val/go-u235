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

	initMixer()
	defer mix.Quit()
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
				if t.State == sdl.PRESSED {
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

	if game.IsOver {
		renderer.SetDrawColor(0, 0, 0, 128)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: SCREEEN_WIDTH, H: SCREEN_HEIGHT})
	}

	renderTexture(renderer, resources, "player.bmp", game.Player.X, game.Player.Y)
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

	renderer.Present()
}

func renderTexture(renderer *sdl.Renderer, resources *Resources, textureKey string, x int, y int) {
	texture := resources.textures[textureKey]
	renderer.Copy(texture.texture, nil, &sdl.Rect{X: int32(CELL_SIZE_PX * x), Y: int32(CELL_SIZE_PX * y), W: texture.w, H: texture.h})
}
