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
	var zero Entity
	if zero.Valid() {
		t.Error("zero entity should not be valid")
	}

	e := Entity{ID: 0, Generation: 1}
	if !e.Valid() {
		t.Error("entity with generation > 0 should be valid")
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
