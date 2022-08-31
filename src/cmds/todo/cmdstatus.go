package main

import (
	"errors"
	"flag"
	"fmt"

	"galuma.net/todo"
)

// commandStatus is the arguments parser of the command status
func commandStatus(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var next todo.TaskIDArray
	flagset.Var(&next, "n", "Change to their next status the specified tasks (comma separated list of indices)")
	var prev todo.TaskIDArray
	flagset.Var(&prev, "p", "Change to their previous status the specified tasks (comma separated list of indices)")
	var info todo.TaskIDArray
	flagset.Var(&info, "i", "Display the complete status of the specified tasks (comma separated list of indices)")

	flagset.Parse(args)

	if len(next) > 0 {
		return modifyStatus(next, modifierNext)
	}
	if len(prev) > 0 {
		return modifyStatus(prev, modifierPrevious)
	}
	if len(info) > 0 {
		return infoStatus(info)
	}

	flagset.Usage()
	return errors.New("ERR: At least one option should be specified (-n or -p)")
}

type statusModifier func(task *todo.Task) error

func modifierNext(task *todo.Task) error {
	return task.Status.Next()
}

func modifierPrevious(task *todo.Task) error {
	return task.Status.Previous()
}

func modifyStatus(indeces todo.TaskIDArray, modifier statusModifier) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	for _, index := range indeces {
		task, err := journal.GetTask(index)
		if err != nil {
			fmt.Println(err)
		} else {
			err = modifier(task)
			if err != nil {
				msg := fmt.Sprintf("WRN: the status of the task %d can not be changed (%s)", index, err)
				fmt.Println(msg)
			} else {
				fmt.Println(task.String())
			}
		}
	}
	return journal.Save()
}

func infoStatus(indeces todo.TaskIDArray) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	fmt.Println()
	for _, uindex := range indeces {
		info, err := journal.GetTaskInfo(uindex)
		if err != nil {
			msg := fmt.Sprintf("WRN: no info for the task %d (%s)", uindex, err)
			fmt.Println(msg)
		} else {
			fmt.Println(info)
		}
		fmt.Println()
	}
	return nil
}
