package data

import (
	"fmt"
	"path/filepath"
	"sort"
	"todogo/conf"
)

const (
	noIndex = -1
	noUID   = 0
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
	ParentID    uint64     // UID of the parent task
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
	task.NotePath = filepath.Join(conf.NotebookDirname, basename)
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

// String implements the stringable interface for a TaskArray
func (tasks TaskArray) String() string {
	s := ""
	for i := 0; i < len(tasks); i++ {
		s += fmt.Sprintf("%s\n", tasks[i].String())
	}
	return s
}

// Remove removes from the array the task of order index (index in the array)
func (tasks *TaskArray) Remove(index int) error {
	if index < 0 || index >= len(*tasks) {
		return fmt.Errorf("ERR: index %d is out of range of tasks", index)
	}
	(*tasks)[index] = (*tasks)[len(*tasks)-1]
	*tasks = (*tasks)[:len(*tasks)-1]
	return nil
}

// Append adds the task to the array. Returns an error if a task with same uid exists
func (tasks *TaskArray) Append(task Task) error {
	ptask, err := tasks.GetTask(task.UIndex)
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
func (tasks TaskArray) indexFromUID(uindex uint64) int {
	filterUID := func(task Task) bool {
		return task.UIndex == uindex
	}
	return tasks.index(filterUID)
}

// GetTask returns a pointer to the task of the sepcified UID
func (tasks *TaskArray) GetTask(uindex uint64) (*Task, error) {
	idx := tasks.indexFromUID(uindex)
	if idx == noIndex {
		err := fmt.Errorf("The task of index %d does not exist", uindex)
		return nil, err
	}
	return &(*tasks)[idx], nil
}

// GetTasksWithFilter returns an array of pointer to the tasks that satisfy the filter
func (tasks TaskArray) GetTasksWithFilter(filter TaskFilter) []*Task {
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

func (tasks *TaskArray) SortByUID() {
	sort.Slice(*tasks, tasks.byUID)
}

func (tasks *TaskArray) SortByGID() {
	sort.Slice(*tasks, tasks.byGID)
}

func (tasks *TaskArray) SortByTimestamp() {
	sort.Slice(*tasks, tasks.byTimestamp)
}

// getFreeUID() returns the first free UID of this tasks list.
func (tasks TaskArray) getFreeUID() uint64 {
	// The free index is determined with the hypothesis that the indeces array
	// is a list of consecutive integer indeces. If the difference between two
	// consecutif indeces is not 1, then it means that there is at least a free
	// index (the index that follows the smallest index of the difference).
	tasks.SortByUID()
	if len(tasks) == 0 {
		return 1
	}
	var freeUID uint64 = 1
	for i := 0; i < len(tasks); i++ {
		if tasks[i].UIndex-freeUID > 0 {
			return freeUID
		}
		freeUID = tasks[i].UIndex + 1
	}
	return freeUID
}

// ancestor returns true if parentId is an ancestor of childID
func (tasks TaskArray) ancestor(childID uint64, parentID uint64) bool {
	if childID == parentID {
		return false
	}
	task, _ := tasks.GetTask(childID)
	if task.ParentID == noUID {
		return false
	}
	if task.ParentID == parentID {
		return true
	}
	return tasks.ancestor(task.ParentID, parentID)
}
