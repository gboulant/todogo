package main

import (
	"errors"
	"flag"
	"fmt"
)

// commandNew is the arguments command of the command new
func commandNew(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var text string
	flagset.StringVar(&text, "t", "", "text of the task")
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
	err = journal.Save()
	if err != nil {
		return err
	}
	fmt.Println(task.String())
	return nil
}
