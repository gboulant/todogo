package conf

// --------------------------------------------------------------
// Implementation of the Stringable interface of  config with pretty
// representations
// --------------------------------------------------------------

import (
	"fmt"
	"todogo/core"
)

// PrettyPrint indicates wether the printable string should be pretty or plain text
const PrettyPrint bool = false
const WithColor bool = true

// String implements the stringable interface for a Config
func (config Config) String() string {
	if PrettyPrint {
		return config.PrettyString()
	} else {
		return config.PlainString()
	}
}

// PlainString implements the stringable interface for a Config
func (config Config) PlainString() string {

	withcolor := func(s string) string {
		return core.ColorString(s, core.ColorMagenta)
	}
	uncolored := func(s string) string {
		return s
	}

	var colored func(s string) string
	if WithColor {
		colored = withcolor
	} else {
		colored = uncolored
	}

	s := "\n"
	for i := 0; i < len(config.ContextList); i++ {
		context := config.ContextList[i]
		if context.Name == config.ContextName {
			s += colored(fmt.Sprintf("* %s\n", context.String()))
		} else {
			s += colored(fmt.Sprintf("  %s\n", context.String()))
		}
	}
	s += fmt.Sprintf("\nLegend: %s\n", colored("* active context"))
	return s
}

// PrettyString is a variant of String for a pretty print of Config on standard output
func (config Config) PrettyString() string {
	s := "\n"
	for i := 0; i < len(config.ContextList); i++ {
		context := config.ContextList[i]
		if context.Name == config.ContextName {
			s += fmt.Sprintf("%s\n", core.ColorString(core.CharacterDisk+" "+context.String(), core.ColorMagenta))
		} else {
			s += fmt.Sprintf("  %s\n", context.String())
		}
	}
	s += fmt.Sprintf("\nLegend: %s", core.ColorString(core.CharacterDisk+" active context\n", core.ColorMagenta))
	return s
}
