package component

// InputControlled marks an entity as controllable by player input.
// Entities with this component will have their velocity updated by InputSystem.
type InputControlled struct {
	Speed float64 // movement speed in units per second
}

func NewInputControlled(speed float64) *InputControlled {
	return &InputControlled{Speed: speed}
}
