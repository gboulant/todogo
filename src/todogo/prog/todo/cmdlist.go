package main

import (
	"flag"
	"fmt"
	"todogo/data"
)

// commandList is the arguments parser of the command list
func commandList(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var board bool
	flagset.BoolVar(&board, "b", false, "List only the tasks on board")

	flagset.Parse(args)

	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	var listing string
	if board {
		listing = journal.ListWithFilter(data.TaskFilterOnBoard)
	} else {
		listing = journal.List()
	}
	fmt.Println(listing)
	return nil
}
