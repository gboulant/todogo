package main

// TODO: encapsulate all that stuff of creating the tree folders and the
// initial not file

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"todogo/core"
)

// commandNote is the arguments parser of the command note
func commandNote(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var editIndex uint64
	flagset.Uint64Var(&editIndex, "e", 0, "Edit the note of the specified task")
	var viewIndex uint64
	flagset.Uint64Var(&viewIndex, "v", 0, "View the note of the specified task")
	var delIndex uint64
	flagset.Uint64Var(&delIndex, "d", 0, "Delete the note of the specified task")

	flagset.Parse(args)

	if editIndex != 0 {
		return editNote(editIndex)
	}
	if viewIndex != 0 {
		return viewNote(viewIndex)
	}

	flagset.Usage()
	return errors.New("ERR: Choose an option (see usage)")
}

func editNote(index uint64) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	task, err := journal.GetTask(index)
	if err != nil {
		return err
	}
	if task.NotePath == "" {
		task.InitNotePath()
	}

	var notepath string
	if filepath.IsAbs(task.NotePath) {
		notepath = task.NotePath
	} else {
		rootdir := filepath.Dir(journal.File())
		notepath = filepath.Join(rootdir, task.NotePath)
	}

	exists, err := core.PathExists(notepath)
	if exists && err != nil {
		return err
	}

	if !exists {
		err := core.CheckAndMakeDir(filepath.Dir(notepath))
		file, err := os.Create(notepath)
		defer file.Close()
		if err != nil {
			return err
		}
		title := fmt.Sprintf("%.2d - %s", task.UIndex, task.Description)
		line := ""
		for i := 0; i < len(title); i++ {
			line += "="
		}
		file.WriteString(fmt.Sprintf("%s\n", title))
		file.WriteString(fmt.Sprintf("%s\n", line))
		file.Sync()
	}

	fmt.Printf("The note of the task %d can be edited in file: %s\n", index, notepath)
	return journal.Save()
}

func viewNote(index uint64) error {
	journal, err := getActiveJournal()
	if err != nil {
		return err
	}
	task, err := journal.GetTask(index)
	if err != nil {
		return err
	}

	if task.NotePath == "" {
		return fmt.Errorf("ERR: the task %d has no associated note", index)
	}

	var notepath string
	if filepath.IsAbs(task.NotePath) {
		notepath = task.NotePath
	} else {
		rootdir := filepath.Dir(journal.File())
		notepath = filepath.Join(rootdir, task.NotePath)
	}

	_, err = core.PathExists(notepath)
	if err != nil {
		return err
	}

	file, err := os.Open(notepath)
	defer file.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return scanner.Err()
}
