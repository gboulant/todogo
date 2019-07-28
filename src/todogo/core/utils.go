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

// CharacterDisk displays a little disk when printed on standard output
const CharacterDisk string = "\u23FA"
