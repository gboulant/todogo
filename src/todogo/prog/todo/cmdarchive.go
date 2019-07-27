package main

import (
	"flag"
	"fmt"
	"todogo/core"
	"todogo/data"
)

// commandArchive is the arguments parser of the command archive
func commandArchive(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var list bool
	flagset.BoolVar(&list, "l", false, "List the tasks of the archive")
	var add core.IndexList
	flagset.Var(&add, "a", "Archive the specified tasks (comma separated list of indeces)")
	var restore core.IndexList
	flagset.Var(&restore, "r", "Restore the specified tasks (comma separated list of indeces)")

	flagset.Parse(args)

	if list {
		return listArchive()
	}
	if len(add) > 0 {
		return moveToArchive(add)
	}
	if len(restore) > 0 {
		return restoreFromArchive(restore)
	}

	return listArchive()
}

func listArchive() error {
	var dba data.Database
	dba.Init(getconfig().GetActiveContext().ArchivePath())
	dba.List()
	return nil
}

func moveToArchive(indeces core.IndexList) error {
	var dba data.Database
	dba.Init(getconfig().GetActiveContext().ArchivePath())

	var dbj data.Database
	dbj.Init(getconfig().GetActiveContext().JournalPath())

	for _, index := range indeces {
		task, err := dbj.Delete(index)
		if err != nil {
			fmt.Println(err)
		} else {
			task.UIndex = task.GIndex
			err = dba.Add(task)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Task %d moved to the archive with a new usage index: %d\n", index, task.UIndex)
			}
		}
	}
	err := dba.Commit()
	if err != nil {
		return err
	}
	return dbj.Commit()
}

func restoreFromArchive(indeces core.IndexList) error {
	var dba data.Database
	dba.Init(getconfig().GetActiveContext().ArchivePath())

	var dbj data.Database
	dbj.Init(getconfig().GetActiveContext().JournalPath())

	for _, index := range indeces {
		task, err := dba.Delete(index)
		if err != nil {
			fmt.Println(err)
		} else {
			task.UIndex = dbj.FreeUsageIndex()
			err = dbj.Add(task)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Task %d restored from archive with a new usage index: %d\n", index, task.UIndex)
			}
		}
	}

	err := dba.Commit()
	if err != nil {
		return err
	}
	return dbj.Commit()
}
