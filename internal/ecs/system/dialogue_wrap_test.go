package system

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// runeWidth measures one unit per rune, so maxWidth is a rune budget per line.
func runeWidth(s string) float64 { return float64(utf8.RuneCountInString(s)) }

func TestWrapTextShortStaysOneLine(t *testing.T) {
	got := wrapText("hello", 10, runeWidth)
	if len(got) != 1 || got[0] != "hello" {
		t.Fatalf("expected [hello], got %q", got)
	}
}

func TestWrapTextWrapsAtSpaces(t *testing.T) {
	got := wrapText("aaa bbb ccc", 7, runeWidth)
	want := []string{"aaa bbb", "ccc"}
	if !equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWrapTextExplicitNewline(t *testing.T) {
	got := wrapText("aaa\nbbb", 10, runeWidth)
	want := []string{"aaa", "bbb"}
	if !equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWrapTextLongWordBreaksByRune(t *testing.T) {
	got := wrapText("aaaaaaaaaa", 4, runeWidth)
	want := []string{"aaaa", "aaaa", "aa"}
	if !equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// CJK has no spaces: the whole run is one word and must break at rune boundaries.
func TestWrapTextCJKBreaksByRune(t *testing.T) {
	got := wrapText("你好世界你好", 2, runeWidth)
	want := []string{"你好", "世界", "你好"}
	if !equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// A single rune wider than the line must not loop forever; it takes its own line.
func TestWrapTextOverWideRuneTerminates(t *testing.T) {
	doubleWidth := func(s string) float64 { return 2 * runeWidth(s) }
	got := wrapText("ab", 1, doubleWidth)
	want := []string{"a", "b"}
	if !equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// A trailing space after an over-wide word must not leave a spurious empty line.
func TestWrapTextTrailingSpaceNoEmptyLine(t *testing.T) {
	got := wrapText("aaaaaaaa ", 4, runeWidth)
	want := []string{"aaaa", "aaaa"}
	if !equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestWrapTextEmpty(t *testing.T) {
	got := wrapText("", 10, runeWidth)
	if len(got) != 1 || got[0] != "" {
		t.Fatalf("expected one empty line, got %q", got)
	}
}

func equal(a, b []string) bool {
	return strings.Join(a, "\x00") == strings.Join(b, "\x00")
}
