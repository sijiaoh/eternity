//go:build !test

package ecs

import "github.com/hajimehoshi/ebiten/v2"

// System processes entities each frame.
type System interface {
	Update(w *World, dt float64)
}

// DrawSystem renders entities to screen.
type DrawSystem interface {
	Draw(w *World, screen *ebiten.Image)
}
