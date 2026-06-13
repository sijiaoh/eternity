package component

import "testing"

func lines(texts ...string) []DialogueLine {
	out := make([]DialogueLine, len(texts))
	for i, t := range texts {
		out[i] = DialogueLine{Text: t}
	}
	return out
}

func TestDialogueStartNonEmpty(t *testing.T) {
	var d Dialogue
	d.Start(lines("a", "b"))

	if !d.Active {
		t.Fatal("expected Active after Start with lines")
	}
	if d.Index != 0 {
		t.Fatalf("expected Index 0, got %d", d.Index)
	}
	cur, ok := d.Current()
	if !ok || cur.Text != "a" {
		t.Fatalf("expected first line \"a\", got %q ok=%v", cur.Text, ok)
	}
}

func TestDialogueStartEmptyStaysInactive(t *testing.T) {
	var d Dialogue
	d.Start(nil)

	if d.Active {
		t.Fatal("expected inactive after Start with no lines")
	}
	if _, ok := d.Current(); ok {
		t.Fatal("expected Current ok=false when inactive")
	}
}

func TestDialogueAdvance(t *testing.T) {
	var d Dialogue
	d.Start(lines("a", "b"))

	d.Advance()
	cur, ok := d.Current()
	if !ok || cur.Text != "b" {
		t.Fatalf("expected second line \"b\", got %q ok=%v", cur.Text, ok)
	}
	if !d.Active {
		t.Fatal("expected still active on last line")
	}
}

func TestDialogueAdvancePastEndDeactivates(t *testing.T) {
	var d Dialogue
	d.Start(lines("a", "b"))

	d.Advance() // -> b
	d.Advance() // past end

	if d.Active {
		t.Fatal("expected inactive after advancing past the last line")
	}
	if _, ok := d.Current(); ok {
		t.Fatal("expected Current ok=false after end")
	}
}

func TestDialogueRestartResetsToFirstLine(t *testing.T) {
	var d Dialogue
	d.Start(lines("a", "b"))
	d.Advance() // Index now 1

	d.Start(lines("x", "y")) // restart must reset the stale cursor

	if !d.Active || d.Index != 0 {
		t.Fatalf("expected restart at Index 0 active, got Active=%v Index=%d", d.Active, d.Index)
	}
	cur, ok := d.Current()
	if !ok || cur.Text != "x" {
		t.Fatalf("expected first line of new script \"x\", got %q ok=%v", cur.Text, ok)
	}
}

func TestDialogueAdvanceWhenInactiveIsNoop(t *testing.T) {
	var d Dialogue
	d.Advance()

	if d.Active || d.Index != 0 {
		t.Fatalf("expected inactive no-op, got Active=%v Index=%d", d.Active, d.Index)
	}
}
