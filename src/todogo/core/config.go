package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var cfgdirpath = filepath.Join(os.Getenv("HOME"), ".todogo")
var cfgfilepath = filepath.Join(cfgdirpath, "config.json")

// Context defines a workspace for a todo list
type Context struct {
	DirPath string
	Name    string
}

// ContextArray defines a list of Context workspaces
type ContextArray []Context

// ContextMap defines a map of Context workspaces whose key is the context name
type ContextMap map[string]*Context

type Parameters struct {
	DefaultCommand string
}

// Config defines the configuration of todo application. A configuration
// contains a list of contexts and the specification of the activze context.
type Config struct {
	ContextName string
	ContextList ContextArray
	Parameters  Parameters
}

func initDefaultConfig() Config {
	config := Config{
		ContextName: "default",
		ContextList: ContextArray{
			{DirPath: cfgdirpath, Name: "default"},
		},
		Parameters: Parameters{
			DefaultCommand: "board",
		},
	}
	return config
}

func contextArray2Map(contextarray ContextArray) ContextMap {
	contextmap := make(ContextMap)
	var pcontext *Context
	for i := 0; i < len(contextarray); i++ {
		pcontext = &contextarray[i]
		contextmap[pcontext.Name] = pcontext
	}
	return contextmap
}

func contextMap2Array(contextmap ContextMap) ContextArray {
	contextarray := make(ContextArray, len(contextmap))
	icontext := 0
	for _, v := range contextmap {
		contextarray[icontext] = *(v)
		icontext++
	}
	return contextarray
}

// --------------------------------------------------------------

// ConfigHandler is a tool to manipulate the configuration (Load, Save, Edit)
type ConfigHandler struct {
	activeContextName string
	contextMap        ContextMap
}

// Load reads the Config from the configuration file
func (configHandler *ConfigHandler) Load() error {
	exists, err := PathExists(cfgfilepath)
	if exists && err != nil {
		return err
	}

	var config Config
	if !exists {
		config = initDefaultConfig()
		config.Save(cfgfilepath)
	} else {
		err = config.Load(cfgfilepath)
		if err != nil {
			return err
		}
	}
	configHandler.activeContextName = config.ContextName
	configHandler.contextMap = contextArray2Map(config.ContextList)
	return nil
}

// getConfig retuns a copy of the current configuration
func (configHandler ConfigHandler) getConfig() Config {
	var config Config
	config.ContextName = configHandler.activeContextName
	config.ContextList = contextMap2Array(configHandler.contextMap)
	return config
}

// List print the Config to the standard output
func (configHandler ConfigHandler) List() {
	config := configHandler.getConfig()
	config.Println()
}

// Save writes the Config to the configuration file
func (configHandler ConfigHandler) Save() error {
	config := configHandler.getConfig()
	return config.Save(cfgfilepath)
}

// AddContext adds the specified Context to the configuration
func (configHandler *ConfigHandler) AddContext(pcontext *Context) error {
	_, exists := configHandler.contextMap[pcontext.Name]
	if exists {
		msg := fmt.Sprintf("ERR: The context %s already exists", pcontext.Name)
		return errors.New(msg)
	}
	configHandler.contextMap[pcontext.Name] = pcontext
	return nil
}

// GetContext returns a pointer to the Context of the specified name
func (configHandler *ConfigHandler) GetContext(contextName string) (*Context, error) {
	pcontext, exists := configHandler.contextMap[contextName]
	if !exists {
		msg := fmt.Sprintf("ERR: The context %s does not exist", contextName)
		return nil, errors.New(msg)
	}
	return pcontext, nil
}

// RemoveContext remove the specified Context from the configuration. It does
// not delete the context workspace.
func (configHandler *ConfigHandler) RemoveContext(contextName string) error {
	if contextName == "default" {
		return errors.New("The default context can not be removed")
	}
	pconfig, err := configHandler.GetContext(contextName)
	if err != nil {
		return err
	}
	delete(configHandler.contextMap, contextName)
	fmt.Printf("The context %s has been removed from the configuration\n", pconfig.Name)
	fmt.Printf("The workspace still exists in folder: %s\n", pconfig.DirPath)

	// If the context was the active context, then we have to change the active
	// context. Set to default
	if pconfig.Name == configHandler.activeContextName {
		fmt.Println("The active context is reset to default")
		configHandler.activeContextName = "default"
	}
	return nil
}

// SetActiveContext defines the active context to the context whose name
// is contextName. The specified context should be defined in the configuration,
// otherwise an error is raised.
func (configHandler *ConfigHandler) SetActiveContext(contextName string) error {
	pcontext, err := configHandler.GetContext(contextName)
	if err != nil {
		return err
	}
	configHandler.activeContextName = contextName
	fmt.Printf("The active context has been set to %s (%s)\n", contextName, pcontext.DirPath)
	return nil
}

// GetActiveContext returns a pointer to the active Context
func (configHandler *ConfigHandler) GetActiveContext() (*Context, error) {
	return configHandler.GetContext(configHandler.activeContextName)
}

// GetActiveContextPath returns the path to the database of the default context
func GetActiveContextPath() string {
	var handler ConfigHandler
	handler.Load()
	pcontext, err := handler.GetActiveContext()
	if err != nil {
		panic(err)
	}
	return pcontext.DirPath
}

// --------------------------------------------------------------
// Implementation of the Jsonable interface

// Load reads a json file and map the data into a Config.
// It implements the jsonable interface.
func (config *Config) Load(filepath string) error {
	bytes, err := LoadBytes(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, config)
}

// Save writes the Config data into a json file.
// It implements the jsonable interface.
func (config *Config) Save(filepath string) error {
	bytes, err := json.MarshalIndent(*config, JsonPrefix, JsonIndent)
	if err != nil {
		return err
	}
	return WriteBytes(filepath, bytes)
}

// --------------------------------------------------------------
// Implementation of the Stringable interface

// String implements the stringable interface for a Context
func (context Context) String() string {
	return fmt.Sprintf("%-8s: %s", context.Name, context.DirPath)
}

// Println implements the stringable interface
func (context Context) Println() {
	fmt.Println(context.String())
}

// PrettyPrint indicates wether the printable string should be pretty or plain text
const PrettyPrint bool = true

// String implements the stringable interface for a Config
func (config Config) String() string {
	if PrettyPrint {
		return config.PrettyString()
	} else {
		return config.PlainString()
	}
}

// Println implements the stringable interface
func (config Config) Println() {
	fmt.Print(config.String())
}

// PlainString implements the stringable interface for a Config
func (config Config) PlainString() string {
	s := "\n"
	for i := 0; i < len(config.ContextList); i++ {
		context := config.ContextList[i]
		if context.Name == config.ContextName {
			s += fmt.Sprintf("* %s\n", context.String())
		} else {
			s += fmt.Sprintf("  %s\n", context.String())
		}
	}
	s += fmt.Sprint("\nLegend: * active context\n")
	return s
}

// PrettyString is a variant of String for a pretty print of Config on standard output
func (config Config) PrettyString() string {
	s := "\n"
	for i := 0; i < len(config.ContextList); i++ {
		context := config.ContextList[i]
		if context.Name == config.ContextName {
			s += fmt.Sprintf("%s\n", ColorString(CharacterDisk+" "+context.String(), ColorMagenta))
		} else {
			s += fmt.Sprintf("  %s\n", context.String())
		}
	}
	s += fmt.Sprintf("\nLegend: %s", ColorString(CharacterDisk+" active context\n", ColorMagenta))
	return s
}
