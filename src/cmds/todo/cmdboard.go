package main

import (
	"flag"
	"fmt"

	"galuma.net/todo"
)

// commandBoard is the arguments parser of the command board
func commandBoard(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var clear bool
	flagset.BoolVar(&clear, "c", false, "Clear all the tasks from the board")
	var list bool
	flagset.BoolVar(&list, "l", false, "List all the tasks on board")
	var add todo.TaskIDArray
	flagset.Var(&add, "a", "Add on board the specified tasks (comma separeted list of indeces)")
	var remove todo.TaskIDArray
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
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	listing := journal.ListWithFilter(todo.TaskFilterOnBoard)
	fmt.Println(listing)
	return nil
}

func clearBoard() error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	tasksOnBoard := journal.GetTasksWithFilter(todo.TaskFilterOnBoard)
	for i := 0; i < len(tasksOnBoard); i++ {
		tasksOnBoard[i].OnBoard = false
	}
	err = journal.Save()
	if err != nil {
		return err
	}
	for i := 0; i < len(tasksOnBoard); i++ {
		fmt.Printf("The task of index %d has been removed from board\n", tasksOnBoard[i].UIndex)
	}
	return nil
}

func addOnBoard(indeces todo.TaskIDArray) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	for _, uindex := range indeces {
		err := journal.AddOnBoard(uindex)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Task of index %d has been added on board\n", uindex)
		}
	}
	return journal.Save()
}

func removeFromBoard(indeces todo.TaskIDArray) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	for _, uindex := range indeces {
		err := journal.RemoveFromBoard(uindex)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Task of index %d has been removed from board\n", uindex)
		}
	}
	return journal.Save()
}
