package main

import (
	"log"

	"ebiten-agent-example/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.New()

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Vampire Survivors Clone")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
