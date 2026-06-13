package system

import "strings"

// wrapText breaks s into lines no wider than maxWidth, measuring with the injected measure.
//
// It lives in this ebiten-free file (no build tag) so it stays testable under -tags test,
// mirroring the spritesheet_calc.go / spritesheet_render.go split. The render system injects
// text.Advance as measure.
//
// Explicit "\n" splits paragraphs. Within a paragraph it wraps at spaces; a word too wide for
// a line on its own falls back to rune-boundary breaking, which also handles CJK runs that have
// no spaces. A single rune wider than maxWidth is kept on its own line rather than looping forever.
func wrapText(s string, maxWidth float64, measure func(string) float64) []string {
	var lines []string
	for _, para := range strings.Split(s, "\n") {
		lines = append(lines, wrapParagraph(para, maxWidth, measure)...)
	}
	return lines
}

func wrapParagraph(s string, maxWidth float64, measure func(string) float64) []string {
	if measure(s) <= maxWidth {
		return []string{s}
	}

	var lines []string
	cur := ""
	for _, word := range strings.Split(s, " ") {
		candidate := word
		if cur != "" {
			candidate = cur + " " + word
		}
		if measure(candidate) <= maxWidth {
			cur = candidate
			continue
		}

		if cur != "" {
			lines = append(lines, cur)
			cur = ""
		}
		if measure(word) <= maxWidth {
			cur = word
			continue
		}

		// Word is wider than a full line: break it at rune boundaries.
		chunks := breakRunes(word, maxWidth, measure)
		lines = append(lines, chunks[:len(chunks)-1]...)
		cur = chunks[len(chunks)-1]
	}

	// Drop a trailing empty remainder left by a trailing space; keep it only if it is
	// the paragraph's sole line (so an empty paragraph still yields one empty line).
	if cur == "" && len(lines) > 0 {
		return lines
	}
	return append(lines, cur)
}

// breakRunes splits a single word into rune-boundary chunks each within maxWidth.
// A rune wider than maxWidth occupies its own chunk, guaranteeing progress.
func breakRunes(word string, maxWidth float64, measure func(string) float64) []string {
	var chunks []string
	line := make([]rune, 0, len(word))
	for _, r := range word {
		if len(line) == 0 || measure(string(line)+string(r)) <= maxWidth {
			line = append(line, r)
			continue
		}
		chunks = append(chunks, string(line))
		line = []rune{r}
	}
	return append(chunks, string(line))
}
