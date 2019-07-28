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
	journal.TaskList.Append(task)
	return &task
}

// Delete removes the task with the specified id. Returns a copy of the deleted
// task on success
func (journal *TaskJournal) Delete(uindex uint64) (Task, error) {
	var task Task
	index := journal.TaskList.indexFromUID(uindex)
	if index == noIndex {
		return task, fmt.Errorf("ERR: The task %d does not exist", uindex)
	}
	task = journal.TaskList[index]
	err := journal.TaskList.Remove(index)
	return task, err
}

func (journal TaskJournal) GetTask(uindex uint64) (*Task, error) {
	return journal.TaskList.GetTask(uindex)
}

func (journal TaskJournal) GetFreeUID() uint64 {
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

func (journal *TaskJournal) Add(task Task) error {
	return journal.TaskList.Append(task)
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
	bytes, err := json.MarshalIndent(journal, core.JsonPrefix, core.JsonIndent)
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

// List prints all tasks (no filter)
func (journal TaskJournal) List() string {
	return journal.ListWithFilter(TaskFilterAll)
}

func (journal TaskJournal) String() string {
	return journal.List()
}

// AddOnBoard adds the specified task on board
func (journal *TaskJournal) AddOnBoard(uindex uint64) error {
	task, err := journal.TaskList.GetTask(uindex)
	if err != nil {
		return err
	}
	task.OnBoard = true
	return nil
}

// RemoveFromBoard removes the specified task from board
func (journal *TaskJournal) RemoveFromBoard(uindex uint64) error {
	task, err := journal.TaskList.GetTask(uindex)
	if err != nil {
		return err
	}
	task.OnBoard = false
	return nil
}

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
