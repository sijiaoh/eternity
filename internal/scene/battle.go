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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

type BattleScene struct {
	clock     *component.Clock
	camera    *component.Camera
	floorTile *ebiten.Image
	tileSize  int

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
	PlayerSpriteSheet   *ebiten.Image
	PlayerSpriteColumns int
	PlayerFrameWidth    int
	PlayerFrameHeight   int
	PlayerAnimFPS       float64
	GoblinSpriteSheet   *ebiten.Image
	GoblinSpriteColumns int
	GoblinFrameWidth    int
	GoblinFrameHeight   int
	GoblinAnimFPS       float64
	DialogueFont        text.Face
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

	playerX := config.PixelsToUnits(config.ScreenWidth / 2)
	playerY := config.PixelsToUnits(config.ScreenHeight / 2)

	camera := component.NewCamera(playerX, playerY, 0.1)

	world := ecs.NewWorld(64)

	// Shared storages
	positions := ecs.NewStorage[component.Position](64)
	velocities := ecs.NewStorage[component.Velocity](64)
	facings := ecs.NewStorage[component.Facing](64)
	animations := ecs.NewStorage[component.Animation](64)
	spriteSheets := ecs.NewStorage[component.SpriteSheet](64)

	// Player-specific storages
	inputs := ecs.NewStorage[component.InputControlled](8)
	cameraTargets := ecs.NewStorage[component.CameraTarget](4)

	// Goblin-specific storages
	aiFollows := ecs.NewStorage[component.AIFollow](32)

	playerComponents := &entity.PlayerComponents{
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

	player := entity.CreatePlayer(world, playerComponents, entity.PlayerFactoryConfig{
		X:           playerX,
		Y:           playerY,
		Speed:       5.0, // units per second
		SizeInUnits: 1.0,
		SpriteSheet: cfg.PlayerSpriteSheet,
		FrameWidth:  cfg.PlayerFrameWidth,
		FrameHeight: cfg.PlayerFrameHeight,
		Columns:     cfg.PlayerSpriteColumns,
		AnimFPS:     cfg.PlayerAnimFPS,
	})

	// Create goblin if sprite is provided
	if cfg.GoblinSpriteSheet != nil {
		entity.CreateGoblin(world, goblinComponents, entity.GoblinFactoryConfig{
			X:           playerX + 3.0,
			Y:           playerY + 3.0,
			Speed:       3.0, // slower than player
			// Body fills ~half the 64px frame; 2.0 keeps it close to player size.
			SizeInUnits: 2.0,
			Target:      player,
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

	return &BattleScene{
		clock:                component.NewClock(),
		camera:               camera,
		floorTile:            tile,
		tileSize:             tile.Bounds().Dx(),
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

	// Demo trigger: press E to start a sample dialogue.
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		s.StartDialogue(sampleDialogue())
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

// sampleDialogue is the demo script for manual verification: a portrait line, a portrait-less
// narration line, and a long line that exercises CJK wrapping.
func sampleDialogue() []component.DialogueLine {
	return []component.DialogueLine{
		{Speaker: "法师", Portrait: "mage", Text: "前方就是哥布林的巢穴了，务必小心。它们成群结队，而且对闯入者毫不留情。"},
		{Speaker: "", Portrait: "", Text: "（四周一片寂静，只有风声掠过枯草，远处似乎有什么东西在移动。）"},
		{Speaker: "法师", Portrait: "mage", Text: "准备好了吗？按下空格、回车或鼠标左键即可继续，我们这就出发。"},
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
