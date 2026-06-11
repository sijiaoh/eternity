package scene

import (
	"image"
	_ "image/png"
	"math"
	"os"

	"ebiten-agent-example/internal/config"

	"github.com/hajimehoshi/ebiten/v2"
)

// BattleScene implements a battle scene with infinite scrolling floor.
type BattleScene struct {
	floorTile *ebiten.Image
	tileSize  int
	offsetX   float64
	offsetY   float64
	scrollX   float64
	scrollY   float64
}

// NewBattleScene creates a new battle scene with the given floor tile image.
func NewBattleScene(floorImagePath string) (*BattleScene, error) {
	f, err := os.Open(floorImagePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	tile := ebiten.NewImageFromImage(img)
	bounds := img.Bounds()

	return &BattleScene{
		floorTile: tile,
		tileSize:  bounds.Dx(),
		scrollX:   50, // pixels per second
		scrollY:   30, // pixels per second
	}, nil
}

func (s *BattleScene) Update() error {
	tps := float64(ebiten.TPS())
	s.offsetX += s.scrollX / tps
	s.offsetY += s.scrollY / tps
	s.offsetX = math.Mod(s.offsetX, float64(s.tileSize))
	s.offsetY = math.Mod(s.offsetY, float64(s.tileSize))
	return nil
}

func (s *BattleScene) Draw(screen *ebiten.Image) {
	ts := float64(s.tileSize)
	tilesX := int(math.Ceil(float64(config.ScreenWidth)/ts)) + 2
	tilesY := int(math.Ceil(float64(config.ScreenHeight)/ts)) + 2

	for y := 0; y < tilesY; y++ {
		for x := 0; x < tilesX; x++ {
			op := &ebiten.DrawImageOptions{}
			drawX := float64(x)*ts - ts + s.offsetX
			drawY := float64(y)*ts - ts + s.offsetY
			op.GeoM.Translate(drawX, drawY)
			screen.DrawImage(s.floorTile, op)
		}
	}
}
