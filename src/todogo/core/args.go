package core

import (
	"errors"
	"fmt"
	"os"
)

type parserFunc func(cmdname string, args []string) error

// Command defines a subcommand of the program. A subcommand is defined by a Name
// (used on the command line to invoke the subcommand), a Description (used for
// printing the usage of the subcommand) and a Parser (that implements the
// processing of the command line arguments, i.e. parsing of arguments and
// execution of the requested subcommand with the given arguments).
type Command struct {
	Name        string
	Description string
	Parser      parserFunc
}

// CommandList is an array of Command
type CommandList []Command

// CommandParser is the general manager of the command line. It gathers the
// progname (name of the program), and the list of the command (list of the
// subcommand of the program)
//of
type CommandParser struct {
	// progname is the name of the executable program
	progname string
	// defaultopt is the list of command line options to consider if no option is given
	defaultopt []string
	// commandList is the list of available commands
	commandList CommandList
}

// Init initialises a CommandParser object
func (commandParser *CommandParser) Init(progname string, commandList CommandList) {
	commandParser.progname = progname
	commandParser.defaultopt = []string{}
	commandParser.commandList = commandList
}

// NewCommandParser creates and initialise a CommandParser
func NewCommandParser(progname string, commandList CommandList) CommandParser {
	commandParser := CommandParser{}
	commandParser.Init(progname, commandList)
	return commandParser
}

// SetDefaultCmdOptions specifies the default command line options
func (commandParser *CommandParser) SetDefaultCmdOptions(cmdopt []string) {
	commandParser.defaultopt = cmdopt
}

// commandNames returns a list of possible command names (from commandList)
func (commandParser CommandParser) commandNames() []string {
	names := make([]string, len(commandParser.commandList))
	for i := 0; i < len(commandParser.commandList); i++ {
		names[i] = commandParser.commandList[i].Name
	}
	return names
}

// usage prints the main usage of the standard output
func (commandParser CommandParser) usage() {
	fmt.Printf("usage: %s <command> [<options>] [<arguments>]\n", commandParser.progname)
	fmt.Printf("\nWith <command> in:\n\n")
	for i := 0; i < len(commandParser.commandList); i++ {
		cmd := commandParser.commandList[i]
		fmt.Printf("* %-10s: %s\n", cmd.Name, cmd.Description)
	}
	fmt.Printf("\nFor a description of possible options, try: %s <command> --help\n", commandParser.progname)
}

// getCommand returns the Command whose name is commandName
func (commandParser CommandParser) getCommand(commandName string) (Command, error) {
	for i := 0; i < len(commandParser.commandList); i++ {
		if commandParser.commandList[i].Name == commandName {
			return commandParser.commandList[i], nil
		}
	}
	var nilcmd = Command{}
	msg := fmt.Sprintf("ERR: the command %s is not defined (should be in: %s)", commandName, commandParser.commandNames())
	err := errors.New(msg)
	return nilcmd, err
}

var helpOptions = []string{"-help", "--help", "-h"}

// ArgParse parses the command line arguments and executed the requested command
func (commandParser CommandParser) ArgParse() error {
	var args []string

	if len(os.Args) < 2 {
		if len(commandParser.defaultopt) > 0 {
			args = make([]string, 1+len(commandParser.defaultopt))
			args[0] = os.Args[0]
			for i := 1; i <= len(commandParser.defaultopt); i++ {
				args[i] = commandParser.defaultopt[i-1]
			}
		} else {
			commandParser.usage()
			msg := fmt.Sprintf("ERR: you should specify a command in: %s", commandParser.commandNames())
			return errors.New(msg)
		}
	} else {
		args = os.Args
	}

	if contains(helpOptions, args[1]) {
		commandParser.usage()
		return nil
	}

	commandName := args[1]
	command, err := commandParser.getCommand(commandName)
	if err != nil {
		return err
	}
	return command.Parser(commandName, args[2:])
}

// -----------------------------------------------------------------------

// contains return true if the array contains the item
func contains(sarray []string, sitem string) bool {
	for i := 0; i < len(sarray); i++ {
		if sarray[i] == sitem {
			return true
		}
	}
	return false
}
