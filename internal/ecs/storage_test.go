package ecs

import "testing"

type testPosition struct {
	X, Y float64
}

func TestStorage_SetGet(t *testing.T) {
	s := NewStorage[testPosition](10)
	e := Entity{ID: 1, Generation: 0}

	s.Set(e, testPosition{X: 10, Y: 20})

	pos, ok := s.Get(e)
	if !ok {
		t.Fatal("component should exist")
	}

	if pos.X != 10 || pos.Y != 20 {
		t.Errorf("position = (%v, %v), want (10, 20)", pos.X, pos.Y)
	}
}

func TestStorage_GetPtr(t *testing.T) {
	s := NewStorage[testPosition](10)
	e := Entity{ID: 1, Generation: 0}

	s.Set(e, testPosition{X: 10, Y: 20})

	ptr := s.GetPtr(e)
	if ptr == nil {
		t.Fatal("GetPtr should return non-nil")
	}

	ptr.X = 100
	pos, _ := s.Get(e)
	if pos.X != 100 {
		t.Error("mutation through GetPtr should be reflected")
	}
}

func TestStorage_Has(t *testing.T) {
	s := NewStorage[testPosition](10)
	e1 := Entity{ID: 1, Generation: 0}
	e2 := Entity{ID: 2, Generation: 0}

	s.Set(e1, testPosition{})

	if !s.Has(e1) {
		t.Error("Has should return true for existing entity")
	}

	if s.Has(e2) {
		t.Error("Has should return false for non-existing entity")
	}
}

func TestStorage_Remove(t *testing.T) {
	s := NewStorage[testPosition](10)
	e := Entity{ID: 1, Generation: 0}

	s.Set(e, testPosition{X: 10, Y: 20})
	s.Remove(e)

	if s.Has(e) {
		t.Error("removed component should not exist")
	}

	if s.Len() != 0 {
		t.Errorf("Len = %d, want 0", s.Len())
	}
}

func TestStorage_Each(t *testing.T) {
	s := NewStorage[testPosition](10)

	entities := []Entity{
		{ID: 1, Generation: 0},
		{ID: 2, Generation: 0},
		{ID: 3, Generation: 0},
	}

	for i, e := range entities {
		s.Set(e, testPosition{X: float64(i), Y: float64(i * 10)})
	}

	count := 0
	s.Each(func(e Entity, c *testPosition) {
		count++
	})

	if count != 3 {
		t.Errorf("Each visited %d, want 3", count)
	}
}

func TestStorage_GenerationCheck(t *testing.T) {
	s := NewStorage[testPosition](10)

	e1 := Entity{ID: 1, Generation: 0}
	s.Set(e1, testPosition{X: 10, Y: 20})

	// Stale reference with different generation
	stale := Entity{ID: 1, Generation: 1}

	if s.Has(stale) {
		t.Error("stale reference should not have component")
	}

	_, ok := s.Get(stale)
	if ok {
		t.Error("Get should return false for stale reference")
	}

	if s.GetPtr(stale) != nil {
		t.Error("GetPtr should return nil for stale reference")
	}
}

func TestStorage_SwapRemove(t *testing.T) {
	s := NewStorage[testPosition](10)

	e1 := Entity{ID: 1, Generation: 0}
	e2 := Entity{ID: 2, Generation: 0}
	e3 := Entity{ID: 3, Generation: 0}

	s.Set(e1, testPosition{X: 1, Y: 1})
	s.Set(e2, testPosition{X: 2, Y: 2})
	s.Set(e3, testPosition{X: 3, Y: 3})

	// Remove middle element
	s.Remove(e2)

	// e1 and e3 should still be accessible
	if !s.Has(e1) || !s.Has(e3) {
		t.Error("remaining entities should exist after swap-remove")
	}

	pos1, _ := s.Get(e1)
	pos3, _ := s.Get(e3)

	if pos1.X != 1 || pos3.X != 3 {
		t.Error("remaining components should retain correct values")
	}

	if s.Len() != 2 {
		t.Errorf("Len = %d, want 2", s.Len())
	}
}

func TestStorage_Update(t *testing.T) {
	s := NewStorage[testPosition](10)
	e := Entity{ID: 1, Generation: 0}

	s.Set(e, testPosition{X: 10, Y: 20})
	s.Set(e, testPosition{X: 100, Y: 200})

	pos, _ := s.Get(e)
	if pos.X != 100 || pos.Y != 200 {
		t.Errorf("position = (%v, %v), want (100, 200)", pos.X, pos.Y)
	}

	if s.Len() != 1 {
		t.Errorf("Len = %d, want 1 (update should not add)", s.Len())
	}
}

func TestStorage_RemoveNonexistent(t *testing.T) {
	s := NewStorage[testPosition](10)
	e := Entity{ID: 1, Generation: 0}

	// Should not panic
	s.Remove(e)

	if s.Len() != 0 {
		t.Errorf("Len = %d, want 0", s.Len())
	}
}

func TestStorage_SetWithStaleReference(t *testing.T) {
	s := NewStorage[testPosition](10)

	old := Entity{ID: 1, Generation: 0}
	s.Set(old, testPosition{X: 10, Y: 20})

	// Simulate entity reuse with new generation
	newEntity := Entity{ID: 1, Generation: 1}
	s.Set(newEntity, testPosition{X: 100, Y: 200})

	// Stale reference Set should add new entry, not update
	s.Set(old, testPosition{X: 999, Y: 999})

	// New entity should be unaffected
	pos, ok := s.Get(newEntity)
	if !ok {
		t.Fatal("new entity should exist")
	}
	if pos.X != 100 || pos.Y != 200 {
		t.Errorf("new entity position = (%v, %v), want (100, 200)", pos.X, pos.Y)
	}
}

func TestStorage_RemoveWithStaleReference(t *testing.T) {
	s := NewStorage[testPosition](10)

	old := Entity{ID: 1, Generation: 0}
	s.Set(old, testPosition{X: 10, Y: 20})

	// Simulate entity reuse with new generation
	newEntity := Entity{ID: 1, Generation: 1}
	s.Set(newEntity, testPosition{X: 100, Y: 200})

	// Stale reference Remove should not affect new entity
	s.Remove(old)

	// New entity should still exist
	if !s.Has(newEntity) {
		t.Error("new entity should still exist after stale Remove")
	}

	pos, ok := s.Get(newEntity)
	if !ok {
		t.Fatal("new entity should have component")
	}
	if pos.X != 100 || pos.Y != 200 {
		t.Errorf("new entity position = (%v, %v), want (100, 200)", pos.X, pos.Y)
	}

	if s.Len() != 1 {
		t.Errorf("Len = %d, want 1", s.Len())
	}
}
