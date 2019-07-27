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
			path = core.DefaultContextPath(newName)
			msg := fmt.Sprintf("WRN: You did't specify the context path. Default to %s", path)
			fmt.Println(msg)
		}
		return createOrUptadeContext(newName, path)
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
	fmt.Println(getconfig().String())
	return nil
}

func createOrUptadeContext(name string, path string) error {
	config := getconfig()
	context := config.GetContext(name)
	if context != nil {
		// The context exists. To be updated
		fmt.Printf("Updating the context %s with path %s\n", name, path)
		context.DirPath = path
	} else {
		// The context does not exists. To be created
		fmt.Printf("Creating the context %s with path %s\n", name, path)
		context := core.Context{
			Name:    name,
			DirPath: path,
		}
		config.AddContext(context)
	}
	// This context become the default context (could be a user choice)
	config.SetActiveContext(name)

	err := config.Save()
	if err == nil {
		fmt.Println(config.String())
	}
	return err
}

func selectContext(name string) error {
	config := getconfig()
	context := config.GetContext(name)
	if context == nil {
		return fmt.Errorf("ERR: the context %s does not exist", name)
	}
	config.SetActiveContext(name)

	err := config.Save()
	if err == nil {
		fmt.Println(config.String())
	}
	return err
}

func removeContext(name string) error {
	config := getconfig()
	err := config.RemoveContext(name)
	if err != nil {
		return err
	}

	err = config.Save()
	if err == nil {
		fmt.Println(config.String())
	}
	return err
}
