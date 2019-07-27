package main

import (
	"flag"
	"todogo/data"
)

// commandList is the arguments parser of the command list
func commandList(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var board bool
	flagset.BoolVar(&board, "b", false, "List only the tasks on board")

	flagset.Parse(args)

	var db data.Database
	db.Init(getconfig().GetActiveContext().JournalPath())
	if board {
		db.ListWithFilter(data.TaskFilterOnBoard)
	} else {
		db.List()
	}
	return nil
}
