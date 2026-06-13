//go:build !test

package scene

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Manager struct {
	current Scene
}

func NewManager(initial Scene) *Manager {
	return &Manager{current: initial}
}

func (m *Manager) SetScene(s Scene) {
	m.current = s
}

func (m *Manager) Update() error {
	if m.current != nil {
		return m.current.Update()
	}
	return nil
}

func (m *Manager) Draw(screen *ebiten.Image) {
	if m.current != nil {
		m.current.Draw(screen)
	}
}
