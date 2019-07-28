package data

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"todogo/core"
)

const (
	noIndex = -1
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

type filterFunction func(task Task) bool

func (tasks TaskArray) index(filter filterFunction) int {
	for i := 0; i < len(tasks); i++ {
		if filter(tasks[i]) {
			return i
		}
	}
	return noIndex
}

func (tasks TaskArray) indeces(filter filterFunction) []int {
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
func (tasks TaskArray) GetTasksWithFilter(filter filterFunction) []*Task {
	results := make([]*Task, 0, len(tasks))
	idxlist := tasks.indeces(filter)
	for i := 0; i < len(idxlist); i++ {
		idx := idxlist[i]
		results = append(results, &tasks[idx])
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

// =========================================================================
// Implementation of a tasks journal (with file persistence)

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
