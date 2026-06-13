//go:build !test

package system

import (
	"image/color"
	"strings"

	"eternity/internal/component"
	"eternity/internal/config"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	dialogueBoxColor     = color.RGBA{R: 12, G: 14, B: 28, A: 220}
	dialogueBorderColor  = color.RGBA{R: 120, G: 130, B: 170, A: 255}
	dialogueTextColor    = color.RGBA{R: 240, G: 240, B: 240, A: 255}
	dialogueSpeakerColor = color.RGBA{R: 250, G: 215, B: 120, A: 255}
)

// DialogueRenderSystem draws the dialogue overlay in screen-pixel space (never camera-relative),
// so it stays fixed on screen above the world.
type DialogueRenderSystem struct {
	dialogue  *component.Dialogue
	portraits map[string]*ebiten.Image // portrait key -> image; nil/missing keys draw no portrait
	face      text.Face
}

func NewDialogueRenderSystem(
	dialogue *component.Dialogue,
	portraits map[string]*ebiten.Image,
	face text.Face,
) *DialogueRenderSystem {
	return &DialogueRenderSystem{
		dialogue:  dialogue,
		portraits: portraits,
		face:      face,
	}
}

// Draw renders the box, optional portrait, speaker name, and wrapped body text.
// w is unused: the overlay reads only the singleton dialogue and never iterates entities.
func (s *DialogueRenderSystem) Draw(w *ecs.World, screen *ebiten.Image) {
	line, ok := s.dialogue.Current()
	if !ok {
		return
	}

	s.drawBox(screen)

	if img := s.portraits[line.Portrait]; line.Portrait != "" && img != nil {
		s.drawPortrait(screen, img)
	}

	if line.Speaker != "" {
		s.drawText(screen, line.Speaker, config.DialogueSpeakerX, config.DialogueSpeakerY, dialogueSpeakerColor)
	}

	body := wrapText(line.Text, config.DialogueTextMaxWidth, s.measure)
	s.drawText(screen, strings.Join(body, "\n"), config.DialogueTextX, config.DialogueTextY, dialogueTextColor)
}

func (s *DialogueRenderSystem) drawBox(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		config.DialogueBoxX, config.DialogueBoxY,
		config.DialogueBoxWidth, config.DialogueBoxHeight,
		dialogueBoxColor, false)
	vector.StrokeRect(screen,
		config.DialogueBoxX, config.DialogueBoxY,
		config.DialogueBoxWidth, config.DialogueBoxHeight,
		2, dialogueBorderColor, false)
}

func (s *DialogueRenderSystem) drawPortrait(screen *ebiten.Image, img *ebiten.Image) {
	src := img.Bounds().Dx()
	scale := float64(config.DialoguePortraitSize) / float64(src)
	x := config.DialoguePortraitCenterX - config.DialoguePortraitSize/2
	y := config.DialogueBoxY - config.DialoguePortraitSize

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scale, scale)
	opts.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, opts)
}

func (s *DialogueRenderSystem) drawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	opts.ColorScale.ScaleWithColor(clr)
	opts.LineSpacing = lineHeight(s.face)
	text.Draw(screen, str, s.face, opts)
}

func (s *DialogueRenderSystem) measure(str string) float64 {
	return text.Advance(str, s.face)
}

// lineHeight returns the baseline-to-baseline distance for the face.
func lineHeight(face text.Face) float64 {
	m := face.Metrics()
	return m.HAscent + m.HDescent + m.HLineGap
}
