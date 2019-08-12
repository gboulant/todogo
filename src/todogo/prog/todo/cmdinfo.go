package main

import (
	"flag"
	"fmt"
	"todogo/data"
)

// commandNote is the arguments parser of the command note
func commandInfo(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var taskid data.TaskID
	help := "Print detailled information concerning this task (value = task ID)"
	flagset.Var(&taskid, "t", help)

	flagset.Parse(args)

	if taskid != data.NoUID {
		return printTaskInfo(taskid)
	}

	flagset.Usage()
	return nil
}

func printTaskInfo(uindex data.TaskID) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	info, err := journal.GetTaskInfo(uindex)
	if err != nil {
		return err
	}
	fmt.Println(info)

	return nil
}
