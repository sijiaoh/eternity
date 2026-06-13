package config

// Title screen layout in screen pixels (the title is fixed UI, drawn directly in
// screen space rather than camera-relative).
const (
	TitleFontSize       = 72
	TitlePromptFontSize = 24

	// Vertical centers for the title and the prompt, placed symmetrically around mid-screen.
	TitleCenterY       = ScreenHeight * 2 / 5 // 240
	TitlePromptCenterY = ScreenHeight * 3 / 5 // 360
)
