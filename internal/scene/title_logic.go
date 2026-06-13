package scene

// titleAdvance invokes onStart when Enter was just pressed. It lives in this ebiten-free
// file — separate from the `//go:build !test` render/input code — so the start condition
// is unit-testable in a headless build.
func titleAdvance(enterPressed bool, onStart func()) {
	if enterPressed {
		onStart()
	}
}
