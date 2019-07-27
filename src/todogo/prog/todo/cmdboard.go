package main

import (
	"errors"
	"flag"
	"fmt"
	"todogo/core"
	"todogo/data"
)

// commandBoard is the arguments parser of the command board
func commandBoard(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var clear bool
	flagset.BoolVar(&clear, "c", false, "Clear all the tasks from the board")
	var list bool
	flagset.BoolVar(&list, "l", false, "List all the tasks on board")
	var add core.IndexList
	flagset.Var(&add, "a", "Add on board the specified tasks (comma separeted list of indeces)")
	var remove core.IndexList
	flagset.Var(&remove, "r", "Remove from board the specified tasks (comma separeted list of indeces)")

	flagset.Parse(args)

	if list {
		return listBoard()
	}
	if clear {
		return clearBoard()
	}

	// At this point the list of indeces is specified
	if len(add) > 0 {
		return addOnBoard(add)
	}
	if len(remove) > 0 {
		return removeFromBoard(remove)
	}

	return listBoard()
}

func listBoard() error {
	var db data.Database
	db.Init(data.JournalPath)
	db.ListWithFilter(data.TaskFilterOnBoard)
	return nil
}

func clearBoard() error {
	var db data.Database
	db.Init(data.JournalPath)
	indeces := db.GetIndeces(data.TaskFilterOnBoard)
	if len(indeces) == 0 {
		return errors.New("WRN: Nothing to clear because there is no task on board")
	}

	for _, index := range indeces {
		err := db.RemoveFromBoard(index)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Task of index %d has been removed from board\n", index)
		}
	}
	return db.Commit()
}

func addOnBoard(indeces core.IndexList) error {
	var db data.Database
	db.Init(data.JournalPath)
	for _, index := range indeces {
		err := db.AddOnBoard(index)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Task of index %d has been added on board\n", index)
		}
	}
	return db.Commit()
}

func removeFromBoard(indeces core.IndexList) error {
	var db data.Database
	db.Init(data.JournalPath)
	for _, index := range indeces {
		err := db.RemoveFromBoard(index)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Task of index %d has been removed from board\n", index)
		}
	}
	return db.Commit()
}
