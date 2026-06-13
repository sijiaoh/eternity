//go:build !test

package main

import (
	"flag"
	"log"
	"os"

	"eternity/internal/config"
	"eternity/internal/game"
	"eternity/internal/i18n"
	"eternity/internal/scenario"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Two developer/debug entry points, both leaving the normal player path untouched when unused:
	//   -scene <name>     boot straight into a scene with its default situation.
	//   -scenario <path>  load a JSON file that names the scene and reproduces an in-scene
	//                     situation (see internal/scenario). It supersedes -scene, so the two
	//                     are mutually exclusive.
	startScene := flag.String("scene", "", "start scene to launch (developer/debug); default is the title screen")
	scenarioPath := flag.String("scenario", "", "path to a debug scenario file (JSON); see internal/scenario")
	flag.Parse()

	sc := scenario.Scenario{Scene: *startScene}
	if *scenarioPath != "" {
		if *startScene != "" {
			log.Fatal("specify either -scene or -scenario, not both")
		}
		data, err := os.ReadFile(*scenarioPath)
		if err != nil {
			log.Fatalf("read scenario file: %v", err)
		}
		if sc, err = scenario.Parse(data); err != nil {
			log.Fatal(err)
		}
	}

	bundle, err := i18n.New()
	if err != nil {
		log.Fatal(err)
	}

	g, err := game.New(bundle, sc)
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
