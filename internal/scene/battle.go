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
	world      *ecs.World
	components *entity.PlayerComponents

	// Systems (update order matters)
	inputSystem          *system.InputSystem
	movementSystem       *system.MovementSystem
	facingSystem         *system.FacingSystem
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

	// Create ECS world and storages
	world := ecs.NewWorld(64)
	components := &entity.PlayerComponents{
		Positions:     ecs.NewStorage[component.Position](64),
		Velocities:    ecs.NewStorage[component.Velocity](64),
		Inputs:        ecs.NewStorage[component.InputControlled](8),
		Facings:       ecs.NewStorage[component.Facing](64),
		Animations:    ecs.NewStorage[component.Animation](64),
		SpriteSheets:  ecs.NewStorage[component.SpriteSheet](64),
		CameraTargets: ecs.NewStorage[component.CameraTarget](4),
	}

	// Create player entity
	entity.CreatePlayer(world, components, entity.PlayerFactoryConfig{
		X:           playerX,
		Y:           playerY,
		Speed:       5.0, // units per second
		SpriteSheet: cfg.PlayerSpriteSheet,
		FrameWidth:  cfg.PlayerFrameWidth,
		FrameHeight: cfg.PlayerFrameHeight,
		Columns:     cfg.PlayerSpriteColumns,
		AnimFPS:     cfg.PlayerAnimFPS,
	})

	// Create systems
	inputSystem := system.NewInputSystem(components.Inputs, components.Velocities)
	movementSystem := system.NewMovementSystem(components.Positions, components.Velocities)
	facingSystem := system.NewFacingSystem(components.Facings, components.Velocities)
	animationStateSystem := system.NewAnimationStateSystem(components.Animations, components.Facings)
	animationSystem := system.NewAnimationSystem(components.Animations)
	cameraSystem := system.NewCameraSystem(components.CameraTargets, components.Positions, camera)
	renderSystem := system.NewSpriteSheetRenderSystem(components.Positions, components.Animations, components.SpriteSheets, camera)

	return &BattleScene{
		clock:                component.NewClock(),
		camera:               camera,
		floorTile:            tile,
		tileSize:             tile.Bounds().Dx(),
		world:                world,
		components:           components,
		inputSystem:          inputSystem,
		movementSystem:       movementSystem,
		facingSystem:         facingSystem,
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
	s.movementSystem.Update(s.world, dt)
	s.facingSystem.Update(s.world, dt)
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

	// Draw entities via ECS render system
	s.renderSystem.Draw(s.world, screen)
}
