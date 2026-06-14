package component

import "testing"

func TestAnimationStateName(t *testing.T) {
	tests := []struct {
		walking bool
		dir     FacingDirection
		want    string
	}{
		{false, FacingDown, "idle_down"},
		{true, FacingDown, "walk_down"},
		{false, FacingLeft, "idle_left"},
		{true, FacingLeft, "walk_left"},
		{false, FacingUp, "idle_up"},
		{true, FacingUp, "walk_up"},
		{false, FacingRight, "idle_right"},
		{true, FacingRight, "walk_right"},
	}
	for _, tt := range tests {
		if got := AnimationStateName(tt.walking, tt.dir); got != tt.want {
			t.Errorf("AnimationStateName(%v, %v) = %s, want %s", tt.walking, tt.dir, got, tt.want)
		}
	}
}

func TestDirectionalAnimation_FourUsesFacingAsIs(t *testing.T) {
	d := NewDirectionalAnimation(DirectionsFour)
	tests := []struct {
		dir  FacingDirection
		want string
	}{
		{FacingUp, "walk_up"},
		{FacingDown, "walk_down"},
		{FacingLeft, "walk_left"},
		{FacingRight, "walk_right"},
	}
	for _, tt := range tests {
		if got := d.ResolveState(true, tt.dir); got != tt.want {
			t.Errorf("ResolveState(true, %v) = %s, want %s", tt.dir, got, tt.want)
		}
	}
}

func TestDirectionalAnimation_HorizontalFallsBackToLastHorizontal(t *testing.T) {
	d := NewDirectionalAnimation(DirectionsHorizontal)

	// Facing right is supported and remembered.
	if got := d.ResolveState(true, FacingRight); got != "walk_right" {
		t.Fatalf("ResolveState(right) = %s, want walk_right", got)
	}
	// Moving up has no artwork → reuse last horizontal (right).
	if got := d.ResolveState(true, FacingUp); got != "walk_right" {
		t.Errorf("ResolveState(up) = %s, want walk_right (fallback)", got)
	}
	// Turn left, then move down → reuse last horizontal (left).
	if got := d.ResolveState(false, FacingLeft); got != "idle_left" {
		t.Fatalf("ResolveState(left) = %s, want idle_left", got)
	}
	if got := d.ResolveState(true, FacingDown); got != "walk_left" {
		t.Errorf("ResolveState(down) = %s, want walk_left (fallback)", got)
	}
}

func TestDirectionalAnimation_HorizontalDefaultsToRight(t *testing.T) {
	d := NewDirectionalAnimation(DirectionsHorizontal)
	// Vertical facing before any horizontal input falls back to the initial right.
	if got := d.ResolveState(false, FacingDown); got != "idle_right" {
		t.Errorf("ResolveState(down) = %s, want idle_right (default)", got)
	}
}

func TestDirectionalSheetSpec_StatesCoverEveryDirection(t *testing.T) {
	spec := DirectionalSheetSpec{
		Directions: DirectionsFour,
		Rows: map[FacingDirection]int{
			FacingDown: 0, FacingLeft: 6, FacingUp: 12, FacingRight: 18,
		},
		IdleFrames: 1, WalkFrames: 6, FPS: 8,
	}

	states := spec.States()
	byName := make(map[string]AnimationState, len(states))
	for _, s := range states {
		byName[s.Name] = s
	}

	if len(states) != 8 {
		t.Fatalf("got %d states, want 8 (idle+walk per direction)", len(states))
	}
	// Idle and walk of a direction share the row's start frame; only the count differs.
	if got := byName["idle_left"]; got.StartFrame != 6 || got.FrameCount != 1 {
		t.Errorf("idle_left = %+v, want StartFrame 6 FrameCount 1", got)
	}
	if got := byName["walk_left"]; got.StartFrame != 6 || got.FrameCount != 6 {
		t.Errorf("walk_left = %+v, want StartFrame 6 FrameCount 6", got)
	}
}

func TestDirectionalSheetSpec_HorizontalGeneratesOnlyLeftRight(t *testing.T) {
	spec := DirectionalSheetSpec{
		Directions: DirectionsHorizontal,
		Rows:       map[FacingDirection]int{FacingLeft: 0, FacingRight: 4},
		IdleFrames: 1, WalkFrames: 4, FPS: 8,
	}

	states := spec.States()
	if len(states) != 4 {
		t.Fatalf("got %d states, want 4 (idle+walk for left/right)", len(states))
	}
	for _, s := range states {
		if s.Name == "idle_up" || s.Name == "walk_up" || s.Name == "idle_down" || s.Name == "walk_down" {
			t.Errorf("horizontal-only spec produced vertical state %q", s.Name)
		}
	}
}

func TestDirectionalSheetSpec_MissingRowPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic when a direction in the set has no row")
		}
	}()

	spec := DirectionalSheetSpec{
		Directions: DirectionsFour,
		Rows:       map[FacingDirection]int{FacingDown: 0}, // missing left/up/right
		IdleFrames: 1, WalkFrames: 6, FPS: 8,
	}
	spec.States()
}
