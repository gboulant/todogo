package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
