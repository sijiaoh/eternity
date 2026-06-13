//go:build !test

package scene

import (
	"image"
	_ "image/png"
	"io/fs"
	"math"

	"eternity/internal/component"
	"eternity/internal/config"
	"eternity/internal/ecs"
	"eternity/internal/ecs/system"
	"eternity/internal/entity"
	"eternity/internal/i18n"
	"eternity/internal/scenario"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

type BattleScene struct {
	clock     *component.Clock
	camera    *component.Camera
	floorTile *ebiten.Image
	tileSize  int
	bundle    *i18n.Bundle

	// ECS
	world *ecs.World

	// Systems (update order matters)
	inputSystem          *system.InputSystem
	aiFollowSystem       *system.AIFollowSystem
	movementSystem       *system.MovementSystem
	facingSystem         *system.FacingSystem
	animationStateSystem *system.AnimationStateSystem
	animationSystem      *system.AnimationSystem
	cameraSystem         *system.CameraSystem

	// Dialogue overlay (scene-owned singleton, gates the game systems while active)
	dialogue       *component.Dialogue
	dialogueSystem *system.DialogueSystem

	// Draw systems
	renderSystem         *system.SpriteSheetRenderSystem
	dialogueRenderSystem *system.DialogueRenderSystem
}

type BattleSceneConfig struct {
	FloorImagePath      string
	MageSpriteSheet     *ebiten.Image
	MageSpriteColumns   int
	MageFrameWidth      int
	MageFrameHeight     int
	MageAnimFPS         float64
	GoblinSpriteSheet   *ebiten.Image
	GoblinSpriteColumns int
	GoblinFrameWidth    int
	GoblinFrameHeight   int
	GoblinAnimFPS       float64
	DialogueFont        text.Face
	Bundle              *i18n.Bundle

	// Situation is the optional debug/test scenario applied at construction. Its zero value
	// reproduces normal play, so the player path is unaffected.
	Situation scenario.Battle
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

	mageX, mageY := resolveMageStart(cfg.Situation)

	camera := component.NewCamera(mageX, mageY, 0.1)

	world := ecs.NewWorld(64)

	// Shared storages
	positions := ecs.NewStorage[component.Position](64)
	velocities := ecs.NewStorage[component.Velocity](64)
	facings := ecs.NewStorage[component.Facing](64)
	animations := ecs.NewStorage[component.Animation](64)
	spriteSheets := ecs.NewStorage[component.SpriteSheet](64)

	// Mage-specific storages
	inputs := ecs.NewStorage[component.InputControlled](8)
	cameraTargets := ecs.NewStorage[component.CameraTarget](4)

	// Goblin-specific storages
	aiFollows := ecs.NewStorage[component.AIFollow](32)

	mageComponents := &entity.MageComponents{
		Positions:     positions,
		Velocities:    velocities,
		Inputs:        inputs,
		Facings:       facings,
		Animations:    animations,
		SpriteSheets:  spriteSheets,
		CameraTargets: cameraTargets,
	}

	goblinComponents := &entity.GoblinComponents{
		Positions:    positions,
		Velocities:   velocities,
		AIFollows:    aiFollows,
		Facings:      facings,
		Animations:   animations,
		SpriteSheets: spriteSheets,
	}

	mage := entity.CreateMage(world, mageComponents, entity.MageFactoryConfig{
		X:           mageX,
		Y:           mageY,
		Speed:       5.0, // units per second
		SizeInUnits: 1.0,
		SpriteSheet: cfg.MageSpriteSheet,
		FrameWidth:  cfg.MageFrameWidth,
		FrameHeight: cfg.MageFrameHeight,
		Columns:     cfg.MageSpriteColumns,
		AnimFPS:     cfg.MageAnimFPS,
	})

	// Create goblin if the sprite is present and the scenario doesn't suppress it.
	if cfg.GoblinSpriteSheet != nil && spawnGoblin(cfg.Situation) {
		goblinX, goblinY := resolveGoblinStart(cfg.Situation, mageX, mageY)
		entity.CreateGoblin(world, goblinComponents, entity.GoblinFactoryConfig{
			X:     goblinX,
			Y:     goblinY,
			Speed: 3.0, // slower than the mage
			// Body fills ~half the 64px frame; 2.0 keeps it close to mage size.
			SizeInUnits: 2.0,
			Target:      mage,
			SpriteSheet: cfg.GoblinSpriteSheet,
			FrameWidth:  cfg.GoblinFrameWidth,
			FrameHeight: cfg.GoblinFrameHeight,
			Columns:     cfg.GoblinSpriteColumns,
			AnimFPS:     cfg.GoblinAnimFPS,
		})
	}

	inputSystem := system.NewInputSystem(inputs, velocities)
	aiFollowSystem := system.NewAIFollowSystem(aiFollows, positions, velocities)
	movementSystem := system.NewMovementSystem(positions, velocities)
	facingSystem := system.NewFacingSystem(facings, velocities)
	animationStateSystem := system.NewAnimationStateSystem(animations, facings)
	animationSystem := system.NewAnimationSystem(animations)
	cameraSystem := system.NewCameraSystem(cameraTargets, positions, camera)
	renderSystem := system.NewSpriteSheetRenderSystem(positions, animations, spriteSheets, camera)

	portraits, err := loadPortraits(fsys)
	if err != nil {
		return nil, err
	}
	dialogue := &component.Dialogue{}
	dialogueSystem := system.NewDialogueSystem(dialogue)
	dialogueRenderSystem := system.NewDialogueRenderSystem(dialogue, portraits, cfg.DialogueFont)

	// Start the scenario's initial dialogue (nil when it doesn't ask, leaving the dialogue
	// inactive). Active dialogue gates the game systems on the first Update (see Update), so a
	// start-in-dialogue scenario opens paused on the first line.
	dialogue.Start(initialDialogue(cfg.Situation, cfg.Bundle))

	clock := component.NewClock()
	clock.SetScale(resolveTimeScale(cfg.Situation))

	return &BattleScene{
		clock:                clock,
		camera:               camera,
		floorTile:            tile,
		tileSize:             tile.Bounds().Dx(),
		bundle:               cfg.Bundle,
		world:                world,
		inputSystem:          inputSystem,
		aiFollowSystem:       aiFollowSystem,
		movementSystem:       movementSystem,
		facingSystem:         facingSystem,
		animationStateSystem: animationStateSystem,
		animationSystem:      animationSystem,
		cameraSystem:         cameraSystem,
		dialogue:             dialogue,
		dialogueSystem:       dialogueSystem,
		renderSystem:         renderSystem,
		dialogueRenderSystem: dialogueRenderSystem,
	}, nil
}

// loadPortraits builds the portrait-key registry from the embedded character images.
func loadPortraits(fsys fs.FS) (map[string]*ebiten.Image, error) {
	keys := []string{"mage", "wolf", "panther"}
	portraits := make(map[string]*ebiten.Image, len(keys))
	for _, key := range keys {
		img, err := loadImage(fsys, "images/characters/"+key+"/portrait.png")
		if err != nil {
			return nil, err
		}
		portraits[key] = img
	}
	return portraits, nil
}

func (s *BattleScene) Update() error {
	s.clock.Update(1.0 / float64(ebiten.TPS()))
	dt := s.clock.DeltaTime()

	// While dialogue is active, gate the game systems off and only advance dialogue,
	// freezing the world so the player can read without characters moving. We gate here
	// rather than via clock.Pause: a zeroed dt only stops dt-scaled movement, but the
	// input-driven systems (velocity/facing/animation) ignore dt, so held keys would still
	// turn the character and play its walk animation in place.
	if s.dialogue.Active {
		s.dialogueSystem.Update(s.world, dt)
		return nil
	}

	// Demo trigger: press L to cycle the locale, so the next dialogue shows another language.
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		s.cycleLocale()
	}

	// Demo trigger: press E to start a sample dialogue.
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		s.StartDialogue(sampleDialogue(s.bundle))
		return nil
	}

	// Update systems in order
	s.inputSystem.Update(s.world, dt)
	s.aiFollowSystem.Update(s.world, dt)
	s.movementSystem.Update(s.world, dt)
	s.facingSystem.Update(s.world, dt)
	s.animationStateSystem.Update(s.world, dt)
	s.animationSystem.Update(s.world, dt)
	s.cameraSystem.Update(s.world, dt)

	return nil
}

// StartDialogue opens a dialogue from the given lines. The trigger source (key, script,
// trigger zone) is the caller's concern and stays decoupled from the dialogue system.
func (s *BattleScene) StartDialogue(lines []component.DialogueLine) {
	s.dialogue.Start(lines)
}

// cycleLocale advances the bundle to the next available locale, wrapping around. Text is
// resolved when a dialogue starts, so switching here changes the language of the next dialogue.
func (s *BattleScene) cycleLocale() {
	locales := s.bundle.Locales()
	current := s.bundle.Locale()
	for i, locale := range locales {
		if locale == current {
			s.bundle.SetLocale(locales[(i+1)%len(locales)])
			return
		}
	}
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
	offsetX, offsetY := system.CameraGetOffset(s.camera)
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

	// Draw entities via ECS render system
	s.renderSystem.Draw(s.world, screen)

	// Dialogue overlay sits on top of the world.
	s.dialogueRenderSystem.Draw(s.world, screen)
}
