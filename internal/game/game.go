//go:build !test

package game

import (
	"image"
	_ "image/png"
	"log"

	"eternity/assets"
	"eternity/internal/config"
	"eternity/internal/i18n"
	"eternity/internal/scenario"
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

// deps holds the shared, ebiten-decoded resources every scene builder draws from. They are
// loaded once up front so building any start scene (or transitioning between scenes) reuses
// the same sprites and fonts rather than reloading assets per scene.
type deps struct {
	bundle       *i18n.Bundle
	battleConfig scene.BattleSceneConfig
	titleFont    text.Face
	promptFont   text.Face
}

// New assembles the game and launches into the scenario's scene. sc.Scene names a registered
// scene (see the scene constants); an empty name keeps the default player path. sc.Locale, when
// set, picks the language before any scene is built so any scene reproduces in that language.
// sc.Battle is the debug situation applied when the battle scene is the target. The build is lazy
// — only the selected start scene is constructed — and scene transitions stay decoupled via the
// Manager and injected callbacks (see docs/ebiten-code-rules.md「场景切换」).
func New(bundle *i18n.Bundle, sc scenario.Scenario) (*Game, error) {
	d, err := loadDeps(bundle)
	if err != nil {
		return nil, err
	}
	// Carry the debug situation into the battle config; the builder picks it up when battle is
	// the chosen scene. Other scenes have no situation to vary, so it's harmlessly ignored.
	d.battleConfig.Situation = sc.Battle

	manager := scene.NewManager(nil)

	// Scene registry: any registered scene can be the launch target. The builders capture the
	// shared deps and the manager so the title→battle transition can swap scenes by itself.
	builders := map[string]func() (scene.Scene, error){
		sceneTitle:  func() (scene.Scene, error) { return buildTitleScene(d, manager), nil },
		sceneBattle: func() (scene.Scene, error) { return buildBattleScene(d) },
	}

	names := make([]string, 0, len(builders))
	for name := range builders {
		names = append(names, name)
	}

	name, err := resolveScene(sc.Scene, DefaultScene, names)
	if err != nil {
		return nil, err
	}
	if err := validateSituation(name, sc.Battle); err != nil {
		return nil, err
	}

	// Apply the locale before building the scene, so its text resolves in the requested language;
	// an invalid locale fails here, before any text is read.
	if err := applyLocale(d.bundle, sc.Locale); err != nil {
		return nil, err
	}

	start, err := builders[name]()
	if err != nil {
		return nil, err
	}
	manager.SetScene(start)

	return &Game{sceneManager: manager}, nil
}

// loadDeps decodes the sprites and fonts shared across scenes and packs the battle scene's
// config, so the per-scene builders only assemble — they don't touch the filesystem.
func loadDeps(bundle *i18n.Bundle) (*deps, error) {
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

	titleFont, err := loadFont("fonts/WenQuanYiMicroHei.ttf", config.TitleFontSize)
	if err != nil {
		return nil, err
	}

	promptFont, err := loadFont("fonts/WenQuanYiMicroHei.ttf", config.TitlePromptFontSize)
	if err != nil {
		return nil, err
	}

	return &deps{
		bundle:     bundle,
		titleFont:  titleFont,
		promptFont: promptFont,
		battleConfig: scene.BattleSceneConfig{
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
			Bundle:              bundle,
		},
	}, nil
}

func buildBattleScene(d *deps) (scene.Scene, error) {
	return scene.NewBattleScene(assets.Images, d.battleConfig)
}

// buildTitleScene constructs the title screen. Pressing Enter builds the battle scene on
// demand and swaps it in via the manager; the title stays decoupled from the battle type.
func buildTitleScene(d *deps, manager *scene.Manager) scene.Scene {
	return scene.NewTitleScene(d.bundle, d.titleFont, d.promptFont, func() {
		battle, err := buildBattleScene(d)
		if err != nil {
			// Battle only loads embedded, build-validated assets, so a failure here means a
			// broken build. Fail loud rather than leaving the player on a dead title screen —
			// this mirrors the startup log.Fatal back when battle was pre-built in New.
			log.Fatalf("launcher: failed to start battle scene: %v", err)
		}
		manager.SetScene(battle)
	})
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
