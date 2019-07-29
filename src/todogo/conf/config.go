package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"todogo/core"
)

const (
	// configDirname is the directory name of the configuration (relative to user HOME)
	configDirname = ".todogo"
	// configFilename is the base name (relative to configDirname) of the configuration file
	configFilename = "config.json"

	// JournalFilename is the base name of the journal of a context
	JournalFilename = "journal.json"
	// ArchiveFilename is the base name of the archive of a context
	ArchiveFilename = "archive.json"
	// NotebookDirname is the base directory name of the notebook of a context
	NotebookDirname = "notes"
	// defaultContextName is the name (and relative path) of the default context
	defaultContextName = "default"
)

var (
	// cfgdirpath is the absolute path of the configuration directory
	cfgdirpath = filepath.Join(os.Getenv("HOME"), configDirname)
	// cfgfilepath is the absolute path to the configuration file
	cfgfilepath = filepath.Join(cfgdirpath, configFilename)

	// userConfig is a pointer to the current configuration (obtained using GetConfig)
	userConfig *Config
)

// DefaultContextPath returns a default path for this context name
func DefaultContextPath(name string) string {
	return filepath.Join(cfgdirpath, name)
}

// =========================================================================
// Implementation of the configuration Config

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

func defaultConfig() Config {
	config := Config{
		ContextName: defaultContextName,
		ContextList: ContextArray{
			{
				DirPath: DefaultContextPath(defaultContextName),
				Name:    defaultContextName,
			},
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
	bytes, err := core.LoadBytes(filepath)
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
	bytes, err := json.MarshalIndent(*config, core.JsonPrefix, core.JsonIndent)
	if err != nil {
		return err
	}
	err = core.WriteBytes(filepath, bytes)
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

// =========================================================================

// GetConfig returns the current configuration (and load it if first call)
func GetConfig() (*Config, error) {
	if userConfig != nil {
		return userConfig, nil
	}
	exists, err := core.PathExists(cfgfilepath)
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

// GetTestConfig creates a test configuration (for test purposes)
func CreateTestConfig() Config {
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
