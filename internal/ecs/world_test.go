package ecs

import "testing"

func TestWorld_Spawn(t *testing.T) {
	w := NewWorld(10)

	e1 := w.Spawn()
	e2 := w.Spawn()

	if e1.ID == e2.ID {
		t.Error("spawned entities should have different IDs")
	}

	if !w.Alive(e1) || !w.Alive(e2) {
		t.Error("spawned entities should be alive")
	}

	if w.Count() != 2 {
		t.Errorf("count = %d, want 2", w.Count())
	}
}

func TestWorld_Despawn(t *testing.T) {
	w := NewWorld(10)

	e := w.Spawn()
	w.Despawn(e)

	if w.Alive(e) {
		t.Error("despawned entity should not be alive")
	}

	if w.Count() != 0 {
		t.Errorf("count = %d, want 0", w.Count())
	}
}

func TestWorld_ReuseSlot(t *testing.T) {
	w := NewWorld(10)

	e1 := w.Spawn()
	oldID := e1.ID
	w.Despawn(e1)

	e2 := w.Spawn()

	if e2.ID != oldID {
		t.Error("new entity should reuse despawned slot")
	}

	if e2.Generation != 1 {
		t.Errorf("generation = %d, want 1", e2.Generation)
	}

	// Stale reference should not be alive
	if w.Alive(e1) {
		t.Error("stale entity reference should not be alive")
	}
}

func TestEntity_Valid(t *testing.T) {
	// Zero value Entity is invalid
	var zero Entity
	if zero.Valid() {
		t.Error("zero entity should not be valid")
	}

	// Entity with ID=0 is always invalid (reserved for zero value)
	e := Entity{ID: 0, Generation: 1}
	if e.Valid() {
		t.Error("entity with ID=0 should not be valid")
	}

	// Entity with ID>0 is valid
	e2 := Entity{ID: 1, Generation: 0}
	if !e2.Valid() {
		t.Error("entity with ID > 0 should be valid")
	}
}

func TestSpawnedEntity_IsValid(t *testing.T) {
	w := NewWorld(10)

	// First spawned entity should be valid
	e := w.Spawn()
	if !e.Valid() {
		t.Errorf("spawned entity %v should be valid", e)
	}

	// First ID should be exactly 1 (not 0, to reserve zero value as invalid)
	if e.ID != 1 {
		t.Errorf("first entity ID = %d, want 1", e.ID)
	}
}

func TestWorld_DespawnTwice(t *testing.T) {
	w := NewWorld(10)

	e := w.Spawn()
	w.Despawn(e)
	w.Despawn(e) // Should not panic or corrupt state

	if w.Count() != 0 {
		t.Errorf("count = %d, want 0", w.Count())
	}
}

func TestWorld_AliveInvalidEntity(t *testing.T) {
	w := NewWorld(10)

	// Entity that was never spawned
	e := Entity{ID: 999, Generation: 0}

	if w.Alive(e) {
		t.Error("unspawned entity should not be alive")
	}
}

func TestWorld_AliveZeroEntity(t *testing.T) {
	w := NewWorld(10)
	w.Spawn() // Ensure world has some state

	// Zero value Entity should safely return false (idx becomes -1)
	var zero Entity
	if w.Alive(zero) {
		t.Error("zero entity should not be alive")
	}
}

func TestWorld_AllAlive(t *testing.T) {
	w := NewWorld(10)

	// Empty world
	if len(w.AllAlive()) != 0 {
		t.Error("empty world should return empty slice")
	}

	// Spawn some entities
	e1 := w.Spawn()
	e2 := w.Spawn()
	e3 := w.Spawn()

	all := w.AllAlive()
	if len(all) != 3 {
		t.Errorf("len(AllAlive()) = %d, want 3", len(all))
	}

	// Verify all returned entities are alive
	for _, e := range all {
		if !w.Alive(e) {
			t.Errorf("entity %v from AllAlive() should be alive", e)
		}
	}

	// Despawn one and verify
	w.Despawn(e2)
	all = w.AllAlive()
	if len(all) != 2 {
		t.Errorf("len(AllAlive()) = %d, want 2", len(all))
	}

	// e2 should not be in the result
	for _, e := range all {
		if e.ID == e2.ID && e.Generation == e2.Generation {
			t.Error("despawned entity should not be in AllAlive()")
		}
	}

	// e1 and e3 should still be present with correct generation
	found := make(map[Entity]bool)
	for _, e := range all {
		found[e] = true
	}
	if !found[e1] {
		t.Errorf("e1 %v should be in AllAlive()", e1)
	}
	if !found[e3] {
		t.Errorf("e3 %v should be in AllAlive()", e3)
	}
}

func TestWorld_SpawnBatch(t *testing.T) {
	w := NewWorld(10)

	// Spawn batch of 5 entities
	entities := w.SpawnBatch(5)

	if len(entities) != 5 {
		t.Errorf("len(SpawnBatch(5)) = %d, want 5", len(entities))
	}

	if w.Count() != 5 {
		t.Errorf("count = %d, want 5", w.Count())
	}

	// All entities should be alive and have unique IDs
	ids := make(map[uint32]bool)
	for _, e := range entities {
		if !w.Alive(e) {
			t.Errorf("entity %v should be alive", e)
		}
		if ids[e.ID] {
			t.Errorf("duplicate ID %d in batch", e.ID)
		}
		ids[e.ID] = true
	}
}

func TestWorld_SpawnBatch_Zero(t *testing.T) {
	w := NewWorld(10)

	entities := w.SpawnBatch(0)
	if entities != nil {
		t.Error("SpawnBatch(0) should return nil")
	}

	if w.Count() != 0 {
		t.Errorf("count = %d, want 0", w.Count())
	}
}

func TestWorld_SpawnBatch_Negative(t *testing.T) {
	w := NewWorld(10)

	entities := w.SpawnBatch(-5)
	if entities != nil {
		t.Error("SpawnBatch(-5) should return nil")
	}

	if w.Count() != 0 {
		t.Errorf("count = %d, want 0", w.Count())
	}
}

func TestWorld_SpawnBatch_ReusesSlots(t *testing.T) {
	w := NewWorld(10)

	// Spawn and despawn some entities to create free slots
	first := w.SpawnBatch(3)
	firstIDs := make(map[uint32]bool)
	for _, e := range first {
		firstIDs[e.ID] = true
	}
	w.DespawnBatch(first)

	// New batch should reuse freed slots (order may differ due to LIFO)
	second := w.SpawnBatch(3)

	// All IDs should be from the first batch, with incremented generations
	for _, e := range second {
		if !firstIDs[e.ID] {
			t.Errorf("expected ID %d to be reused from first batch", e.ID)
		}
		if e.Generation != 1 {
			t.Errorf("expected generation 1, got %d for ID %d", e.Generation, e.ID)
		}
	}
}

func TestWorld_DespawnBatch(t *testing.T) {
	w := NewWorld(10)

	entities := w.SpawnBatch(5)
	w.DespawnBatch(entities)

	if w.Count() != 0 {
		t.Errorf("count = %d, want 0", w.Count())
	}

	for _, e := range entities {
		if w.Alive(e) {
			t.Errorf("entity %v should not be alive", e)
		}
	}
}

func TestWorld_DespawnBatch_Partial(t *testing.T) {
	w := NewWorld(10)

	entities := w.SpawnBatch(5)
	w.DespawnBatch(entities[:3])

	if w.Count() != 2 {
		t.Errorf("count = %d, want 2", w.Count())
	}

	// First 3 should be dead
	for i := range 3 {
		if w.Alive(entities[i]) {
			t.Errorf("entity %d should not be alive", i)
		}
	}

	// Last 2 should be alive
	for i := 3; i < 5; i++ {
		if !w.Alive(entities[i]) {
			t.Errorf("entity %d should be alive", i)
		}
	}
}

func TestWorld_DespawnBatch_Empty(t *testing.T) {
	w := NewWorld(10)
	w.Spawn()

	// Should not panic
	w.DespawnBatch(nil)
	w.DespawnBatch([]Entity{})

	if w.Count() != 1 {
		t.Errorf("count = %d, want 1", w.Count())
	}
}

func TestWorld_DespawnBatch_StaleEntities(t *testing.T) {
	w := NewWorld(10)

	entities := w.SpawnBatch(3)
	w.DespawnBatch(entities)

	// Spawn new entities to get same IDs with new generations
	newEntities := w.SpawnBatch(3)

	// Trying to despawn old stale entities should have no effect
	w.DespawnBatch(entities)

	// New entities should still be alive
	for _, e := range newEntities {
		if !w.Alive(e) {
			t.Errorf("new entity %v should still be alive", e)
		}
	}

	if w.Count() != 3 {
		t.Errorf("count = %d, want 3", w.Count())
	}
}
