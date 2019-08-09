package main

import (
	"flag"
	"fmt"
	"todogo/data"
)

// commandArchive is the arguments parser of the command archive
func commandArchive(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var list bool
	flagset.BoolVar(&list, "l", false, "List the tasks of the archive")
	var add data.TaskIDArray
	flagset.Var(&add, "a", "Archive the specified tasks (comma separated list of indeces)")
	var restore data.TaskIDArray
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
	archive, err := getActiveArchive()
	if err != nil {
		return err
	}
	fmt.Println(archive.List())
	return nil
}

func moveToArchive(indeces data.TaskIDArray) error {
	archive, err := getActiveArchive()
	if err != nil {
		return err
	}

	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	for _, index := range indeces {
		task, err := journal.Delete(index)
		if err != nil {
			fmt.Println(err)
		} else {
			task.UIndex = task.GIndex
			err = archive.Add(task)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Task %d moved to the archive with a new usage index: %d\n", index, task.UIndex)
			}
		}
	}
	err = archive.Save()
	if err != nil {
		return err
	}
	return journal.Save()
}

func restoreFromArchive(indeces data.TaskIDArray) error {
	archive, err := getActiveArchive()
	if err != nil {
		return err
	}

	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	for _, index := range indeces {
		task, err := archive.Delete(index)
		if err != nil {
			fmt.Println(err)
		} else {
			task.UIndex = journal.GetFreeUID()
			err = journal.Add(task)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Task %d restored from archive with a new usage index: %d\n", index, task.UIndex)
			}
		}
	}
	err = archive.Save()
	if err != nil {
		return err
	}
	return journal.Save()
}
