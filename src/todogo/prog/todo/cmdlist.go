package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	var printlist func(text string) (int, error)
	colorflag := config.Parameters.WithColor
	if filepath == "" {
		// Print listing on the standard output
		printlist = func(text string) (int, error) {
			return fmt.Println(text)
		}
	} else {
		// Print listing in a file. We add a header with the date, and
		// deactivate the color rendering
		printlist = func(text string) (int, error) {
			dlabel := time.Unix(time.Now().Unix(), 0).Format("Monday, January 2, 2006")
			header := fmt.Sprintf("TODO list at %s:", dlabel)
			underline := ""
			for range header {
				underline += "-"
			}
			header = fmt.Sprintf("%s\n%s\n", header, underline)
			return printfile(filepath, header+text)
		}
		config.Parameters.WithColor = false
	}

	var listing string
	if board {
		listing = journal.ListWithFilter(data.TaskFilterOnBoard)
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

// printfile prints the text content to a file with path fpath. The output format
// is a plain text. If fpath has a pdf extention then the content is printed
// also in a pdf file (in this case, the function results in two files: the plain
// text file and a pdf file)
func printfile(fpath string, text string) (int, error) {

	var txtpath string
	var pdfpath string
	pdf := false
	ext := filepath.Ext(fpath)
	if ext == ".pdf" {
		noxpath := strings.TrimSuffix(fpath, ext)
		txtpath = fmt.Sprintf("%s.txt", noxpath)
		pdfpath = fpath
		pdf = true
	} else {
		txtpath = fpath
	}

	fw, err := os.Create(txtpath)
	defer fw.Close()
	if err != nil {
		return 0, err
	}
	n, err := fw.WriteString(text)
	if err != nil {
		return 0, err
	}
	fmt.Printf("The todo list has been printed in the txt file: %s\n", txtpath)

	if pdf {
		err = txt2pdf(txtpath, pdfpath)
		if err != nil {
			err = fmt.Errorf("ERR: an error occurs during pdf creation: %s", err.Error())
		} else {
			fmt.Printf("The todo list has been printed in the pdf file: %s\n", pdfpath)
		}
	}

	return n, err
}

func txt2pdf(txtpath, pdfpath string) error {
	// To create a pdf file from a txt file, we use the external program
	// cupsfilter, available on every linux system. The typical usage is:
	//
	// $ cupsfilter <txtpath> -o cpi=14 -o lpi=8 -o landscape > <pdfpath>
	//
	cmd := exec.Command(
		"cupsfilter", txtpath,
		"-o", "cpi=14",
		"-o", "lpi=8",
		"-o", "landscape")

	fw, err := os.Create(pdfpath)
	defer fw.Close()
	if err != nil {
		return err
	}

	cmd.Stdout = fw
	cmd.Stderr = debugFilter{writer: os.Stderr}
	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

// debugFilter is an implementation of io.Writer used to filter the debug output
// of the function txt2pdf.
type debugFilter struct {
	writer io.Writer
}

// Write implements the io.Writer interface so that debugFilter is a Writer
func (f debugFilter) Write(p []byte) (int, error) {
	s := string(p[:])
	if strings.HasPrefix(s, "DEBUG: ") {
		// We return n as if the good number of characters was written but we
		// don't print the line in fact
		return len(p), nil
	}
	return f.writer.Write(p)
}
