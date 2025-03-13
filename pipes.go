package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
	"time"
)

type pipes struct {
	mu sync.RWMutex

	pipes   []*pipe
	texture *sdl.Texture
	speed   int32
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(r, "res/img/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("Could not load pipe texture: %v", err)
	}

	ps := &pipes{
		texture: texture,
		speed:   2,
	}

	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	return ps, nil
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var remaining []*pipe
	for _, p := range ps.pipes {
		p.update(ps.speed)

		if p.x > -p.w {
			remaining = append(remaining, p)
		}
	}

	ps.pipes = remaining
}

func (ps *pipes) touch(b *bird) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, p := range ps.pipes {
		p.touch(b)
	}
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}

	return nil
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}
