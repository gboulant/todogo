package main

import (
	"errors"
	"flag"
	"fmt"
	"todogo/core"
)

// commandDelete is the arguments parser of the command delete
func commandDelete(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var delete core.IndexList
	flagset.Var(&delete, "d", "Delete (definitively) the specified tasks (comma separated list of indeces)")

	var archive core.IndexList
	flagset.Var(&archive, "a", "Move to the archive the specified tasks (comma separated list of indeces)")

	flagset.Parse(args)

	if len(delete) > 0 {
		return deleteFromJournal(delete)
	}
	if len(archive) > 0 {
		return moveToArchive(archive)
	}

	flagset.Usage()
	return errors.New("ERR: At least one option should be specified (-d or -a)")
}

func deleteFromJournal(indeces core.IndexList) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	for _, index := range indeces {
		_, err := journal.Delete(index)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Task of index %d has been deleted\n", index)
		}
	}
	return journal.Save()
}
