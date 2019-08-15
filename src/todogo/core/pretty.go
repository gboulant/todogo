package core

import (
	"fmt"
)

// ColorIndex defines a color index
type ColorIndex int

// Enumeration of possible ColorIndex
const (
	ColorRed ColorIndex = iota + 31
	ColorGreen
	ColorOrange
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// ColorString returns the text with color when printed on standard output
func ColorString(text string, color ColorIndex) string {
	return fmt.Sprintf("\033[1;%dm%s\033[1;0m", color, text)
}

// Enumeration of possible pretty characters
const (
	PrettyDiskVoid string = "\u25cb"
	PrettyDiskHalf string = "\u25d0"
	PrettyDiskFull string = "\u25cf"
	PrettyDisk     string = PrettyDiskFull

	PrettyTriangleRight string = "\u25b6"
	PrettyCross         string = "\u274c"

	PrettyArrowRightShort string = "\u27f6"
	PrettyArrowRightLong  string = "\u27ff"
	PrettyArrowRight      string = PrettyArrowRightShort
)
