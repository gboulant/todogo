package main

import (
	"errors"
	"flag"
	"fmt"
	"todogo/core"
)

// commandConfig is the arguments parser of the command config
func commandConfig(cmdname string, args []string) error {
	flagset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var help string
	var newName string
	help = "Create (or update if exists) the context with name = string"
	flagset.StringVar(&newName, "n", "", help)

	var path string
	help = "path of the root directory of the context (to be used with option -n)"
	flagset.StringVar(&path, "p", "", help)

	var removeName string
	help = "Remove the context of name=string from the configuration"
	flagset.StringVar(&removeName, "r", "", help)

	var list bool
	help = "List the existing contexts (default option)"
	flagset.BoolVar(&list, "l", false, help)

	var selectName string
	help = "Select the context with name=string as the active context"
	flagset.StringVar(&selectName, "s", "", help)

	flagset.Parse(args)

	if newName != "" {
		if path == "" {
			msg := fmt.Sprintf("ERR: a path should be specified using the -p option")
			return errors.New(msg)
		} else {
			return createOrUptadeContext(newName, path)
		}
	}

	if selectName != "" {
		return selectContext(selectName)
	}

	if removeName != "" {
		return removeContext(removeName)
	}

	if len(flagset.Args()) > 0 {
		msg := fmt.Sprintf("ERR: the arguments %v are not valid", flagset.Args())
		return errors.New(msg)
	}

	return listContexts()
}

func listContexts() error {
	var handler core.ConfigHandler
	err := handler.Load()
	if err != nil {
		return err
	}
	handler.List()
	return nil
}

func createOrUptadeContext(name string, path string) error {
	var handler core.ConfigHandler
	err := handler.Load()
	if err != nil {
		return err
	}

	var pcontext *core.Context
	pcontext, err = handler.GetContext(name)
	if err != nil {
		// The context does not exists. To be created
		fmt.Printf("Creating the context %s with path %s\n", name, path)
		context := core.Context{
			Name:    name,
			DirPath: path,
		}
		pcontext = &context
		handler.AddContext(pcontext)
	} else {
		// The context exists. To be updated
		fmt.Printf("Updating the context %s with path %s\n", name, path)
		pcontext.DirPath = path
	}
	return handler.Save()
}

func selectContext(name string) error {
	var handler core.ConfigHandler
	err := handler.Load()
	if err != nil {
		return err
	}
	err = handler.SetActiveContext(name)
	if err != nil {
		return err
	}
	return handler.Save()
}

func removeContext(name string) error {
	var handler core.ConfigHandler
	err := handler.Load()
	if err != nil {
		return err
	}
	err = handler.RemoveContext(name)
	if err != nil {
		return err
	}
	return handler.Save()
}
