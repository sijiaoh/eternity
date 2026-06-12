package ecs

// Entity is a unique identifier for a game object.
// It uses a generation counter to detect stale references.
type Entity struct {
	ID         uint32
	Generation uint32
}

// Valid returns true if the entity has been assigned.
func (e Entity) Valid() bool {
	return e.ID > 0 || e.Generation > 0
}
