package assets

import (
	"embed"
)

//go:embed images
var Images embed.FS

//go:embed audio
var Audio embed.FS

//go:embed fonts
var Fonts embed.FS
