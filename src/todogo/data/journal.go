package data

// Implementation of a tasks journal (with file persistence)

import (
	"encoding/json"
	"fmt"
	"todogo/core"
)

// TaskJournal defines the structure to manage a task journal. A tasks journal
// could be the current collection of tasks (called journal) or the archive
// collection of tasks (called archive).
type TaskJournal struct {
	TaskList TaskArray
	filepath string
}

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
	return &task
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

// GetTasksWithFilter returns an array of pointer to the tasks that satisfy the
// given filter.
func (journal TaskJournal) GetTasksWithFilter(filter TaskFilter) []*Task {
	return journal.TaskList.getTasksWithFilter(filter)
}

// GetFreeUID returns the next free usage index in this journal
func (journal TaskJournal) GetFreeUID() TaskID {
	return journal.TaskList.getFreeUID()
}

// Load reads a journal of tasks from the given file. Returns an error if the
// file does not exist. Use LoadOrCreate to make sure to initialise a joournal
// whatever the starting situation (inn the case of the first usage of todo for
// example). It implements the jsonable interface.
func (journal *TaskJournal) Load(filepath string) error {
	bytes, err := core.LoadBytes(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, journal)
	if err == nil {
		journal.filepath = filepath
	}
	return err

}

// Add adds the given task to this journal
func (journal *TaskJournal) Add(task Task) error {
	return journal.TaskList.append(task)
}

// LoadOrCreate tries to load a journal from the given file, and create a void
// journal if the file does not exist.
func (journal *TaskJournal) LoadOrCreate(filepath string) error {
	exists, err := core.PathExists(filepath)
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
	bytes, err := json.MarshalIndent(journal, core.JSONPrefix, core.JSONIndent)
	if err != nil {
		return err
	}
	err = core.WriteBytes(filepath, bytes)
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
		s += fmt.Sprint("No tasks. Go have a drink\n\n")
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
		return fmt.Sprint("\nNo tasks. Go have a drink\n\n")
	}
	s := fmt.Sprintln()
	s += TreeString(journal.TaskList)
	s += fmt.Sprintf("\nLegend: %s  %s  %s\n", StatusTodo.legend(), StatusDoing.legend(), StatusDone.legend())
	return s
}

func (journal TaskJournal) String() string {
	return journal.List()
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
