package game

import (
	"image"
	_ "image/png"

	"eternity/assets"
	"eternity/internal/config"
	"eternity/internal/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

// Mage sprite sheet layout: 6 columns x 5 rows
// Row order: down, left, up, right (+ extra row)
// Per row: 6 walk frames (frame 0 used as idle)
const (
	mageFrameWidth    = 48
	mageFrameHeight   = 48
	mageSpriteColumns = 6
	mageAnimFPS       = 8
)

type Game struct {
	sceneManager *scene.Manager
}

func New() (*Game, error) {
	playerSprite, err := loadImage("images/characters/mage/sprite.png")
	if err != nil {
		return nil, err
	}

	cfg := scene.BattleSceneConfig{
		FloorImagePath:      "images/battle/floor/stone.png",
		PlayerSpriteSheet:   playerSprite,
		PlayerFrameWidth:    mageFrameWidth,
		PlayerFrameHeight:   mageFrameHeight,
		PlayerSpriteColumns: mageSpriteColumns,
		PlayerAnimFPS:       mageAnimFPS,
	}

	battleScene, err := scene.NewBattleScene(assets.Images, cfg)
	if err != nil {
		return nil, err
	}

	return &Game{
		sceneManager: scene.NewManager(battleScene),
	}, nil
}

func loadImage(path string) (*ebiten.Image, error) {
	f, err := assets.Images.Open(path)
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

func (g *Game) Update() error {
	return g.sceneManager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
