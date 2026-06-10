package render

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Image   *ebiten.Image
	OffsetX float64
	OffsetY float64
}
