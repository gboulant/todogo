package core

import (
	"fmt"
)

// FreeIndex returns the first free index of this list. The free index is
// determined with the hypothesis that the indeces array is a list of
// consecutive integer indeces. If the difference between two consecutif indeces
// is not 1, then it means that there is at least a free index (the index that
// follows the smallest index of the difference)
func FreeIndex(indeces []uint64) uint64 {
	if len(indeces) == 0 {
		return 1
	}
	var freeIndex uint64 = 1
	for i := 0; i < len(indeces); i++ {
		if indeces[i]-freeIndex > 0 {
			return freeIndex
		}
		freeIndex = indeces[i] + 1
	}
	return freeIndex
}

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
