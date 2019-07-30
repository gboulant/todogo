package conf

import (
	"fmt"
	"path/filepath"
	"sort"
)

const (
	// JournalFilename is the base name of the journal of a context
	JournalFilename = "journal.json"
	// ArchiveFilename is the base name of the archive of a context
	ArchiveFilename = "archive.json"
	// NotebookDirname is the base directory name of the notebook of a context
	NotebookDirname = "notes"

	// noIndex is used to specify an undefined ContextArray index
	noIndex = -1
)

// =========================================================================
// Implementation of the a context Context

// Context defines a workspace for a todo list
type Context struct {
	DirPath string
	Name    string
}

// String implements the stringable interface for a Context
func (context Context) String() string {
	return fmt.Sprintf("%-8s: %s", context.Name, context.absDirPath())
}

// absDirPath returns the absolute path to the root directory of this context.
// If DirPath is a relative path, then it is considered as relative to the
// configuration root directory. That is the case generally when the option
// -p was not specified at context creation. In such a case, the context
// workspace is created as a subdirectory of the configuration directory
// with name equal to the context name.
func (context Context) absDirPath() string {
	if filepath.IsAbs(context.DirPath) {
		return context.DirPath
	}
	return filepath.Join(cfgdirpath, context.DirPath)
}

// JournalPath returns the absolute path of the journal of this context
func (context Context) JournalPath() string {
	return filepath.Join(context.absDirPath(), JournalFilename)
}

// ArchivePath returns the absolute path of the archive of this context
func (context Context) ArchivePath() string {
	return filepath.Join(context.absDirPath(), ArchiveFilename)
}

// NotesPath returns the absolute path of the notes directory of this context
func (context Context) NotesPath() string {
	return filepath.Join(context.absDirPath(), NotebookDirname)
}

// =========================================================================
// Implementation of the collection of contexts ContextArray

// ContextArray defines a list of Context workspaces
type ContextArray []Context

// String implements the stringable interface for a ContextArray
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
