package component

// Health tracks current and maximum health for an entity.
type Health struct {
	Current int
	Max     int
}

func NewHealth(max int) *Health {
	return &Health{Current: max, Max: max}
}
