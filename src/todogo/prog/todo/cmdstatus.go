package main

import (
	"errors"
	"flag"
	"fmt"
	"todogo/core"
	"todogo/data"
)

// commandStatus is the arguments parser of the command status
func commandStatus(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var next core.IndexList
	flagset.Var(&next, "n", "Change to their next status the specified tasks (comma separated list of indices)")
	var prev core.IndexList
	flagset.Var(&prev, "p", "Change to their previous status the specified tasks (comma separated list of indices)")

	flagset.Parse(args)

	if len(next) > 0 {
		return modifyStatus(next, modifierNext)
	}
	if len(prev) > 0 {
		return modifyStatus(prev, modifierPrevious)
	}

	flagset.Usage()
	return errors.New("ERR: At least one option should be specified (-n or -p)")
}

type statusModifier func(task *data.Task) error

func modifierNext(task *data.Task) error {
	return task.Status.Next()
}

func modifierPrevious(task *data.Task) error {
	return task.Status.Next()
}

func modifyStatus(indeces core.IndexList, modifier statusModifier) error {
	var db data.Database
	db.Init(data.JournalPath)
	for _, index := range indeces {
		task, err := db.Get(index)
		if err != nil {
			fmt.Println(err)
		} else {
			err = modifier(&task)
			if err != nil {
				msg := fmt.Sprintf("WRN: the status of the task %d can not be changed (%s)", index, err)
				fmt.Println(msg)
			} else {
				db.Set(index, task)
				core.Println(task)
			}
		}
	}
	return db.Commit()
}
