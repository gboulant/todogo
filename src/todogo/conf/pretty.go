package conf

// --------------------------------------------------------------
// Implementation of the Stringable interface of  config with pretty
// representations
// --------------------------------------------------------------

import (
	"fmt"
	"todogo/core"
)

var dotSymbolMap = map[bool]string{
	true:  core.PrettyDisk, // Pretty
	false: "*",             // Plain
}

type colorFunction func(s string) string

var colorFunctionMap = map[bool]colorFunction{
	true:  func(s string) string { return core.ColorString(s, core.ColorMagenta) },
	false: func(s string) string { return s },
}

// String implements the stringable interface for a Config
func (config Config) String() string {
	symbol := dotSymbolMap[PrettyPrint]
	clrfun := colorFunctionMap[WithColor]
	return config.createString(symbol, clrfun)
}

func (config Config) createString(dotSymbol string, coloredstr colorFunction) string {
	s := "\n"
	for i := 0; i < len(config.ContextList); i++ {
		context := config.ContextList[i]
		if context.Name == config.ContextName {
			s += coloredstr(fmt.Sprintf("%s %s\n", dotSymbol, context.String()))
		} else {
			s += fmt.Sprintf("  %s\n", context.String())
		}
	}
	s += fmt.Sprintf("\nLegend: %s\n", coloredstr(dotSymbol+" active context"))
	return s
}
