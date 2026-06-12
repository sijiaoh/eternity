package ecs

// Storage is a type-safe container for components of a specific type.
// It uses a sparse set for O(1) add/remove/lookup and cache-friendly iteration.
type Storage[T any] struct {
	// dense holds the actual component data
	dense []T
	// entities maps dense index to entity
	entities []Entity
	// sparse maps entity ID to dense index (or -1 if not present)
	sparse []int
}

// NewStorage creates a storage with pre-allocated capacity.
func NewStorage[T any](capacity int) *Storage[T] {
	return &Storage[T]{
		dense:    make([]T, 0, capacity),
		entities: make([]Entity, 0, capacity),
		sparse:   make([]int, 0),
	}
}

// Set adds or updates a component for an entity.
// If the entity's generation is older than the stored one, the call is ignored.
func (s *Storage[T]) Set(e Entity, c T) {
	s.ensureSparse(int(e.ID))

	idx := s.sparse[e.ID]
	if idx >= 0 && idx < len(s.entities) {
		stored := s.entities[idx]
		if e.Generation < stored.Generation {
			return // Ignore stale reference
		}
		if e.Generation == stored.Generation {
			s.dense[idx] = c
			return
		}
		// e.Generation > stored.Generation: replace old data
		s.dense[idx] = c
		s.entities[idx] = e
		return
	}

	s.sparse[e.ID] = len(s.dense)
	s.dense = append(s.dense, c)
	s.entities = append(s.entities, e)
}

// Get returns the component and true if found, zero value and false otherwise.
func (s *Storage[T]) Get(e Entity) (T, bool) {
	if int(e.ID) >= len(s.sparse) {
		var zero T
		return zero, false
	}

	idx := s.sparse[e.ID]
	if idx < 0 || idx >= len(s.entities) || s.entities[idx].Generation != e.Generation {
		var zero T
		return zero, false
	}

	return s.dense[idx], true
}

// GetPtr returns a pointer to the component for direct mutation.
func (s *Storage[T]) GetPtr(e Entity) *T {
	if int(e.ID) >= len(s.sparse) {
		return nil
	}

	idx := s.sparse[e.ID]
	if idx < 0 || idx >= len(s.entities) || s.entities[idx].Generation != e.Generation {
		return nil
	}

	return &s.dense[idx]
}

// Has returns true if the entity has this component.
func (s *Storage[T]) Has(e Entity) bool {
	if int(e.ID) >= len(s.sparse) {
		return false
	}

	idx := s.sparse[e.ID]
	return idx >= 0 && idx < len(s.entities) && s.entities[idx].Generation == e.Generation
}

// Remove removes the component from an entity.
func (s *Storage[T]) Remove(e Entity) {
	if int(e.ID) >= len(s.sparse) {
		return
	}

	idx := s.sparse[e.ID]
	if idx < 0 || idx >= len(s.entities) || s.entities[idx].Generation != e.Generation {
		return
	}

	// Swap with last element
	lastIdx := len(s.dense) - 1
	if idx != lastIdx {
		s.dense[idx] = s.dense[lastIdx]
		s.entities[idx] = s.entities[lastIdx]
		s.sparse[s.entities[idx].ID] = idx
	}

	s.dense = s.dense[:lastIdx]
	s.entities = s.entities[:lastIdx]
	s.sparse[e.ID] = -1
}

// Each iterates over all components with their entities.
func (s *Storage[T]) Each(fn func(e Entity, c *T)) {
	for i := range s.dense {
		fn(s.entities[i], &s.dense[i])
	}
}

// Len returns the number of stored components.
func (s *Storage[T]) Len() int {
	return len(s.dense)
}

// ensureSparse grows sparse array to accommodate entity ID.
func (s *Storage[T]) ensureSparse(id int) {
	if id < len(s.sparse) {
		return
	}

	newSize := id + 1
	if newSize < len(s.sparse)*2 {
		newSize = len(s.sparse) * 2
	}
	if newSize < 64 {
		newSize = 64
	}

	grown := make([]int, newSize)
	copy(grown, s.sparse)
	for i := len(s.sparse); i < newSize; i++ {
		grown[i] = -1
	}
	s.sparse = grown
}
