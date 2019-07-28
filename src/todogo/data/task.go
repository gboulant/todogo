package data

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"todogo/core"
)

// =========================================================================
// Implementation of the Task concept

// Task is the data structure for a single task
type Task struct {
	UIndex      uint64     // Usage Index (could be recycled)
	GIndex      uint64     // General Index (invariant and unique)
	Timestamp   int64      // Date of the task (unix format)
	Description string     // Description of the Task
	Status      TaskStatus // Status of the task
	OnBoard     bool       // True if the task is on board
	NotePath    string     // Path to the note file (relative to the db root)
}

// initGlobalIndex initialises the global index of this task.
// We create a global index (unique and invariant ever) by creating a hash integer
// from a string representation of the task and its timestamp. The string
// representation is a composition of the usage id, the timestamp and the
// description. The objective is to make it impossible ever to have two
// tasks with the same global index.
func (task *Task) initGlobalIndex() {
	taskstr := fmt.Sprintf("%d [%d]: %s", task.UIndex, task.Timestamp, task.Description)
	task.GIndex = hashdate(taskstr, task.Timestamp)
}

// InitNotePath initializes the NotePath with the default value
func (task *Task) InitNotePath() {
	basename := fmt.Sprintf("%d.rst", task.GIndex)
	task.NotePath = filepath.Join(core.NotebookDirname, basename)
}

// String implements the stringable interface
func (task Task) String() string {
	dtlabel := datelabel(task.Timestamp)
	template := "%2d [%s] %s : %s"
	s := fmt.Sprintf(template, task.UIndex, dtlabel, task.Status.String(), task.Description)
	return s
}

// CreateTestTask creates a dummy task for test purposes
func CreateTestTask(uindex int, text string) Task {
	var task = Task{
		UIndex:      uint64(uindex),
		Description: text,
		Timestamp:   timestamp(),
		Status:      StatusTodo,
		OnBoard:     false,
	}
	task.initGlobalIndex()
	return task
}

// =========================================================================
// Implementation of the collection of Tasks TaskArray

// TaskArray is the data structure for a list (array) of Tasks
type TaskArray []Task

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// WARNING: TO BE DELETE

func (ta *TaskArray) Load(filepath string) error {
	bytes, err := core.LoadBytes(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, ta)
}

// It implements the jsonable interface.
func (ta *TaskArray) SaveTo(filepath string) error {
	bytes, err := json.MarshalIndent(*ta, core.JsonPrefix, core.JsonIndent)
	if err != nil {
		return err
	}
	return core.WriteBytes(filepath, bytes)
}

func (ta *TaskArray) File() string {
	return ""
}

func (ta *TaskArray) Save() error {
	return ta.SaveTo(ta.File())
}

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// =========================================================================
// Implementation of a tasks journal (with file persistence)

// TaskJournal defines the structure to manage a task journal. A tasks journal
// could be the current collection of tasks (called journal) or the archive
// collection of tasks (called archive).
type TaskJournal struct {
	TaskList TaskArray
	filepath string
}

// Load reads a journal of tasks from the given file.
//It implements the jsonable interface.
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

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

// TaskMap is the data structure for a map of indexed Tasks
type TaskMap map[uint64]Task

// taskArray2Map creates a map from the taskarray. The keys are the task indeces
//and value the corresponding task
// task index
func taskArray2Map(taskarray TaskArray) TaskMap {
	taskmap := make(TaskMap)
	var task Task
	for i := 0; i < len(taskarray); i++ {
		task = taskarray[i]
		taskmap[task.UIndex] = task
	}
	return taskmap
}

// sortedIndeces returns an ordered list of the indeces of the tasks contained
// in the specified TaskMap.
func sortedIndeces(taskmap TaskMap) []uint64 {
	indeces := make([]uint64, len(taskmap))
	i := 0
	for index := range taskmap {
		indeces[i] = index
		i++
	}
	sort.Slice(indeces, func(i, j int) bool { return indeces[i] < indeces[j] })
	return indeces
}

// taskMap2Array create an array from the taskmap. The task array is ordered by
// the task indeces
func taskMap2Array(taskmap TaskMap) TaskArray {
	// We first have to sort the map keys (task indeces)
	indeces := sortedIndeces(taskmap)
	// Then we create the array in this order
	taskarray := make(TaskArray, len(indeces))
	for i := 0; i < len(indeces); i++ {
		taskarray[i] = taskmap[indeces[i]]
	}
	return taskarray
}

// --------------------------------------------------------------
// Implementation of the Jsonable interface

// --------------------------------------------------------------
// Implementation of the Stringable interface
