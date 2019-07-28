package main

import (
	"flag"
	"fmt"
	"todogo/core"
)

// commandChild is the arguments parser of the command child
func commandChild(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var children core.IndexList
	flagset.Var(&children, "c", "List of children tasks (comma separated list of indeces)")
	var parent uint64
	flagset.Uint64Var(&parent, "p", 0, "Index of the parent task")

	flagset.Parse(args)

	if len(children) > 0 {
		return addChildren(parent, children)
	}

	return nil
}

func addChildren(parentUID uint64, children core.IndexList) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	parent, err := journal.GetTask(parentUID)
	if err != nil {
		return err
	}

	for _, index := range children {
		child, err := journal.GetTask(index)
		if err != nil {
			fmt.Println(err)
		}
		child.ParentID = parent.UIndex
	}

	return journal.Save()
}
