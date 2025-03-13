package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
}

func newScene(r *sdl.Renderer) (*scene, error) {
	texture, err := img.LoadTexture(r, "res/img/background.jpg")
	if err != nil {
		return nil, fmt.Errorf("Could not load background: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, fmt.Errorf("Could not create bird: %v", err)
	}

	pipes, err := newPipes(r)
	if err != nil {
		return nil, fmt.Errorf("Could not create pipe: %v", err)
	}

	return &scene{
		bg:    texture,
		bird:  bird,
		pipes: pipes,
	}, nil
}

func (s *scene) run(events chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)

		ticks := time.Tick(10 * time.Millisecond)

		for {
			select {
			case event := <-events:
				if done := s.handleEvent(event); done {
					return
				}
			case <-ticks:
				s.update()

				if s.bird.isDead() {
					fmt.Println("Bird is dead")
					drawTitle(r, "Game Over")
					time.Sleep(2 * time.Second)
					s.restart()
				}

				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background texture: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return fmt.Errorf("Could not paint bird texture: %v", err)
	}

	if err := s.pipes.paint(r); err != nil {
		return fmt.Errorf("Could not paint pipe texture: %v", err)
	}

	r.Present()

	return nil
}

func (s *scene) handleEvent(event sdl.Event) bool {
	// fmt.Printf("Event: %T\n", event)

	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		key := e.Keysym.Sym
		t := e.Type

		if key == sdl.K_SPACE && t == sdl.KEYDOWN {
			s.bird.jump()
		}
	}

	return false
}

func (s *scene) restart() {
	fmt.Println("Restarting")
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}
