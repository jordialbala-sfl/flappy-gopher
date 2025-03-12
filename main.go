package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("Could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Could not create window: %v", err)
	}
	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		return fmt.Errorf("Could not draw title: %v", err)
	}

	sdl.PumpEvents()
	sdl.PumpEvents()

	fmt.Println("Print title")
	time.Sleep(5 * time.Second)

	if err := drawBackground(r); err != nil {
		return fmt.Errorf("Could not draw background: %v", err)
	}
	time.Sleep(5 * time.Second)

	return nil
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	font, err := ttf.OpenFont("res/fonts/flappy-font.ttf", 20)
	if err != nil {
		return fmt.Errorf("Could not open font: %v", err)
	}
	defer font.Close()

	color := sdl.Color{R: 255, G: 0, B: 0, A: 255}
	surface, err := font.RenderUTF8Solid("Flappy Gopher", color)
	if err != nil {
		return fmt.Errorf("Could not render text: %v", err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer texture.Destroy()

	if err := r.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	r.Present()

	return nil
}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()

	texture, err := img.LoadTexture(r, "res/img/background.jpg")
	if err != nil {
		return fmt.Errorf("Could not load background: %v", err)
	}
	defer texture.Destroy()

	if err := r.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	r.Present()

	return nil
}
