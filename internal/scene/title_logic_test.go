package scene

import "testing"

// The title screen advances to the next scene only when Enter is pressed. Pinning both
// directions guards the real failure modes: firing without Enter would auto-skip the
// title, and not firing on Enter would strand the player there.
func TestTitleAdvance(t *testing.T) {
	starts := 0
	onStart := func() { starts++ }

	titleAdvance(false, onStart)
	if starts != 0 {
		t.Fatalf("expected no transition without Enter, got %d", starts)
	}

	titleAdvance(true, onStart)
	if starts != 1 {
		t.Fatalf("expected one transition on Enter, got %d", starts)
	}
}
