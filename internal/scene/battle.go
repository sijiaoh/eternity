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
	spriteFlipSystem     *system.SpriteFlipSystem
	animationStateSystem *system.AnimationStateSystem
	animationSystem      *system.AnimationSystem
	cameraSystem         *system.CameraSystem

	// Draw systems
	renderSystem *system.SpriteSheetRenderSystem
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
	spriteFlipSystem := system.NewSpriteFlipSystem(facings, spriteSheets)
	animationStateSystem := system.NewAnimationStateSystem(animations, facings)
	animationSystem := system.NewAnimationSystem(animations)
	cameraSystem := system.NewCameraSystem(cameraTargets, positions, camera)
	renderSystem := system.NewSpriteSheetRenderSystem(positions, animations, spriteSheets, camera)

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
		spriteFlipSystem:     spriteFlipSystem,
		animationStateSystem: animationStateSystem,
		animationSystem:      animationSystem,
		cameraSystem:         cameraSystem,
		renderSystem:         renderSystem,
	}, nil
}

func (s *BattleScene) Update() error {
	s.clock.Update(1.0 / float64(ebiten.TPS()))
	dt := s.clock.DeltaTime()

	// Update systems in order
	s.inputSystem.Update(s.world, dt)
	s.aiFollowSystem.Update(s.world, dt)
	s.movementSystem.Update(s.world, dt)
	s.facingSystem.Update(s.world, dt)
	s.spriteFlipSystem.Update(s.world, dt)
	s.animationStateSystem.Update(s.world, dt)
	s.animationSystem.Update(s.world, dt)
	s.cameraSystem.Update(s.world, dt)

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
}
