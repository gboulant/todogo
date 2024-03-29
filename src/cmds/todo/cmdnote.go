package main

import (
	"errors"
	"flag"
	"fmt"

	"galuma.net/todo"
)

// commandNote is the arguments parser of the command note
func commandNote(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var editIndex todo.TaskID
	flagset.Var(&editIndex, "e", "Edit the note of the specified task")
	var viewIndex todo.TaskID
	flagset.Var(&viewIndex, "v", "View the note of the specified task")
	var delIndex todo.TaskID
	flagset.Var(&delIndex, "d", "Delete the note of the specified task")

	flagset.Parse(args)

	if editIndex != 0 {
		return editNote(editIndex)
	}
	if viewIndex != 0 {
		return viewNote(viewIndex)
	}
	if delIndex != 0 {
		return deleteNote(delIndex)
	}

	flagset.Usage()
	return errors.New("ERR: Choose an option (see usage)")
}

func editNote(index todo.TaskID) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	notepath, err := journal.GetOrCreateNoteFile(index)
	if err != nil {
		return err
	}

	fmt.Printf("The note of the task %d can be edited in file: %s\n", index, notepath)
	return journal.Save()
}

func viewNote(index todo.TaskID) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	notepath, err := journal.GetNoteFile(index)
	if err != nil {
		return err
	}

	if notepath == "" {
		return fmt.Errorf("ERR: the task %d has no associated note", index)
	}

	content, err := todo.LoadString(notepath)
	if err == nil {
		fmt.Println(content)
	}
	return err
}

func deleteNote(index todo.TaskID) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	err = journal.DeleteNoteFile(index)
	if err != nil {
		return err
	}
	return journal.Save()
}
