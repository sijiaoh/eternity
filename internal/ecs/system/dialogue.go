//go:build !test

package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// DialogueSystem advances the singleton dialogue on a press edge.
type DialogueSystem struct {
	dialogue *component.Dialogue
}

func NewDialogueSystem(dialogue *component.Dialogue) *DialogueSystem {
	return &DialogueSystem{dialogue: dialogue}
}

// Update advances to the next line when an advance key/button is just pressed.
// dt is unused: advancement is edge-driven with no typewriter or auto-play, so it needs
// no time base; the parameter only satisfies the System interface.
func (s *DialogueSystem) Update(w *ecs.World, dt float64) {
	if !s.dialogue.Active {
		return
	}
	if advanceJustPressed() {
		s.dialogue.Advance()
	}
}

// advanceJustPressed reports a press edge of Space, Enter, or the left mouse button.
func advanceJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}
