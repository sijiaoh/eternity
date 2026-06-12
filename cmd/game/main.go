//go:build !test

package main

import (
	"log"

	"eternity/internal/config"
	"eternity/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g, err := game.New()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Vampire Survivors Clone")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
