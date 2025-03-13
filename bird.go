package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"sync"
)

const gravity = 0.1
const jumpSpeed = 7

type bird struct {
	mu sync.RWMutex

	time   int
	frames []*sdl.Texture

	y, speed float64

	dead bool
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var frames []*sdl.Texture

	for i := 1; i <= 4; i++ {
		bird, err := img.LoadTexture(r, fmt.Sprintf("res/img/frame-%d.png", i))
		if err != nil {
			return nil, fmt.Errorf("Could not load bird frame %d: %v", i, err)
		}

		frames = append(frames, bird)
	}

	return &bird{frames: frames, y: 300}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	//fmt.Printf("%v\t%v\n", int32(b.y), int32(b.speed))

	b.time++

	b.y -= b.speed
	b.speed += gravity

	if b.y <= 0 {
		b.dead = true
	}
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	rect := &sdl.Rect{X: 10, Y: (600 - int32(b.y)) - 43/2, W: 50, H: 43}
	frameIndex := (b.time / 10) % len(b.frames)

	if err := r.Copy(b.frames[frameIndex], nil, rect); err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	return nil
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = math.Abs(b.speed)
	b.speed -= jumpSpeed
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300 - 43/2
	b.speed = 0
	b.dead = false
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, frame := range b.frames {
		frame.Destroy()
	}
}
