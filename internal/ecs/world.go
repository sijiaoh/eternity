package ecs

// World manages entity lifecycle (spawn/despawn).
type World struct {
	nextID     uint32
	free       []Entity // recycled entity slots
	generation []uint32 // current generation for each ID
	alive      []bool   // true if entity at index is alive
}

// NewWorld creates an empty world with pre-allocated capacity.
func NewWorld(capacity int) *World {
	return &World{
		nextID:     0,
		free:       make([]Entity, 0, capacity/4),
		generation: make([]uint32, 0, capacity),
		alive:      make([]bool, 0, capacity),
	}
}

// Spawn creates a new entity.
func (w *World) Spawn() Entity {
	// Reuse recycled slot if available
	if len(w.free) > 0 {
		e := w.free[len(w.free)-1]
		w.free = w.free[:len(w.free)-1]
		w.alive[e.ID] = true
		return e
	}

	// Create new slot
	id := w.nextID
	w.nextID++

	w.generation = append(w.generation, 0)
	w.alive = append(w.alive, true)

	return Entity{ID: id, Generation: 0}
}

// Despawn removes an entity and recycles its slot.
func (w *World) Despawn(e Entity) {
	if !w.Alive(e) {
		return
	}

	w.alive[e.ID] = false
	w.generation[e.ID]++

	w.free = append(w.free, Entity{
		ID:         e.ID,
		Generation: w.generation[e.ID],
	})
}

// Alive returns true if the entity exists and matches current generation.
func (w *World) Alive(e Entity) bool {
	if int(e.ID) >= len(w.alive) {
		return false
	}
	return w.alive[e.ID] && w.generation[e.ID] == e.Generation
}

// Count returns the number of alive entities.
func (w *World) Count() int {
	count := 0
	for _, a := range w.alive {
		if a {
			count++
		}
	}
	return count
}
