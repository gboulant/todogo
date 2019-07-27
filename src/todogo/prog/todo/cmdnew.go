package main

import (
	"errors"
	"flag"
	"todogo/core"
	"todogo/data"
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

	var db data.Database
	db.Init(data.JournalPath)
	task := db.New(text)
	err := db.Commit()
	if err != nil {
		return err
	}
	core.Println(task)
	return nil
}
