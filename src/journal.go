package todo

// Implementation of a tasks journal (with file persistence)

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"
)

// TaskJournal defines the structure to manage a task journal. A tasks journal
// could be the current collection of tasks (called journal) or the archive
// collection of tasks (called archive).
type TaskJournal struct {
	TaskList TaskArray
	filepath string
}

// =========================================================================
// Implementation of the edition functions

// New creates a new task in the database
func (journal *TaskJournal) New(text string) *Task {
	uindex := journal.TaskList.getFreeUID()
	var task = Task{
		UIndex:      uindex,
		Description: text,
		Timestamp:   timestamp(),
		Status:      StatusTodo,
		OnBoard:     false,
	}
	task.initGlobalIndex()
	journal.TaskList.append(task)
	ptask, _ := journal.GetTask(task.UIndex)
	return ptask
}

// Add adds the given task to this journal
func (journal *TaskJournal) Add(task Task) error {
	return journal.TaskList.append(task)
}

// Delete removes the task with the specified id. Returns a copy of the deleted
// task on success
func (journal *TaskJournal) Delete(uindex TaskID) (Task, error) {
	var task Task
	index := journal.TaskList.indexFromUID(uindex)
	if index == noIndex {
		return task, fmt.Errorf("ERR: The task %d does not exist", uindex)
	}
	task = journal.TaskList[index]
	err := journal.TaskList.remove(index)
	return task, err
}

// GetTask returns a pointer to the task whose usage ID is uindex
func (journal TaskJournal) GetTask(uindex TaskID) (*Task, error) {
	return journal.TaskList.getTask(uindex)
}

// GetTaskInfo returns a string representation of the charateristics of this
// task. This is the journal that creates this string and not the task itself,
// so that we can use absolute paths
func (journal TaskJournal) GetTaskInfo(uindex TaskID) (string, error) {
	task, err := journal.TaskList.getTask(uindex)
	if err != nil {
		return "", err
	}

	datelabel := func(timestamp int64) string {
		return time.Unix(timestamp, 0).Format("Monday 2006-January-02 at 15:04:05")
	}

	notepath, _ := journal.getNoteFile(uindex, false)

	s := ""
	s += fmt.Sprintf("Task               : %s\n", task.Description)
	s += fmt.Sprintf("Usage Index  (UID) : %d\n", task.UIndex)
	s += fmt.Sprintf("Global Index (GID) : %d\n", task.GIndex)
	s += fmt.Sprintf("Creation Date      : %s\n", datelabel(task.Timestamp))
	s += fmt.Sprintf("Status             : %s\n", task.Status.Label())
	s += fmt.Sprintf("Is on board        : %v\n", task.OnBoard)
	s += fmt.Sprintf("Note filepath      : %s\n", notepath)
	s += fmt.Sprintf("Parent UID         : %d", task.ParentID)

	return s, nil
}

// GetTasksWithFilter returns an array of pointer to the tasks that satisfy the
// given filter.
func (journal TaskJournal) GetTasksWithFilter(filter TaskFilter) []*Task {
	return journal.TaskList.getTasksWithFilter(filter)
}

// GetFreeUID returns the next free usage index in this journal
func (journal TaskJournal) GetFreeUID() TaskID {
	return journal.TaskList.getFreeUID()
}

// AddOnBoard adds the specified task on board
func (journal *TaskJournal) AddOnBoard(uindex TaskID) error {
	task, err := journal.TaskList.getTask(uindex)
	if err != nil {
		return err
	}
	task.OnBoard = true
	return nil
}

// RemoveFromBoard removes the specified task from board
func (journal *TaskJournal) RemoveFromBoard(uindex TaskID) error {
	task, err := journal.TaskList.getTask(uindex)
	if err != nil {
		return err
	}
	task.OnBoard = false
	return nil
}

// =========================================================================
// Implementation of the serialization functions

// Load reads a journal of tasks from the given file. Returns an error if the
// file does not exist. Use LoadOrCreate to make sure to initialise a joournal
// whatever the starting situation (inn the case of the first usage of todo for
// example). It implements the jsonable interface.
func (journal *TaskJournal) Load(filepath string) error {
	bytes, err := LoadBytes(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, journal)
	if err == nil {
		journal.filepath = filepath
	}
	return err

}

// LoadOrCreate tries to load a journal from the given file, and create a void
// journal if the file does not exist.
func (journal *TaskJournal) LoadOrCreate(filepath string) error {
	exists, err := PathExists(filepath)
	if exists && err != nil {
		return err
	}

	if !exists {
		journal.TaskList = make(TaskArray, 0)
	} else {
		err = journal.Load(filepath)
		if err != nil {
			return err
		}
	}
	journal.filepath = filepath
	return nil
}

// SaveTo writes the journal data to the given file.
// It implements the jsonable interface.
func (journal *TaskJournal) SaveTo(filepath string) error {
	bytes, err := json.MarshalIndent(journal, JSONPrefix, JSONIndent)
	if err != nil {
		return err
	}
	err = WriteBytes(filepath, bytes)
	if err == nil {
		journal.filepath = filepath
	}
	return err
}

// File returns the persistance filepath (if journal is created by Load)
func (journal *TaskJournal) File() string {
	return journal.filepath
}

// Save writes the journal data to the persistence file
func (journal *TaskJournal) Save() error {
	return journal.SaveTo(journal.File())
}

// =========================================================================
// Implementation of the stringable function (function creating string
// representations of a journal)

const notasks = "No tasks. Go have a drink"

// ListWithFilter returns a string representation of the list of tasks that
// satisfy the given filter (tasks are included in the list if the taskFilter
// returns true).
func (journal TaskJournal) ListWithFilter(taskFilter TaskFilter) string {
	s := fmt.Sprintln()
	nlisted := 0
	for i := 0; i < len(journal.TaskList); i++ {
		task := journal.TaskList[i]
		if taskFilter(task) {
			s += fmt.Sprintf("%s\n", task.String())
			nlisted++
		}
	}
	if nlisted == 0 {
		s += fmt.Sprintf("%s\n\n", notasks)
	} else {
		s += fmt.Sprintf("\nLegend: %s  %s  %s\n", StatusTodo.legend(), StatusDoing.legend(), StatusDone.legend())
	}
	return s
}

// List returns a string representation of the list of all tasks (no filter)
func (journal TaskJournal) List() string {
	return journal.ListWithFilter(TaskFilterAll)
}

// Tree returns a string representation of the tree structure of tasks (parent relations)
func (journal TaskJournal) Tree() string {
	if len(journal.TaskList) == 0 {
		return fmt.Sprintf("\n%s\n\n", notasks)
	}
	s := ""
	tree := TreeString(journal.TaskList)
	if tree[0] != '\n' {
		// We add a line return for a pretty look
		s += fmt.Sprintln()
	}
	s += tree
	s += fmt.Sprintf("\nLegend: %s  %s  %s\n", StatusTodo.legend(), StatusDoing.legend(), StatusDone.legend())
	return s
}

func (journal TaskJournal) String() string {
	return journal.List()
}

// =========================================================================
// Implementation of the functions to edit task features

func (journal *TaskJournal) getNoteFile(uindex TaskID, create bool) (string, error) {
	task, err := journal.GetTask(uindex)
	if err != nil {
		return "", err
	}

	if task.NotePath == "" {
		if !create {
			return task.NotePath, nil
		}
		// initializes the NotePath with the default value (path relative to the
		// context root directory)
		basename := fmt.Sprintf("%d.rst", task.GIndex)
		task.NotePath = filepath.Join(NotebookDirname, basename)
	}

	var notepath string
	if filepath.IsAbs(task.NotePath) {
		notepath = task.NotePath
	} else {
		rootdir := filepath.Dir(journal.File())
		notepath = filepath.Join(rootdir, task.NotePath)
	}

	exists, err := PathExists(notepath)
	if exists && err != nil {
		return notepath, err
	}

	if !create {
		return notepath, nil
	}

	if !exists {
		err := CheckAndMakeDir(filepath.Dir(notepath))
		if err != nil {
			return notepath, err
		}
		file, err := os.Create(notepath)
		if err != nil {
			return notepath, err
		}
		defer file.Close()
		title := fmt.Sprintf("%.2d - %s", task.UIndex, task.Description)
		size := utf8.RuneCountInString(title) // nb runes, and not nb bytes (because of accents)
		line := ""
		for i := 0; i < size; i++ {
			line += "="
		}
		file.WriteString(fmt.Sprintf("%s\n", title))
		file.WriteString(fmt.Sprintf("%s\n", line))
		file.Sync()
	}

	return notepath, nil
}

// GetNoteFile returns the filepath to the note associated to this task. Returns
// a blank string ("") if no note is associated to this task.
func (journal *TaskJournal) GetNoteFile(uindex TaskID) (string, error) {
	return journal.getNoteFile(uindex, false)
}

// GetOrCreateNoteFile returns the filepath to the note associated to this task.
// It ensures that this note exists. If it is not defined, then the function
// creates it and return the absolute path to this note file.
func (journal *TaskJournal) GetOrCreateNoteFile(uindex TaskID) (string, error) {
	return journal.getNoteFile(uindex, true)
}

func (journal *TaskJournal) DeleteNoteFile(uindex TaskID) error {
	task, err := journal.GetTask(uindex)
	if err != nil {
		return err
	}

	if task.NotePath == "" {
		return errors.New("the task has no associated note")
	}

	var notepath string
	if filepath.IsAbs(task.NotePath) {
		notepath = task.NotePath
	} else {
		rootdir := filepath.Dir(journal.File())
		notepath = filepath.Join(rootdir, task.NotePath)
	}

	err = os.Remove(notepath)
	if err != nil {
		return err
	}
	task.NotePath = ""
	return nil
}

// ListNotes returns the list of all task notes as a single concatenated string
func (journal TaskJournal) ListNotes() string {
	listNotes := ""
	var notepath string
	for _, task := range journal.TaskList {

		if task.NotePath == "" {
			continue
		}

		if filepath.IsAbs(task.NotePath) {
			notepath = task.NotePath
		} else {
			rootdir := filepath.Dir(journal.File())
			notepath = filepath.Join(rootdir, task.NotePath)
		}
		content, err := LoadString(notepath)
		if err != nil {
			listNotes += err.Error() + "\n"
		} else {
			listNotes += content + "\n"
		}
	}
	return listNotes
}

// Report returns a complete report including (1) the task tree, (2) the task
// board and (3) the notes contents
func (journal TaskJournal) Report() string {
	report := journal.Tree()
	report += "\n------------------------------------------------------\n"
	report += "Board:\n"
	report += journal.ListWithFilter(TaskFilterOnBoard)
	report += "\n------------------------------------------------------\n"
	report += "Notes:\n\n"
	report += journal.ListNotes()
	return report
}

// =========================================================================
// Helper functions for testing purpose

// CreateTestJournal creates a dummy journal for test purposes
func CreateTestJournal() TaskJournal {
	journal := TaskJournal{
		TaskList: TaskArray{
			CreateTestTask(1, "Write documentation for todogo"),
			CreateTestTask(2, "Create unit test for todogo"),
			CreateTestTask(3, "Add a function to print a tasks journal"),
			CreateTestTask(4, "Organize a code review of todogo"),
		},
	}
	return journal
}
