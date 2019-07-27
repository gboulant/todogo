package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"todogo/core"
)

// Database holds the data describing the tasks
type Database struct {
	taskmap  TaskMap
	filepath string
}

var dbdirpath = core.GetActiveContextPath()

// JournalPath points to the default journal task file
var JournalPath = filepath.Join(dbdirpath, "journal.json")

// ArchivePath points to the default archive task file
var ArchivePath = filepath.Join(dbdirpath, "archive.json")

// Init loads the data files and initialize the memory database
func (db *Database) Init(taskFilePath string) error {
	exists, err := core.PathExists(taskFilePath)
	if exists && err != nil {
		return err
	}

	if !exists {
		db.taskmap = make(TaskMap)
	} else {
		var tasklist TaskArray
		err = core.Load(taskFilePath, &tasklist)
		if err != nil {
			return err
		}
		db.taskmap = taskArray2Map(tasklist)
	}
	db.filepath = taskFilePath
	return nil // not an error
}

// FreeUsageIndex returns the first free usage index of this database for
// initializing the usage index of a newly created task
func (db Database) FreeUsageIndex() uint64 {
	if len(db.taskmap) == 0 {
		return 1
	}
	indeces := sortedIndeces(db.taskmap)
	return core.FreeIndex(indeces)
}

// New creates a new task in the database
func (db *Database) New(text string) Task {
	var task = Task{
		UIndex:      uint64(db.FreeUsageIndex()),
		Description: text,
		Timestamp:   timestamp(),
		Status:      StatusTodo,
		OnBoard:     false,
	}
	task.initGlobalIndex()
	db.taskmap[task.UIndex] = task
	return task
}

// Delete removes the task with the specified id. Returns a copy of the deleted
// task on success
func (db *Database) Delete(taskIndex uint64) (Task, error) {
	task, exists := db.taskmap[taskIndex]
	if !exists {
		msg := fmt.Sprintf("WRN: The task %d does not exist", taskIndex)
		return task, errors.New(msg)
	}
	delete(db.taskmap, taskIndex)
	return task, nil
}

// List prints all tasks (no filter)
func (db Database) List() {
	db.ListWithFilter(TaskFilterAll)
}

// ListWithFilter prints the list of tasks restricted to the given filter
func (db Database) ListWithFilter(taskFilter TaskFilter) {
	fmt.Println()
	indeces := sortedIndeces(db.taskmap)
	nlisted := 0
	for i := 0; i < len(indeces); i++ {
		task := db.taskmap[indeces[i]]
		if taskFilter(task) {
			core.Println(task)
			nlisted++
		}
	}
	if nlisted == 0 {
		fmt.Printf("No tasks. Go have a drink\n\n")
	} else {
		fmt.Printf("\nLegend: %s  %s  %s\n", StatusTodo.legend(), StatusDoing.legend(), StatusDone.legend())
	}
}

// GetIndeces returns the list of indeces satisfaying the given filter
func (db Database) GetIndeces(taskFilter TaskFilter) []uint64 {
	indeces := sortedIndeces(db.taskmap)
	buffer := ""
	for i := 0; i < len(indeces); i++ {
		task := db.taskmap[indeces[i]]
		if taskFilter(task) {
			buffer += fmt.Sprintf("%d,", task.UIndex)
		}
	}
	if buffer == "" {
		return []uint64{}
	}
	buffer = strings.Trim(buffer, ",")
	indexList := core.IndexList{}
	indexList.Set(buffer)
	return indexList
}

// Get returns a copy of the specified task
func (db Database) Get(taskIndex uint64) (Task, error) {
	_, exists := db.taskmap[taskIndex]
	if !exists {
		msg := fmt.Sprintf("ERR: The task %d does not exist", taskIndex)
		return Task{}, errors.New(msg)
	}
	return db.taskmap[taskIndex], nil
}

// Set replaces the specified task by the given task
func (db *Database) Set(taskIndex uint64, task Task) error {
	_, exists := db.taskmap[taskIndex]
	if !exists {
		msg := fmt.Sprintf("ERR: The task %d does not exist", taskIndex)
		return errors.New(msg)
	}
	if task.UIndex != taskIndex {
		msg := fmt.Sprintf("ERR: the task index is %d (should be %d)", task.UIndex, taskIndex)
		return errors.New(msg)
	}
	db.taskmap[taskIndex] = task
	return nil
}

// Add adds the specified task to the database. It is up to you to check before
// that the usage index is consistent (as for the global index).
func (db *Database) Add(task Task) error {
	_, exists := db.taskmap[task.UIndex]
	if exists {
		msg := fmt.Sprintf("ERR: The task %d already exist", task.UIndex)
		return errors.New(msg)
	}
	db.taskmap[task.UIndex] = task
	return nil
}

// AddOnBoard adds the specified task on board
func (db *Database) AddOnBoard(taskIndex uint64) error {
	task, err := db.Get(taskIndex)
	if err != nil {
		return err
	}
	task.OnBoard = true
	return db.Set(taskIndex, task)
}

// RemoveFromBoard removes the specified task from board
func (db *Database) RemoveFromBoard(taskIndex uint64) error {
	task, err := db.Get(taskIndex)
	if err != nil {
		return err
	}
	task.OnBoard = false
	return db.Set(taskIndex, task)
}

// Commit saves the memory database to the data files
func (db Database) Commit() error {
	tasklist := taskMap2Array(db.taskmap)
	return core.Save(db.filepath, &tasklist)
}
