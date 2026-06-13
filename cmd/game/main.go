//go:build !test

package main

import (
	"log"

	"eternity/internal/config"
	"eternity/internal/game"
	"eternity/internal/i18n"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	bundle, err := i18n.New()
	if err != nil {
		log.Fatal(err)
	}

	g, err := game.New(bundle)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle(bundle.Get("window.title"))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
