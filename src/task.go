package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

// =========================================================================
// Implementation of the TaskID concept.
//
// The task index (or identifier) is an positive integer that identifies a task
// in its context.

// TaskID is the data type of a task index (Usage ID or General ID)
type TaskID uint64

// TaskIDArray is a list of TaskID
type TaskIDArray []TaskID

const (
	// noIndex is used to specify that there is no array index (whatever the array is)
	noIndex int = -1
	// NoUID is used to specify that there is no task index (task identifier)
	NoUID TaskID = 0
)

//
// The task indeces are used one the command lines to specify the target of
// actions, then we give here an implementation of the flag.Value interface for
// a list of task indeces.
//
func (taskID *TaskID) String() string {
	return fmt.Sprintf("%d", *taskID)
}

// Set implement the flag.Value interface
func (taskID *TaskID) Set(value string) error {
	index, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	(*taskID) = TaskID(index)
	return nil
}

// String implement the flag.Value interface
func (il *TaskIDArray) String() string {
	return fmt.Sprintf("%v", *il)
}

// Set implement the flag.Value interface
func (il *TaskIDArray) Set(value string) error {
	sl := strings.Split(value, ",")
	*il = make(TaskIDArray, len(sl))
	for i := 0; i < len(sl); i++ {
		index, err := strconv.ParseUint(sl[i], 10, 64)
		if err != nil {
			return err
		}
		(*il)[i] = TaskID(index)
	}
	return nil
}

// =========================================================================
// Implementation of the Task concept

// Task is the data structure for a single task
type Task struct {
	UIndex      TaskID     // Usage Index (could be recycled)
	GIndex      TaskID     // Global Index (invariant and unique)
	Timestamp   int64      // Date of the task (unix format)
	Description string     // Description of the Task
	Status      TaskStatus // Status of the task
	OnBoard     bool       // True if the task is on board
	NotePath    string     // Path to the note file (relative to the db root)
	ParentID    TaskID     // UID of the parent task
}

// initGlobalIndex initialises the global index of this task.
// We create a global index (unique and invariant ever) by creating a hash integer
// from a string representation of the task and its timestamp. The string
// representation is a composition of the usage id, the timestamp and the
// description. The objective is to make it impossible ever to have two
// tasks with the same global index.
func (task *Task) initGlobalIndex() {
	taskstr := fmt.Sprintf("%d [%d]: %s", task.UIndex, task.Timestamp, task.Description)
	task.GIndex = TaskID(hashdate(taskstr, task.Timestamp))
}

// String returns a string representation of this task
func (task Task) String() string {
	return task.OnelineString()
}

// OnelineString returns a string representation of this task on one signe line.
// This shouldbe used for a pretty presentation of task lists.
func (task Task) OnelineString() string {
	t := "%2d %s %s : %s"
	s := fmt.Sprintf(t, task.UIndex, task.getTaskIndicators(), task.Status.String(), task.Description)
	return s
}

func (task Task) getTaskIndicators() string {
	cfg, _ := GetConfig() // unused to test the err, we can not arrive here in case of config error
	indicatorsTemplate := cfg.Parameters.Indicators
	tmpl, err := template.New("indicators").Parse(indicatorsTemplate)
	if err != nil {
		tmpl, _ = template.New("indicators").Parse(DefaultIndicatorsTemplate)
	}

	dtlabel := datelabel(task.Timestamp)

	var hasNote string
	if task.NotePath != "" {
		hasNote = "n"
	} else {
		hasNote = "-"
	}

	var onBoard string
	if task.OnBoard {
		onBoard = "b"
	} else {
		onBoard = "-"
	}

	indicators := struct {
		Date  string
		Note  string
		Board string
	}{
		Date:  dtlabel,
		Note:  hasNote,
		Board: onBoard,
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, indicators)
	if err != nil {
		return "no indicators"
	}
	return buffer.String()
}

// JSONString returns a json string representation of this task
func (task Task) JSONString() string {
	bytes, err := json.MarshalIndent(task, JSONPrefix, JSONIndent)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// InfoString returns a string representation of the task attributes
func (task Task) InfoString() string {
	return task.JSONString()
}

// CreateTestTask creates a dummy task for test purposes
func CreateTestTask(uindex TaskID, text string) Task {
	var task = Task{
		UIndex:      uindex,
		Description: text,
		Timestamp:   timestamp(),
		Status:      StatusTodo,
		OnBoard:     false,
	}
	task.initGlobalIndex()
	return task
}

// =========================================================================
// Implementation of the collection of Tasks TaskArray. The TaskArray is the
// underlying data structure of a task journal. It should not be exposed output
// from the data package.

// TaskArray is the data structure for a list (array) of Tasks
type TaskArray []Task

// String implements the stringable interface for a TaskArray
func (tasks TaskArray) String() string {
	s := ""
	for i := 0; i < len(tasks); i++ {
		s += fmt.Sprintf("%s\n", tasks[i].String())
	}
	return s
}

// Remove removes from the array the task of order index (index in the array)
func (tasks *TaskArray) remove(index int) error {
	if index < 0 || index >= len(*tasks) {
		return fmt.Errorf("ERR: index %d is out of range of tasks", index)
	}
	(*tasks)[index] = (*tasks)[len(*tasks)-1]
	*tasks = (*tasks)[:len(*tasks)-1]
	return nil
}

// Append adds the task to the array. Returns an error if a task with same uid exists
func (tasks *TaskArray) append(task Task) error {
	ptask, err := tasks.getTask(task.UIndex)
	if ptask != nil && err == nil {
		return fmt.Errorf("ERR: a task with UID %d already exists", task.UIndex)
	}
	*tasks = append(*tasks, task)
	return nil
}

func (tasks TaskArray) index(filter TaskFilter) int {
	for i := 0; i < len(tasks); i++ {
		if filter(tasks[i]) {
			return i
		}
	}
	return noIndex
}

func (tasks TaskArray) indeces(filter TaskFilter) []int {
	results := make([]int, 0, len(tasks))
	for i := 0; i < len(tasks); i++ {
		if filter(tasks[i]) {
			results = append(results, i)
		}
	}
	return results
}

// IndexFromUID returns the array index of the task with the given UID
func (tasks TaskArray) indexFromUID(uindex TaskID) int {
	filterUID := func(task Task) bool {
		return task.UIndex == uindex
	}
	return tasks.index(filterUID)
}

// GetTask returns a pointer to the task of the sepcified UID
func (tasks *TaskArray) getTask(uindex TaskID) (*Task, error) {
	idx := tasks.indexFromUID(uindex)
	if idx == noIndex {
		err := fmt.Errorf("The task of index %d does not exist", uindex)
		return nil, err
	}
	return &(*tasks)[idx], nil
}

// GetTasksWithFilter returns an array of pointer to the tasks that satisfy the filter
func (tasks TaskArray) getTasksWithFilter(filter TaskFilter) []*Task {
	results := make([]*Task, 0, len(tasks))
	for i := 0; i < len(tasks); i++ {
		if filter(tasks[i]) {
			results = append(results, &tasks[i])
		}
	}
	return results
}

func (tasks TaskArray) byUID(i int, j int) bool {
	return tasks[i].UIndex < tasks[j].UIndex
}
func (tasks TaskArray) byGID(i int, j int) bool {
	return tasks[i].GIndex < tasks[j].GIndex
}
func (tasks TaskArray) byTimestamp(i int, j int) bool {
	return tasks[i].Timestamp < tasks[j].Timestamp
}

func (tasks *TaskArray) sortByUID() {
	sort.Slice(*tasks, tasks.byUID)
}

func (tasks *TaskArray) sortByGID() {
	sort.Slice(*tasks, tasks.byGID)
}

func (tasks *TaskArray) sortByTimestamp() {
	sort.Slice(*tasks, tasks.byTimestamp)
}

// getFreeUID() returns the first free UID of this tasks list.
func (tasks TaskArray) getFreeUID() TaskID {
	// The free index is determined with the hypothesis that the indeces array
	// is a list of consecutive integer indeces. If the difference between two
	// consecutif indeces is not 1, then it means that there is at least a free
	// index (the index that follows the smallest index of the difference).
	tasks.sortByUID()
	if len(tasks) == 0 {
		return 1
	}
	var freeUID TaskID = 1
	for i := 0; i < len(tasks); i++ {
		if tasks[i].UIndex-freeUID > 0 {
			return freeUID
		}
		freeUID = tasks[i].UIndex + 1
	}
	return freeUID
}

// ancestor returns true if parentId is an ancestor of childID
func (tasks TaskArray) ancestor(childID TaskID, parentID TaskID) bool {
	if childID == parentID {
		return false
	}
	task, _ := tasks.getTask(childID)
	if task.ParentID == NoUID {
		return false
	}
	if task.ParentID == parentID {
		return true
	}
	return tasks.ancestor(task.ParentID, parentID)
}
