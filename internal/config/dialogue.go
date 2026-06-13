package config

// Dialogue box layout in screen pixels (screen is 800x600).
// All coordinates are screen pixels (not camera-relative): the overlay is fixed UI.
const (
	DialogueBoxX      = 32
	DialogueBoxY      = 416
	DialogueBoxWidth  = 736
	DialogueBoxHeight = 160
	DialoguePadding   = 24

	// DialogueSpeakerY sits above the body so the name does not overlap the first text line.
	DialogueSpeakerX = DialogueBoxX + DialoguePadding
	DialogueSpeakerY = DialogueBoxY + 16

	DialogueTextX        = DialogueBoxX + DialoguePadding
	DialogueTextY        = DialogueBoxY + 56
	DialogueTextMaxWidth = DialogueBoxWidth - 2*DialoguePadding

	// DialogueFontSize sets the face size used for both glyph rendering and line spacing.
	DialogueFontSize = 22

	// DialoguePortrait is drawn centered horizontally, bottom edge resting on the box top.
	DialoguePortraitSize    = 192
	DialoguePortraitCenterX = ScreenWidth / 2
)
