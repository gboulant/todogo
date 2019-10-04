package main

import (
	"fmt"
	"os"
	"strings"
	"todogo/conf"
	"todogo/core"
)

var commands = core.CommandList{
	{Name: "new", Description: "Create a new task", Parser: commandNew},
	{Name: "list", Description: "Print the list of tasks", Parser: commandList},
	{Name: "status", Description: "Change the status of tasks", Parser: commandStatus},
	{Name: "board", Description: "Append/Remove tasks on/from the board", Parser: commandBoard},
	{Name: "note", Description: "Edit/View the note associated to a task", Parser: commandNote},
	{Name: "child", Description: "Make tasks be children of a parent task", Parser: commandChild},
	{Name: "delete", Description: "Delete tasks (definitely or in archive)", Parser: commandDelete},
	{Name: "archive", Description: "Archive/Restore tasks", Parser: commandArchive},
	{Name: "config", Description: "Manage de configuration", Parser: commandConfig},
}

func getConfig() *conf.Config {
	cfg, err := conf.GetConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfg
}

func main() {
	app := core.NewCommandParser("todo", commands)
	app.SetDefaultCmdOptions(strings.Fields(getConfig().Parameters.DefaultCommand))
	err := app.ArgParse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
