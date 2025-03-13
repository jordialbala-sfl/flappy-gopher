package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
)

type pipe struct {
	mu sync.RWMutex

	x, h, w  int32
	inverted bool
}

func newPipe() *pipe {
	return &pipe{
		x:        800,
		h:        200 + rand.Int31n(200),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}
}

func (p *pipe) update(speed int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x -= speed
}

func (p *pipe) touch(b *bird) {
	b.touch(p)
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := r.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("Could not copy pipe: %v", err)
	}

	return nil
}
