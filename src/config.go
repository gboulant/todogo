package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// configDirname is the directory name of the configuration (relative to user HOME)
	configDirname = ".config/galuma/todogo"
	// configFilename is the base name (relative to configDirname) of the configuration file
	configFilename = "config.json"

	// defaultContextName is the name (and relative path) of the default context
	defaultContextName = "default"

	// DefaultIndicatorsTemplate is the template for the task indicators string
	DefaultIndicatorsTemplate = "[{{.Date}}:{{.Note}}{{.Board}}]"
)

var (
	// cfgdirpath is the absolute path of the configuration directory
	cfgdirpath = filepath.Join(os.Getenv("HOME"), configDirname)
	// cfgfilepath is the absolute path to the configuration file
	cfgfilepath = filepath.Join(cfgdirpath, configFilename)

	// userConfig is a pointer to the current configuration (obtained using GetConfig)
	userConfig *Config
)

// DefaultContextPath returns a default path for this context name. The default
// context workspace is a subdirectory of the configuration root folder with a
// dirname equals to the context name.
func DefaultContextPath(name string) string {
	// Practically we set the default path to the context name itself so that it
	// is interpreted as a relative path (relative to the configuration root
	// directory.)
	return name
}

// =========================================================================
// Implementation of the configuration Config

// Parameters is a structure holding various configuration parameters in
// addition to the list of working contexts
type Parameters struct {
	// DefaultCommand specifies the default todo command (no options on the cmd line)
	DefaultCommand string
	// PrettyPrint indicates wether the printable string should be pretty or plain text
	PrettyPrint bool
	// WithColor indicates wether the printable string should be with color or not
	WithColor bool
	// Indicators is the template of indicators of a task string representation
	Indicators string
}

func (parameters Parameters) String() string {
	bytes, err := json.MarshalIndent(parameters, JSONPrefix, JSONIndent)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// Config defines the configuration of todo application. A configuration
// contains a list of contexts and the specification of the activze context.
type Config struct {
	ContextName string
	ContextList ContextArray
	Parameters  Parameters
	filepath    string // WRN: no jsonified (on purpose) because starts with minus letter
}

// defaultConfig() creates and returns a default configuration when no
// configuration exists. The default configuration is a configuration with one
// single context named "default".
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
			PrettyPrint:    true,
			WithColor:      true,
			Indicators:     DefaultIndicatorsTemplate,
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
	bytes, err := json.MarshalIndent(*config, JSONPrefix, JSONIndent)
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

// Save writes the config to the file it has been loaded from. The file path is
// given by the methode File.
func (config *Config) Save() error {
	return config.SaveTo(config.File())
}

// AddContext adds a new context in the configuration
func (config *Config) AddContext(context Context) error {
	return config.ContextList.append(context)
}

// GetContext returns the context whose name is the given name
func (config *Config) GetContext(name string) *Context {
	return config.ContextList.getContext(name)
}

// SetActiveContext selects the context with the given name as the active context
func (config *Config) SetActiveContext(name string) error {
	context := config.GetContext(name)
	if context == nil {
		return fmt.Errorf("ERR: The context %s does not exists", name)
	}
	config.ContextName = name
	return nil
}

// GetActiveContext returns the currently active context
func (config *Config) GetActiveContext() *Context {
	return config.ContextList.getContext(config.ContextName)
}

// RemoveContext removes the context of given name from the configuration.
// This operation does not delete the workspace associated to the context.
func (config *Config) RemoveContext(name string) error {
	if name == defaultContextName {
		return errors.New("ERR: The default context can not be removed")
	}
	context := config.GetContext(name)
	if context == nil {
		return fmt.Errorf("ERR: The context %s does not exists", name)
	}
	dir := context.DirPath
	idx := config.ContextList.indexFromName(name)
	err := config.ContextList.remove(idx)

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
// Implementation of the Stringable interface of  config with pretty
// representations
// --------------------------------------------------------------

var dotSymbolMap = map[bool]string{
	true:  PrettyDisk, // Pretty
	false: "*",        // Plain
}

type renderingFunction func(s string) string

var renderingFunctionMap = map[bool]renderingFunction{
	true:  func(s string) string { return ColorString(s, ColorMagenta) },
	false: func(s string) string { return s },
}

// String implements the stringable interface for a Config
func (config Config) String() string {
	return config.InfoString()
}

// ContextsString returns a pretty string representation of the list of contexts
func (config Config) ContextsString() string {
	symbol := dotSymbolMap[config.Parameters.PrettyPrint]
	clrfun := renderingFunctionMap[config.Parameters.WithColor]
	return config.createContextsString(symbol, clrfun)
}

func (config Config) createContextsString(dotSymbol string, renderingFunc renderingFunction) string {
	s := "\n"
	for i := 0; i < len(config.ContextList); i++ {
		context := config.ContextList[i]
		if context.Name == config.ContextName {
			s += renderingFunc(fmt.Sprintf("%s %s\n", dotSymbol, context.String()))
		} else {
			s += fmt.Sprintf("  %s\n", context.String())
		}
	}
	s += fmt.Sprintf("\nLegend: %s\n", renderingFunc(dotSymbol+" active context"))
	return s
}

// InfoString returns a plain string representation of this configuration
func (config Config) InfoString() string {
	s := "\n"

	// Configuration files
	s += "Configuration files:\n"
	s += "-------------------\n\n"
	s += fmt.Sprintf(" Configuration root directory: %s\n", cfgdirpath)
	s += fmt.Sprintf(" Configuration file path     : %s\n", cfgfilepath)
	s += "\n"

	// Configuration Parameters
	s += "Configuration parameters:\n"
	s += "------------------------\n\n"
	s += fmt.Sprintf("%s\n", config.Parameters.String())
	s += "\n"

	// List of contexts
	s += "List of contexts:\n"
	s += "----------------\n"
	s += config.createContextsString(dotSymbolMap[false], renderingFunctionMap[false])
	s += "\n"

	s += "!! NOTE: apart form the management of the contexts (create, remove, select),\n"
	s += "!! there is no todo command to edit the configuration (up to now). If you need\n"
	s += "!! to modify the configuration, you should directly edit the configuration file.\n"
	return s
}

// =========================================================================

// GetConfig returns the current configuration (and load it if first call)
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

// CreateTestConfig creates a test configuration (for test purposes)
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
