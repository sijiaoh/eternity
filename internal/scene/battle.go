package scene

import (
	"image"
	_ "image/png"
	"io/fs"
	"math"

	"ebiten-agent-example/internal/component"
	"ebiten-agent-example/internal/config"
	"ebiten-agent-example/internal/entity"

	"github.com/hajimehoshi/ebiten/v2"
)

type BattleScene struct {
	clock     *component.Clock
	floorTile *ebiten.Image
	tileSize  int
	offsetX   float64
	offsetY   float64
	scrollX   float64
	scrollY   float64
	player    *entity.Player
}

type BattleSceneConfig struct {
	FloorImagePath      string
	PlayerSpriteSheet   *ebiten.Image
	PlayerSpriteColumns int
	PlayerFrameWidth    int
	PlayerFrameHeight   int
	PlayerAnimFPS       float64
}

func loadImage(fsys fs.FS, path string) (*ebiten.Image, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func NewBattleScene(fsys fs.FS, cfg BattleSceneConfig) (*BattleScene, error) {
	tile, err := loadImage(fsys, cfg.FloorImagePath)
	if err != nil {
		return nil, err
	}

	playerCfg := entity.PlayerConfig{
		SpriteSheet: cfg.PlayerSpriteSheet,
		FrameWidth:  cfg.PlayerFrameWidth,
		FrameHeight: cfg.PlayerFrameHeight,
		Columns:     cfg.PlayerSpriteColumns,
		AnimFPS:     cfg.PlayerAnimFPS,
	}

	return &BattleScene{
		clock:     component.NewClock(),
		floorTile: tile,
		tileSize:  tile.Bounds().Dx(),
		scrollX:   50, // pixels per second (visual offset, not world position)
		scrollY:   30, // pixels per second (visual offset, not world position)
		player:    entity.NewPlayer(config.PixelsToUnits(config.ScreenWidth/2), config.PixelsToUnits(config.ScreenHeight/2), playerCfg),
	}, nil
}

func (s *BattleScene) Update() error {
	s.clock.Update(1.0 / float64(ebiten.TPS()))
	dt := s.clock.DeltaTime()
	s.offsetX += s.scrollX * dt
	s.offsetY += s.scrollY * dt
	s.offsetX = math.Mod(s.offsetX, float64(s.tileSize))
	s.offsetY = math.Mod(s.offsetY, float64(s.tileSize))
	s.player.Update(dt)
	return nil
}

func (s *BattleScene) SetTimeScale(scale float64) {
	s.clock.SetScale(scale)
}

func (s *BattleScene) Pause() {
	s.clock.Pause()
}

func (s *BattleScene) Resume() {
	s.clock.Resume()
}

func (s *BattleScene) IsPaused() bool {
	return s.clock.IsPaused()
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
	s.player.Draw(screen)
}
