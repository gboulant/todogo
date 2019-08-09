package main

import (
	"flag"
	"fmt"
	"todogo/data"
)

// commandChild is the arguments parser of the command child
func commandChild(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var children data.TaskIDArray
	flagset.Var(&children, "c", "List of children tasks (comma separated list of indeces)")
	var parent data.TaskID
	flagset.Var(&parent, "p", "Index of the parent task")

	flagset.Parse(args)

	if len(children) > 0 {
		return addChildren(parent, children)
	}

	return nil
}

func addChildren(parentUID data.TaskID, children data.TaskIDArray) error {
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
