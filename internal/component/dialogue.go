package component

// DialogueLine is one line of dialogue, shown all at once (no typewriter).
type DialogueLine struct {
	Speaker  string // display name; "" hides the name
	Text     string
	Portrait string // portrait asset key resolved by the render registry; "" shows no portrait
}

// Dialogue is the scene-owned singleton dialogue state. Only one dialogue is ever active,
// so it mirrors Camera (a plain struct passed to systems) rather than living in per-entity
// ECS storage. It is pure data with no ebiten dependency, so it can be unit tested under -tags test.
type Dialogue struct {
	Lines  []DialogueLine
	Index  int
	Active bool
}

// Start begins a new dialogue. Empty lines leave the dialogue inactive.
func (d *Dialogue) Start(lines []DialogueLine) {
	if len(lines) == 0 {
		d.Lines = nil
		d.Index = 0
		d.Active = false
		return
	}
	d.Lines = lines
	d.Index = 0
	d.Active = true
}

// Current returns the line at the cursor; ok is false when inactive or out of range.
func (d *Dialogue) Current() (DialogueLine, bool) {
	if !d.Active || d.Index < 0 || d.Index >= len(d.Lines) {
		return DialogueLine{}, false
	}
	return d.Lines[d.Index], true
}

// Advance moves to the next line, deactivating once the cursor passes the last line.
func (d *Dialogue) Advance() {
	if !d.Active {
		return
	}
	d.Index++
	if d.Index >= len(d.Lines) {
		d.Active = false
	}
}
