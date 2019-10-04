package main

import (
	"flag"
	"fmt"
	"time"
	"todogo/conf"
	"todogo/data"
)

// commandList is the arguments parser of the command list
func commandList(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var board bool
	flagset.BoolVar(&board, "b", false, "List only the tasks on board")
	var tree bool
	flagset.BoolVar(&tree, "t", false, "Tree representation of parent relations")
	var report bool
	flagset.BoolVar(&report, "r", false, "List a complete report (list, board, and notes")

	var filepath string
	flagset.StringVar(&filepath, "f", "", "Print the listing in the specified file")

	flagset.Parse(args)

	journal, err := getActiveJournal()
	if err != nil {
		return err
	}

	config, err := conf.GetConfig()
	if err != nil {
		return err
	}

	// If the output is a file, then we deactivate temporarely the color
	// rendering
	var printlist printer
	colorflag := config.Parameters.WithColor
	if filepath == "" {
		// Print listing on the standard output
		printlist = stdOutPrinter()
	} else {
		// Print listing in a file. We add a header with the date, and
		// deactivate the color rendering
		printlist = filePrinter(filepath)
		config.Parameters.WithColor = false
	}

	var listing string
	if board {
		listing = journal.ListWithFilter(data.TaskFilterOnBoard)
	} else if report {
		listing = journal.Report()
	} else {
		if tree {
			listing = journal.Tree()
		} else {
			listing = journal.List()
		}
	}

	_, err = printlist(listing)
	config.Parameters.WithColor = colorflag
	return err
}

type printer func(text string) (int, error)

func stdOutPrinter() printer {
	return func(text string) (int, error) {
		return fmt.Println(text)
	}
}

func filePrinter(filepath string) printer {
	return func(text string) (int, error) {
		dlabel := time.Unix(time.Now().Unix(), 0).Format("Monday, January 2, 2006")
		header := fmt.Sprintf("TODO list at %s:", dlabel)
		underline := ""
		for range header {
			underline += "-"
		}
		header = fmt.Sprintf("%s\n%s\n", header, underline)
		return printfile(filepath, header+text)
	}
}
