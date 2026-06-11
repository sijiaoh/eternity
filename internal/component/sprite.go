package component

import "github.com/hajimehoshi/ebiten/v2"

type Anchor struct {
	X, Y float64 // 0-1 range: (0,0) = top-left, (0.5,0.5) = center, (1,1) = bottom-right
}

func AnchorTopLeft() Anchor     { return Anchor{X: 0, Y: 0} }
func AnchorCenter() Anchor      { return Anchor{X: 0.5, Y: 0.5} }
func AnchorBottomRight() Anchor { return Anchor{X: 1, Y: 1} }

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

func (s *Sprite) Draw(screen *ebiten.Image, x, y float64) {
	s.DrawFrame(screen, s.Image, x, y)
}

func (s *Sprite) CalcDrawPosition(x, y float64, imageWidth, imageHeight int) (drawX, drawY float64) {
	offsetX := float64(imageWidth) * s.Anchor.X
	offsetY := float64(imageHeight) * s.Anchor.Y
	return x - offsetX, y - offsetY
}

// DrawFrame renders a frame image, allowing Animation component to pass the current frame.
func (s *Sprite) DrawFrame(screen *ebiten.Image, frame *ebiten.Image, x, y float64) {
	if frame == nil {
		return
	}

	w, h := frame.Bounds().Dx(), frame.Bounds().Dy()
	drawX, drawY := s.CalcDrawPosition(x, y, w, h)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(drawX, drawY)
	screen.DrawImage(frame, opts)
}
