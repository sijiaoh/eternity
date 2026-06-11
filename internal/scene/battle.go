package scene

import (
	"image"
	_ "image/png"
	"io/fs"
	"math"

	"eternity/internal/component"
	"eternity/internal/config"
	"eternity/internal/entity"

	"github.com/hajimehoshi/ebiten/v2"
)

type BattleScene struct {
	clock     *component.Clock
	camera    *component.Camera
	floorTile *ebiten.Image
	tileSize  int
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

	playerX := config.PixelsToUnits(config.ScreenWidth / 2)
	playerY := config.PixelsToUnits(config.ScreenHeight / 2)

	return &BattleScene{
		clock:     component.NewClock(),
		camera:    component.NewCamera(playerX, playerY, 0.1),
		floorTile: tile,
		tileSize:  tile.Bounds().Dx(),
		player:    entity.NewPlayer(playerX, playerY, playerCfg),
	}, nil
}

func (s *BattleScene) Update() error {
	s.clock.Update(1.0 / float64(ebiten.TPS()))
	dt := s.clock.DeltaTime()
	s.player.Update(dt)
	s.camera.Update(s.player.Position.X, s.player.Position.Y, dt)
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
	offsetX, offsetY := s.camera.GetOffset()
	ts := float64(s.tileSize)

	// Calculate tile offset for seamless scrolling (always positive)
	tileOffsetX := offsetX - math.Floor(offsetX/ts)*ts
	tileOffsetY := offsetY - math.Floor(offsetY/ts)*ts

	tilesX := int(math.Ceil(float64(config.ScreenWidth)/ts)) + 2
	tilesY := int(math.Ceil(float64(config.ScreenHeight)/ts)) + 2

	for y := 0; y < tilesY; y++ {
		for x := 0; x < tilesX; x++ {
			op := &ebiten.DrawImageOptions{}
			drawX := float64(x)*ts - ts - tileOffsetX
			drawY := float64(y)*ts - ts - tileOffsetY
			op.GeoM.Translate(drawX, drawY)
			screen.DrawImage(s.floorTile, op)
		}
	}

	// Draw player relative to camera
	screenX, screenY := s.camera.WorldToScreen(s.player.Position.X, s.player.Position.Y)
	s.player.DrawAt(screen, screenX, screenY)
}
