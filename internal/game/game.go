package game

import (
	"ebiten-agent-example/internal/config"
	"ebiten-agent-example/internal/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	sceneManager *scene.Manager
}

func New() (*Game, error) {
	battleScene, err := scene.NewBattleScene("assets/images/battle/floor/stone.png")
	if err != nil {
		return nil, err
	}

	return &Game{
		sceneManager: scene.NewManager(battleScene),
	}, nil
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
