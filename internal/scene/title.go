//go:build !test

package scene

import (
	"image/color"

	"eternity/internal/config"
	"eternity/internal/i18n"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	titleBackgroundColor = color.RGBA{R: 12, G: 14, B: 28, A: 255}
	titleTextColor       = color.RGBA{R: 240, G: 240, B: 240, A: 255}
	titlePromptColor     = color.RGBA{R: 180, G: 188, B: 210, A: 255}
)

// TitleScene is the startup screen: it shows the game's title and a "press Enter"
// prompt, then runs onStart when the player presses Enter. The transition target is
// injected as the onStart callback, so TitleScene stays decoupled from the Manager
// and the concrete next scene.
type TitleScene struct {
	bundle     *i18n.Bundle
	titleFace  text.Face
	promptFace text.Face
	onStart    func()
}

// NewTitleScene builds the title screen. titleFace renders the brand name, promptFace
// the smaller hint; onStart is invoked once when the player presses Enter.
func NewTitleScene(bundle *i18n.Bundle, titleFace, promptFace text.Face, onStart func()) *TitleScene {
	return &TitleScene{
		bundle:     bundle,
		titleFace:  titleFace,
		promptFace: promptFace,
		onStart:    onStart,
	}
}

func (s *TitleScene) Update() error {
	titleAdvance(inpututil.IsKeyJustPressed(ebiten.KeyEnter), s.onStart)
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(titleBackgroundColor)

	// The on-screen title is the game's brand name; reuse window.title rather than
	// duplicating the string across locale files (DRY).
	drawCenteredText(screen, s.bundle.Get("window.title"), s.titleFace,
		config.ScreenWidth/2, config.TitleCenterY, titleTextColor)
	drawCenteredText(screen, s.bundle.Get("title.press_enter"), s.promptFace,
		config.ScreenWidth/2, config.TitlePromptCenterY, titlePromptColor)
}

// drawCenteredText draws str centered horizontally on centerX and vertically on centerY.
func drawCenteredText(screen *ebiten.Image, str string, face text.Face, centerX, centerY int, clr color.Color) {
	spacing := faceLineHeight(face)
	w, h := text.Measure(str, face, spacing)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(centerX)-w/2, float64(centerY)-h/2)
	op.ColorScale.ScaleWithColor(clr)
	op.LineSpacing = spacing
	text.Draw(screen, str, face, op)
}

// faceLineHeight returns the baseline-to-baseline distance for the face.
func faceLineHeight(face text.Face) float64 {
	m := face.Metrics()
	return m.HAscent + m.HDescent + m.HLineGap
}
