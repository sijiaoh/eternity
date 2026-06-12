package ecs

// World manages entity lifecycle (spawn/despawn).
type World struct {
	nextID     uint32
	free       []Entity // recycled entity slots
	generation []uint32 // current generation for each ID
	alive      []bool   // true if entity at index is alive
	count      int      // number of alive entities
}

// NewWorld creates an empty world with pre-allocated capacity.
func NewWorld(capacity int) *World {
	return &World{
		nextID:     1, // Start from 1 so that Entity{} (zero value) is invalid
		free:       make([]Entity, 0, capacity/4),
		generation: make([]uint32, 0, capacity),
		alive:      make([]bool, 0, capacity),
	}
}

// Spawn creates a new entity.
func (w *World) Spawn() Entity {
	w.count++

	// Reuse recycled slot if available
	if len(w.free) > 0 {
		e := w.free[len(w.free)-1]
		w.free = w.free[:len(w.free)-1]
		w.alive[w.index(e.ID)] = true
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

	idx := w.index(e.ID)
	w.count--
	w.alive[idx] = false
	w.generation[idx]++

	w.free = append(w.free, Entity{
		ID:         e.ID,
		Generation: w.generation[idx],
	})
}

// Alive returns true if the entity exists and matches current generation.
func (w *World) Alive(e Entity) bool {
	idx := w.index(e.ID)
	if idx < 0 || idx >= len(w.alive) {
		return false
	}
	return w.alive[idx] && w.generation[idx] == e.Generation
}

// index converts entity ID to array index.
func (w *World) index(id uint32) int {
	return int(id) - 1
}

// Count returns the number of alive entities.
func (w *World) Count() int {
	return w.count
}

// AllAlive returns a slice containing all currently alive entities.
func (w *World) AllAlive() []Entity {
	result := make([]Entity, 0, w.count)
	for i, alive := range w.alive {
		if alive {
			result = append(result, Entity{
				ID:         uint32(i + 1), // Convert index back to ID
				Generation: w.generation[i],
			})
		}
	}
	return result
}

// SpawnBatch creates multiple entities at once.
func (w *World) SpawnBatch(count int) []Entity {
	if count <= 0 {
		return nil
	}

	entities := make([]Entity, count)
	for i := range count {
		entities[i] = w.Spawn()
	}
	return entities
}

// DespawnBatch removes multiple entities at once.
func (w *World) DespawnBatch(entities []Entity) {
	for _, e := range entities {
		w.Despawn(e)
	}
}
