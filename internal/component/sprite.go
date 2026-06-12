//go:build !test

package component

import "github.com/hajimehoshi/ebiten/v2"

// Anchor defines the sprite's origin point relative to its bounds.
type Anchor struct {
	X, Y float64 // 0-1 range: (0,0) = top-left, (0.5,0.5) = center, (1,1) = bottom-right
}

func AnchorTopLeft() Anchor     { return Anchor{X: 0, Y: 0} }
func AnchorCenter() Anchor      { return Anchor{X: 0.5, Y: 0.5} }
func AnchorBottomRight() Anchor { return Anchor{X: 1, Y: 1} }

// Sprite holds a static image for rendering.
type Sprite struct {
	Image  *ebiten.Image
	Anchor Anchor
}

func NewSprite(image *ebiten.Image) *Sprite {
	return &Sprite{
		Image:  image,
		Anchor: AnchorCenter(),
	}
}

func NewSpriteWithAnchor(image *ebiten.Image, anchor Anchor) *Sprite {
	return &Sprite{
		Image:  image,
		Anchor: anchor,
	}
}
