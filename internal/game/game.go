//go:build !test

package game

import (
	"image"
	_ "image/png"

	"eternity/assets"
	"eternity/internal/config"
	"eternity/internal/scene"

	"github.com/hajimehoshi/ebiten/v2"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	// Mage: 6 columns x 4 rows, rows = down/left/up/right (frame 0 = idle).
	mageFrameWidth    = 48
	mageFrameHeight   = 48
	mageSpriteColumns = 6
	mageAnimFPS       = 8

	// LPC goblin: 11 columns x 5 rows, rows = down/left/up/right/death.
	goblinFrameWidth    = 64
	goblinFrameHeight   = 64
	goblinSpriteColumns = 11
	goblinAnimFPS       = 8
)

type Game struct {
	sceneManager *scene.Manager
}

func New() (*Game, error) {
	playerSprite, err := loadImage("images/characters/mage/sprite.png")
	if err != nil {
		return nil, err
	}

	goblinSprite, err := loadImage("images/characters/goblin/sprite.png")
	if err != nil {
		return nil, err
	}

	dialogueFont, err := loadFont("fonts/WenQuanYiMicroHei.ttf", config.DialogueFontSize)
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
		GoblinSpriteSheet:   goblinSprite,
		GoblinFrameWidth:    goblinFrameWidth,
		GoblinFrameHeight:   goblinFrameHeight,
		GoblinSpriteColumns: goblinSpriteColumns,
		GoblinAnimFPS:       goblinAnimFPS,
		DialogueFont:        dialogueFont,
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

// loadFont parses an embedded TTF into a text.Face at the given pixel size.
func loadFont(path string, size float64) (text.Face, error) {
	data, err := assets.Fonts.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tt, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	return text.NewGoXFace(face), nil
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
