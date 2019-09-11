package main

import (
	"errors"
	"flag"
	"fmt"
	"todogo/data"
)

// commandNew is the arguments command of the command new
func commandNew(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var text string
	flagset.StringVar(&text, "t", "", "text of the task")
	var parentUID data.TaskID
	flagset.Var(&parentUID, "p", "parent task (default is: no parent)")
	flagset.Parse(args)

	if text == "" {
		flagset.Usage()
		return errors.New("ERR: The text should be specified")
	}

	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	task := journal.New(text)

	if parentUID != data.NoUID {
		parentTask, err := journal.GetTask(parentUID)
		if err != nil {
			return err
		}
		task.ParentID = parentTask.UIndex
	}

	err = journal.Save()
	if err != nil {
		return err
	}
	fmt.Println(task.String())
	return nil
}
