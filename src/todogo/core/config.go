package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var cfgdirpath = filepath.Join(os.Getenv("HOME"), ".todogo")
var cfgfilepath = filepath.Join(cfgdirpath, "config.json")

// Context defines a workspace for a todo list
type Context struct {
	DirPath string
	Name    string
}

const (
	JournalFilename = "journal.json"
	ArchiveFilename = "archive.json"
	NotebookDirname = "notes"
)

// String implements the stringable interface for a Context
func (context Context) String() string {
	return fmt.Sprintf("%-8s: %s", context.Name, context.DirPath)
}
func (context Context) JournalPath() string {
	return filepath.Join(context.DirPath, JournalFilename)
}
func (context Context) ArchivePath() string {
	return filepath.Join(context.DirPath, ArchiveFilename)
}
func (context Context) NotesPath() string {
	return filepath.Join(context.DirPath, NotebookDirname)
}

// ContextArray defines a list of Context workspaces
type ContextArray []Context

// String implements the stringable interface for a Context
func (contexts ContextArray) String() string {
	s := ""
	for i := 0; i < len(contexts); i++ {
		s += fmt.Sprintf("%s\n", contexts[i].String())
	}
	return s
}

func (contexts *ContextArray) Remove(index int) error {
	if index < 0 || index >= len(*contexts) {
		return fmt.Errorf("ERR: index %d is out of range of contexts", index)
	}
	(*contexts)[index] = (*contexts)[len(*contexts)-1]
	*contexts = (*contexts)[:len(*contexts)-1]
	return nil
}

func (contexts *ContextArray) Append(context Context) error {
	if contexts.GetContext(context.Name) != nil {
		return fmt.Errorf("ERR: a context with name %s already exists", context.Name)
	}
	*contexts = append(*contexts, context)
	return nil
}

func (contexts *ContextArray) SortByName() {
	byName := func(i int, j int) bool {
		return (*contexts)[i].Name < (*contexts)[j].Name
	}
	sort.Slice(*contexts, byName)
}

const noIndex = -1

type filterFunction func(context Context) bool

func (contexts ContextArray) index(filter filterFunction) int {
	for i := 0; i < len(contexts); i++ {
		if filter(contexts[i]) {
			return i
		}
	}
	return noIndex
}

func (contexts ContextArray) IndexFromName(name string) int {
	filter := func(context Context) bool {
		return context.Name == name
	}
	return contexts.index(filter)
}

func (contexts *ContextArray) GetContext(name string) *Context {
	idx := contexts.IndexFromName(name)
	if idx == noIndex {
		return nil
	}
	return &(*contexts)[idx]
}

func GetTestConfig() Config {
	config := Config{
		ContextName: "toto",
		ContextList: ContextArray{
			Context{DirPath: "/tmp/toto", Name: "toto"},
			Context{DirPath: "/tmp/tutu", Name: "tutu"},
			Context{DirPath: "/tmp/titi", Name: "titi"},
			Context{DirPath: "/tmp/tata", Name: "tata"},
		},
		Parameters: Parameters{
			DefaultCommand: "board",
		},
	}
	return config
}

// Parameters is a structure holding various configuration parameters in
// addition to the list of working contexts
type Parameters struct {
	DefaultCommand string
}

// Config defines the configuration of todo application. A configuration
// contains a list of contexts and the specification of the activze context.
type Config struct {
	ContextName string
	ContextList ContextArray
	Parameters  Parameters
	filepath    string // WRN: no jsonified (on purpose) because starts with minus letter
}

const defaultContextName = "default"

var defaultContextPath = cfgdirpath

func defaultConfig() Config {
	config := Config{
		ContextName: defaultContextName,
		ContextList: ContextArray{
			{DirPath: defaultContextPath, Name: defaultContextName},
		},
		Parameters: Parameters{
			DefaultCommand: "board",
		},
	}
	return config
}

// Load reads a json file and map the data into a Config.
// It implements the jsonable interface.
func (config *Config) Load(filepath string) error {
	bytes, err := LoadBytes(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	config.filepath = filepath
	return nil
}

// SaveTo writes the Config data into a json file.
// It implements the jsonable interface.
func (config *Config) SaveTo(filepath string) error {
	bytes, err := json.MarshalIndent(*config, JsonPrefix, JsonIndent)
	if err != nil {
		return err
	}
	err = WriteBytes(filepath, bytes)
	if err != nil {
		return err
	}
	config.filepath = filepath
	return nil
}

// File returns the source file (if Load was used)
// It implements the jsonable interface.
func (config *Config) File() string {
	return config.filepath
}

func (config *Config) Save() error {
	return config.SaveTo(config.File())
}

func (config *Config) AddContext(context Context) error {
	return config.ContextList.Append(context)
}

func (config *Config) GetContext(name string) *Context {
	return config.ContextList.GetContext(name)
}

func (config *Config) SetActiveContext(name string) error {
	context := config.GetContext(name)
	if context == nil {
		return fmt.Errorf("ERR: The context %s does not exists", name)
	}
	config.ContextName = name
	return nil
}

func (config *Config) GetActiveContext() *Context {
	return config.ContextList.GetContext(config.ContextName)
}

func (config *Config) RemoveContext(name string) error {
	if name == defaultContextName {
		return errors.New("ERR: The default context can not be removed")
	}
	context := config.GetContext(name)
	if context == nil {
		return fmt.Errorf("ERR: The context %s does not exists", name)
	}
	dir := context.DirPath
	idx := config.ContextList.IndexFromName(name)
	err := config.ContextList.Remove(idx)

	if err != nil {
		return err
	}

	fmt.Printf("The context %s has been removed from the configuration\n", name)
	fmt.Printf("The workspace still exists in folder: %s\n", dir)

	// If the context was the active context, then we have to change the active
	// context. Reset to default
	if name == config.ContextName {
		fmt.Println("The active context is reset to default")
		config.ContextName = defaultContextName
	}
	return nil
}

// --------------------------------------------------------------
// Implementation of the Stringable interface

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

// =========================================================================

var userConfig *Config

func GetConfig() (*Config, error) {
	if userConfig != nil {
		return userConfig, nil
	}
	exists, err := PathExists(cfgfilepath)
	if exists && err != nil {
		return nil, err
	}

	var config Config
	if !exists {
		config = defaultConfig()
		config.SaveTo(cfgfilepath)
	} else {
		err = config.Load(cfgfilepath)
		if err != nil {
			return nil, err
		}
	}
	userConfig = &config
	return userConfig, nil
}
