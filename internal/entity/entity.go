package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	Update(deltaTime float64)
	Draw(screen *ebiten.Image)
}
